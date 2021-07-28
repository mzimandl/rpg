package ui

import (
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) checkGroundItems(level *game.Level) *game.Item {
	for i, item := range level.Items[level.Player.Pos] {
		itemDstRect := ui.getGroundItemRect(i)
		if ui.mouseState.onRect(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkInventoryItems(level *game.Level) *game.Item {
	for i, item := range level.Player.Items {
		itemDstRect := ui.getInventoryItemRect(i)
		if ui.mouseState.onRect(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkEquippedItems(level *game.Level) *game.Item {
	if ui.mouseState.onRect(ui.placements.invCharHelmet) {
		return level.Player.Helmet
	} else if ui.mouseState.onRect(ui.placements.invCharWeapon) {
		return level.Player.Weapon
	} else if ui.mouseState.onRect(ui.placements.invCharArmor) {
		return level.Player.Armor
	}
	return nil
}

func (ui *ui) checkInventoryDrag() *game.Item {
	if ui.mouseState.onRect(ui.placements.inv) {
		return ui.draggedItem
	}
	return nil
}

func (ui *ui) checkDropDrag() *game.Item {
	if !ui.mouseState.onRect(ui.placements.inv) {
		return ui.draggedItem
	}
	return nil
}

func (ui *ui) checkEquipDrag() *game.Item {
	var slot *sdl.Rect

	switch ui.draggedItem.Typ {
	case game.Weapon:
		slot = ui.placements.invCharWeapon
	case game.Helmet:
		slot = ui.placements.invCharHelmet
	case game.Armor:
		slot = ui.placements.invCharArmor
	default:
		return nil
	}

	if ui.mouseState.onRect(slot) {
		return ui.draggedItem
	}
	return nil
}
