package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

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

	runtime.GOMAXPROCS(runtime.NumCPU())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	server := getty.NewServer("")
	server.OnConnect = OnConnect
	server.OnMessage = OnMessage
	server.OnDisconnect = OnDisconnect
	server.BeforeAccept = BeforeAccept
	handle := &http.Server{Addr: ":" + port, Handler: nil}
	go func() {
		http.HandleFunc("/", server.Listen)
		if err := handle.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	//go func() {
	//	for {
	//		fmt.Println(runtime.NumGoroutine())
	//		time.Sleep(2 * time.Second)
	//	}
	//}()
	log.Println("PORT: " + port + ", GOMAXPROCS: " + strconv.Itoa(runtime.GOMAXPROCS(0)) + " - Run the Godori server.")
	sig := <-sigChan
	switch sig {
	case os.Interrupt:
		log.Println("The Godori server has been shut down.")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()
		if err := handle.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed: %+v", err)
		}
		log.Print("Server Exited Properly.")
	}
}

func BeforeAccept() bool {
	return connections < maxAcceptCnt
}

func Login(u *game.User) {
	u.Send(toClient.UserData(u.GetUserdata()))
	// TODO
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
		connections++
		log.Printf("클라이언트 %s - %s 접속 (동시접속자: %d/%d명)\n", u.UserData.Name, c.RemoteAddr(), connections, maxAcceptCnt)
	}
}

func OnDisconnect(c *getty.Client) {
	if u, ok := game.Users[c]; ok {
		name := u.UserData.Name
		u.Disconnect()
		connections--
		log.Printf("클라이언트 %s - %s 종료 (동시접속자: %d/%d명)\n", name, c.RemoteAddr(), connections, maxAcceptCnt)
	}
}

func DataToInt(d *getty.Data) int {
	var val int32
	buf := bytes.NewBuffer(d.Buffers)
	err := binary.Read(buf, binary.BigEndian, &val)
	CheckError(err)
	return int(val)
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
			u.Result(int(d.Buffers[0]))
		case toServer.ESCAPE:
			u.Leave()
		case toServer.CHAT:
			u.Chat(string(d.Buffers))
		case toServer.CHANGE_USERNAME:
			u.ChangeName(string(d.Buffers))
		case toServer.CREATE_CLAN:
			u.CreateClan(string(d.Buffers))
		case toServer.GET_CLAN:
			u.GetClan()
		case toServer.LEAVE_CLAN:
			u.LeaveClan()
		case toServer.JOIN_CLAN:
			u.JoinClan(DataToInt(d))
		case toServer.CANCEL_CLAN:
			u.CancelClan(DataToInt(d))
		case toServer.INVITE_CLAN:
			u.InviteClan(string(d.Buffers))
		case toServer.KICK_CLAN:
			u.KickClan(DataToInt(d))
		case toServer.SET_OPTION_CLAN:
			u.SetOptionClan(d.Buffers)
		case toServer.PAY_CLAN:
			u.PayClan(DataToInt(d))
		case toServer.DONATE_CLAN:
			u.DonateClan(DataToInt(d))
		case toServer.WITHDRAW_CLAN:
			u.WithdrawClan(DataToInt(d))
		case toServer.LEVEL_UP_CLAN:
			u.LevelUpClan()
		case toServer.MEMBER_INFO_CLAN:
			u.Send(toClient.MemberInfoClan(DataToInt(d)))
		case toServer.SET_UP_MEMBER_LEVEL_CLAN:
			u.SetUpMemberLevelClan(DataToInt(d))
		case toServer.SET_DOWN_MEMBER_LEVEL_CLAN:
			u.SetDownMemberLevelClan(DataToInt(d))
		case toServer.CHANGE_MASTER_CLAN:
			u.ChangeMasterClan(DataToInt(d))
		case toServer.GET_BILLING:
			u.GetBilling()
		case toServer.USE_BILLING:
			u.UseBilling(DataToInt(d))
		case toServer.REFUND_BILLING:
			u.RefundBilling(DataToInt(d))
		case toServer.GET_SHOP:
			u.GetShop(DataToInt(d))
		case toServer.GET_INFO_ITEM:
			u.GetInfoItem(DataToInt(d))
		case toServer.BUY_ITEM:
			u.BuyItem(d.Buffers)
		case toServer.GET_SKIN_LIST:
			u.GetSkinList()
		case toServer.SET_SKIN:
			u.GetSkinList()
		case toServer.GET_PAY_INFO_ITEM:
			u.GetPayInfoItem(DataToInt(d))
		case toServer.GET_RANK:
			u.GetRank(DataToInt(d))
		case toServer.GET_USER_INFO_RANK:
			u.GetUserInfoRank(DataToInt(d))
		case toServer.GET_USER_INFO_RANK_BY_USERNAME:
			u.GetUserInfoRankByUserName(string(d.Buffers))
		case toServer.GET_NOTICE_MESSAGE_COUNT:
			u.GetNoticeMessageCount()
		case toServer.GET_NOTICE_MESSAGE:
			u.GetNoticeMessage(int(d.Buffers[0]))
		case toServer.GET_INFO_NOTICE_MESSAGE:
			u.GetInfoNoticeMessage(DataToInt(d))
		case toServer.WITHDRAW_NOTICE_MESSAGE:
			u.WithdrawNoticeMessage(DataToInt(d))
		case toServer.DELETE_NOTICE_MESSAGE:
			u.DeleteNoticeMessage(DataToInt(d))
		case toServer.RESTORE_NOTICE_MESSAGE:
			u.RestoreNoticeMessage(DataToInt(d))
		case toServer.CLEAR_NOTICE_MESSAGE:
			u.ClearNoticeMessage()
		case toServer.ADD_NOTICE_MESSAGE:
			u.AddNoticeMessage(d.Buffers)
		case toServer.ADD_USER_REPORT:
			u.AddUserReport(d.Buffers)
		case toServer.USE_ITEM:
			u.UseItem()
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
