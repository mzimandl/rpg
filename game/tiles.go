package game

type Tile struct {
	Rune        rune
	OverlayRune rune
	Visible     bool
	Visited     bool
}

const (
	StoneWall  rune = '#'
	DirtFloor       = '.'
	ClosedDoor      = '|'
	OpenedDoor      = '/'
	UpStair         = 'u'
	DownStair       = 'd'
	Blank           = 0
	Pending         = -1
)
