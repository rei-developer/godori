package getty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
	"unicode/utf8"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"godori.com/db"
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

const VERSION = 123

// auth.go
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

func VerifyByGoogle(token string) {

}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	CheckError(err)
	verify := GetJwtToken(data["token"].(string))
	if u, ok := db.GetUserByOAuth(verify, 1); ok {
		if u.Verify.Int32 == 1 {
			fmt.Fprint(w, "FAILED")
			return
		}
		name := data["name"].(string)
		nameLen := utf8.RuneCountInString(name)
		if nameLen < 1 || nameLen > 6 {
			fmt.Fprint(w, "FAILED")
			return
		}
		if match, _ := regexp.MatchString("[^가-힣]", name); match {
			fmt.Fprint(w, "FAILED")
			return
		}
		if cFilter.Check(name) {
			fmt.Fprint(w, "UNAVAILABLE_NAME")
			return
		}
		go db.UpdateUserVerify(name, verify, 1)
		fmt.Fprint(w, "LOGIN_SUCCESS")
	}
}

func (s *Server) HandleAuthByGoogle(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}

	fmt.Println(body)
	fmt.Println(string(body))

	err := json.Unmarshal(body, &data)
	CheckError(err)
	// TODO : version check
	// TODO : check block
	token, err := OAuthConf.Exchange(oauth2.NoContext, data["token"].(string))
	CheckError(err)
	client := OAuthConf.Client(oauth2.NoContext, token)
	// UserInfoAPIEndpoint는 유저 정보 API URL을 담고 있음
	userInfoResp, err := client.Get(UserInfoAPIEndpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer userInfoResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(userInfo)
	fmt.Println(string(userInfo))
	//var authUser User
	//json.Unmarshal(userInfo, &authUser)

	//token := data["token"]
	//uuid := data["uuid"]
	//version := int(data["version"].(float64))
	//var jData []byte
	//jData, err = json.Marshal(struct {
	//	State int
	//	Token string
	//}{123, "sdafdsf"})
	//CheckError(err)
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(jData)
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
			wg.Add(1)
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
		panic(err)
	}
}
