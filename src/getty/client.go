package getty

import (
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const HEADER_SIZE = 2

type Client struct {
	conn     *websocket.Conn
	server   *Server
	token    string
	sendChan chan []byte
	Run      bool
}

func NewClient(c *websocket.Conn, s *Server, t string) *Client {
	return &Client{
		conn:     c,
		server:   s,
		token:    t,
		sendChan: make(chan []byte),
		Run:      true,
	}
}

func BytesToInt(b []byte) int {
	var n int
	addr := uint((len(b) - 1) * 8)
	for i, _ := range b {
		n += int(b[i]) << addr
		addr -= 8
	}
	return n
}

func (c *Client) GetToken() string {
	return c.token
}

func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) Send(d []byte) {
	if c.Run {
		c.sendChan <- d
	}
}

func (c *Client) Broadcast(d []byte) {
	for _, ic := range c.server.Clients {
		ic.Send(d)
	}
}

func (c *Client) BroadcastAnother(d []byte) {
	for _, ic := range c.server.Clients {
		if c == ic {
			continue
		}
		ic.Send(d)
	}
}

func (c *Client) Close() {
	c.server = nil
	c.conn.Close()
	c.conn = nil
}

func (c *Client) Request() {
	defer func() {
		c.server.DisConnChan <- c
	}()
	for c.Run {
		_, message, err := c.conn.ReadMessage()
		if e, ok := err.(*websocket.CloseError); ok {
			switch e.Code {
			case 1001, 1005, 1006:
				return
			default:
				log.Println(e)
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		pSize := len(message)
		pType := BytesToInt(message[:HEADER_SIZE])
		c.server.PacketChan <- &Message{c, &Data{pType, message[HEADER_SIZE:pSize]}}
	}
}

func (c *Client) Response() {
	//ticker := time.NewTicker(3 * time.Second)
	//defer func() {
	//	ticker.Stop()
	//}()
	for c.Run {
		data := <-c.sendChan
		log.Println(string(data))
		err := c.conn.WriteMessage(websocket.TextMessage, data)
		CheckError(err)
		//case tick := <-ticker.C:
		//	log.Println("ping:", c.RemoteAddr(), tick.Second())
		//	err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
		//	CheckError(err)
	}
}

func (c *Client) Handle() {
	go c.Request()
	go c.Response()
}
