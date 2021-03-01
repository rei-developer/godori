package getty

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"time"
)

type Client struct {
	conn     *websocket.Conn
	server   *Server
	sendChan chan string
	done     chan struct{}
	buffer   []byte
}

func NewClient(conn *websocket.Conn, server *Server) *Client {
	return &Client{
		conn:     conn,
		server:   server,
		sendChan: make(chan string),
		done:     make(chan struct{}),
	}
}

func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) Send(data []byte) {
	c.sendChan <- string(data)
}

func (c *Client) Broadcast(data []byte) {
	for _, ic := range c.server.clients {
		ic.Send(data)
	}
}

func (c *Client) BroadcastAnother(data []byte) {
	for _, ic := range c.server.clients {
		if c == ic {
			continue
		}
		ic.Send(data)
	}
}

func (c *Client) Close() {
	close(c.done)
	c.server = nil
	c.conn.Close()
	c.conn = nil
}

func (c *Client) request() {
	defer func() {
		c.server.disconnChan <- c
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if c, k := err.(*websocket.CloseError); k {
				if c.Code == 1005 {
					return
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			var data map[string]interface{}
			json.Unmarshal(message, &data)
			c.server.packetChan <- &Message{c, &Data{data["type"].(float64), data["data"].(interface{})}}
		}
	}
}

func (c *Client) response() {
	for {
		select {
		case <-c.done:
			return
		case data := <-c.sendChan:
			log.Println(data, "response입니다")
			//c.writer.WriteString(data)
			//c.writer.Flush()
		}
	}
}

func (c *Client) Handle() {
	go c.request()
	go c.response()
}
