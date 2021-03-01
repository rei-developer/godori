package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"sandspoon.com/getty"
	"sandspoon.com/packet"
)

const (
	maxAcceptCnt = 3
)

var (
	connections = 0
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		fmt.Println(string(message))

		var randomData string = `{ "event": "Bye" }`

		var data map[string]interface{}
		json.Unmarshal(message, &data)
		fmt.Println(data["event"].(string))

		var sendData map[string]interface{}
		json.Unmarshal([]byte(randomData), &sendData)
		fmt.Println(sendData["event"].(string))

		refineSendData, err := json.Marshal(sendData)
		err = c.WriteMessage(mt, refineSendData)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup
	server := getty.NewServer("")
	server.OnConnect = onConnect
	server.OnMessage = onMessage
	server.OnDisconnect = onDisconnect
	server.BeforeAccept = beforeAccept
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/", server.Listen)
		log.Fatal(http.ListenAndServe(":50000", nil))
	}()
	fmt.Println("슈퍼호옹호옹이서버를 실행합니다.")
	wg.Wait()
}

func beforeAccept() bool {
	// 최대 접속 인원 제한
	return connections < maxAcceptCnt
}

func onConnect(c *getty.Client) {
	connections++
	fmt.Printf("클라이언트 %s 접속 (동시접속자: %d/%d명)\n", c.RemoteAddr(), connections, maxAcceptCnt)
}

func onDisconnect(c *getty.Client) {
	connections--
	fmt.Printf("클라이언트 %s 종료 (동시접속자: %d/%d명)\n", c.RemoteAddr(), connections, maxAcceptCnt)
}

func onMessage(c *getty.Client, d *getty.Data) {
	b, _ := json.MarshalIndent(d.JsonData, "", "  ")
	var data map[string]interface{}
	json.Unmarshal(b, &data)
	fmt.Println(data["test"])

	switch d.Type {
	case packet.USER_LOGIN:

	case packet.USER_CHAT:
		//message, err := packet.ReadChat(d.Buffers)
		//if err != nil {
		//	return
		//}
		//message = c.RemoteAddr().String() + ":" + message
		//chat := packet.MakeChat(builder, 100, message)
		//c.Broadcast(1, chat)
		//fmt.Println(message)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
