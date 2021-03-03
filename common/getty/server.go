package getty

import (
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
	host         string
	clients      map[string]*Client
	packetChan   chan *Message
	connChan     chan *Client
	disconnChan  chan *Client
	OnConnect    func(*Client)
	OnMessage    func(*Client, *Data)
	OnDisconnect func(*Client)
	BeforeAccept func() bool
}

func NewServer(host string) *Server {
	server := &Server{host: host}

	server.OnConnect = func(*Client) {}
	server.OnMessage = func(*Client, *Data) {}
	server.OnDisconnect = func(*Client) {}
	server.BeforeAccept = func() bool { return true }

	server.clients = make(map[string]*Client)
	server.packetChan = make(chan *Message)
	server.connChan = make(chan *Client)
	server.disconnChan = make(chan *Client)
	return server
}

func (s *Server) onConnect(c *Client) {
	key := c.RemoteAddr().String()
	if client, ok := s.clients[key]; ok {
		client.Close()
	}
	s.clients[key] = c
	c.Handle()

	s.OnConnect(c)
}

func (s *Server) onDisconnect(c *Client) {
	s.OnDisconnect(c)
	delete(s.clients, c.RemoteAddr().String())
	c.Close()
}

func (s *Server) onMessage(c *Client, data *Data) {
	s.OnMessage(c, data)
}

func (s *Server) Listen(w http.ResponseWriter, r *http.Request) {
	fork := NewServer(s.host)
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
			if s.BeforeAccept() == false {
				continue
			}
			token := r.URL.Query().Get("token")
			fork.connChan <- NewClient(conn, fork, token)
			wg.Wait()
		}
	}()

	for {
		select {
		case conn := <-fork.connChan:
			fork.onConnect(conn)
			defer wg.Done()
		case disconn := <-fork.disconnChan:
			fork.onDisconnect(disconn)
		case packet := <-fork.packetChan:
			fork.onMessage(packet.Client, packet.Data)
		}
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
