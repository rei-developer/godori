package game

import (
	"fmt"

	toClient "godori.com/packet/toClient"
	cMath "godori.com/util/math"
)

type NoneMode struct {
	Room  *Room
	Count int
}

func (m *NoneMode) ChangeMode() {
	rand := cMath.Rand(1) + 1
	m.Room.Mode.ChangeMode(rand, true)
}

func (m *NoneMode) MoveToBase(u *User) {
	u.Teleport(42, 9, 7)
}

func (m *NoneMode) Join(u *User) {
	u.SetGraphics(u.character.Graphics.Image)
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
		u.Send(toClient.CreateGameObject(self.GetCreateGameObject(false)))
		self.Send(toClient.CreateGameObject(u.GetCreateGameObject(false)))
	}
}

func (m *NoneMode) Hit(self *User, target *User) bool {
	fmt.Println("째ㅑㅂ쨉!")
	return true
}

func (m *NoneMode) UseItem(u *User) {

}

func (m *NoneMode) Result(winner int) {

}

func (m *NoneMode) Update() {
	if len(m.Room.Users) >= 1 {
		m.ChangeMode()
	} else {
		if m.Count%100 == 0 {
			m.Room.Publish(toClient.NoticeMessage("4명부터 시작합니다."))
		}
		m.Count++
		if m.Count == 10000 {
			m.Count = 0
		}
	}
}
