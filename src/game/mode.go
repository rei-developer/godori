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
	mode.ChangeMode(r.RoomType, modeType.NONE)
	return mode
}

func (m *GameMode) ChangeMode(rType int, mType int) {
	if rType == roomType.PLAYGROUND {
		m.Mode = &PlaygroundMode{
			Name: "플레이그라운드",
			Room: m.Room,
		}
	} else if rType == roomType.GAME {
		switch mType {
		case modeType.NONE:
			m.Mode = &NoneMode{
				Name: "없다네",
			}
		case modeType.RESCUE:
			m.Mode = &RescueMode{
				Name: "dd",
			}
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
