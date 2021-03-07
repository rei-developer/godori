package game

import (
	toClient "godori.com/packet/toClient"
)

type NoneMode struct {
	Room *Room
}

func (m *NoneMode) MoveToBase(u *User) {
	u.Teleport(42, 9, 7)
}

func (m *NoneMode) Join(u *User) {
	u.SetGraphics(string(u.character.Graphics))
	m.MoveToBase(u)
}

func (m *NoneMode) Leave(u *User) {

}

func (m *NoneMode) DrawEvents(u *User) {

}

func (m *NoneMode) DrawUsers(self *User) {
	for _, u := range m.Room.SameMapUsers(self.place) {
		if u == self {
			return
		}
		u.Send(toClient.CreateGameObject(self.GetCreateGameObject()))
		self.Send(toClient.CreateGameObject(u.GetCreateGameObject()))
	}
}

func (m *NoneMode) Hit(u *User) {

}

func (m *NoneMode) UseItem(u *User) {

}

func (m *NoneMode) Update() {

}
