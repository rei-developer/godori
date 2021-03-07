package game

import (
	"godori.com/db"
	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	cMath "godori.com/util/math"
)

type UserData struct {
	Id     int
	Uuid   string
	Name   string
	Team   int
	Level  int
	Exp    int
	MaxExp int
}

type User struct {
	Index     int
	client    *getty.Client
	character Character
	room      int
	place     int
	userdata  *UserData
}

var nextUserIndex int
var Users map[*getty.Client]*User = make(map[*getty.Client]*User)

func NewUser(c *getty.Client, uid string, loginType int) (*User, bool) {
	if result, ok := db.GetUserByOAuth(uid, loginType); ok {
		nextUserIndex++
		user := &User{
			Index:  nextUserIndex,
			client: c,
			room:   0,
			place:  0,
			userdata: &UserData{
				int(result.Id.Int32),
				result.Uuid.String,
				result.Name.String,
				0,
				0,
				0,
				0,
			},
		}
		Users[c] = user
		user.character.SetPosition(0, -1)
		user.character.Graphics = "Yuzuha"
		return user, true
	}
	return nil, false
}

func (u *User) Remove() bool {
	_, ok := Users[u.client]
	if ok {
		delete(Users, u.client)
	}
	return ok
}

func (u *User) GetUserdata() *UserData {
	return u.userdata
}

func (u *User) GetCreateGameObject() (model int, index int, name string, clanName string, team int, level int, image string, x int, y int, dirX int, dirY int, collider bool) {
	model = 1
	index = u.Index
	name = u.userdata.Name
	clanName = ""
	team = u.userdata.Team
	level = u.userdata.Level
	image = string(u.character.Graphics)
	x = u.character.x
	y = u.character.y
	dirX = u.character.dirX
	dirY = u.character.dirY
	collider = false
	return
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

func (u *User) SetGraphics(image string) {
	u.character.Graphics = Graphics(image)
	u.PublishMap(toClient.SetGraphics(1, u.Index, image))
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
	u.character.Turn(x, y)
	dir := u.character.GetDirection(x, y)
	if r, ok := Rooms[u.room]; ok {
		if r.Passable(u.place, u.character.x, u.character.y, dir, false) && r.Passable(u.place, u.character.x+x, u.character.y-y, 10-dir, true) {
			u.character.Move(x, -y)
			r.Portal(u)
		} else {
			u.Teleport(u.place, u.character.x, u.character.y)
		}
	}
}

func (u *User) Entry(rType int) {
	if u.room > 0 {
		return
	}
	// TODO : set state, send
	r := AvailableRoom(rType)
	r.Join(u)
}

func (u *User) Leave() {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.Leave(u)
	}
}

func (u *User) Hit() {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.UseItem(u)
	}
}

func (u *User) Portal(place int, x int, y int, dirX int, dirY int) {
	// TODO : broadcast
	u.place = place
	u.character.SetPosition(x, y)
	if !(dirX == dirY && dirX == 0) {
		u.Turn(dirX, dirY)
	}
	u.Send(toClient.Portal(place, x, y, u.character.dirX, u.character.dirY))
}

func (u *User) Teleport(place int, x int, y int) {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.Teleport(u, place, x, y, 0, -1)
	}
}

func (u *User) Disconnect() {
	u.Leave()
	u.Remove()
	// TODO : db 저장
}

func (u *User) Send(d []byte) {
	u.client.Send(d)
}

// TODO : notice는 clients의 broadcast 사용할 것.

func (u *User) Publish(d []byte) {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.Publish(d)
	}
}

func (u *User) PublishMap(d []byte) {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.PublishMap(u.place, d)
	}
}

func (u *User) Broadcast(d []byte) {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.Broadcast(u, d)
	}
}

func (u *User) BroadcastMap(d []byte) {
	if u.room < 1 {
		return
	}
	if r, ok := Rooms[u.room]; ok {
		r.BroadcastMap(u, d)
	}
}
