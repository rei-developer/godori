package game

import (
	modeType "godori.com/util/constant/modeType"
	roomType "godori.com/util/constant/roomType"
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
	Hit(*User)
	UseItem(*User)
	Update()
}

func NewMode(r *Room) *GameMode {
	mode := &GameMode{Room: r}
	mode.ChangeMode(modeType.NONE)
	return mode
}

func (m *GameMode) ChangeMode(mType int) {
	rType := m.Room.RoomType
	if rType == roomType.PLAYGROUND {
		m.Mode = &PlaygroundMode{Room: m.Room}
	} else if rType == roomType.GAME {
		switch mType {
		case modeType.NONE:
			m.Mode = &NoneMode{Room: m.Room}
		case modeType.RESCUE:
			m.Mode = &RescueMode{Room: m.Room}
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

func (m *GameMode) Hit(u *User) {
	m.Mode.Hit(u)
}

func (m *GameMode) UseItem(u *User) {
	m.Mode.UseItem(u)
}

func (m *GameMode) Update() {
	m.Mode.Update()
}
