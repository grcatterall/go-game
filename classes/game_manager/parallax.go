package game_manager

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParallaxLayer struct {
	Texture   rl.Texture2D
	Speed     float32
	PositionX float32
}

type ParallaxBackground struct {
	Layers []ParallaxLayer
}

var (
	backgroundLayers []rl.Texture2D
	parallaxFactors  []float32
	playerPosition   rl.Vector2
)

func NewParallaxBackground(layerFiles []string, speeds []float32) *ParallaxBackground {
	layers := make([]ParallaxLayer, len(layerFiles))
	for i, file := range layerFiles {
		texture := rl.LoadTexture(file)
		layers[i] = ParallaxLayer{
			Texture:   texture,
			Speed:     speeds[i],
			PositionX: 0,
		}
	}
	return &ParallaxBackground{Layers: layers}
}

func (pb *ParallaxBackground) Update(cameraX float32) {
	for i := range pb.Layers {
		pb.Layers[i].PositionX = -cameraX * (pb.Layers[i].Speed * 0.5)
	}
}

func (pb *ParallaxBackground) Draw() {
	screenWidth := float32(rl.GetScreenWidth())
	for _, layer := range pb.Layers {
		// Draw the initial texture
		rl.DrawTexture(layer.Texture, int32(layer.PositionX), 0, rl.White)

		// Draw the additional textures to fill the screen
		for x := layer.PositionX + float32(layer.Texture.Width); x < screenWidth; x += float32(layer.Texture.Width) {
			rl.DrawTexture(layer.Texture, int32(x), 0, rl.White)
		}
		for x := layer.PositionX - float32(layer.Texture.Width); x > -float32(layer.Texture.Width); x -= float32(layer.Texture.Width) {
			rl.DrawTexture(layer.Texture, int32(x), 0, rl.White)
		}
	}
}

func LoadBackgrounds(playerSpeed float32) {
	backgroundLayers = []rl.Texture2D{
		rl.LoadTexture("assets/world/2 Background/Day/1.png"), // Farthest layer
		rl.LoadTexture("assets/world/2 Background/Day/2.png"), // Middle layer
		rl.LoadTexture("assets/world/2 Background/Day/3.png"), // Closest layer
		rl.LoadTexture("assets/world/2 Background/Day/4.png"), // Closest layer
	}
	parallaxFactors = []float32{
		playerSpeed,
		0.2,
		0.4,
		0.8,
	}
}

func UnloadBackgrounds() {
	for _, texture := range backgroundLayers {
		rl.UnloadTexture(texture)
	}
}
