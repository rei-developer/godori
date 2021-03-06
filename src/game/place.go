package game

import (
	"godori.com/getty"
)

type Place struct {
	Index     int
	RoomIndex int
	Users     map[*getty.Client]*User
	Events    []string
	Light     bool
}

func NewPlace(index int, rIndex int) *Place {
	return &Place{
		Index:     index,
		RoomIndex: rIndex,
		Users:     make(map[*getty.Client]*User),
	}
}

func (p *Place) AddUser(u *User) {
	p.Users[u.client] = u
}

func (p *Place) RemoveUser(u *User) {
	delete(p.Users, u.client)
}

func (p *Place) RemoveAllUser() {
	p.Users = make(map[*getty.Client]*User)
}

func (p *Place) AddEvent() {
	// TODO
}

func (p *Place) RemoveEvent() {
	// TODO
}

func (p *Place) RemoveAllEvent() {
	// TODO
}

func (p *Place) Update() {
	if len(p.Users) < 1 {
		return
	}

	// TODO : 동기화
}
