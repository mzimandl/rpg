package game

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Game struct {
	LevelChan    chan *Level
	InputChan    chan *Input
	Player       *Player
	Levels       map[string]*Level
	CurrentLevel *Level
}

func NewGame() *Game {
	levelChan := make(chan *Level)
	inputChan := make(chan *Input)

	levels := make(map[string]*Level)
	filenames, err := filepath.Glob("game/maps/*.map")
	if err != nil {
		panic(err)
	}

	player := NewPlayer(Pos{0, 0})
	for _, filename := range filenames {
		extIndex := strings.LastIndex(filename, ".map")
		lastSlashIndex := strings.LastIndex(filename, "/")
		levelName := filename[lastSlashIndex+1 : extIndex]
		levels[levelName] = NewLevelFromFile(filename, player)
	}
	game := &Game{levelChan, inputChan, player, levels, nil}
	game.loadWorldFile()
	return game
}

type InputType int

const (
	INone InputType = iota
	IMove
	IAction
	ITakeAll
	ITakeItem
	IDropItem
	IEquipItem
	IStripItem
	IQuitGame
)

type DirectionType int

const (
	DNone DirectionType = iota
	DUp
	DDown
	DLeft
	DRight
)

type Input struct {
	Typ       InputType
	Item      *Item
	Direction DirectionType
}

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
	Rune rune
	Name string
}

func (game *Game) loadWorldFile() {
	file, err := os.Open("game/maps/world.txt")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for rowIndex, row := range rows {
		// first level
		if rowIndex == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			continue
		}

		// portal entry
		level := game.Levels[row[0]]
		x, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			panic(err)
		}
		pos := Pos{int(x), int(y)}

		// portal destination
		dstLevel := game.Levels[row[3]]
		x, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err = strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			panic(err)
		}
		dstPos := Pos{int(x), int(y)}

		// link
		level.Portals[pos] = &LevelPos{dstLevel, dstPos}
	}
}

func (game *Game) resolveMovement(pos Pos) {
	monster, exists := game.CurrentLevel.AliveMonstersPos[pos]
	if exists {
		event := game.Player.Attack(&monster.Character)
		game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, Attack)
		game.CurrentLevel.addEvent(event)
		if !monster.IsAlive() {
			monster.Kill(game.CurrentLevel)
		}
		if !game.Player.IsAlive() {
			game.CurrentLevel.addEvent("DED")
		}
	} else if game.CurrentLevel.canWalk(pos) {
		game.CurrentLevel.Player.Move(pos, game.CurrentLevel)
		game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, Move)

		portal, portalExists := game.CurrentLevel.Portals[game.Player.Pos]
		if portalExists {
			game.CurrentLevel = portal.level
			game.Player.Pos = portal.pos
			game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, Portal)
		}
		game.CurrentLevel.resetVisibility()
		game.CurrentLevel.resolveVisibility()
	} else {
		game.CurrentLevel.checkClosedDoor(pos)
		game.CurrentLevel.resetVisibility()
		game.CurrentLevel.resolveVisibility()
	}
}

func (game *Game) resolveAction(pos Pos) {
	monster, exists := game.CurrentLevel.AliveMonstersPos[pos]
	if exists {
		event := game.Player.Attack(&monster.Character)
		game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, Attack)
		game.CurrentLevel.addEvent(event)
		if !monster.IsAlive() {
			monster.Kill(game.CurrentLevel)
		}
		if !game.Player.IsAlive() {
			game.CurrentLevel.addEvent("DED")
		}
	} else {
		if game.CurrentLevel.checkClosedDoor(pos) || game.CurrentLevel.checkOpenedDoor(pos) {
			game.CurrentLevel.resetVisibility()
			game.CurrentLevel.resolveVisibility()
		}
	}
}

func (game *Game) handleInput(input *Input) {
	p := game.Player
	switch input.Typ {
	case IMove, IAction:
		var newPos Pos
		switch input.Direction {
		case DUp:
			newPos = Pos{p.X, p.Y - 1}
		case DDown:
			newPos = Pos{p.X, p.Y + 1}
		case DLeft:
			newPos = Pos{p.X - 1, p.Y}
		case DRight:
			newPos = Pos{p.X + 1, p.Y}
		}

		switch input.Typ {
		case IMove:
			game.resolveMovement(newPos)
		case IAction:
			game.resolveAction(newPos)
		}
	case ITakeItem:
		if game.Player.TakeItem(game.CurrentLevel, input.Item) {
			game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, PickUp)
		}
	case IDropItem:
		game.Player.Strip(input.Item)
		if game.Player.DropItem(game.CurrentLevel, input.Item) {
			game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, DropDown)
		}
	case ITakeAll:
		for _, item := range game.CurrentLevel.Items[game.Player.Pos] {
			game.Player.TakeItem(game.CurrentLevel, item)
		}
		game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, PickUp)
	case IEquipItem:
		if game.Player.Equip(input.Item) {
			game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, Equip)
		}
	case IStripItem:
		if game.Player.Strip(input.Item) {
			game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, TakeOff)
		}
	}
}

func (game *Game) Run() {
	game.CurrentLevel.resolveVisibility()
	game.LevelChan <- game.CurrentLevel

	for input := range game.InputChan {
		game.CurrentLevel.LastEvents = make([]GameEvent, 0)
		if input.Typ == IQuitGame {
			return
		}

		game.handleInput(input)
		switch input.Typ {
		case IAction, IMove:
			for _, monster := range game.CurrentLevel.Monsters {
				if monster.IsAlive() {
					monster.Update(game.CurrentLevel)
				}
			}
		}

		game.LevelChan <- game.CurrentLevel
	}
}
