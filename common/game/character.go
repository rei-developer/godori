package game

type characterPos struct {
	x int
	y int
	d int
}

type graphics string

type character struct {
	characterPos
	graphics
}

func (c *character) move(d int) {
}
