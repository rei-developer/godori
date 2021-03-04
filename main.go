package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	_ "godori.com/game/shop"
	user "godori.com/game/user"
	"godori.com/getty"
	toserver "godori.com/packet/toserver"
)

const maxAcceptCnt = 3

var connections = 0

//var randomData string = `{ "event": "Bye" }`
//var sendData map[string]interface{}
//json.Unmarshal([]byte(randomData), &sendData)
//fmt.Println(sendData["event"].(string))
//refineSendData, err := json.Marshal(sendData)
//err = c.WriteMessage(mt, refineSendData)

func main() {

	//id, uuid := db.GetUser(1)
	//fmt.Println(id, uuid, "입니다")

	//item, _ := db.GetUserById(1)
	//fmt.Println(item.Name)
	//
	//item2, ok := db.GetUserByOAuth("test", 0)
	//if ok {
	//	fmt.Println(item2)
	//} else {
	//	fmt.Println(" 없구나")
	//}

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
	server.OnConnect = OnConnect
	server.OnMessage = OnMessage
	server.OnDisconnect = OnDisconnect
	server.BeforeAccept = BeforeAccept
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/", server.Listen)
		log.Fatal(http.ListenAndServe(":50000", nil))
	}()
	fmt.Println("슈퍼호옹호옹이서버를 실행합니다.")
	wg.Wait()
}

func BeforeAccept() bool {
	return connections < maxAcceptCnt
}

func Login(u *user.User) {
	//fmt.Println(u.GetUserdata().Name + "이라네~~~~ 로그인 성공이라네")
}

func OnConnect(c *getty.Client) {
	token := c.GetToken()
	var uid string
	var loginType int
	if token == "debug" {
		uid, loginType = "110409668035092753325", 0
	} else {
		uid, loginType = "110409668035092753325", 0
	}
	if u, ok := user.New(c, uid, loginType); ok {
		data := u.GetUserdata()
		//Login(u)
		connections++
		fmt.Printf("클라이언트 %s - %s 접속 (동시접속자: %d/%d명)\n", data.Name, c.RemoteAddr(), connections, maxAcceptCnt)
	}
}

func OnDisconnect(c *getty.Client) {
	if ok := user.RemoveByClient(c); ok {
		connections--
		fmt.Printf("클라이언트 %s 종료 (동시접속자: %d/%d명)\n", c.RemoteAddr(), connections, maxAcceptCnt)
	}
}

func OnMessage(c *getty.Client, d *getty.Data) {
	u := user.Users[c]
	switch d.Type {
	case toserver.HELLO:
		Login(u)
	case toserver.INPUT_ARROW:
		var pos [2]int8
		var dir [1]uint8
		buf := bytes.NewBuffer(d.Buffers)
		binary.Read(buf, binary.BigEndian, &pos)
		binary.Read(buf, binary.BigEndian, &dir)
		fmt.Println(pos[0], pos[1])
		fmt.Println(dir[0])
	case toserver.INPUT_HIT:
	case toserver.ENTER_ROOM:
		u.Entry()
	case toserver.REWARD:
	case toserver.ESCAPE:
	case toserver.CHAT:
	case toserver.CHANGE_USERNAME:
	case toserver.CREATE_CLAN:
	case toserver.GET_CLAN:
	case toserver.LEAVE_CLAN:
	case toserver.JOIN_CLAN:
	case toserver.CANCEL_CLAN:
	case toserver.KICK_CLAN:
	case toserver.SET_OPTION_CLAN:
	case toserver.PAY_CLAN:
	case toserver.DONATE_CLAN:
	case toserver.WITHDRAW_CLAN:
	case toserver.LEVEL_UP_CLAN:
	case toserver.MEMBER_INFO_CLAN:
	case toserver.SET_UP_MEMBER_LEVEL_CLAN:
	case toserver.SET_DOWN_MEMBER_LEVEL_CLAN:
	case toserver.CHANGE_MASTER_CLAN:
	case toserver.GET_BILLING:
	case toserver.USE_BILLING:
	case toserver.REFUND_BILLING:
	case toserver.GET_SHOP:
	case toserver.GET_INFO_ITEM:
	case toserver.BUY_ITEM:
	case toserver.GET_SKIN_LIST:
	case toserver.SET_SKIN:
	case toserver.GET_PAY_INFO_ITEM:
	case toserver.GET_RANK:
	case toserver.GET_USER_INFO_RANK:
	case toserver.GET_USER_INFO_RANK_BY_USERNAME:
	case toserver.GET_NOTICE_MESSAGE_COUNT:
	case toserver.GET_NOTICE_MESSAGE:
	case toserver.GET_INFO_NOTICE_MESSAGE:
	case toserver.WITHDRAW_NOTICE_MESSAGE:
	case toserver.DELETE_NOTICE_MESSAGE:
	case toserver.RESTORE_NOTICE_MESSAGE:
	case toserver.CLEAR_NOTICE_MESSAGE:
	case toserver.ADD_NOTICE_MESSAGE:
	case toserver.ADD_USER_REPORT:
	case toserver.USE_ITEM:
		//case toserver.ADD_USER_REPORT:
		//	b := []byte(string(d.Buffers))
		//	var data map[string]interface{}
		//	err := json.Unmarshal(b, &data)
		//	CheckError(err)
		//	num := int(data["number"].(float64))
		//	fmt.Println(num)
		//	fmt.Println(data["string"])
		//case toserver.BUY_ITEM:
		//	fmt.Println("하하 채팅이네")
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

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
