package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	db "godori.com/database"
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

type Users struct {
	Users []User `json:"users"`
}

// User struct which contains a name
// a type and a list of social links
type User struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Age    int    `json:"Age"`
	Social Social `json:"social"`
}

// Social struct which contains a
// list of links
type Social struct {
	Facebook string `json:"facebook"`
	Twitter  string `json:"twitter"`
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("./lib/shop.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var users Users

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(users.Users); i++ {
		fmt.Println("User Type: " + users.Users[i].Type)
		fmt.Println("User Age: " + strconv.Itoa(users.Users[i].Age))
		fmt.Println("User Name: " + users.Users[i].Name)
		fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	}

	id, uuid := db.GetUser(1)
	fmt.Println(id, uuid, "입니다")

	result := db.GetUsers()
	for i, v := range result {
		fmt.Println(v.Id, i)
		fmt.Println(v.Name, i)
		fmt.Println(v.Uuid, i)
	}

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
