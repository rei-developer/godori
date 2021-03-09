package game

import (
	"fmt"
	"log"
	"reflect"
)

type Action struct{}

func (a *Action) Rescue(r *Room, self *User, target *Event) {
	fmt.Println("ㅁㄴㅇㄻㄴㅇㄹ")
}

func CallFuncByName(class interface{}, funcName string, params ...interface{}) []reflect.Value {
	myClassValue := reflect.ValueOf(class)
	m := myClassValue.MethodByName(funcName)
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
