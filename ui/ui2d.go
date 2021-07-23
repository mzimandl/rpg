package ui2d

import (
	"math/rand"
	"rpg/game"

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
	centerX, centerY    int

	renderer  *sdl.Renderer
	window    *sdl.Window
	levelChan chan *game.Level
	inputChan chan *game.Input

	tileRandomizer *rand.Rand
	textureAtlas   *sdl.Texture
	whiteDot       *sdl.Texture
	textureIndex   map[rune][]sdl.Rect
	fonts          map[FontType]*ttf.Font
	textCache      map[TextCacheKey]*sdl.Texture
	draggedItem    *game.Item

	music  *mix.Music
	sounds *sounds

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

	ui.mouseState = NewMouseState()
	ui.keyboardState = NewKeyboardState()

	ui.loadFonts()
	ui.loadTextures()
	ui.loadAudio()

	ui.recalculatePlacements()
	return ui
}

func (ui *ui) Destroy() {
	ui.sounds.Free()
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
	input := game.Input{game.INone, nil, game.DNone}
	currentLevel := <-ui.levelChan
	ui.music.Play(-1)

	for {
		ui.mouseState.update()
		ui.keyboardState.update()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				input.Typ = game.IQuitGame
			case *sdl.WindowEvent:
				switch e.Event {
				case sdl.WINDOWEVENT_CLOSE:
					input.Typ = game.IQuitGame
				case sdl.WINDOWEVENT_RESIZED:
					ui.winWidth, ui.winHeight = e.Data1, e.Data2
					ui.recalculatePlacements()
				}
			}
		}

		// inventory dragging
		if ui.state == UIInventory {
			if ui.mouseState.leftDoubleClicked() {
				item := ui.checkInventoryItems(currentLevel)
				if item != nil {
					input.Typ = game.IEquipItem
					input.Item = item
				}
			} else if ui.mouseState.leftClicked() {
				item := ui.checkInventoryItems(currentLevel)
				if item == nil {
					item = ui.checkEquippedItems(currentLevel)
				}
				if item != nil {
					ui.draggedItem = item
				}
			} else if ui.mouseState.leftUnclicked() && ui.draggedItem != nil {
				item := ui.checkDropDrag()
				if item != nil {
					input.Typ = game.IDropItem
					input.Item = item
				}
				if item == nil {
					item = ui.checkEquipDrag()
					if item != nil {
						input.Typ = game.IEquipItem
						input.Item = item
					}
				}
				if item == nil {
					item = ui.checkInventoryDrag()
					if item != nil {
						input.Typ = game.ITakeOffItem
						input.Item = item
					}
				}
				ui.draggedItem = nil
			} else if ui.mouseState.rightClicked() {
				item := ui.checkEquippedItems(currentLevel)
				if item != nil {
					input.Typ = game.ITakeOffItem
					input.Item = item
				}
				if item == nil {
					item = ui.checkInventoryItems(currentLevel)
					if item != nil {
						input.Typ = game.IDropItem
						input.Item = item
					}
				}
			}
		}

		// take item from the ground
		if ui.mouseState.leftClicked() {
			item := ui.checkGroundItems(currentLevel)
			if item != nil {
				input.Typ = game.ITakeItem
				input.Item = item
			}
		}

		if ui.keyboardState.pressed(sdl.SCANCODE_ESCAPE) {
			if ui.state == UIInventory {
				ui.state = UIMain
			} else {
				input.Typ = game.IQuitGame
			}
		} else if ui.keyboardState.pressed(sdl.SCANCODE_UP) {
			input.Direction = game.DUp
			if ui.keyboardState.hold(sdl.SCANCODE_SPACE) {
				input.Typ = game.IAction
			} else {
				input.Typ = game.IMove
			}
		} else if ui.keyboardState.pressed(sdl.SCANCODE_DOWN) {
			input.Direction = game.DDown
			if ui.keyboardState.hold(sdl.SCANCODE_SPACE) {
				input.Typ = game.IAction
			} else {
				input.Typ = game.IMove
			}
		} else if ui.keyboardState.pressed(sdl.SCANCODE_LEFT) {
			input.Direction = game.DLeft
			if ui.keyboardState.hold(sdl.SCANCODE_SPACE) {
				input.Typ = game.IAction
			} else {
				input.Typ = game.IMove
			}
		} else if ui.keyboardState.pressed(sdl.SCANCODE_RIGHT) {
			input.Direction = game.DRight
			if ui.keyboardState.hold(sdl.SCANCODE_SPACE) {
				input.Typ = game.IAction
			} else {
				input.Typ = game.IMove
			}
		} else if ui.keyboardState.pressed(sdl.SCANCODE_T) {
			input.Typ = game.ITakeAll
		} else if ui.keyboardState.pressed(sdl.SCANCODE_I) {
			if ui.state != UIInventory {
				ui.state = UIInventory
			} else {
				ui.state = UIMain
			}
		}

		if input.Typ != game.INone {
			ui.inputChan <- &input
			switch input.Typ {
			case game.IQuitGame:
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
					case game.DoorClose:
						playRandomSound(ui.sounds.doorClose, 10)
					}
				}
			}
			input.Typ = game.INone
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
