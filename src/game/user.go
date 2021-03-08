package game

import (
	"godori.com/db"
	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	cMath "godori.com/util/math"
)

type UserData struct {
	Id           int
	Uid          string
	Uuid         string
	Name         string
	ClanName     string
	Rank         int
	Sex          int
	Level        int
	Exp          int
	MaxExp       int
	Coin         int
	Cash         int
	Point        int
	Win          int
	Lose         int
	Kill         int
	Death        int
	Assist       int
	Blast        int
	Rescue       int
	Survive      int
	Escape       int
	Graphics     string
	RedGraphics  string
	BlueGraphics string
	Memo         string
	Admin        int
}

type User struct {
	Index     int
	client    *getty.Client
	character Character
	room      int
	place     int
	UserData  *UserData
	GameData  map[string]interface{}
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
			UserData: &UserData{
				int(result.Id.Int32),
				result.Uid.String,
				result.Uuid.String,
				result.Name.String,
				"",
				1,
				int(result.Sex.Int32),
				int(result.Level.Int32),
				int(result.Exp.Int32),
				0,
				int(result.Coin.Int32),
				int(result.Cash.Int32),
				int(result.Point.Int32),
				int(result.Win.Int32),
				int(result.Lose.Int32),
				0, //int(result.Kill.Int32),
				int(result.Death.Int32),
				int(result.Assist.Int32),
				int(result.Blast.Int32),
				int(result.Rescue.Int32),
				int(result.Survive.Int32),
				int(result.Escape.Int32),
				"",
				result.RedGraphics.String,
				result.BlueGraphics.String,
				result.Memo.String,
				int(result.Admin.Int32),
			},
		}
		Users[c] = user
		user.character.Setting(1, user.UserData.RedGraphics, user.UserData.BlueGraphics)
		user.UserData.MaxExp = cMath.GetMaxExp(user.UserData.Level)
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
	return u.UserData
}

func (u *User) GetCreateGameObject(hide bool) (model int, index int, name string, clanName string, team int, level int, image string, x int, y int, dirX int, dirY int, collider bool) {
	model = 1
	index = u.Index
	name = u.UserData.Name
	clanName = ""
	team = 0
	if t, ok := u.GameData["team"]; ok {
		team = t.(int)
	}
	level = u.UserData.Level
	image = u.character.Graphics.Image
	x = u.character.x
	y = u.character.y
	dirX = u.character.dirX
	dirY = u.character.dirY
	collider = false
	if hide {
		name = ""
		clanName = ""
		level = 0
	}
	return
}

func (u *User) SetUpLevel(v int) {
	u.UserData.Level += v
}

func (u *User) SetUpExp(v int) {
	if u.UserData.Level > 200 {
		return
	}
	u.UserData.Exp = cMath.Max(u.UserData.Exp+v, 0)
	for u.UserData.Exp >= u.UserData.MaxExp {
		u.UserData.Exp -= u.UserData.MaxExp
		u.SetUpLevel(1)
	}
}

func (u *User) SetUpCash(v int) {
	u.UserData.Cash += v
	// TODO : 개발중
}

func (u *User) GetMaxExp(v int) int {
	return (cMath.Pow(v, 2) * (v * 5)) + 200
}

func (u *User) SetGraphics(image string) {
	u.character.Graphics.Image = image
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
		r.Hit(u)
	}
}

func (u *User) UseItem() {
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
