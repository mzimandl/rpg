package ui

import (
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) drawGroundItems(level *game.Level, x, y int32) {
	storage := level.Storages[level.Player.Pos]
	if storage != nil {
		var srcRect *sdl.Rect
		if ui.usedRepository == &storage.Repository {
			srcRect = &ui.textureIndex[storage.Rune][len(ui.textureIndex[storage.Rune])-1]
		} else {
			srcRect = &ui.textureIndex[storage.Rune][0]
		}
		dstRect := ui.getGroundItemRect(0)
		ui.renderer.Copy(ui.textureAtlas, srcRect, dstRect)
	} else {
		for i, item := range level.Items[level.Player.Pos] {
			itemSrcRect := &ui.textureIndex[item.Rune][0]
			itemDstRect := ui.getGroundItemRect(i)
			ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
		}
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

func (ui *ui) drawInventory(level *game.Level) {
	ui.drawBox(ui.placements.inv, sdl.Color{149, 84, 19, 128})
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, ui.placements.invChar)
	ui.drawBox(ui.placements.invCharHelmet, sdl.Color{0, 0, 0, 128})
	ui.drawBox(ui.placements.invCharWeapon, sdl.Color{0, 0, 0, 128})
	ui.drawBox(ui.placements.invCharArmor, sdl.Color{0, 0, 0, 128})

	for i, item := range level.Player.Items {
		if item != ui.draggedItem {
			itemSrcRect := &ui.textureIndex[item.Rune][0]
			itemDstRect := ui.getInventoryItemRect(i)
			ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
		}
	}

	if level.Player.Helmet != nil && level.Player.Helmet != ui.draggedItem {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Helmet.Rune][0], ui.placements.invCharHelmet)
	}
	if level.Player.Weapon != nil && level.Player.Weapon != ui.draggedItem {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Weapon.Rune][0], ui.placements.invCharWeapon)
	}
	if level.Player.Armor != nil && level.Player.Armor != ui.draggedItem {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Armor.Rune][0], ui.placements.invCharArmor)
	}
}

func (ui *ui) drawExchange() {
	ui.drawBox(ui.placements.exch, sdl.Color{149, 84, 19, 128})

	for i, item := range ui.usedRepository.Items {
		if item != ui.draggedItem {
			itemSrcRect := &ui.textureIndex[item.Rune][0]
			itemDstRect := ui.getExchangeItemRect(i)
			ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
		}
	}
}

func (ui *ui) drawDraggedItem() {
	if ui.draggedItem != nil {
		itemSrcRect := &ui.textureIndex[ui.draggedItem.Rune][0]
		itemDstRect := &sdl.Rect{ui.mouseState.x, ui.mouseState.y, ui.placements.itemSize, ui.placements.itemSize}
		ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
	}
}
