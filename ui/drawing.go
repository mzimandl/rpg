package ui2d

import "github.com/veandco/go-sdl2/sdl"

func (ui *ui) drawBox(rect *sdl.Rect, color sdl.Color) {
	ui.whiteDot.SetColorMod(color.R, color.G, color.B)
	ui.whiteDot.SetAlphaMod(color.A)
	ui.renderer.Copy(ui.whiteDot, nil, rect)
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X, rect.Y, 1, rect.H})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X, rect.Y, rect.W, 1})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X, rect.Y + rect.H - 1, rect.W, 1})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X + rect.W - 1, rect.Y, 1, rect.H})
}
