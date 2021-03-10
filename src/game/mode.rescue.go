package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	mapType "godori.com/util/constant/mapType"
	modeType "godori.com/util/constant/modeType"
	teamType "godori.com/util/constant/teamType"
	pix "godori.com/util/pix"
)

type RescueMode struct {
	Room      *Room
	MapType   int
	ModeType  int
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
		MapType:   mapType.TATAMI, //pType,
		ModeType:  modeType.RESCUE,
		RedScore:  0,
		BlueScore: 0,
		RedUsers:  make(map[*getty.Client]*User),
		BlueUsers: make(map[*getty.Client]*User),
		State:     STATE_READY,
		Tick:      0,
		Count:     201,
		MaxCount:  230,
	}
}

func (m *RescueMode) InitEvent() {
	type ED struct {
		Id    int `json: "id"`
		Place int `json: "place"`
		X     int `json: "x"`
		Y     int `json: "y"`
	}
	type EPD struct {
		Events []ED `json: "events"`
	}
	var events map[int]ED = make(map[int]ED)
	jsonFile, _ := os.Open("./lib/events/rescue.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var eventPureData EPD
	json.Unmarshal(byteValue, &eventPureData)
	for i := 0; i < len(eventPureData.Events); i++ {
		id := eventPureData.Events[i].Id
		events[id] = ED{
			Id:    id,
			Place: eventPureData.Events[i].Place,
			X:     eventPureData.Events[i].X,
			Y:     eventPureData.Events[i].Y,
		}
	}
	for _, e := range events {
		fmt.Println(e)
		m.Room.AddEvent(NewEvent(m.Room, e.Id, e.Place, e.X, e.Y))
	}
}

func (m *RescueMode) AddUser(u *User) {
	if tType, ok := u.GameData["team"]; ok {
		if tType == teamType.RED {
			m.RedUsers[u.Client] = u
		} else if tType == teamType.BLUE {
			m.BlueUsers[u.Client] = u
		}
	}
}

func (m *RescueMode) RemoveUser(u *User) {
	if tType, ok := u.GameData["team"]; ok {
		if tType == teamType.RED {
			delete(m.RedUsers, u.Client)
		} else if tType == teamType.BLUE {
			delete(m.BlueUsers, u.Client)
		}
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
	if tType, ok := u.GameData["team"]; ok {
		if tType == teamType.RED {
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
		} else if tType == teamType.BLUE {
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
		u.SetGraphics(u.BlueImage)
		m.AddUser(u)
		m.MoveToBase(u)
	case STATE_GAME:
		u.GameData["caught"] = true
		u.SetGraphics(u.BlueImage)
		m.AddUser(u)
		m.MoveToPrison(u)
		m.RedScore++
		u.Send(toClient.NoticeMessage("감옥에 갇힌 인질을 전원 구출하라."))
	}
	u.PublishMap(toClient.SetGameTeam(u.Index, u.GameData["team"].(int)))
	u.Publish(toClient.ModeData(m.ModeType))
}

func (m *RescueMode) Leave(u *User) {
	m.RemoveUser(u)
	if caught, ok := u.GameData["caught"]; ok && caught.(bool) {
		m.RedScore--
		fmt.Println("레드 스코어 삭감") // 테스트 로그
	}
	u.GameData = nil
	u.SetGraphics(u.BlueImage)
	u.Publish(toClient.UpdateModeCount(m.RedScore))
}

func (m *RescueMode) DrawEvents(u *User) {
	fmt.Println("호출 1")
	for _, e := range m.Room.GetPlace(u.Place).Events {
		fmt.Println(e)
		u.Send(toClient.CreateGameObject(e.GetCreateGameObject()))
	}
}

func (m *RescueMode) DrawUsers(self *User) {
	selfHide := false
	for _, u := range m.Room.SameMapUsers(self.Place) {
		if u == self {
			return
		}
		userHide := false
		if u.GameData["team"] != self.GameData["team"] {
			if !(self.UserData.Admin > 1 && u.UserData.Admin > 1) {
				if self.UserData.Admin > 1 {
					selfHide = true
				} else if u.UserData.Admin > 1 {
					userHide = true
				} else {
					selfHide, userHide = true, true
				}
			}
		}
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
	target.Send(toClient.DeadAnimation())
	self.Send(toClient.NoticeMessage(pix.Maker(target.UserData.Name, "를", "을") + " 인질로 붙잡았다."))
	self.Send(toClient.PlaySound("Eat"))
	self.Broadcast(toClient.NoticeMessage(pix.Maker(target.UserData.Name, "가", "이") + " 인질로 붙잡혔다!"))
	self.Broadcast(toClient.PlaySound("Shock"))
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
	fmt.Println("끝", winner)
	m.State = STATE_RESULT
	for _, u := range m.Room.Users {
		u.Room = nil
		u.GameData["result"] = true
	}
	m.Room.Remove()
	//for _, u := range m.RedUsers {
	//
	//} TODO : 점수 가산
}

func (m *RescueMode) Update() {
	m.Tick++
	if m.Tick%10 != 0 {
		return
	}
	m.Tick = 0
	switch m.State {
	case STATE_READY:
		if m.Count <= 230 && m.Count > 200 {
			if m.Count == 210 {
				m.Room.Publish(toClient.PlaySound("GhostsTen"))
			}
			m.Room.Publish(toClient.NoticeMessage(strconv.Itoa(m.Count - 200)))
		} else if m.Count == 200 {
			m.Room.Lock = true
			m.State = STATE_GAME
			for _, u := range m.Room.Mode.Sample(m.BlueUsers, (len(m.BlueUsers)/5)+1) {
				m.RemoveUser(u)
				u.GameData["team"] = teamType.RED
				m.AddUser(u)
				u.SetGraphics(u.RedImage)
				// TODO : 장농
				//u.Send()
			}
			for _, u := range m.RedUsers {
				u.Send(toClient.NoticeMessage("단 한 명의 인간이라도 감옥에 가둬라."))
			}
			for _, u := range m.BlueUsers {
				u.Send(toClient.NoticeMessage("감옥에 갇힌 인질을 전원 구출하라."))
			}
			m.Room.Publish(toClient.PlaySound("A4"))
		}
	case STATE_GAME:
		for _, u := range m.RedUsers {
			if GameMaps[u.Place].RangePortal(u.X, u.Y, 2) {
				if count, ok := u.GameData["count"]; ok {
					u.GameData["count"] = count.(int) + 1
					count = u.GameData["count"]
					if count.(int) >= 3 && count.(int) <= 5 {
						u.Send(toClient.InformMessage("<color=red>경고!!! 포탈 주변을 막지 마십시오.</color>"))
						u.Send(toClient.PlaySound("Warn"))
					} else if count.(int) > 5 {
						u.GameData["count"] = 0
						m.MoveToBase(u)
						u.Send(toClient.InformMessage("<color=red>지속적인 게임 플레이 방해로 인해 본진으로 추방되었습니다.</color>"))
					}
				}
			} else {
				if count, ok := u.GameData["count"]; ok {
					if count.(int) > 0 {
						u.GameData["count"] = count.(int) - 1
					} else if count.(int) < 0 {
						u.GameData["count"] = 0
					}
				}
			}
		}
		if m.Count == 15 || m.Count%40 == 5 {
			m.Room.Publish(toClient.InformMessage("<color=#B5E61D>잠시 후 인질 구출이 가능해집니다...</color>"))
		} else if m.Count == 10 || m.Count%40 == 0 {
			m.Room.Publish(toClient.InformMessage("<color=#B5E61D>인질 구출이 가능합니다!</color>"))
			m.Room.Publish(toClient.PlaySound("thump"))
		}
		//if len(m.RedUsers) == 0 {
		//	m.Result(teamType.BLUE)
		//} else if len(m.BlueUsers) == 0 || m.RedScore == len(m.BlueUsers) {
		//	m.Result(teamType.RED)
		//} else if m.Count == 5 {
		//	m.Room.Publish(toClient.PlaySound("Second"))
		//} else if m.Count == 0 {
		//	if m.RedScore > 0 {
		//		m.Result(teamType.RED)
		//	} else {
		//		m.Result(teamType.BLUE)
		//	}
		//}
	}
	m.Count--
}
