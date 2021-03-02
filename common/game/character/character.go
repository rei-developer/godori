package character

import (
	"fmt"
)

type CharacterPos struct {
	x int
	y int
	d int
}

type Graphics string

type Character struct {
	CharacterPos
	Graphics
}

func (c *Character) Move(d int) {
	fmt.Println("ㅇㅇ")
}
