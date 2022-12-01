package game

import (
	"log"
	"reflect"
	"strconv"

	toClient "godori/packet/toClient"
	teamType "godori/util/constant/teamType"
	cMath "godori/util/math"
)

type Action struct{}

func (a *Action) Door(r *Room, self *User, target *Event) {
	if team, ok := self.GameData["team"]; ok {
		if open, ok := target.GameData["open"]; ok {
			const (
				OPEN_SOUND  = "door03"
				CLOSE_SOUND = "door04"
				KNOCK_SOUND = "door06"
			)
			if open.(bool) {
				if team == teamType.RED {
					return
				}
				self.PublishMap(toClient.PlaySound(CLOSE_SOUND))
				target.Move(-1, 0)
				target.GameData["open"] = false
			} else {
				rand := cMath.Rand(10)
				if team == teamType.BLUE || rand == 0 {
					self.PublishMap(toClient.PlaySound(OPEN_SOUND))
					target.Move(1, 0)
					target.GameData["open"] = true
				} else {
					self.PublishMap(toClient.PlaySound(KNOCK_SOUND))
				}
			}
		} else {
			target.GameData["open"] = false
			a.Door(r, self, target)
		}
	}
}

func (a *Action) Rescue(r *Room, self *User, target *Event) {
	if team, ok := self.GameData["team"]; ok {
		if caught, ok := self.GameData["caught"]; ok {
			if team == teamType.RED || caught.(bool) {
				return
			}
			m := r.Mode.Mode.(*RescueMode)
			if !m.Caught {
				self.Send(toClient.InformMessage("<color=#B5E61D>아직 인질을 구출할 수 없습니다.</color>"))
				return
			} else if m.RedScore < 1 {
				self.Send(toClient.InformMessage("<color=#B5E61D>붙잡힌 인질이 없습니다.</color>"))
				return
			}
			for _, u := range m.RedUsers {
				m.MoveToOutside(u)
			}
			count := 0
			for _, u := range m.BlueUsers {
				if caught, ok := u.GameData["caught"]; ok {
					if caught.(bool) {
						u.Teleport(self.Place, self.X, self.Y)
						u.GameData["caught"] = false
						count++
					}
				}
			}
			m.RedScore = 0
			m.Caught = false
			self.Publish(toClient.NoticeMessage(self.UserData.Name + " 인질 " + strconv.Itoa(count) + "명 구출!"))
			self.Publish(toClient.ComboMessage(self.UserData.Name + " (" + strconv.Itoa(count) + "명)"))
			self.Publish(toClient.PlaySound("Rescue"))
			self.Publish(toClient.UpdateModeCount(0))
			self.Score.Rescue += count
			self.Score.RescueCombo++
		}
	}
}

func (a *Action) Light(r *Room, self *User, target *Event) {

}

func CallFuncByName(funcName string, params ...interface{}) []reflect.Value {
	m := reflect.ValueOf(&Action{}).MethodByName(funcName)
	if !m.IsValid() {
		log.Println("Method not found", funcName)
		return nil
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	return m.Call(in)
}
