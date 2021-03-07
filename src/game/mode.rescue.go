package game

import (
	"fmt"

	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	mapType "godori.com/util/constant/mapType"
	teamType "godori.com/util/constant/teamType"
)

type RescueMode struct {
	Room      *Room
	MapType   int
	RedScore  int
	BlueScore int
	RedUsers  map[*getty.Client]*User
	BlueUsers map[*getty.Client]*User
	State     int
	Tick      int
	Count     int
	MaxCount  int
}

func NewRescueMode(r *Room, pType int) *RescueMode {
	return &RescueMode{
		Room:      r,
		MapType:   pType,
		RedScore:  0,
		BlueScore: 0,
		RedUsers:  make(map[*getty.Client]*User),
		BlueUsers: make(map[*getty.Client]*User),
		State:     STATE_READY,
		Tick:      0,
		Count:     230,
		MaxCount:  230,
	}
}

func (m *RescueMode) AddUser(u *User) {
	tType := u.GameData["team"]
	if tType == teamType.RED {
		m.RedUsers[u.client] = u
	} else if tType == teamType.BLUE {
		m.BlueUsers[u.client] = u
	}
}

func (m *RescueMode) RemoveUser(u *User) {
	tType := u.GameData["team"]
	if tType == teamType.RED {
		delete(m.RedUsers, u.client)
	} else if tType == teamType.BLUE {
		delete(m.BlueUsers, u.client)
	}
}

func (m *RescueMode) SetUserGameData(u *User) {
	u.GameData = make(map[string]interface{})
	u.GameData["team"] = teamType.BLUE
	u.GameData["state"] = 0
	u.GameData["hp"] = 100
	u.GameData["spawn"] = 10
	u.GameData["count"] = 0
	u.GameData["caught"] = false
	u.GameData["result"] = false
}

func (m *RescueMode) MoveToBase(u *User) {
	if u.UserData.Team == teamType.RED {
		switch m.MapType {
		case mapType.ASYLUM:
			u.Teleport(29, 9, 19)
		case mapType.TATAMI:
			u.Teleport(54, 10, 5)
		case mapType.GON:
			u.Teleport(75, 20, 26)
		case mapType.LABORATORY:
			u.Teleport(86, 9, 11)
		case mapType.SCHOOL:
			u.Teleport(115, 13, 9)
		case mapType.MINE:
			u.Teleport(172, 6, 8)
		case mapType.ISLAND:
			u.Teleport(189, 7, 7)
		case mapType.MANSION:
			u.Teleport(226, 10, 9)
		case mapType.DESERT:
			u.Teleport(244, 9, 11)
		}
	} else if u.UserData.Team == teamType.BLUE {
		switch m.MapType {
		case mapType.ASYLUM:
			u.Teleport(2, 8, 13)
		case mapType.TATAMI:
			u.Teleport(42, 9, 7)
		case mapType.GON:
			u.Teleport(60, 16, 11)
		case mapType.LABORATORY:
			u.Teleport(99, 10, 8)
		case mapType.SCHOOL:
			u.Teleport(149, 14, 8)
		case mapType.MINE:
			u.Teleport(154, 9, 8)
		case mapType.ISLAND:
			u.Teleport(199, 10, 8)
		case mapType.MANSION:
			u.Teleport(238, 17, 8)
		case mapType.DESERT:
			u.Teleport(249, 7, 17)
		}
	}
}

func (m *RescueMode) MoveToPrison(u *User) {
	switch m.MapType {
	case mapType.ASYLUM:
		u.Teleport(13, 11, 15)
	case mapType.TATAMI:
		u.Teleport(57, 21, 6)
	case mapType.GON:
		u.Teleport(74, 14, 12)
	case mapType.LABORATORY:
		u.Teleport(96, 7, 30)
	case mapType.SCHOOL:
		u.Teleport(122, 6, 12)
	case mapType.MINE:
		u.Teleport(169, 13, 6)
	case mapType.ISLAND:
		u.Teleport(191, 11, 7)
	case mapType.MANSION:
		u.Teleport(217, 25, 7)
	case mapType.DESERT:
		u.Teleport(255, 20, 17)
	}
}

func (m *RescueMode) MoveToOutside(u *User) {
	switch m.MapType {
	case mapType.ASYLUM:
		u.Teleport(19, 9, 8)
	case mapType.TATAMI:
		u.Teleport(47, 17, 6)
	case mapType.GON:
		u.Teleport(72, 15, 8)
	case mapType.LABORATORY:
		u.Teleport(89, 16, 12)
	case mapType.SCHOOL:
		u.Teleport(118, 5, 15)
	case mapType.MINE:
		u.Teleport(166, 34, 31)
	case mapType.ISLAND:
		u.Teleport(174, 12, 7)
	case mapType.MANSION:
		u.Teleport(218, 19, 8)
	case mapType.DESERT:
		u.Teleport(243, 13, 22)
	}
}

func (m *RescueMode) Join(u *User) {
	m.SetUserGameData(u)
	switch m.State {
	case STATE_READY:
		u.SetGraphics(u.character.Graphics.BlueImage)
		m.AddUser(u)
		m.MoveToBase(u)
	case STATE_GAME:
		u.GameData["caught"] = true
		u.SetGraphics(u.character.Graphics.BlueImage)
		m.AddUser(u)
		m.MoveToPrison(u)
		m.RedScore++
		u.Send(toClient.NoticeMessage("감옥에 갇힌 인질을 전원 구출하라."))
	}
	//u.PublishMap(toClient.SetGameTeam()) TODO :
}

func (m *RescueMode) Leave(u *User) {
	m.RemoveUser(u)
	if u.GameData["caught"] == true {
		m.RedScore--
	}
	u.SetGraphics(u.character.Graphics.BlueImage)
	// TODO : score publish
}

func (m *RescueMode) DrawEvents(u *User) {

}

func (m *RescueMode) DrawUsers(self *User) {
	selfHide := false
	for _, u := range m.Room.SameMapUsers(self.place) {
		if u == self {
			return
		}
		userHide := false
		if self.GameData["team"] != u.GameData["team"] {
			selfHide = true
			userHide = true
		} // TODO
		u.Send(toClient.CreateGameObject(self.GetCreateGameObject(userHide)))
		self.Send(toClient.CreateGameObject(u.GetCreateGameObject(selfHide)))
	}
}

func (m *RescueMode) Hit(self *User, target *User) bool {
	if self.GameData["team"] == teamType.BLUE {
		return true
	}
	if self.GameData["team"] == target.GameData["team"] {
		return false
	}
	if target.GameData["caught"] == true {
		return true
	}
	m.MoveToPrison(target)
	target.GameData["caught"] = true
	// TODO : dead animation
	self.Send(toClient.NoticeMessage(target.UserData.Name + "을 인질로 붙잡았다."))
	// TODO
	self.Broadcast(toClient.NoticeMessage(target.UserData.Name + "을 인질로 붙잡혔다!"))
	// TODO
	switch target.GameData["state"] {
	case 1:
		// TODO : 장농
	default:
		// TODO : 기본
	}
	m.RedScore++
	//self.Publish()
	return true
}

func (m *RescueMode) UseItem(u *User) {

}

func (m *RescueMode) Result(winner int) {

}

func (m *RescueMode) Update() {
	fmt.Println("흐으음 레스큐다 레스큐")
}
