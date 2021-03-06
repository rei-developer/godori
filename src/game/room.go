package game

import (
	"godori.com/getty"
	roomType "godori.com/util/constant/roomType"
)

type Room struct {
	Index          int
	RoomType       int
	Max            int
	Mode           int
	NextEventIndex int
	Users          map[*getty.Client]*User
	Places         map[int]*Place
	Run            bool
	Lock           bool
}

var nextRoomIndex int
var Rooms map[int]*Room = make(map[int]*Room)

func NewRoom(rType int) *Room {
	nextRoomIndex++
	room := &Room{
		Index:    nextRoomIndex,
		RoomType: rType,
		Max:      30,
		Users:    make(map[*getty.Client]*User),
		Places:   make(map[int]*Place),
	}
	Rooms[nextRoomIndex] = room
	return room
}

func GetRoom(index int) *Room {
	return Rooms[index]
}

func AvailableRoom(rType int) (*Room, bool) {
	for index := range Rooms {
		var r = Rooms[index]
		if r.RoomType == rType && r.CheckJoin() {
			return r, true
		}
	}
	return nil, false
}

func RemoveRoom(r *Room) {
	r.Stop()
	Rooms[r.Index] = &Room{} // TODO : 삭제해도 되나 테스트
	delete(Rooms, r.Index)
}

func (r *Room) Setting() {
	switch r.RoomType {
	case roomType.GAME:
		// TODO : r.mode = new GameMode
	case roomType.PLAYGROUND:
	}
}

func (r *Room) AddEvent() {
	// TODO : add event
}

func (r *Room) RemoveEvent() {
	// TODO : remove event
}

func (r *Room) AddUser(u *User) {
	u.room = r.Index
	r.Users[u.client] = u
	r.Places[u.place].AddUser(u)
}

func (r *Room) RemoveUser(u *User) {
	delete(r.Users, u.client)
	r.Places[u.place].RemoveUser(u)
	u.room = 0
}

func (r *Room) GetPlace(place int) *Place {
	if p, ok := r.Places[place]; ok {
		return p
	} else {
		r.Places[place] = NewPlace(place, r.Index)
		return r.Places[place]
	}
}

func (r *Room) ChangeGameMode(mode int) {
	// TODO : change mode
}

func (r *Room) Publish(d []byte) {
	for _, u := range r.Users {
		u.Send(d)
	}
}

func (r *Room) PublishMap(place int, d []byte) {
	for _, u := range r.GetPlace(place).Users {
		u.Send(d)
	}
}

func (r *Room) Broadcast(self *User, d []byte) {
	for _, u := range r.Users {
		if u == self {
			continue
		}
		u.Send(d)
	}
}

func (r *Room) BroadcastMap(self *User, d []byte) {
	for _, u := range r.GetPlace(self.place).Users {
		if u == self {
			continue
		}
		u.Send(d)
	}
}

func (r *Room) SameMapUsers(place int) map[*getty.Client]*User {
	return r.GetPlace(place).Users
}

func (r *Room) Passable(place int, x int, y int, dir int, collider bool) bool {
	if collider {
		// TODO : event
		//r.GetPlace(place).Events
	}
	return Maps[place].Passable(x, y, dir)
}

func (r *Room) Portal(u *User) {
	if p, ok := Maps[u.place].GetPortal(u.character.x, u.character.y); ok {
		r.Teleport(u, p.NextPlace, p.NextX, p.NextY, p.NextDirX, p.NextDirY)
		if p.Sound != "" {
			// TODO : 사운드 재생
			//r.PublishMap(u.place, )
		}
	}
}

func (r *Room) Teleport(u *User, place int, x int, y int, dirX int, dirY int) {
	r.GetPlace(u.place).RemoveUser(u)
	u.Portal(place, x, y, dirX, dirY)
	r.GetPlace(u.place).AddUser(u)
	r.Draw(u)
}

func (r *Room) Hit(u *User) {
	// TODO : hit
}

func (r *Room) UseItem(u *User) {
	// TODO : use item
}

func (r *Room) CheckJoin() bool {
	return len(r.Users) < r.Max && !r.Lock
}

func (r *Room) Draw(u *User) {
	// TODO : draw
}

func (r *Room) Join(u *User) {
	// TODO : join
	r.AddUser(u)
}

func (r *Room) Leave(u *User) {
	// TODO : leave
	r.RemoveUser(u)
	if len(r.Users) <= 0 {
		RemoveRoom(r)
	}
}

func (r *Room) Start() {
	if r.Run {
		return
	}
	r.Run = true
}

func (r *Room) Pause() {
	r.Run = false
}

func (r *Room) Stop() {
	if !r.Run {
		return
	}
	r.Run = false
}

func (r *Room) Update() {
	// TODO : update
}
