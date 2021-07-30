package ui

import (
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

type UIArea int

const (
	UIAInv UIArea = iota
	UIASlot
	UIAExch
)

func (ui *ui) checkGroundItems(level *game.Level) *game.Item {
	indexShift := 0
	storage := level.Storages[level.Player.Pos]
	if storage != nil {
		indexShift++
	}
	for i, item := range level.Items[level.Player.Pos] {
		itemDstRect := ui.getGroundItemRect(i + indexShift)
		if ui.mouseState.onRect(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkGroundStorage(level *game.Level) *game.Repository {
	storage := level.Storages[level.Player.Pos]
	if storage != nil && !storage.Locked {
		itemDstRect := ui.getGroundItemRect(0)
		if ui.mouseState.onRect(itemDstRect) {
			return &storage.Repository
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

func (ui *ui) checkExchangeItems(level *game.Level) *game.Item {
	for i, item := range ui.usedRepository.Items {
		itemDstRect := ui.getExchangeItemRect(i)
		if ui.mouseState.onRect(itemDstRect) {
			return item
		}
	}
	return nil
}

func (ui *ui) checkInventoryDrag() *game.Item {
	if ui.mouseState.onRect(ui.placements.inv) {
		return ui.draggedItem
	}
	return nil
}

func (ui *ui) checkExchangeDrag() *game.Item {
	if ui.mouseState.onRect(ui.placements.exch) {
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
