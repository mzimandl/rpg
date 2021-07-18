package game

type Tile rune

const (
	StoneWall  Tile = '#'
	DirtFloor  Tile = '.'
	ClosedDoor Tile = '|'
	OpenedDoor Tile = '/'
	Blank      Tile = 0
	Pending    Tile = -1
)
