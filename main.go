package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
	db "godori.com/database"
	_ "godori.com/game/shop"
	user "godori.com/game/user"
	"godori.com/getty"
	toserver "godori.com/packet/toserver"
)

const maxAcceptCnt = 3

var (
	connections = 0
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

//var randomData string = `{ "event": "Bye" }`
//var sendData map[string]interface{}
//json.Unmarshal([]byte(randomData), &sendData)
//fmt.Println(sendData["event"].(string))
//refineSendData, err := json.Marshal(sendData)
//err = c.WriteMessage(mt, refineSendData)

func main() {

	//id, uuid := db.GetUser(1)
	//fmt.Println(id, uuid, "입니다")

	items := db.GetClanMembers(29)
	for i, v := range items {
		fmt.Println(v, i, "ㅇㅇ")
	}

	//result := db.GetUsers()
	//for i, v := range result {
	//	var index int32 = v.Id.Int32
	//	fmt.Println(index, i)
	//	fmt.Println(v.Uid, i)
	//	if v.Uuid.Valid {
	//		fmt.Println(v.Uuid, "uuid")
	//	}
	//	if v.Name.Valid {
	//		fmt.Println(v.Name, "name")
	//	}
	//}

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))

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
	return connections < maxAcceptCnt
}

func onConnect(c *getty.Client) {
	connections++
	fmt.Printf("클라이언트 %s 접속 (동시접속자: %d/%d명)\n", c.RemoteAddr(), connections, maxAcceptCnt)
}

func onDisconnect(c *getty.Client) {
	if _, ok := user.Users[c]; ok {
		delete(user.Users, c)
	}
	connections--
	fmt.Printf("클라이언트 %s 종료 (동시접속자: %d/%d명)\n", c.RemoteAddr(), connections, maxAcceptCnt)
}

func onMessage(c *getty.Client, d *getty.Data) {
	fmt.Println(d.Type)
	switch d.Type {
	case toserver.HELLO:
		user.Users[c] = *user.New(c, user.UserData{Id: "아이디데스", Uuid: "유유아이디데스"})
		u := user.Users[c]
		fmt.Println(u.GetName())
		fmt.Println(u.GetUserdata())
		for key, val := range user.Users {
			fmt.Println(key, val)
		}
		fmt.Println(len(user.Users))
	case toserver.ADD_USER_REPORT:
		b := []byte(string(d.Buffers))
		var data map[string]interface{}
		err := json.Unmarshal(b, &data)
		checkError(err)
		num := int(data["number"].(float64))
		fmt.Println(num)
		fmt.Println(data["string"])
	case toserver.BUY_ITEM:
		fmt.Println("하하 채팅이네")
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
