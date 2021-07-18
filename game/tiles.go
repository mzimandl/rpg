package game

type Tile struct {
	Rune    rune
	Visible bool
	Visited bool
}

const (
	StoneWall  rune = '#'
	DirtFloor       = '.'
	ClosedDoor      = '|'
	OpenedDoor      = '/'
	Blank           = 0
	Pending         = -1
)
