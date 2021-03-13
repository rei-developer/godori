package game

import toClient "godori.com/packet/toClient"

type PlaygroundMode struct {
	Room *Room
}

func (m *PlaygroundMode) MoveToBase(u *User) {
	u.Teleport(79, 36, 24)
}

func (m *PlaygroundMode) Join(u *User) {
	u.SetGraphics(u.BlueImage)
	m.MoveToBase(u)
}

func (m *PlaygroundMode) Leave(u *User) {

}

func (m *PlaygroundMode) DrawEvents(u *User) {

}

func (m *PlaygroundMode) DrawUsers(self *User) {
	for _, u := range m.Room.SameMapUsers(self.Place) {
		if u == self {
			return
		}
		u.Send(toClient.CreateGameObject(self.GetCreateGameObject(false)))
		self.Send(toClient.CreateGameObject(u.GetCreateGameObject(false)))
	}
}

func (m *PlaygroundMode) Hit(self *User, target *User) bool {
	return true
}

func (m *PlaygroundMode) UseItem(u *User) {

}

func (m *PlaygroundMode) Update() {

}
