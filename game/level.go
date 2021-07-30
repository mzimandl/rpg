package game

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type Level struct {
	Map      [][]Tile
	Player   *Player
	Monsters []*Monster

	AliveMonstersPos map[Pos]*Monster
	Items            map[Pos][]*Item
	Portals          map[Pos]*LevelPos
	Storages         map[Pos]*Storage

	Log        []string
	Debug      map[Pos]bool
	LastEvents []GameEvent
}

type GameEvent int

const (
	Move GameEvent = iota
	DoorOpen
	DoorClose
	Attack
	Portal
	PickUp
	DropDown
	Equip
	TakeOff
)

type LevelPos struct {
	level *Level
	pos   Pos
}

func NewLevelFromFile(filename string, player *Player) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	entityLines := make([]string, 0)
	longestRow := 0
	index := 0
	end := false
	for scanner.Scan() {
		text := scanner.Text()
		if text == "ENTITIES:" {
			end = true
			continue
		}

		if end {
			entityLines = append(entityLines, text)
		} else {
			levelLines = append(levelLines, text)
			if len(levelLines[index]) > longestRow {
				longestRow = len(levelLines[index])
			}
			index++
		}
	}

	level := &Level{}
	level.Map = make([][]Tile, len(levelLines))
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}
	level.Player = player
	level.AliveMonstersPos = make(map[Pos]*Monster)
	level.Portals = make(map[Pos]*LevelPos)
	level.Storages = make(map[Pos]*Storage)
	level.Items = make(map[Pos][]*Item)
	level.Debug = make(map[Pos]bool)

	for y := range level.Map {
		line := levelLines[y]
		for x, c := range line {
			level.generateTile(x, y, c)
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune == Pending {
				level.Map[y][x].Rune = level.BfsFloor(Pos{x, y})
			}
		}
	}

	for _, line := range entityLines {
		splitCXY := strings.Split(line, ",")
		if len(splitCXY) < 3 {
			continue
		}
		c := rune(splitCXY[0][0])
		x, err := strconv.Atoi(splitCXY[1])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(splitCXY[2])
		if err != nil {
			panic(err)
		}
		level.generateEntity(x, y, c)
	}

	return level
}

func (level *Level) inRange(pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func (level *Level) canWalk(pos Pos) bool {
	if level.inRange(pos) {
		t := level.Map[pos.Y][pos.X]
		return t.canWalk
	}
	return false
}

func (level *Level) canSeeThrough(pos Pos) bool {
	if level.inRange(pos) {
		t := level.Map[pos.Y][pos.X]
		return t.canSee
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

func (level *Level) checkClosedDoor(pos Pos) bool {
	t := level.Map[pos.Y][pos.X]
	switch t.OverlayRune {
	case ClosedDoor:
		t.OverlayRune = OpenedDoor
		t.canSee = true
		t.canWalk = true
		level.Map[pos.Y][pos.X] = t
		level.LastEvents = append(level.LastEvents, DoorOpen)
		return true
	}
	return false
}

func (level *Level) checkOpenedDoor(pos Pos) bool {
	t := level.Map[pos.Y][pos.X]
	switch t.OverlayRune {
	case OpenedDoor:
		t.OverlayRune = ClosedDoor
		t.canSee = false
		t.canWalk = false
		level.Map[pos.Y][pos.X] = t
		level.LastEvents = append(level.LastEvents, DoorClose)
		return true
	}
	return false
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
	dist := level.Player.SightRange + 2
	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			if xDelta*xDelta+yDelta*yDelta < dist*dist {
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

func (level *Level) BfsFloor(start Pos) rune {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile.Rune {
		case DirtFloor:
			return DirtFloor
		case StoneFloor:
			return StoneFloor
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
	return DirtFloor
}

func (level *Level) addEvent(s string) {
	level.Log = append(level.Log, s)
	if len(level.Log) > 25 {
		level.Log = level.Log[len(level.Log)-25:]
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
