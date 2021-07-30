package game

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
	item := level.generateItem(pos, c)
	if item != nil {
		level.Items[pos] = append(level.Items[pos], item)
	} else {
		switch c {
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

		case '=':
			level.Storages[pos] = NewChest(pos, &StorageConf{items: level.Items[pos]})
			delete(level.Items, pos)

		default:
			panic("Invalid rune: " + string(c))
		}
	}
}

func (level *Level) generateItem(pos Pos, c rune) *Item {
	switch c {
	case 's':
		return NewSword(pos)
	case 'h':
		return NewHelmet(pos)
	case 'a':
		return NewArmor(pos)
	default:
		return nil
	}
}
