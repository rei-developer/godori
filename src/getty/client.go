package getty

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"sync"
)

type Client struct {
	Conn     *websocket.Conn
	Server   *Server
	Token    string
	SendChan chan []byte
	Run      chan bool
	Lock     sync.RWMutex
}

const HEADER_SIZE = 2

var isDev = false

func init() {
	err := godotenv.Load()
	CheckError(err)
	mode := os.Getenv("MODE")
	if mode == "dev" {
		isDev = true
	}
}

func NewClient(c *websocket.Conn, s *Server, t string) *Client {
	return &Client{
		Conn:     c,
		Server:   s,
		Token:    t,
		SendChan: make(chan []byte),
		Run:      make(chan bool),
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
	return c.Token
}

func (c *Client) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Client) Send(d []byte) {
	c.SendChan <- d
}

func (c *Client) Broadcast(d []byte) {
	for _, ic := range c.Server.Clients {
		ic.Send(d)
	}
}

func (c *Client) BroadcastAnother(d []byte) {
	for _, ic := range c.Server.Clients {
		if c == ic {
			continue
		}
		ic.Send(d)
	}
}

func (c *Client) Close() {
	c.Server = nil
	c.Conn.Close()
	c.Conn = nil
}

func (c *Client) Request() {
	defer func() {
		c.Server.DisConnChan <- c
	}()
	for {
		select {
		case <-c.Run:
			fmt.Println("RUN ")
			return
		default:
			//c.Conn.SetReadLimit(512)
			//c.Conn.SetReadDeadline(time.Now().Add(600 * time.Second))
			//c.Conn.SetPongHandler(func(string) error {
			//	c.Conn.SetReadDeadline(time.Now().Add(600 * time.Second))
			//	return nil
			//})
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				if e, ok := err.(*websocket.CloseError); ok {
					switch e.Code {
					case 1001, 1005, 1006:
						return
					default:
						log.Println(e, " : close error")
						return
					}
				}
				log.Println(err, " - 비정상 종료?")
				return
			}
			pSize := len(message)
			if pSize >= HEADER_SIZE {
				pType := BytesToInt(message[:HEADER_SIZE])
				c.Server.PacketChan <- &Message{c, &Data{pType, message[HEADER_SIZE:pSize]}}
			}
		}
	}
}

func (c *Client) Response() {
	//ticker := time.NewTicker(3 * time.Second)
	//defer func() {
	//	ticker.Stop()
	//}()
	//case tick := <-ticker.C:
	//	log.Println("ping:", c.RemoteAddr(), tick.Second())
	//	err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
	//	CheckError(err)
	for {
		select {
		case <-c.Run:
			fmt.Println("response RUN false ")
			return
		default:
			data := <-c.SendChan
			var err error
			var head uint8
			buf := bytes.NewBuffer(data)
			err = binary.Read(buf, binary.BigEndian, &head)
			CheckError(err)
			if head == 0 {
				err = c.Conn.WriteMessage(websocket.BinaryMessage, data)
			} else {
				if isDev {
					log.Println(string(data))
				}
				err = c.Conn.WriteMessage(websocket.TextMessage, data)
			}
			CheckError(err)
		}
	}
}

func (c *Client) Handle() {
	go c.Request()
	go c.Response()
}
