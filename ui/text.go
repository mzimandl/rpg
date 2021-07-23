package ui2d

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

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
