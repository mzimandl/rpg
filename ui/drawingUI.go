package ui2d

import (
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) drawGroundItems(level *game.Level, x, y int32) {
	for i, item := range level.Items[level.Player.Pos] {
		itemSrcRect := &ui.textureIndex[item.Rune][0]
		itemDstRect := ui.getGroundItemRect(i)
		ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
	}
}

func (ui *ui) drawLog(level *game.Level) {
	var textPosY int32 = 0
	ui.drawBox(ui.placements.log, sdl.Color{64, 64, 64, 192})
	for i := len(level.Log) - 1; i >= 0; i-- {
		text := ui.stringToTexture(level.Log[i], FontSmall)
		text.SetColorMod(255, 0, 0)
		_, _, w, h, err := text.Query()

		if textPosY+h > int32(ui.winHeight/4) {
			break
		}

		if err != nil {
			panic(err)
		}
		ui.renderer.Copy(text, nil, &sdl.Rect{ui.placements.log.X + 4, ui.placements.log.Y + ui.placements.log.H - textPosY - h, w, h})
		textPosY += h
	}
}
