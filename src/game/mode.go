package game

import (
	"fmt"
	"godori.com/getty"
	mapType "godori.com/util/constant/mapType"
	modeType "godori.com/util/constant/modeType"
	roomType "godori.com/util/constant/roomType"
	cMath "godori.com/util/math"
)

type GameMode struct {
	Room *Room
	Mode IGameMode
}

type IGameMode interface {
	InitEvent()
	MoveToBase(*User)
	Join(*User)
	Leave(*User)
	DrawEvents(*User)
	DrawUsers(*User)
	Hit(*User, *User) bool
	UseItem(*User)
	Update()
}

const (
	STATE_READY = iota
	STATE_GAME
	STATE_RESULT
)

func NewMode(r *Room) *GameMode {
	mode := &GameMode{Room: r}
	mode.ChangeMode(modeType.NONE, false)
	return mode
}

func (m *GameMode) ChangeMode(mType int, join bool) {
	rType := m.Room.RoomType
	if rType == roomType.PLAYGROUND {
		m.Mode = &PlaygroundMode{Room: m.Room}
	} else if rType == roomType.GAME {
		pType := cMath.Rand(mapType.DESERT) + 1
		switch mType {
		case modeType.NONE:
			m.Mode = &NoneMode{Room: m.Room}
		case modeType.RESCUE:
			m.Mode = NewRescueMode(m.Room, pType)
		}
	}
	m.Mode.InitEvent()
	if join {
		for _, u := range m.Room.Users {
			m.Mode.Join(u)
		}
	}
}

func (m *GameMode) MoveToBase(u *User) {
	m.Mode.MoveToBase(u)
}

func (m *GameMode) Join(u *User) {
	m.Mode.Join(u)
}

func (m *GameMode) Leave(u *User) {
	m.Mode.Leave(u)
}

func (m *GameMode) DrawEvents(u *User) {
	m.Mode.DrawEvents(u)
}

func (m *GameMode) DrawUsers(u *User) {
	m.Mode.DrawUsers(u)
}

func (m *GameMode) Hit(self *User, target *User) bool {
	return m.Mode.Hit(self, target)
}

func (m *GameMode) UseItem(u *User) {
	m.Mode.UseItem(u)
}

func (m *GameMode) Sample(target map[*getty.Client]*User, count int) map[*getty.Client]*User {
	users := make(map[*getty.Client]*User)
	pickers := make(map[*getty.Client]*User)
	for _, u := range target {
		users[u.Client] = u
	}
	for count > 0 {
		c := 0
		fmt.Println(len(users))
		pick := cMath.Rand(len(users) + 1)
		for _, u := range users {
			if c == pick {
				if _, ok := pickers[u.Client]; !ok {
					pickers[u.Client] = u
					delete(users, u.Client)
					break
				}
			}
			c++
		}
		count--
	}
	return pickers
}

func (m *GameMode) Update() {
	m.Mode.Update()
}
