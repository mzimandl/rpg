package ui2d

import (
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) getInventoryRectangle() *sdl.Rect {
	invWidth := int32(float64(ui.winWidth) * 0.4)
	invHeight := int32(float64(ui.winHeight) * 0.75)
	offsetX := (ui.winWidth - invWidth) / 2
	offsetY := (ui.winHeight - invHeight) / 2
	return &sdl.Rect{offsetX, offsetY, invWidth, invHeight}
}

func (ui *ui) checkInventoryItems(level *game.Level, itemSize int32) *game.Item {
	for i, item := range level.Player.Items {
		itemDstRect := ui.getInventoryItemRect(i, itemSize)
		if ui.mouseState.onArea(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkGroundItems(level *game.Level, itemSize int32) *game.Item {
	for i, item := range level.Items[level.Player.Pos] {
		itemDstRect := ui.getGroundItemRect(i, 0, 3*ui.winHeight/4, itemSize)
		if ui.mouseState.onArea(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) drawInventory(level *game.Level) {
	invRect := ui.getInventoryRectangle()
	ui.drawBox(invRect, sdl.Color{149, 84, 19, 200})
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{invRect.X + invRect.W/4, invRect.Y + invRect.H/20, invRect.W / 2, invRect.H / 2})
	itemSize := int32(float32(ui.winWidth) * itemSizeRatio)
	for i, item := range level.Player.Items {
		itemSrcRect := &ui.textureIndex[item.Rune][0]
		var itemDstRect *sdl.Rect
		if item == ui.draggedItem {
			itemDstRect = &sdl.Rect{ui.mouseState.x, ui.mouseState.y, itemSize, itemSize}
		} else {
			itemDstRect = ui.getInventoryItemRect(i, itemSize)
		}
		ui.renderer.Copy(ui.textureAtlas, itemSrcRect, itemDstRect)
	}

	helmetSlot := ui.getHelmetSlotRect(itemSize)
	ui.drawBox(helmetSlot, sdl.Color{0, 0, 0, 128})
	weaponSlot := ui.getWeaponSlotRect(itemSize)
	ui.drawBox(weaponSlot, sdl.Color{0, 0, 0, 128})

	if level.Player.Helmet != nil {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Helmet.Rune][0], helmetSlot)
	}
	if level.Player.Weapon != nil {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Weapon.Rune][0], weaponSlot)
	}
}

func (ui *ui) getInventoryItemRect(index int, itemSize int32) *sdl.Rect {
	invRect := ui.getInventoryRectangle()
	return &sdl.Rect{int32(index)*itemSize + invRect.X, invRect.Y + invRect.H - itemSize, itemSize, itemSize}
}

func (ui *ui) getHelmetSlotRect(itemSize int32) *sdl.Rect {
	invRect := ui.getInventoryRectangle()
	return &sdl.Rect{invRect.X + invRect.W/2 + invRect.W/40 - itemSize/2, invRect.Y + invRect.H/20, itemSize, itemSize}
}

func (ui *ui) getWeaponSlotRect(itemSize int32) *sdl.Rect {
	invRect := ui.getInventoryRectangle()
	return &sdl.Rect{invRect.X + 2*invRect.W/7 - itemSize/2, invRect.Y + 2*invRect.H/9, itemSize, itemSize}
}

func (ui *ui) CheckDroppedItem() *game.Item {
	invRect := ui.getInventoryRectangle()
	if !ui.mouseState.onArea(invRect) {
		return ui.draggedItem
	}
	return nil
}

func (ui *ui) CheckEquippedItem(itemSize int32) *game.Item {
	var slot *sdl.Rect

	switch ui.draggedItem.Typ {
	case game.Weapon:
		slot = ui.getWeaponSlotRect(itemSize)
	case game.Helmet:
		slot = ui.getHelmetSlotRect(itemSize)
	default:
		return nil
	}

	if ui.mouseState.onArea(slot) {
		return ui.draggedItem
	}
	return nil
}
