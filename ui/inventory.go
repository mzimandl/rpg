package ui2d

import (
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) drawInventory(level *game.Level) {
	ui.drawBox(ui.placements.inv, sdl.Color{149, 84, 19, 200})
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, ui.placements.invChar)

	ui.drawBox(ui.placements.invCharHelmet, sdl.Color{0, 0, 0, 128})
	ui.drawBox(ui.placements.invCharWeapon, sdl.Color{0, 0, 0, 128})

	for i, item := range level.Player.Items {
		itemSrcRect := &ui.textureIndex[item.Rune][0]
		var itemDstRect *sdl.Rect
		if item == ui.draggedItem {
			itemDstRect = &sdl.Rect{ui.mouseState.x, ui.mouseState.y, ui.placements.itemSize, ui.placements.itemSize}
		} else {
			itemDstRect = ui.getInventoryItemRect(i)
		}
		ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
	}

	if level.Player.Helmet != nil {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Helmet.Rune][0], ui.placements.invCharHelmet)
	}
	if level.Player.Weapon != nil {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Weapon.Rune][0], ui.placements.invCharWeapon)
	}
}

func (ui *ui) checkInventoryItems(level *game.Level) *game.Item {
	for i, item := range level.Player.Items {
		itemDstRect := ui.getInventoryItemRect(i)
		if ui.mouseState.onArea(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkGroundItems(level *game.Level) *game.Item {
	for i, item := range level.Items[level.Player.Pos] {
		itemDstRect := ui.getGroundItemRect(i)
		if ui.mouseState.onArea(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkDroppedItem() *game.Item {
	if !ui.mouseState.onArea(ui.placements.inv) {
		return ui.draggedItem
	}
	return nil
}

func (ui *ui) checkEquippedItem() *game.Item {
	var slot *sdl.Rect

	switch ui.draggedItem.Typ {
	case game.Weapon:
		slot = ui.placements.invCharWeapon
	case game.Helmet:
		slot = ui.placements.invCharHelmet
	default:
		return nil
	}

	if ui.mouseState.onArea(slot) {
		return ui.draggedItem
	}
	return nil
}

func (ui *ui) checkTakeOffItem(level *game.Level) *game.Item {
	if ui.mouseState.onArea(ui.placements.invCharHelmet) {
		return level.Player.Helmet
	} else if ui.mouseState.onArea(ui.placements.invCharWeapon) {
		return level.Player.Weapon
	}
	return nil
}
