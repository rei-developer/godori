package game

type CharacterPos struct {
	x    int
	y    int
	dirX int
	dirY int
}

type Graphics struct {
	Image     string
	RedImage  string
	BlueImage string
}

type Character struct {
	Model int
	CharacterPos
	Graphics Graphics
	Dirty    bool
}

var dirTable = [3][3]int{
	{0, 4, 0},
	{2, 0, 8},
	{0, 6, 0},
}

func (c *Character) GetDirection(x int, y int) int {
	return dirTable[x+1][y+1]
}

func (c *Character) Setting(model int, redImage string, blueImage string) {
	c.Model = model
	c.SetPosition(0, -1)
	c.Graphics = Graphics{Image: blueImage, RedImage: redImage, BlueImage: blueImage}
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
