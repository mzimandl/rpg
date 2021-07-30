package ui

import (
	"math"
	"rpg/game"

	"github.com/veandco/go-sdl2/sdl"
)

const tileSize = 32

func (ui *ui) getRandomTile(r rune) sdl.Rect {
	srcRects := ui.textureIndex[r]
	return srcRects[ui.tileRandomizer.Intn(len(srcRects))]
}

func (ui *ui) calculateOffset(level *game.Level) (int32, int32) {
	if ui.centerX == -1 || ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	} else {
		if level.Player.X > ui.centerX+5 {
			ui.centerX = level.Player.X - 5
		} else if level.Player.X < ui.centerX-5 {
			ui.centerX = level.Player.X + 5
		}

		if level.Player.Y > ui.centerY+5 {
			ui.centerY = level.Player.Y - 5
		} else if level.Player.Y < ui.centerY-5 {
			ui.centerY = level.Player.Y + 5
		}
	}

	offsetX := (ui.winWidth / 2) - int32(ui.centerX)*tileSize
	offsetY := (ui.winHeight / 2) - int32(ui.centerY)*tileSize

	return offsetX, offsetY
}

func (ui *ui) drawTiles(level *game.Level, offsetX, offsetY int32) {
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRect := ui.getRandomTile(tile.Rune)
				var srcOverlayRect sdl.Rect
				if tile.OverlayRune != game.Blank {
					srcOverlayRect = ui.getRandomTile(tile.OverlayRune)
				}
				if tile.Visible || tile.Visited {
					dstRect := sdl.Rect{offsetX + int32(x)*tileSize, offsetY + int32(y)*tileSize, tileSize, tileSize}
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
						ui.renderer.Copy(ui.textureAtlas, &srcOverlayRect, &dstRect)
					}
				}
			}
		}
	}
	ui.textureAtlas.SetColorMod(255, 255, 255)
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
				monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*tileSize, offsetY + int32(monster.Y)*tileSize, tileSize, tileSize}
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
			monsterDstRect := sdl.Rect{offsetX + int32(monster.X)*tileSize, offsetY + int32(monster.Y)*tileSize, tileSize, tileSize}
			ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &monsterDstRect)
		}
	}
}

func (ui *ui) drawItemsTile(level *game.Level, offsetX, offsetY int32) {
	for _, items := range level.Items {
		side := int32(tileSize / math.Sqrt(float64(len(items))))
		diff := float64(tileSize-side) / float64(len(items))
		for i, item := range items {
			if level.Map[item.Y][item.X].Visible {
				itemSrcRect := ui.textureIndex[item.Rune][0]
				itemDstRect := sdl.Rect{offsetX + int32(item.X)*tileSize + int32(float64(i)*diff), offsetY + int32(item.Y)*tileSize + int32(float64(i)*diff), side, side}
				ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &itemDstRect)
			}
		}
	}
}

func (ui *ui) drawPlayer(level *game.Level, offsetX, offsetY int32) {
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	playerDstRect := sdl.Rect{offsetX + int32(level.Player.X)*tileSize, offsetY + int32(level.Player.Y)*tileSize, tileSize, tileSize}
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &playerDstRect)
}

func (ui *ui) drawStorages(level *game.Level, offsetX, offsetY int32) {
	for pos, storage := range level.Storages {
		if level.Map[pos.Y][pos.X].Visited {
			var srcRect sdl.Rect
			if storage == ui.usedStorage {
				srcRect = ui.textureIndex[storage.Rune][len(ui.textureIndex[storage.Rune])-1]
			} else {
				srcRect = ui.textureIndex[storage.Rune][0]
			}

			dstRect := sdl.Rect{offsetX + int32(pos.X)*tileSize, offsetY + int32(pos.Y)*tileSize, tileSize, tileSize}

			if !level.Map[pos.Y][pos.X].Visible {
				ui.textureAtlas.SetColorMod(128, 128, 128)
			}
			ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
			ui.textureAtlas.SetColorMod(255, 255, 255)
		}
	}
}
