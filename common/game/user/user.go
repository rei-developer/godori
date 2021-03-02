package user

import (
	//"godori.com/game/character"
	"godori.com/getty"
)

type UserData struct {
	Id   string
	Uuid string
}

type User struct {
	client    *getty.Client
	//character character.Character
	room      int // `tag: "채널"`
	place     int
	name      string
	userdata  UserData
}

var Users map[*getty.Client]User = make(map[*getty.Client]User)

func New(
	client *getty.Client,
	userdata UserData,
) *User {
	return &User{
		client:   client,
		name:     "호옹이",
		userdata: userdata,
	}
}

func (u *User) Move(d int) {
	//u.character.Moves(d)
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetUserdata() UserData {
	return u.userdata
}

func UserLength() int {
	return len(Users)
}
