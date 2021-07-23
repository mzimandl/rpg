package ui2d

import (
	"math"
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) getRandomTile(tile game.Tile) sdl.Rect {
	srcRects := ui.textureIndex[tile.Rune]
	return srcRects[ui.tileRandomizer.Intn(len(srcRects))]
}

func (ui *ui) drawTiles(level *game.Level, offsetX, offsetY int32) {
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRect := ui.getRandomTile(tile)
				if tile.Visible || tile.Visited {
					dstRect := sdl.Rect{offsetX + int32(x)*32, offsetY + int32(y)*32, 32, 32}
					pos := game.Pos{x, y}
					if level.Debug[pos] {
						ui.textureAtlas.SetColorMod(128, 0, 0)
					} else if tile.Visited && !tile.Visible {
						ui.textureAtlas.SetColorMod(128, 128, 128)
					} else {
						ui.textureAtlas.SetColorMod(255, 255, 255)
					}
					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)

					if tile.OverlayRune != game.Blank {
						srcRect = ui.textureIndex[tile.OverlayRune][0]
						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					}
				}
			}
		}
	}
}

func (ui *ui) drawDeadMonsters(level *game.Level, offsetX, offsetY int32) {
	for _, monster := range level.Monsters {
		if !monster.IsAlive() {
			if level.Map[monster.Y][monster.X].Visited {
				if level.Map[monster.Y][monster.X].Visible {
					ui.textureAtlas.SetColorMod(255, 64, 64)
				} else {
					ui.textureAtlas.SetColorMod(128, 32, 32)
				}

				monsterSrcRect := ui.textureIndex[monster.Rune][0]
				monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*32, offsetY + int32(monster.Y)*32, 32, 32}
				ui.renderer.CopyEx(ui.textureAtlas, &monsterSrcRect, &monsterDstRect, 0, nil, sdl.FLIP_VERTICAL)
			}
		}
	}
	ui.textureAtlas.SetColorMod(255, 255, 255)
}

func (ui *ui) drawMonsters(level *game.Level, offsetX, offsetY int32) {
	for _, monster := range level.Monsters {
		if monster.IsAlive() && level.Map[monster.Y][monster.X].Visible {
			monsterSrcRect := ui.textureIndex[monster.Rune][0]
			monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*32, offsetY + int32(monster.Y)*32, 32, 32}
			ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &monsterDstRect)
		}
	}
}

func (ui *ui) drawItemsTile(level *game.Level, offsetX, offsetY int32) {
	for _, items := range level.Items {
		side := int32(32 / math.Sqrt(float64(len(items))))
		diff := float64(32-side) / float64(len(items))
		for i, item := range items {
			if level.Map[item.Y][item.X].Visible {
				itemSrcRect := ui.textureIndex[item.Rune][0]
				itemDstRect := sdl.Rect{offsetX + int32(item.X)*32 + int32(float64(i)*diff), offsetY + int32(item.Y)*32 + int32(float64(i)*diff), side, side}
				ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &itemDstRect)
			}
		}
	}
}

func (ui *ui) drawPlayer(level *game.Level, offsetX, offsetY int32) {
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	playerDstRect := sdl.Rect{offsetX + int32(level.Player.X)*32, offsetY + int32(level.Player.Y)*32, 32, 32}
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &playerDstRect)
}
