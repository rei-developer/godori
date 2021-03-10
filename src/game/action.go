package game

import (
	"fmt"
	"log"
	"reflect"

	toClient "godori.com/packet/toClient"
	teamType "godori.com/util/constant/teamType"
)

type Action struct{}

func (a *Action) Door(r *Room, self *User, target *Event) {
	if team, ok := self.GameData["team"]; ok {
		if open, ok := target.GameData["open"]; ok {
			const (
				openSound  = "door03"
				closeSound = "door04"
				knockSound = "door06"
			)
			if open.(bool) {
				//if team == teamType.RED {
				//	return
				//}
				self.PublishMap(toClient.PlaySound(closeSound))
				target.Move(-1, 0)
				target.GameData["open"] = false
			} else {
				rand := 0// cMath.Rand(10)
				if team == teamType.BLUE || rand == 0 {
					self.PublishMap(toClient.PlaySound(openSound))
					target.Move(1, 0)
					target.GameData["open"] = true
				} else {
					self.PublishMap(toClient.PlaySound(knockSound))
				}
			}
		} else {
			target.GameData["open"] = false
			a.Door(r, self, target)
		}
	}
}

func (a *Action) Rescue(r *Room, self *User, target *Event) {
	fmt.Println("ㅁㄴㅇㄻㄴㅇㄹ")
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
