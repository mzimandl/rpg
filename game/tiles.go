package game

type Tile struct {
	Rune    rune
	visible bool
	// visited bool
}

const (
	StoneWall  rune = '#'
	DirtFloor       = '.'
	ClosedDoor      = '|'
	OpenedDoor      = '/'
	Blank           = 0
	Pending         = -1
)
