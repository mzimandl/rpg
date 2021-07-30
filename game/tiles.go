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

func (level *Level) generateTile(x, y int, c rune) {
	var t Tile
	t.OverlayRune = Blank
	t.canSee = true
	t.canWalk = true

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
		t.Rune = Pending
		t.canSee = false
		t.canWalk = false
	case OpenedDoor:
		t.OverlayRune = OpenedDoor
		t.Rune = Pending
	case UpStair:
		t.OverlayRune = UpStair
		t.Rune = Pending
	case DownStair:
		t.OverlayRune = DownStair
		t.Rune = Pending
	case StonePillar:
		t.OverlayRune = StonePillar
		t.Rune = Pending
		t.canWalk = false
	default:
		level.generateEntity(x, y, c)
		t.Rune = Pending
	}
	level.Map[y][x] = t
}

func (level *Level) generateEntity(x, y int, c rune) {
	pos := Pos{x, y}
	switch c {
	case 's':
		level.Items[pos] = append(level.Items[pos], NewSword(pos))
	case 'h':
		level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
	case 'a':
		level.Items[pos] = append(level.Items[pos], NewArmor(pos))

	case '@':
		level.Player.Pos = pos
	case 'R':
		m := NewRat(pos)
		level.Monsters = append(level.Monsters, m)
		level.AliveMonstersPos[pos] = m
	case 'S':
		m := NewSpider(pos)
		level.Monsters = append(level.Monsters, m)
		level.AliveMonstersPos[pos] = m
	default:
		panic("Invalid rune: " + string(c))
	}
}
