package getty

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

func (s *Server) HandleAuthByGoogle(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	CheckError(err)
	//token := data["token"]
	//uuid := data["uuid"]
	//version := int(data["version"].(float64))
	var jData []byte
	jData, err = json.Marshal(struct {
		State int
		Token string
	}{123, "sdafdsf"})
	CheckError(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
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
