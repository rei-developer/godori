package packet

import (
	"encoding/json"
)

func UserData(index int, id int, name string) []byte {
	type Packet struct {
		Header int
		Index  int
		Id     int
		Name   string
	}
	packet := Packet{USER_DATA, index, id, name}
	bytes, err := json.Marshal(packet)
	CheckError(err)
	return bytes
}

func NoticeMessage(text string) []byte {
	type Packet struct {
		Header int
		Text   string
	}
	packet := Packet{NOTICE_MESSAGE, text}
	bytes, err := json.Marshal(packet)
	CheckError(err)
	return bytes
}

func Portal(place int, x int, y int, dirX int, dirY int) []byte {
	type Packet struct {
		Header int
		Place  int
		X      int
		Y      int
		DirX   int
		DirY   int
	}
	packet := Packet{PORTAL, place, x, y, dirX, dirY}
	bytes, err := json.Marshal(packet)
	CheckError(err)
	return bytes
}

func CreateGameObject(model int, index int, name string, clanName string, team int, level int, image string, x int, y int, dirX int, dirY int, collider bool) []byte {
	type Packet struct {
		Header   int
		Model    int
		Index    int
		Name     string
		ClanName string
		Team     int
		Level    int
		Image    string
		X        int
		Y        int
		DirX     int
		DirY     int
		Collider bool
	}
	packet := Packet{CREATE_GAME_OBJECT, model, index, name, clanName, team, level, image, x, y, dirX, dirY, collider}
	bytes, err := json.Marshal(packet)
	CheckError(err)
	return bytes
}

func SetGraphics(model int, index int, image string) []byte {
	type Packet struct {
		Header int
		Model  int
		Index  int
		Image  string
	}
	packet := Packet{SET_GRAPHICS, model, index, image}
	bytes, err := json.Marshal(packet)
	CheckError(err)
	return bytes
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	USER_DATA = iota
	VIBRATE
	CONNECTION_COUNT
	SYSTEM_MESSAGE
	INFORM_MESSAGE
	NOTICE_MESSAGE
	CHAT_MESSAGE
	PORTAL
	CREATE_GAME_OBJECT
	REMOVE_GAME_OBJECT
	SET_GRAPHICS
	PLAY_SOUND
	UPDATE_ROOM_USER_COUNT
	UPDATE_MODE_COUNT
	UPDATE_GAME_ITEM
	REMOVE_GAME_ITEM
	SET_GAME_TEAM
	MODE_DATA
	GET_CLAN
	INVITE_CLAN
	DEAD_ANIMATION
	RESULT_GAME
	ENTER_WARDROBE
	LEAVE_WARDROBE
	SWITCH_LIGHT
	QUIT_GAME
	UPDATE_MODE_SCORE
	SET_OPTION_CLAN
	MEMBER_INFO_CLAN
	UPDATE_CLAN
	MESSAGE_CLAN
	GET_BILLING
	UPDATE_BILLING
	GET_SHOP
	GET_SKIN_ITEM
	MESSAGE_SHOP
	MESSAGE_LOBBY
	GET_SKIN_LIST
	UPDATE_CASH_AND_COIN
	GET_PAY_INFO_ITEM
	GET_RANK
	GET_USER_INFO_RANK
	MESSAGE_RANK
	GET_NOTICE_MESSAGE_COUNT
	GET_NOTICE_MESSAGE
	GET_INFO_NOTICE_MESSAGE
	DELETE_NOTICE_MESSAGE
	MESSAGE_GAME
	SET_ANIMATION
	COMBO_MESSAGE
)
