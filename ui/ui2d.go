package ui2d

import (
	"math/rand"
	"rpg/game"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type uiState int

const (
	UIMain uiState = iota
	UIInventory
)

type ui struct {
	state               uiState
	winWidth, winHeight int32
	placements          placements

	draggedItem      *game.Item
	renderer         *sdl.Renderer
	window           *sdl.Window
	textureAtlas     *sdl.Texture
	whiteDot         *sdl.Texture
	textureIndex     map[rune][]sdl.Rect
	centerX, centerY int
	tileRandomizer   *rand.Rand
	levelChan        chan *game.Level
	inputChan        chan *game.Input
	fonts            map[FontType]*ttf.Font
	textCache        map[TextCacheKey]*sdl.Texture
	music            *mix.Music
	sounds           sounds

	keyboardState *keyboardState
	mouseState    *mouseState
}

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	err = img.Init(img.INIT_PNG)
	if err != nil {
		panic(err)
	}
	err = ttf.Init()
	if err != nil {
		panic(err)
	}
	err = mix.Init(mix.INIT_OGG)
	if err != nil {
		panic(err)
	}
	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}
}

func Destroy() {
	mix.Quit()
	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {
	var err error = nil
	ui := &ui{
		state:          UIMain,
		winWidth:       1280,
		winHeight:      720,
		inputChan:      inputChan,
		levelChan:      levelChan,
		tileRandomizer: rand.New(rand.NewSource(1)),
		centerX:        -1,
		centerY:        -1,
	}

	ui.window, err = sdl.CreateWindow("RPG!!!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, ui.winWidth, ui.winHeight, sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}

	ui.renderer, err = sdl.CreateRenderer(ui.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.fonts = make(map[FontType]*ttf.Font)
	ui.fonts[FontSmall], err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(16*float64(ui.winWidth)/1280))
	if err != nil {
		panic(err)
	}
	ui.fonts[FontMedium], err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(32*float64(ui.winWidth)/1280))
	if err != nil {
		panic(err)
	}
	ui.fonts[FontLarge], err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(64*float64(ui.winWidth)/1280))
	if err != nil {
		panic(err)
	}
	ui.textCache = make(map[TextCacheKey]*sdl.Texture)

	ui.textureAtlas, err = img.LoadTexture(ui.renderer, "ui/assets/tiles.png")
	if err != nil {
		panic(err)
	}
	ui.textureIndex = loadTextureIndex()

	ui.whiteDot = getSinglePixelTexture(ui.renderer, sdl.Color{255, 255, 255, 255})
	ui.whiteDot.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.mouseState = NewMouseState()
	ui.keyboardState = NewKeyboardState()

	ui.music, err = mix.LoadMUS("ui/assets/dungeon002.ogg")
	if err != nil {
		panic(err)
	}
	ui.music.Play(-1)

	footstepBase := "ui/assets/sounds/footstep0"
	for i := 0; i <= 9; i++ {
		footstepFile := footstepBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(footstepFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.footstep = append(ui.sounds.footstep, chunk)
	}
	doorOpenBase := "ui/assets/sounds/doorOpen_"
	for i := 1; i <= 2; i++ {
		doorOpenFile := doorOpenBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(doorOpenFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.doorOpen = append(ui.sounds.doorOpen, chunk)
	}

	ui.recalculatePlacements()
	return ui
}

func (ui *ui) Destroy() {
	for _, chunk := range ui.sounds.doorOpen {
		chunk.Free()
	}
	for _, chunk := range ui.sounds.footstep {
		chunk.Free()
	}
	ui.music.Free()
	ui.textureAtlas.Destroy()
	for _, texture := range ui.textCache {
		texture.Destroy()
	}
	for _, font := range ui.fonts {
		font.Close()
	}
	ui.renderer.Destroy()
	ui.window.Destroy()
}

func (ui *ui) calculateOffset(level *game.Level) (int32, int32) {
	if ui.centerX == -1 || ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	} else {
		if level.Player.X > ui.centerX+5 {
			ui.centerX = level.Player.X - 5
		} else if level.Player.X < ui.centerX-5 {
			ui.centerX = level.Player.X + 5
		}

		if level.Player.Y > ui.centerY+5 {
			ui.centerY = level.Player.Y - 5
		} else if level.Player.Y < ui.centerY-5 {
			ui.centerY = level.Player.Y + 5
		}
	}

	offsetX := (ui.winWidth / 2) - int32(ui.centerX)*32
	offsetY := (ui.winHeight / 2) - int32(ui.centerY)*32

	return offsetX, offsetY
}

func (ui *ui) drawLevel(level *game.Level) {
	offsetX, offsetY := ui.calculateOffset(level)
	ui.tileRandomizer.Seed(1)

	ui.drawTiles(level, offsetX, offsetY)
	ui.drawDeadMonsters(level, offsetX, offsetY)
	ui.drawItemsTile(level, offsetX, offsetY)
	ui.drawMonsters(level, offsetX, offsetY)
	ui.drawPlayer(level, offsetX, offsetY)
}

func (ui *ui) drawUI(level *game.Level) {
	ui.drawGroundItems(level, 0, 3*ui.winHeight/4)
	ui.drawLog(level)
}

func (ui *ui) Run() {
	input := game.Input{game.None, nil}
	currentLevel := <-ui.levelChan

	for {
		ui.mouseState.update()
		ui.keyboardState.update()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				input.Typ = game.QuitGame
			case *sdl.WindowEvent:
				switch e.Event {
				case sdl.WINDOWEVENT_CLOSE:
					input.Typ = game.QuitGame
				case sdl.WINDOWEVENT_RESIZED:
					ui.winWidth, ui.winHeight = e.Data1, e.Data2
					ui.recalculatePlacements()
				}
			}
		}

		// inventory dragging
		if ui.state == UIInventory {
			if ui.mouseState.leftClicked() {
				item := ui.checkInventoryItems(currentLevel)
				if item != nil {
					ui.draggedItem = item
				}
			} else if ui.mouseState.leftUnclicked() && ui.draggedItem != nil {
				item := ui.checkDroppedItem()
				if item != nil {
					input.Typ = game.DropItem
					input.Item = item
				}
				item = ui.checkEquippedItem()
				if item != nil {
					input.Typ = game.EquipItem
					input.Item = item
				}
				ui.draggedItem = nil
			} else if ui.mouseState.rightClicked() {
				item := ui.checkTakeOffItem(currentLevel)
				if item != nil {
					input.Typ = game.TakeOffItem
					input.Item = item
				}
				item = ui.checkInventoryItems(currentLevel)
				if item != nil {
					input.Typ = game.DropItem
					input.Item = item
				}
			}
		}

		// take item from the ground
		if ui.mouseState.leftClicked() {
			item := ui.checkGroundItems(currentLevel)
			if item != nil {
				input.Typ = game.TakeItem
				input.Item = item
			}
		}

		if ui.keyboardState.pressed(sdl.SCANCODE_ESCAPE) {
			if ui.state == UIInventory {
				ui.state = UIMain
			} else {
				input.Typ = game.QuitGame
			}
		} else if ui.keyboardState.pressed(sdl.SCANCODE_UP) {
			input.Typ = game.Up
		} else if ui.keyboardState.pressed(sdl.SCANCODE_DOWN) {
			input.Typ = game.Down
		} else if ui.keyboardState.pressed(sdl.SCANCODE_LEFT) {
			input.Typ = game.Left
		} else if ui.keyboardState.pressed(sdl.SCANCODE_RIGHT) {
			input.Typ = game.Right
		} else if ui.keyboardState.pressed(sdl.SCANCODE_T) {
			input.Typ = game.TakeAll
		} else if ui.keyboardState.pressed(sdl.SCANCODE_I) {
			if ui.state != UIInventory {
				ui.state = UIInventory
			} else {
				ui.state = UIMain
			}
		}

		if input.Typ != game.None {
			ui.inputChan <- &input
			switch input.Typ {
			case game.QuitGame:
				return
			default:
				currentLevel = <-ui.levelChan
				for _, lastEvent := range currentLevel.LastEvents {
					switch lastEvent {
					case game.Portal:
						ui.centerX, ui.centerY = -1, -1
					case game.Move:
						playRandomSound(ui.sounds.footstep, 10)
					case game.DoorOpen:
						playRandomSound(ui.sounds.doorOpen, 10)
					}
				}
			}
			input.Typ = game.None
		}

		ui.renderer.Clear()
		ui.drawLevel(currentLevel)
		ui.drawUI(currentLevel)
		if ui.state == UIInventory {
			ui.drawInventory(currentLevel)
		}
		ui.renderer.Present()

		sdl.Delay(10)
	}
}
