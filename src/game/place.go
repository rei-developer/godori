package game

import (
	"bytes"
	"encoding/binary"
	"godori.com/getty"
)

type Place struct {
	Index  int
	Room   *Room
	Events map[int]*Event
	Users  map[*getty.Client]*User
	Light  bool
}

func NewPlace(index int, r *Room) *Place {
	return &Place{
		Index:  index,
		Room:   r,
		Events: make(map[int]*Event),
		Users:  make(map[*getty.Client]*User),
	}
}

func (p *Place) AddEvent(e *Event) {
	p.Events[e.EventData.Id] = e
}

func (p *Place) RemoveEvent(e *Event) {
	delete(p.Events, e.EventData.Id)
}

func (p *Place) RemoveAllEvent() {
	p.Events = make(map[int]*Event)
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

func (p *Place) Update() {
	if len(p.Users) < 1 {
		return
	}
	var users map[*getty.Client]*User = make(map[*getty.Client]*User)
	for c, u := range p.Users {
		if u.character.Dirty {
			users[c] = u
		}
	}
	if len(users) <= 0 {
		return
	}
	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(0),
		uint8(len(users)),
	}
	for _, u := range users {
		u.character.Dirty = false
		data = append(data, int8(u.character.Model))
		data = append(data, int32(u.Index))
		data = append(data, int16(u.character.CharacterPos.x))
		data = append(data, int16(u.character.CharacterPos.y))
		data = append(data, int8(u.character.CharacterPos.dirX))
		data = append(data, int8(u.character.CharacterPos.dirY))
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		CheckError(err)
	}
	//fmt.Println(buf.Bytes())
	for _, u := range p.Users {
		u.Send(buf.Bytes())
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
