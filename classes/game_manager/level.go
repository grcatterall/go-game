package game_manager

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileMap struct {
	Tiles [][]*Tile
}

func LoadLevel(level [][]int, tileTextures map[int]rl.Texture2D) *TileMap {
	var tileMap TileMap
	for y, row := range level {
		var tileRow []*Tile
		for x, tileID := range row {
			if tileID != 0 {
				tile := &Tile{
					Texture:  tileTextures[tileID],
					Position: rl.Vector2{X: float32(x * 32), Y: float32(y * 32)},
				}
				tileRow = append(tileRow, tile)
			} else {
				tileRow = append(tileRow, nil)
			}
		}
		tileMap.Tiles = append(tileMap.Tiles, tileRow)
	}
	return &tileMap
}

func (tileMap *TileMap) Draw() {
	for _, row := range tileMap.Tiles {
		for _, tile := range row {
			if tile != nil {
				rl.DrawTextureV(tile.Texture, tile.Position, rl.White)
			}
		}
	}
}
