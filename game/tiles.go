package game

type Tile struct {
	Rune        rune
	OverlayRune rune
	Visible     bool
	Visited     bool
	canWalk     bool
	canSee      bool
}

const (
	StoneWall    rune = '#'
	OldStoneWall      = '%'
	StoneFloor        = '_'
	DirtFloor         = '.'
	ClosedDoor        = '|'
	OpenedDoor        = '/'
	UpStair           = 'u'
	DownStair         = 'd'
	StonePillar       = 'I'
	Blank             = 0
	Pending           = -1
)
