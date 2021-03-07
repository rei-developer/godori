package game

import (
	"fmt"
	"godori.com/getty"
)

type Place struct {
	Index  int
	Room   *Room
	Users  map[*getty.Client]*User
	Events []string
	Light  bool
}

func NewPlace(index int, r *Room) *Place {
	return &Place{
		Index: index,
		Room:  r,
		Users: make(map[*getty.Client]*User),
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
	fmt.Println("루룰랄라")
	
	if len(p.Users) < 1 {
		return
	}
}
