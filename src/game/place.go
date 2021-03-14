package game

import (
	"bytes"
	"encoding/binary"
	"log"
	"sync"

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
	p.Users[u.Client] = u
}

func (p *Place) RemoveUser(u *User) {
	delete(p.Users, u.Client)
}

func (p *Place) RemoveAllUser() {
	p.Users = make(map[*getty.Client]*User)
}

var lock = sync.RWMutex{}

func (p *Place) Update() {
	if len(p.Users) < 1 {
		return
	}
	var events map[int]*Event = make(map[int]*Event)
	var users map[*getty.Client]*User = make(map[*getty.Client]*User)
	for i, e := range p.Events {
		lock.RLock()
		e.Update()
		if e.Dirty {
			events[i] = e
		}
		lock.RUnlock()
	}
	for c, u := range p.Users {
		lock.RLock()
		if u.Dirty {
			users[c] = u
		}
		lock.RUnlock()
	}
	if len(events)+len(users) <= 0 {
		return
	}
	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint8(0),
		uint8(len(events) + len(users)),
	}
	for _, e := range events {
		lock.RLock()
		e.Dirty = false
		data = append(data, int8(e.Model))
		data = append(data, int32(e.Index))
		data = append(data, int16(e.X))
		data = append(data, int16(e.Y))
		data = append(data, int8(e.DirX))
		data = append(data, int8(e.DirY))
		lock.RUnlock()
	}
	for _, u := range users {
		lock.RLock()
		u.Dirty = false
		data = append(data, int8(u.Model))
		data = append(data, int32(u.Index))
		data = append(data, int16(u.X))
		data = append(data, int16(u.Y))
		data = append(data, int8(u.DirX))
		data = append(data, int8(u.DirY))
		lock.RUnlock()
	}
	for _, v := range data {
		lock.RLock()
		err := binary.Write(buf, binary.LittleEndian, v)
		CheckError(err)
		lock.RUnlock()
	}
	for _, u := range p.Users {
		lock.RLock()
		u.Client.Send(buf.Bytes())
		lock.RUnlock()
	}
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}
