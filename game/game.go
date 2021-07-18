package game

type Game struct {
	LevelChans []chan *Level
	InputChan  chan *Input
	Level      *Level
}

func NewGame(numWindows int, levelPath string) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)

	return &Game{levelChans, inputChan, NewLevelFromFile(levelPath)}
}

type InputType int

const (
	None InputType = iota

	Up
	Down
	Left
	Right

	CloseWindow
	QuitGame
)

type Input struct {
	Typ          InputType
	LevelChannel chan *Level
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
	p := game.Level.Player
	switch input.Typ {
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.Level.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.Level.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.Level.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.Level.resolveMovement(newPos)
	}
}

func (game *Game) Run() {

	for _, lchan := range game.LevelChans {
		lchan <- game.Level
	}

	for input := range game.InputChan {
		if input.Typ == QuitGame {
			return
		}

		game.handleInput(input)
		if len(game.LevelChans) == 0 {
			return
		}

		for _, monster := range game.Level.Monsters {
			if monster.IsAlive() {
				monster.Update(game.Level)
			}
		}

		for _, lchan := range game.LevelChans {
			lchan <- game.Level
		}
	}
}
