package ui2d

import (
	"math"
	"math/rand"
	"rpg/game"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const itemSizeRatio = 0.033

type uiState int

const (
	UIMain uiState = iota
	UIInventory
)

type ui struct {
	state             uiState
	draggedItem       *game.Item
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
	tileRandomizer    *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
	fonts             map[FontType]*ttf.Font
	textCache         map[TextCacheKey]*sdl.Texture
	music             *mix.Music
	sounds            sounds
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

type sounds struct {
	doorOpen []*mix.Chunk
	footstep []*mix.Chunk
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
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

	ui.window, err = sdl.CreateWindow("RPG!!!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, ui.winWidth, ui.winHeight, sdl.WINDOW_SHOWN)
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

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}

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

func (ui *ui) getRandomTile(tile game.Tile) sdl.Rect {
	srcRects := ui.textureIndex[tile.Rune]
	return srcRects[ui.tileRandomizer.Intn(len(srcRects))]
}

func (ui *ui) drawTiles(level *game.Level, offsetX, offsetY int32) {
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRect := ui.getRandomTile(tile)
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

					if tile.OverlayRune != game.Blank {
						srcRect = ui.textureIndex[tile.OverlayRune][0]
						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					}
				}
			}
		}
	}
}

func (ui *ui) drawDeadMonsters(level *game.Level, offsetX, offsetY int32) {
	for _, monster := range level.Monsters {
		if !monster.IsAlive() {
			if level.Map[monster.Y][monster.X].Visited {
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

func (ui *ui) drawMonsters(level *game.Level, offsetX, offsetY int32) {
	for _, monster := range level.Monsters {
		if monster.IsAlive() && level.Map[monster.Y][monster.X].Visible {
			monsterSrcRect := ui.textureIndex[monster.Rune][0]
			monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*32, offsetY + int32(monster.Y)*32, 32, 32}
			ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &monsterDstRect)
		}
	}
}

func (ui *ui) drawItemsTile(level *game.Level, offsetX, offsetY int32) {
	for _, items := range level.Items {
		side := int32(32 / math.Sqrt(float64(len(items))))
		diff := float64(32-side) / float64(len(items))
		for i, item := range items {
			if level.Map[item.Y][item.X].Visible {
				itemSrcRect := ui.textureIndex[item.Rune][0]
				itemDstRect := sdl.Rect{offsetX + int32(item.X)*32 + int32(float64(i)*diff), offsetY + int32(item.Y)*32 + int32(float64(i)*diff), side, side}
				ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &itemDstRect)
			}
		}
	}
}

func (ui *ui) drawPlayer(level *game.Level, offsetX, offsetY int32) {
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	playerDstRect := sdl.Rect{offsetX + int32(level.Player.X)*32, offsetY + int32(level.Player.Y)*32, 32, 32}
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &playerDstRect)
}

func (ui *ui) drawUI(level *game.Level) {
	ui.drawGroundItems(level, 0, 3*ui.winHeight/4)
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
}

func (ui *ui) getGroundItemRect(index int, x, y, itemSize int32) *sdl.Rect {
	return &sdl.Rect{int32(index)*itemSize + x, y - itemSize, itemSize, itemSize}
}

func (ui *ui) drawGroundItems(level *game.Level, x, y int32) {
	itemSize := int32(float32(ui.winWidth) * itemSizeRatio)
	for i, item := range level.Items[level.Player.Pos] {
		itemSrcRect := &ui.textureIndex[item.Rune][0]
		itemDstRect := ui.getGroundItemRect(i, x, y, itemSize)
		ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
	}
}

func (ui *ui) drawInventory(level *game.Level, mouseState *mouseState) {
	invRect := ui.getInventoryRectangle()
	ui.drawBox(invRect.X, invRect.Y, invRect.W, invRect.H, sdl.Color{149, 84, 19, 200})
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{invRect.X + invRect.W/4, invRect.Y, invRect.W / 2, invRect.H / 2})
	itemSize := int32(float32(ui.winWidth) * itemSizeRatio)
	for i, item := range level.Player.Items {
		itemSrcRect := &ui.textureIndex[item.Rune][0]
		var itemDstRect *sdl.Rect
		if item == ui.draggedItem {
			itemDstRect = &sdl.Rect{mouseState.x, mouseState.y, itemSize, itemSize}
		} else {
			itemDstRect = ui.getInventoryItemRect(i, invRect.X, invRect.Y+invRect.H, itemSize)
		}
		ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
	}
}

func (ui *ui) getInventoryItemRect(index int, x, y, itemSize int32) *sdl.Rect {
	return &sdl.Rect{int32(index)*itemSize + x, y - itemSize, itemSize, itemSize}
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

func (ui *ui) checkGroundItems(level *game.Level, mouseState *mouseState, itemSize int32) *game.Item {
	for i, item := range level.Items[level.Player.Pos] {
		itemDstRect := ui.getGroundItemRect(i, 0, 3*ui.winHeight/4, itemSize)
		if mouseState.onArea(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkInventoryItems(level *game.Level, mouseState *mouseState, itemSize int32) *game.Item {
	invRect := ui.getInventoryRectangle()
	for i, item := range level.Player.Items {
		itemDstRect := ui.getInventoryItemRect(i, invRect.X, invRect.Y+invRect.H, itemSize)
		if mouseState.onArea(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) getInventoryRectangle() *sdl.Rect {
	invWidth := int32(float64(ui.winWidth) * 0.4)
	invHeight := int32(float64(ui.winHeight) * 0.75)
	offsetX := (ui.winWidth - invWidth) / 2
	offsetY := (ui.winHeight - invHeight) / 2
	return &sdl.Rect{offsetX, offsetY, invWidth, invHeight}
}

func (ui *ui) Run() {
	mouseState := NewMouseState()
	input := game.Input{game.None, nil}
	currentLevel := <-ui.levelChan

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				input.Typ = game.QuitGame
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					input.Typ = game.QuitGame
				}
			}
		}

		mouseState.update()
		itemSize := int32(float32(ui.winWidth) * itemSizeRatio)
		switch ui.state {
		case UIInventory:
			if mouseState.leftClicked() {
				item := ui.checkInventoryItems(currentLevel, mouseState, itemSize)
				if item != nil {
					ui.draggedItem = item
				}
			}
			if ui.draggedItem != nil && !mouseState.leftButton {
				invRect := ui.getInventoryRectangle()
				if !invRect.HasIntersection(&sdl.Rect{mouseState.x, mouseState.y, 1, 1}) {
					input.Typ = game.DropItem
					input.Item = ui.draggedItem
				}
				ui.draggedItem = nil
			}

		default:
			if mouseState.leftClicked() {
				item := ui.checkGroundItems(currentLevel, mouseState, itemSize)
				if item != nil {
					input.Typ = game.TakeItem
					input.Item = item
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
		} else if ui.keyPressed(sdl.SCANCODE_T) {
			input.Typ = game.TakeAll
		} else if ui.keyPressed(sdl.SCANCODE_I) {
			if ui.state != UIInventory {
				ui.state = UIInventory
			} else {
				ui.state = UIMain
			}
		}
		for i, v := range ui.keyboardState {
			ui.prevKeyboardState[i] = v
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
			ui.drawInventory(currentLevel, mouseState)
		}
		ui.renderer.Present()

		sdl.Delay(10)
	}
}
