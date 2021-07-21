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
	None InputType = iota
	Up
	Down
	Left
	Right
	TakeAll
	TakeItem
	QuitGame
)

type Input struct {
	Typ  InputType
	Item *Item
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
		game.CurrentLevel.checkDoor(pos)
		game.CurrentLevel.resetVisibility()
		game.CurrentLevel.resolveVisibility()
	}
}

func (game *Game) handleInput(input *Input) {
	p := game.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.resolveMovement(newPos)
	case TakeItem:
		game.CurrentLevel.MoveItem(input.Item, &game.Player.Character)
		game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, PickUp)
	case TakeAll:
		for _, item := range game.CurrentLevel.Items[game.Player.Pos] {
			game.CurrentLevel.MoveItem(item, &game.Player.Character)
			game.CurrentLevel.addEvent("You took item: " + item.Name)
		}
		game.CurrentLevel.LastEvents = append(game.CurrentLevel.LastEvents, PickUp)
	}
}

func (game *Game) Run() {
	game.CurrentLevel.resolveVisibility()
	game.LevelChan <- game.CurrentLevel

	for input := range game.InputChan {
		game.CurrentLevel.LastEvents = make([]GameEvent, 0)
		if input.Typ == QuitGame {
			return
		}

		game.handleInput(input)
		for _, monster := range game.CurrentLevel.Monsters {
			if monster.IsAlive() {
				monster.Update(game.CurrentLevel)
			}
		}

		game.LevelChan <- game.CurrentLevel
	}
}
