package game

type PlaygroundMode struct {
	Room *Room
}

func (m *PlaygroundMode) MoveToBase(u *User) {

}

func (m *PlaygroundMode) Join(u *User) {

}

func (m *PlaygroundMode) Leave(u *User) {

}

func (m *PlaygroundMode) DrawEvents(u *User) {

}

func (m *PlaygroundMode) DrawUsers(self *User) {

}

func (m *PlaygroundMode) Hit(self *User, target *User) bool {
	return true
}

func (m *PlaygroundMode) UseItem(u *User) {

}

func (m *PlaygroundMode) Result(winner int) {

}

func (m *PlaygroundMode) Update() {

}
