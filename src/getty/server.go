package getty

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"godori.com/db"
	lType "godori.com/util/constant/loginType"
	cFilter "godori.com/util/filter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Claims struct {
	Token string
	jwt.StandardClaims
}

type Server struct {
	Host         string
	Clients      map[string]*Client
	PacketChan   chan *Message
	ConnChan     chan *Client
	DisConnChan  chan *Client
	OnConnect    func(*Client)
	OnMessage    func(*Client, *Data)
	OnDisconnect func(*Client)
	BeforeAccept func() bool
}

const VERSION = 4

const (
	CallBackURL         = "http://localhost:50001/verify/google" //auth/callback"
	UserInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	ScopeEmail          = "https://www.googleapis.com/auth/userinfo.email"
	ScopeProfile        = "https://www.googleapis.com/auth/userinfo.profile"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var expirationTime = 5 * time.Minute
var JwtKey = []byte("apple")
var OAuthConf *oauth2.Config

func init() {
	OAuthConf = &oauth2.Config{
		ClientID:     "112494846092-ar8ml4nm16mr7bhd3cekb87846fr5k0e.apps.googleusercontent.com",
		ClientSecret: "PYpSubbdJdIbzSUy4mqxcpVf",
		RedirectURL:  CallBackURL,
		Scopes:       []string{ScopeEmail, ScopeProfile},
		Endpoint:     google.Endpoint,
	}
}

func NewServer(h string) *Server {
	server := &Server{Host: h}
	server.OnConnect = func(*Client) {}
	server.OnMessage = func(*Client, *Data) {}
	server.OnDisconnect = func(*Client) {}
	server.BeforeAccept = func() bool { return true }
	server.Clients = make(map[string]*Client)
	server.PacketChan = make(chan *Message)
	server.ConnChan = make(chan *Client)
	server.DisConnChan = make(chan *Client)
	return server
}

func (s *Server) onConnect(c *Client) {
	key := c.RemoteAddr().String()
	if client, ok := s.Clients[key]; ok {
		client.Close()
	}
	s.Clients[key] = c
	c.Handle()
	s.OnConnect(c)
}

func (s *Server) onDisconnect(c *Client) {
	c.Run = false
	s.OnDisconnect(c)
	delete(s.Clients, c.RemoteAddr().String())
	c.Close()
}

func (s *Server) onMessage(c *Client, d *Data) {
	s.OnMessage(c, d)
}

func GetJwtToken(token string) string {
	expirationTime := time.Now().Add(expirationTime)
	claims := &Claims{
		Token: token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	verify := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	verifyString, err := verify.SignedString(JwtKey)
	CheckError(err)
	return verifyString
}

func (s *Server) ParseJwtToken(receivedToken string) string {
	token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return JwtKey, nil
	})
	CheckError(err)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["Token"].(string)
	}
	return ""
}

func VerifyByGoogle(token string) []byte {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token)
	CheckError(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	return body
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.FormValue("token")
	name := r.FormValue("name")
	//recommend := r.FormValue("recommend")
	var state string
	verify := strings.Split(s.ParseJwtToken(token), " ")
	if len(verify) < 2 {
		state = "FAILED"
	} else {
		loginType, err := strconv.Atoi(verify[0])
		CheckError(err)
		uid := verify[1]
		if u, ok := db.GetUserByOAuth(uid, loginType); ok {
			nameLen := utf8.RuneCountInString(name)
			if u.Verify.Int32 == 1 {
				state = "FAILED"
			} else if nameLen < 1 || nameLen > 6 {
				state = "FAILED"
			} else if match, _ := regexp.MatchString("[^가-힣]", name); match {
				state = "FAILED"
			} else if cFilter.Check(name) {
				state = "UNAVAILABLE_NAME"
			} else {
				if _, ok := db.GetUserByName(name); ok {
					state = "RE_REQUEST"
				} else {
					go db.UpdateUserVerify(name, uid, loginType)
					state = "LOGIN_SUCCESS"
				}
			}
		}
	}
	fmt.Fprint(w, state)
}

func (s *Server) GetEnvValue(key string) string {
	err := godotenv.Load()
	CheckError(err)
	return os.Getenv(key)
}

func (s *Server) HandleAuthByGoogle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.FormValue("token")
	//uuid := r.FormValue("uuid")
	regex := regexp.MustCompile("[^0-9]")
	version, err := strconv.Atoi(regex.ReplaceAllString(r.FormValue("version"), ""))
	CheckError(err)
	var state, verify string
	if version < VERSION {
		state = "NOT_UPDATED"
	} else {
		body := VerifyByGoogle(token)
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		CheckError(err)
		if clientId := s.GetEnvValue("GOOGLE_CLIENT_ID"); clientId != data["aud"].(string) {
			return
		}
		uid := data["sub"].(string)
		verify = GetJwtToken(strconv.Itoa(lType.GOOGLE) + " " + uid)
		if u, ok := db.GetUserByOAuth(uid, lType.GOOGLE); ok {
			if u.Verify.Int32 == 1 {
				state = "LOGIN_SUCCESS"
			} else {
				state = "REGISTER_SUCCESS"
			}
		} else {
			go db.InsertUser(uid, lType.GOOGLE)
			state = "REGISTER_SUCCESS"
		}
	}
	jsonData, err := json.Marshal(struct {
		State string
		Token string
	}{state, verify})
	CheckError(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) Listen(w http.ResponseWriter, r *http.Request) {
	fork := NewServer(s.Host)
	fork.OnMessage = s.OnMessage
	fork.OnConnect = s.OnConnect
	fork.OnDisconnect = s.OnDisconnect
	conn, err := upgrader.Upgrade(w, r, nil)
	CheckError(err)
	defer conn.Close()
	var wg sync.WaitGroup
	go func() {
		for {
			wg.Add(-1)
			if !s.BeforeAccept() {
				continue
			}
			token := r.URL.Query().Get("token")
			fork.ConnChan <- NewClient(conn, fork, token)
			wg.Wait()
		}
	}()
	for {
		select {
		case conn := <-fork.ConnChan:
			fork.onConnect(conn)
			defer wg.Done()
		case disconn := <-fork.DisConnChan:
			fork.onDisconnect(disconn)
		case packet := <-fork.PacketChan:
			fork.onMessage(packet.Client, packet.Data)
		}
	}
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}
