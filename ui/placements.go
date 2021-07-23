package ui2d

import "github.com/veandco/go-sdl2/sdl"

const itemSizeRatio = 0.033

type placements struct {
	itemSize int32

	log *sdl.Rect

	inv           *sdl.Rect
	invChar       *sdl.Rect
	invCharHelmet *sdl.Rect
	invCharWeapon *sdl.Rect
	invCharArmor  *sdl.Rect
}

func (ui *ui) recalculatePlacements() {
	ui.placements.itemSize = int32(float32(ui.winWidth) * itemSizeRatio)
	ui.placements.log = &sdl.Rect{
		0,
		3 * ui.winHeight / 4,
		ui.winWidth / 4,
		ui.winHeight / 4,
	}
	ui.placements.inv = ui.getInventoryRectangle()
	ui.placements.invChar = &sdl.Rect{
		ui.placements.inv.X + ui.placements.inv.W/4,
		ui.placements.inv.Y + ui.placements.inv.H/20,
		ui.placements.inv.W / 2,
		ui.placements.inv.H / 2,
	}
	ui.placements.invCharHelmet = &sdl.Rect{
		ui.placements.invChar.X + ui.placements.invChar.W/2 + ui.placements.invChar.W/20 - ui.placements.itemSize/2,
		ui.placements.invChar.Y,
		ui.placements.itemSize,
		ui.placements.itemSize,
	}
	ui.placements.invCharWeapon = &sdl.Rect{
		ui.placements.invChar.X + ui.placements.invChar.W/10 - ui.placements.itemSize/2,
		ui.placements.invChar.Y + ui.placements.invChar.H/3,
		ui.placements.itemSize,
		ui.placements.itemSize,
	}
	ui.placements.invCharArmor = &sdl.Rect{
		ui.placements.invChar.X + ui.placements.invChar.W/2 + ui.placements.invChar.W/20 - ui.placements.itemSize/2,
		ui.placements.invChar.Y + ui.placements.invChar.H/3,
		ui.placements.itemSize,
		ui.placements.itemSize,
	}
}

func (ui *ui) getGroundItemRect(index int) *sdl.Rect {
	return &sdl.Rect{
		int32(index) * ui.placements.itemSize,
		3*ui.winHeight/4 - ui.placements.itemSize,
		ui.placements.itemSize,
		ui.placements.itemSize,
	}
}

func (ui *ui) getInventoryRectangle() *sdl.Rect {
	invWidth := int32(float64(ui.winWidth) * 0.4)
	invHeight := int32(float64(ui.winHeight) * 0.75)
	offsetX := (ui.winWidth - invWidth) / 2
	offsetY := (ui.winHeight - invHeight) / 2
	return &sdl.Rect{offsetX, offsetY, invWidth, invHeight}
}

func (ui *ui) getInventoryItemRect(index int) *sdl.Rect {
	return &sdl.Rect{
		int32(index)*ui.placements.itemSize + ui.placements.inv.X,
		ui.placements.inv.Y + ui.placements.inv.H - ui.placements.itemSize,
		ui.placements.itemSize,
		ui.placements.itemSize,
	}
}
