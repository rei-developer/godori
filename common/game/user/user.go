package user

import (
	"godori.com/game/character"
	"godori.com/getty"
)

type UserData struct {
	Id   string
	Uuid string
}

type User struct {
	client    *getty.Client
	character character.Character
	room      int // `tag: "채널"`
	place     int
	userdata  UserData
}

var Users map[*getty.Client]User = make(map[*getty.Client]User)

func New(client *getty.Client) *User {
	return &User{
		client:   client,
		//userdata: userdata,
	}
}

func (u *User) Move(d int) {
	u.character.Move(d)
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetUserdata() UserData {
	return u.userdata
}