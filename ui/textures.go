package ui2d

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) loadTextures() {
	var err error
	ui.textCache = make(map[TextCacheKey]*sdl.Texture)
	ui.textureAtlas, err = img.LoadTexture(ui.renderer, "ui/assets/tiles.png")
	if err != nil {
		panic(err)
	}
	ui.textureIndex = loadTextureIndex()

	ui.whiteDot = getSinglePixelTexture(ui.renderer, sdl.Color{255, 255, 255, 255})
	ui.whiteDot.SetBlendMode(sdl.BLENDMODE_BLEND)
}

func loadTextureIndex() map[rune][]sdl.Rect {
	textureIndex := make(map[rune][]sdl.Rect)

	infile, err := os.Open("ui/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		tileRune := rune(line[0])

		xyc := line[1:]
		splitXYC := strings.Split(xyc, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXYC[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYC[1]), 10, 64)
		if err != nil {
			panic(err)
		}

		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYC[2]), 10, 64)
		for i := 0; i < int(variationCount); i++ {
			textureIndex[tileRune] = append(textureIndex[tileRune], sdl.Rect{int32(x) * 32, int32(y) * 32, 32, 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
	}
	return textureIndex
}

func getSinglePixelTexture(renderer *sdl.Renderer, color sdl.Color) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}
	pixel := make([]byte, 4)
	pixel[0] = color.R
	pixel[1] = color.G
	pixel[2] = color.B
	pixel[3] = color.A
	tex.Update(nil, pixel, 4)
	return tex
}
