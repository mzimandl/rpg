package game

import (
	"path/filepath"
	"strings"
)

type Game struct {
	LevelChan    chan *Level
	InputChan    chan *Input
	Player       *Player
	Level        map[string]*Level
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
	currentLevel := levels["level1"]

	return &Game{levelChan, inputChan, player, levels, currentLevel}
}

type InputType int

const (
	None InputType = iota

	Up
	Down
	Left
	Right

	QuitGame
)

type Input struct {
	Typ InputType
}

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
	Rune rune
	Name string
}

func (game *Game) handleInput(input *Input) {
	p := game.CurrentLevel.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.CurrentLevel.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.CurrentLevel.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.CurrentLevel.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.CurrentLevel.resolveMovement(newPos)
	}
}

func (game *Game) Run() {

	game.LevelChan <- game.CurrentLevel

	for input := range game.InputChan {
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
