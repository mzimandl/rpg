package ui2d

import (
	"math/rand"
	"rpg/game"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type ui struct {
	winWidth          int32
	winHeight         int32
	renderer          *sdl.Renderer
	window            *sdl.Window
	textureAtlas      *sdl.Texture
	whiteDot          *sdl.Texture
	textureIndex      map[rune][]sdl.Rect
	keyboardState     []uint8
	prevKeyboardState []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
	fonts             map[FontType]*ttf.Font
	textCache         map[TextCacheKey]*sdl.Texture
}

type FontType int

const (
	FontSmall FontType = iota
	FontMedium
	FontLarge
)

type TextCacheKey struct {
	fontType FontType
	text     string
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
}

func Destroy() {
	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {
	var err error = nil
	ui := &ui{
		winWidth:  1280,
		winHeight: 720,
		inputChan: inputChan,
		levelChan: levelChan,
		r:         rand.New(rand.NewSource(1)),
		centerX:   -1,
		centerY:   -1,
	}

	ui.window, err = sdl.CreateWindow("RPG!!!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, ui.winWidth, ui.winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	ui.renderer, err = sdl.CreateRenderer(ui.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

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

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}

	return ui
}

func (ui *ui) stringToTexture(s string, fontType FontType) *sdl.Texture {
	font, exists := ui.fonts[fontType]
	if exists {
		textKey := TextCacheKey{fontType, s}
		texture, exists := ui.textCache[textKey]
		if exists {
			return texture
		}

		fontSurface, err := font.RenderUTF8Blended(s, sdl.Color{255, 255, 255, 0})
		if err != nil {
			panic(err)
		}
		fontTexture, err := ui.renderer.CreateTextureFromSurface(fontSurface)
		if err != nil {
			panic(err)
		}

		ui.textCache[textKey] = fontTexture
		return fontTexture
	} else {
		panic("Font type not found: " + strconv.Itoa(int(fontType)))
	}
}

func (ui *ui) Destroy() {
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

func (ui *ui) handleSrolling(level *game.Level) (int32, int32) {
	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

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

	offsetX := (ui.winWidth / 2) - int32(ui.centerX)*32
	offsetY := (ui.winHeight / 2) - int32(ui.centerY)*32

	return offsetX, offsetY
}

func (ui *ui) draw(level *game.Level) {
	offsetX, offsetY := ui.handleSrolling(level)

	ui.r.Seed(1)
	ui.renderer.Clear()

	ui.drawTiles(level, offsetX, offsetY)
	ui.drawMonsters(level, offsetX, offsetY)
	ui.drawPlayer(level, offsetX, offsetY)

	var textPosY int32 = 0
	ui.drawBox(0, 3*ui.winHeight/4, ui.winWidth/4, ui.winHeight/4, sdl.Color{64, 64, 64, 192})
	for i := len(level.Events) - 1; i >= 0; i-- {
		text := ui.stringToTexture(level.Events[i], FontSmall)
		text.SetColorMod(255, 0, 0)
		_, _, w, h, err := text.Query()

		if textPosY+h > int32(ui.winHeight/4) {
			break
		}

		if err != nil {
			panic(err)
		}
		ui.renderer.Copy(text, nil, &sdl.Rect{4, ui.winHeight - textPosY - h, w, h})
		textPosY += h
	}

	ui.renderer.Present()
}

func (ui *ui) getTile(tile game.Tile) sdl.Rect {
	srcRects := ui.textureIndex[tile.Rune]
	return srcRects[ui.r.Intn(len(srcRects))]
}

func (ui *ui) drawTiles(level *game.Level, offsetX, offsetY int32) {
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRect := ui.getTile(tile)
				if tile.Visible || tile.Visited {
					dstRect := sdl.Rect{offsetX + int32(x)*32, offsetY + int32(y)*32, 32, 32}
					pos := game.Pos{x, y}
					if level.Debug[pos] {
						ui.textureAtlas.SetColorMod(128, 0, 0)
					} else if tile.Visited && !tile.Visible {
						ui.textureAtlas.SetColorMod(128, 128, 128)
					} else {
						ui.textureAtlas.SetColorMod(255, 255, 255)
					}
					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
				}
			}
		}
	}
}

func (ui *ui) drawMonsters(level *game.Level, offsetX, offsetY int32) {
	for _, monster := range level.Monsters {
		if level.Map[monster.Y][monster.X].Visited {
			if monster.IsAlive() && level.Map[monster.Y][monster.X].Visible {
				ui.textureAtlas.SetColorMod(255, 255, 255)
				monsterSrcRect := ui.textureIndex[monster.Rune][0]
				monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*32, offsetY + int32(monster.Y)*32, 32, 32}
				ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &monsterDstRect)
			} else {
				if level.Map[monster.Y][monster.X].Visible {
					ui.textureAtlas.SetColorMod(255, 64, 64)
				} else {
					ui.textureAtlas.SetColorMod(128, 32, 32)
				}
				monsterSrcRect := ui.textureIndex[monster.Rune][0]
				monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*32, offsetY + int32(monster.Y)*32, 32, 32}
				ui.renderer.CopyEx(ui.textureAtlas, &monsterSrcRect, &monsterDstRect, 0, nil, sdl.FLIP_VERTICAL)
			}
		}
	}

	ui.textureAtlas.SetColorMod(255, 255, 255)
}

func (ui *ui) drawPlayer(level *game.Level, offsetX, offsetY int32) {
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	playerDstRect := sdl.Rect{offsetX + int32(level.Player.X)*32, offsetY + int32(level.Player.Y)*32, 32, 32}
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &playerDstRect)
}

func (ui *ui) drawBox(x, y, w, h int32, color sdl.Color) {
	ui.whiteDot.SetColorMod(color.R, color.G, color.B)
	ui.whiteDot.SetAlphaMod(color.A)
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{x, y, w, h})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{x, y, 1, h})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{x, y, w, 1})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{x, y + h - 1, w, 1})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{x + w - 1, y, 1, h})
}

func (ui *ui) keyPressed(scancode int) bool {
	return ui.keyboardState[scancode] != 0 && ui.prevKeyboardState[scancode] == 0
}

func (ui *ui) Run() {
	input := game.Input{game.None, ui.levelChan}
	ui.draw(<-ui.levelChan)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				input.Typ = game.QuitGame
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					input.Typ = game.CloseWindow
				}
			}
		}

		ui.keyboardState = sdl.GetKeyboardState()
		if ui.keyPressed(sdl.SCANCODE_ESCAPE) {
			input.Typ = game.QuitGame
		} else if ui.keyPressed(sdl.SCANCODE_UP) {
			input.Typ = game.Up
		} else if ui.keyPressed(sdl.SCANCODE_DOWN) {
			input.Typ = game.Down
		} else if ui.keyPressed(sdl.SCANCODE_LEFT) {
			input.Typ = game.Left
		} else if ui.keyPressed(sdl.SCANCODE_RIGHT) {
			input.Typ = game.Right
		}
		for i, v := range ui.keyboardState {
			ui.prevKeyboardState[i] = v
		}

		if input.Typ != game.None {
			ui.inputChan <- &input
			switch input.Typ {
			case game.QuitGame, game.CloseWindow:
				return
			default:
				// after input wait for modified level
				ui.draw(<-ui.levelChan)
			}
			input.Typ = game.None
		}

		sdl.Delay(10)
	}
}
