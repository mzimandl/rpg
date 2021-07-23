package ui2d

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) drawBox(rect *sdl.Rect, color sdl.Color) {
	ui.whiteDot.SetColorMod(color.R, color.G, color.B)
	ui.whiteDot.SetAlphaMod(color.A)
	ui.renderer.Copy(ui.whiteDot, nil, rect)
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X, rect.Y, 1, rect.H})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X, rect.Y, rect.W, 1})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X, rect.Y + rect.H - 1, rect.W, 1})
	ui.renderer.Copy(ui.whiteDot, nil, &sdl.Rect{rect.X + rect.W - 1, rect.Y, 1, rect.H})
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
