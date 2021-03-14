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
	Lock         sync.RWMutex
}

const VERSION = 4

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var expirationTime = 5 * time.Minute
var JwtKey = []byte("apple")

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
	fmt.Println(key, " 접속 시도")
	if client, ok := s.Clients[key]; ok {
		s.onDisconnect(client)
	}
	s.Clients[key] = c
	c.Handle()
	s.OnConnect(c)
}

func (s *Server) onDisconnect(c *Client) {
	key := c.RemoteAddr().String()
	c.Run = false
	s.OnDisconnect(c)
	fmt.Println(key, " 입니다용 키값은!")
	delete(s.Clients, key)
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
	if !s.BeforeAccept() {
		return
	}
	token := r.URL.Query().Get("token")
	go func() {
		s.Lock.Lock()
		fork.ConnChan <- NewClient(conn, fork, token)
		s.Lock.Unlock()
	}()
	for {
		select {
		case conn := <-fork.ConnChan:
			fork.onConnect(conn)
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
