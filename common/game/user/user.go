package user

import (
	db "godori.com/database"
	"godori.com/game/character"
	"godori.com/getty"
)

type UserData struct {
	Uuid string
	Name string
}

type User struct {
	client    *getty.Client
	character character.Character
	room      int
	place     int
	userdata  *UserData
}

var Users map[*getty.Client]*User = make(map[*getty.Client]*User)

func New(client *getty.Client, uid string, loginType int) (*User, bool) {
	result, ok := db.GetUserByOAuth(uid, loginType)
	return &User{
		client: client,
		userdata: &UserData{
			result.Uuid.String,
			result.Name.String,
		},
	}, ok
}

func (u *User) Move(d int) {
	u.character.Move(d)
}

func (u *User) GetUserdata() UserData {
	return *u.userdata
}
