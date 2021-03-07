package game

import (
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

func (m *GameMode) Hit(self *User, target *User) {
	m.Mode.Hit(self, target)
}

func (m *GameMode) UseItem(u *User) {
	m.Mode.UseItem(u)
}

func (m *GameMode) Update() {
	m.Mode.Update()
}
