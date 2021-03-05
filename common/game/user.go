package game

import (
	db "godori.com/database"
	"godori.com/getty"
	cMath "godori.com/util/math"
)

type UserData struct {
	Uuid   string
	Name   string
	Level  int
	Exp    int
	MaxExp int
}

type User struct {
	client    *getty.Client
	character Character
	room      int
	place     int
	userdata  *UserData
}

var Users map[*getty.Client]*User = make(map[*getty.Client]*User)

func NewUser(client *getty.Client, uid string, loginType int) (*User, bool) {
	result, ok := db.GetUserByOAuth(uid, loginType)
	user := &User{
		client: client,
		userdata: &UserData{
			result.Uuid.String,
			result.Name.String,
			0,
			0,
			0,
		},
	}
	Users[client] = user
	return user, ok
}

func (u *User) GetUserdata() UserData {
	return *u.userdata
}

func (u *User) SetUpLevel(val int) {
	u.userdata.Level += val
}

func (u *User) SetUpExp(val int) {
	if u.userdata.Level > 200 {
		return
	}
	u.userdata.Exp = cMath.Max(u.userdata.Exp+val, 0)
	for u.userdata.Exp >= u.userdata.MaxExp {
		u.userdata.Exp -= u.userdata.MaxExp
		u.SetUpLevel(1)
	}
}

func (u *User) SetUpCash(val int) {
	// TODO : 개발중
}

func (u *User) GetMaxExp(val int) int {
	return (cMath.Pow(val, 2) * (val * 5)) + 200
}

func (u *User) Move(d int) {
	u.character.Move(d)
}

func (u *User) Entry(roomType int) {
	if u.room > 0 {
		return
	}
	// TODO : set state, send

}

func RemoveByClient(c *getty.Client) bool {
	_, ok := Users[c]
	if ok {
		delete(Users, c)
	}
	return ok
}

func RemoveByUser(u *User) bool {
	_, ok := Users[u.client]
	if ok {
		delete(Users, u.client)
	}
	return ok
}
