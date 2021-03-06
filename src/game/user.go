package game

import (
	"godori.com/db"
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

func NewUser(c *getty.Client, uid string, loginType int) (*User, bool) {
	result, ok := db.GetUserByOAuth(uid, loginType)
	user := &User{
		client: c,
		userdata: &UserData{
			result.Uuid.String,
			result.Name.String,
			0,
			0,
			0,
		},
	}
	Users[c] = user
	return user, ok
}

func RemoveUser(u *User) bool {
	_, ok := Users[u.client]
	if ok {
		delete(Users, u.client)
	}
	return ok
}

func RemoveUserByClient(c *getty.Client) bool {
	_, ok := Users[c]
	if ok {
		delete(Users, c)
	}
	return ok
}

func (u *User) GetUserdata() UserData {
	return *u.userdata
}

func (u *User) SetUpLevel(v int) {
	u.userdata.Level += v
}

func (u *User) SetUpExp(v int) {
	if u.userdata.Level > 200 {
		return
	}
	u.userdata.Exp = cMath.Max(u.userdata.Exp+v, 0)
	for u.userdata.Exp >= u.userdata.MaxExp {
		u.userdata.Exp -= u.userdata.MaxExp
		u.SetUpLevel(1)
	}
}

func (u *User) SetUpCash(v int) {
	// TODO : 개발중
}

func (u *User) GetMaxExp(v int) int {
	return (cMath.Pow(v, 2) * (v * 5)) + 200
}

func (u *User) Turn(dirX int, dirY int) {
	if u.room < 1 {
		return
	}
	u.character.Turn(dirX, dirY)
}

func (u *User) Move(x int, y int) {
	if u.room < 1 {
		return
	}
	u.character.Turn(x, -y)
	dir := u.character.GetDirection(x, -y)
	r := GetRoom(u.room)
	if r.Passable(u.place, u.character.x, u.character.y, dir, false) && r.Passable(u.place, u.character.x+x, u.character.y+y, 10-dir, true) {
		u.character.Move(x, y)
		r.Portal(u)
	} else {
		u.Teleport(u.place, u.character.x, u.character.y)
	}
}

func (u *User) Entry(rType int) {
	if u.room > 0 {
		return
	}
	// TODO : set state, send
	if r, ok := AvailableRoom(rType); ok {
		r.Join(u)
	} else {
		r := NewRoom(rType)
		r.Join(u)
	}
}

func (u *User) Leave() {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).Leave(u)
}

func (u *User) Hit() {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).UseItem(u)
}

func (u *User) Portal(place int, x int, y int, dirX int, dirY int) {
	// TODO : broadcast
	u.place = place
	u.character.SetPosition(x, y)
	if !(dirX == dirY && dirX == 0) {
		u.Turn(dirX, dirY)
	}
	// TODO : send
}

func (u *User) Teleport(place int, x int, y int) {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).Teleport(u, place, x, y, 0, 0)
}

func (u *User) Send(d []byte) {
	u.client.Send(d)
}

// TODO : notice는 clients의 broadcast 사용할 것.

func (u *User) Publish(d []byte) {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).Publish(d)
}

func (u *User) PublishMap(d []byte) {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).PublishMap(u.place, d)
}

func (u *User) Broadcast(d []byte) {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).Broadcast(u, d)
}

func (u *User) BroadcastMap(d []byte) {
	if u.room < 1 {
		return
	}
	GetRoom(u.room).BroadcastMap(u, d)
}
