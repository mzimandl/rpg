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
	Blank             = 0
	Pending           = -1
)

func (level *Level) generateTile(x, y int, c rune) {
	var t Tile
	t.OverlayRune = Blank
	t.canSee = true
	t.canWalk = true
	pos := Pos{x, y}

	switch c {
	case ' ', '\t', '\n', '\r':
		t.Rune = Blank
	case DirtFloor:
		t.Rune = DirtFloor
	case StoneFloor:
		t.Rune = StoneFloor
	case StoneWall:
		t.Rune = StoneWall
		t.canSee = false
		t.canWalk = false
	case OldStoneWall:
		t.Rune = OldStoneWall
		t.canSee = false
		t.canWalk = false
	case ClosedDoor:
		t.OverlayRune = ClosedDoor
		t.canSee = false
		t.canWalk = false
		t.Rune = Pending
	case OpenedDoor:
		t.OverlayRune = OpenedDoor
		t.Rune = Pending
	case UpStair:
		t.OverlayRune = UpStair
		t.Rune = Pending
	case DownStair:
		t.OverlayRune = DownStair
		t.Rune = Pending

	case 's':
		level.Items[pos] = append(level.Items[pos], NewSword(pos))
		t.Rune = Pending
	case 'h':
		level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
		t.Rune = Pending

	case '@':
		level.Player.Pos = pos
		t.Rune = Pending
	case 'R':
		m := NewRat(pos)
		level.Monsters = append(level.Monsters, m)
		level.AliveMonstersPos[pos] = m
		t.Rune = Pending
	case 'S':
		m := NewSpider(pos)
		level.Monsters = append(level.Monsters, m)
		level.AliveMonstersPos[pos] = m
		t.Rune = Pending
	default:
		panic("Invalid character in the map: " + string(c))
	}
	level.Map[y][x] = t
}
