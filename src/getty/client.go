package getty

import (
	"fmt"
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
	sendChan chan string
	done     chan struct{}
}

func NewClient(c *websocket.Conn, s *Server, t string) *Client {
	return &Client{
		conn:     c,
		server:   s,
		token:    t,
		sendChan: make(chan string),
		done:     make(chan struct{}),
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
	c.sendChan <- string(d)
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
	close(c.done)
	c.server = nil
	c.conn.Close()
	c.conn = nil
}

func (c *Client) Request() {
	defer func() {
		c.server.DisConnChan <- c
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if e, ok := err.(*websocket.CloseError); ok {
				switch e.Code {
				case 1001, 1005:
					return
				default:
					fmt.Println(e)
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			pSize := len(message)
			pType := BytesToInt(message[:HEADER_SIZE])
			c.server.PacketChan <- &Message{c, &Data{pType, message[HEADER_SIZE:pSize]}}
			//far := []byte{2,166}
			//fmt.Println(far)
			//fmt.Println(BytesToInt(far))

			//var data map[string]interface{}
			//json.Unmarshal(message, &data)
		}
	}
}

func (c *Client) Response() {
	for {
		select {
		case <-c.done:
			return
		case data := <-c.sendChan:
			log.Println(data, "Response입니다")
			//c.writer.WriteString(data)
			//c.writer.Flush()
		}
	}
}

func (c *Client) Handle() {
	go c.Request()
	go c.Response()
}
