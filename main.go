package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	//"godori.com/db"
	"godori.com/game"
	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	toServer "godori.com/packet/toServer"
)

const (
	port         = "50000"
	maxAcceptCnt = 3
)

var connections int

func main() {
	//id, uuid := db.GetUserById(1)
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

	var wg sync.WaitGroup
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := getty.NewServer("")
	server.OnConnect = OnConnect
	server.OnMessage = OnMessage
	server.OnDisconnect = OnDisconnect
	server.BeforeAccept = BeforeAccept
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/", server.Listen)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	log.Println("PORT: " + port + ", GOMAXPROCS: " + strconv.Itoa(runtime.GOMAXPROCS(0)) + " - 서버를 실행합니다.")
	wg.Wait()
}

func BeforeAccept() bool {
	return connections < maxAcceptCnt
}

func Login(u *game.User) {
	uData := u.GetUserdata()
	u.Send(toClient.UserData(u.Index, uData.Id, uData.Name))
	// TODO :
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
	if u, ok := game.NewUser(c, uid, loginType); ok {
		data := u.GetUserdata()
		connections++
		fmt.Printf("클라이언트 %s - %s 접속 (동시접속자: %d/%d명)\n", data.Name, c.RemoteAddr(), connections, maxAcceptCnt)
	}
}

func OnDisconnect(c *getty.Client) {
	if u, ok := game.Users[c]; ok {
		u.Disconnect()
		connections--
		fmt.Printf("클라이언트 %s 종료 (동시접속자: %d/%d명)\n", c.RemoteAddr(), connections, maxAcceptCnt)
	}
}

func OnMessage(c *getty.Client, d *getty.Data) {
	if u, ok := game.Users[c]; ok {
		switch d.Type {
		case toServer.HELLO:
			Login(u)
		case toServer.INPUT_ARROW:
			var err error
			var pos [2]int8
			var dir [1]uint8
			buf := bytes.NewBuffer(d.Buffers)
			err = binary.Read(buf, binary.BigEndian, &pos)
			err = binary.Read(buf, binary.BigEndian, &dir)
			CheckError(err)
			x, y := int(pos[0]), int(pos[1])
			if dir[0] == 0 {
				u.Turn(x, y)
			} else {
				u.Move(x, y)
			}
		case toServer.INPUT_HIT:
			u.Hit()
		case toServer.ENTER_ROOM:
			u.Entry(int(d.Buffers[0]))
		case toServer.REWARD:

		case toServer.ESCAPE:
			u.Leave()
		case toServer.CHAT:

		case toServer.CHANGE_USERNAME:
		case toServer.CREATE_CLAN:
		case toServer.GET_CLAN:
		case toServer.LEAVE_CLAN:
		case toServer.JOIN_CLAN:
		case toServer.CANCEL_CLAN:
		case toServer.KICK_CLAN:
		case toServer.SET_OPTION_CLAN:
		case toServer.PAY_CLAN:
		case toServer.DONATE_CLAN:
		case toServer.WITHDRAW_CLAN:
		case toServer.LEVEL_UP_CLAN:
		case toServer.MEMBER_INFO_CLAN:
		case toServer.SET_UP_MEMBER_LEVEL_CLAN:
		case toServer.SET_DOWN_MEMBER_LEVEL_CLAN:
		case toServer.CHANGE_MASTER_CLAN:
		case toServer.GET_BILLING:
		case toServer.USE_BILLING:
		case toServer.REFUND_BILLING:
		case toServer.GET_SHOP:
		case toServer.GET_INFO_ITEM:
		case toServer.BUY_ITEM:
		case toServer.GET_SKIN_LIST:
		case toServer.SET_SKIN:
		case toServer.GET_PAY_INFO_ITEM:
		case toServer.GET_RANK:
		case toServer.GET_USER_INFO_RANK:
		case toServer.GET_USER_INFO_RANK_BY_USERNAME:
		case toServer.GET_NOTICE_MESSAGE_COUNT:
		case toServer.GET_NOTICE_MESSAGE:
		case toServer.GET_INFO_NOTICE_MESSAGE:
		case toServer.WITHDRAW_NOTICE_MESSAGE:
		case toServer.DELETE_NOTICE_MESSAGE:
		case toServer.RESTORE_NOTICE_MESSAGE:
		case toServer.CLEAR_NOTICE_MESSAGE:
		case toServer.ADD_NOTICE_MESSAGE:
		case toServer.ADD_USER_REPORT:
		case toServer.USE_ITEM:
			//case toServer.ADD_USER_REPORT:
			//	b := []byte(string(d.Buffers))
			//	var data map[string]interface{}
			//	err := json.Unmarshal(b, &data)
			//	CheckError(err)
			//	num := int(data["number"].(float64))
			//	fmt.Println(num)
			//	fmt.Println(data["string"])
			//case toServer.BUY_ITEM:
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
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
