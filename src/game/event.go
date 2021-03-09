package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type EventPureData struct {
	Events []EventData `json: "events"`
}

type EventData struct {
	Id       int    `json: "id"`
	Name     string `json: "name"`
	PureHp   int    `json: "pureHp"`
	Image    string `json: "image"`
	Command  string `json: "command"`
	Collider bool   `json: "collider"`
}

type Event struct {
	Model int
	Index int
	Name  string
	Room  *Room
	Place int
	Hp    int
	Character
	EventData *EventData
}

var nextEventIndex int
var EventDatas map[int]*EventData = make(map[int]*EventData)
var Events map[int]*Event = make(map[int]*Event)

func init() {
	fmt.Println("이벤트 로딩중...")
	jsonFile, err := os.Open("./lib/event.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var eventPureData EventPureData
	json.Unmarshal(byteValue, &eventPureData)
	for i := 0; i < len(eventPureData.Events); i++ {
		id := eventPureData.Events[i].Id
		EventDatas[id] = &EventData{
			Id:       id,
			Name:     eventPureData.Events[i].Name,
			PureHp:   eventPureData.Events[i].PureHp,
			Image:    eventPureData.Events[i].Image,
			Command:  eventPureData.Events[i].Command,
			Collider: eventPureData.Events[i].Collider,
		}
	}
	fmt.Println("이벤트 로딩 완료")
}

func NewEvent(r *Room, id int, place int, x int, y int) *Event {
	nextEventIndex++
	eventData := EventDatas[id]
	event := &Event{
		Model:     2,
		Index:     nextEventIndex,
		Name:      eventData.Name,
		Room:      r,
		Place:     place,
		Hp:        eventData.PureHp,
		EventData: eventData,
	}
	Events[nextEventIndex] = event
	event.Character.Setting(event.Model, event.EventData.Image, event.EventData.Image)
	event.SetPosition(x, y)
	return event
}

func (e *Event) Remove() bool {
	_, ok := Events[e.Index]
	if ok {
		delete(Events, e.Index)
	}
	return ok
}

func (e *Event) GetCreateGameObject() (model int, index int, name string, clanName string, team int, level int, image string, x int, y int, dirX int, dirY int, collider bool) {
	model = e.Model
	index = e.Index
	name = e.Name
	clanName = ""
	team = 0
	level = 0
	image = e.Graphics.Image
	x = e.x
	y = e.y
	dirX = e.CharacterPos.dirX
	dirY = e.CharacterPos.dirY
	collider = e.EventData.Collider
	return
}

func (e *Event) Do(r *Room, u *User) bool {
	return CallFuncByName(&Action{}, e.EventData.Command, r, u, e) != nil
}
