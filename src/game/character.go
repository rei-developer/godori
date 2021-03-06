package game

type CharacterPos struct {
	x    int
	y    int
	dirX int
	dirY int
}

type Graphics string

type Character struct {
	CharacterPos
	Graphics
	Dirty bool
}

var dirTable = [3][3]int{
	{0, 4, 0},
	{2, 0, 8},
	{0, 6, 0},
}

func (c *Character) GetDirection(x int, y int) int {
	return dirTable[x+1][y+1]
}

func (c *Character) SetPosition(x int, y int) {
	c.CharacterPos.x = x
	c.CharacterPos.y = y
	c.Dirty = true
}

func (c *Character) Turn(dirX int, dirY int) {
	c.CharacterPos.dirX = dirX
	c.CharacterPos.dirY = dirY
	c.Dirty = true
}

func (c *Character) Move(x int, y int) {
	c.CharacterPos.x += x
	c.CharacterPos.y += y
	c.Dirty = true
}
