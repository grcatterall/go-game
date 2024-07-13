package helpers

import rl "github.com/gen2brain/raylib-go/raylib"

type Animation struct {
	Texture      rl.Texture2D
	Frames       int32
	FrameRec     rl.Rectangle
	CurrentFrame int32
	FrameSpeed   float32
	FrameCounter float32
}

// loadAnimation loads an animation from a texture file
func LoadAnimation(filePath string, frameSpeed float32, frameSize int32) Animation {
	texture := rl.LoadTexture(filePath)
	frames := texture.Width / frameSize
	frameWidth := float32(texture.Width) / float32(frames)
	frameHeight := float32(texture.Height)

	return Animation{
		Texture:      texture,
		Frames:       frames,
		FrameRec:     rl.Rectangle{X: 0, Y: 0, Width: frameWidth, Height: frameHeight},
		CurrentFrame: 0,
		FrameSpeed:   frameSpeed,
		FrameCounter: 0,
	}
}
