package ui2d

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

func (ui *ui) loadFonts() {
	var err error
	ui.fonts = make(map[FontType]*ttf.Font)
	ui.fonts[FontSmall], err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(16*float64(ui.winWidth)/1280))
	if err != nil {
		panic(err)
	}
	ui.fonts[FontMedium], err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(32*float64(ui.winWidth)/1280))
	if err != nil {
		panic(err)
	}
	ui.fonts[FontLarge], err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(64*float64(ui.winWidth)/1280))
	if err != nil {
		panic(err)
	}
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
