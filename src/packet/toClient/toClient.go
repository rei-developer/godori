package packet

import (
	"encoding/json"
)

func UserData(name string) []byte {
	type Packet struct {
		Header int
		Name   string
	}
	packet := Packet{
		USER_DATA,
		name,
	}
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
