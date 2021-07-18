package game

import (
	"bufio"
	"math"
	"os"
)

type Level struct {
	Map              [][]Tile
	Player           *Player
	Monsters         []*Monster
	AliveMonstersPos map[Pos]*Monster
	Events           []string
	Debug            map[Pos]bool
}

func NewLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}
	level := &Level{}
	level.Map = make([][]Tile, len(levelLines))
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}
	level.Debug = make(map[Pos]bool)
	level.AliveMonstersPos = make(map[Pos]*Monster)
	level.Monsters = make([]*Monster, 0)
	level.Events = make([]string, 0)

	for y := range level.Map {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t.Rune = Blank
			case '#':
				t.Rune = StoneWall
			case '|':
				t.Rune = ClosedDoor
			case '/':
				t.Rune = ClosedDoor
			case '.':
				t.Rune = DirtFloor
			case '@':
				level.Player = NewPlayer(Pos{x, y})
				t.Rune = Pending
			case 'R':
				m := NewRat(Pos{x, y})
				level.Monsters = append(level.Monsters, m)
				level.AliveMonstersPos[Pos{x, y}] = m
				t.Rune = Pending
			case 'S':
				m := NewSpider(Pos{x, y})
				level.Monsters = append(level.Monsters, m)
				level.AliveMonstersPos[Pos{x, y}] = m
				t.Rune = Pending
			default:
				panic("Invalid character in the map: " + string(c))
			}
			level.Map[y][x] = t
		}
	}

	if level.Player == nil {
		panic("Missing player in the map!")
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune == Pending {
				level.Map[y][x] = level.BfsFloor(Pos{x, y})
			}
		}
	}
	level.resolveVisibility()

	return level
}

func (level *Level) inRange(pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func (level *Level) canWalk(pos Pos) bool {
	if level.inRange(pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoneWall, ClosedDoor:
			return false
		}
		return true
	}
	return false
}

func (level *Level) canSeeThrough(pos Pos) bool {
	if level.inRange(pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoneWall, ClosedDoor, Blank:
			return false
		default:
			return true
		}
	}
	return false
}

func (level *Level) bresenhamVisibility(start Pos, end Pos) {
	steep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))
	if steep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}

	deltaX := int(math.Abs(float64(end.X - start.X)))
	deltaY := int(math.Abs(float64(end.Y - start.Y)))
	err := 0
	y := start.Y
	ystep := 1
	if start.Y >= end.Y {
		ystep = -1
	}

	reversed := start.X > end.X
	if reversed {
		for x := start.X; x > end.X; x-- {
			var pos Pos
			if steep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Map[pos.Y][pos.X].Visible = true
			level.Map[pos.Y][pos.X].Visited = true
			if !level.canSeeThrough(pos) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	} else {
		for x := start.X; x < end.X; x++ {
			var pos Pos
			if steep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Map[pos.Y][pos.X].Visible = true
			level.Map[pos.Y][pos.X].Visited = true
			if !level.canSeeThrough(pos) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	}
}

func (level *Level) checkDoor(pos Pos) {
	t := level.Map[pos.Y][pos.X]
	switch t.Rune {
	case ClosedDoor:
		level.Map[pos.Y][pos.X].Rune = OpenedDoor
	}
}

func (level *Level) resolveMovement(pos Pos) {
	monster, exists := level.AliveMonstersPos[pos]
	if exists {
		event := level.Player.Attack(&monster.Character)
		level.addEvent(event)
		if !monster.IsAlive() {
			monster.Die(level)
		}
		if !level.Player.IsAlive() {
			level.addEvent("DED")
		}
	} else if level.canWalk(pos) {
		level.Player.Move(pos, level)
		level.resetVisibility()
		level.resolveVisibility()
	} else {
		level.checkDoor(pos)
		level.resetVisibility()
		level.resolveVisibility()
	}
}

func (level *Level) resetVisibility() {
	for y, row := range level.Map {
		for x := range row {
			level.Map[y][x].Visible = false
		}
	}
}

func (level *Level) resolveVisibility() {
	pos := level.Player.Pos
	dist := level.Player.SightRange
	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			if xDelta*xDelta+yDelta*yDelta <= dist*dist {
				level.bresenhamVisibility(pos, Pos{x, y})
			}
		}
	}
	level.bresenhamVisibility(level.Player.Pos, Pos{level.Player.X, level.Player.Y + level.Player.SightRange})
}

func (level *Level) getNeighbors(pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if level.canWalk(left) {
		neighbors = append(neighbors, left)
	}
	if level.canWalk(right) {
		neighbors = append(neighbors, right)
	}
	if level.canWalk(up) {
		neighbors = append(neighbors, up)
	}
	if level.canWalk(down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func (level *Level) BfsFloor(start Pos) Tile {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile.Rune {
		case DirtFloor:
			return Tile{DirtFloor, false, false}
		default:
		}

		frontier = frontier[1:]
		for _, next := range level.getNeighbors(current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
	return Tile{DirtFloor, false, false}
}

func (level *Level) addEvent(s string) {
	level.Events = append(level.Events, s)
	if len(level.Events) > 25 {
		level.Events = level.Events[len(level.Events)-25:]
	}
}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {
		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)

			for current != start {
				path = append(path, current)
				current = cameFrom[current]
			}
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return path
		}

		for _, next := range level.getNeighbors(current) {
			newCost := costSoFar[current]

			_, exists := level.AliveMonstersPos[next]
			if exists {
				newCost += 10
			} else {
				newCost += 1
			}

			_, exists = costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				cameFrom[next] = current
			}
		}
	}

	return nil
}
