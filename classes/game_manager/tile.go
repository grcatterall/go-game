package game_manager

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	Texture  rl.Texture2D
	Position rl.Vector2
}

func LoadTile(filepath string) rl.Texture2D {
	return rl.LoadTexture(filepath)
}
