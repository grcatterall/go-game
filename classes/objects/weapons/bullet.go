package weapons

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Active   bool
	Width    float32
	Height   float32
}

// NewBullet creates a new bullet instance
func NewBullet(position, speed rl.Vector2, width, height float32) *Bullet {
	return &Bullet{
		Position: position,
		Speed:    speed,
		Active:   true,
		Width:    width,
		Height:   height,
	}
}

// Update updates the bullet's position and active status
func (b *Bullet) Update() {
	if b.Active {
		b.Position.X += b.Speed.X
		b.Position.Y += b.Speed.Y
		fmt.Println("Bullet position:", b.Position)

		// Deactivate the bullet if it moves off-screen (example logic)
		if b.Position.X < 0 || b.Position.X > float32(rl.GetScreenWidth()) ||
			b.Position.Y < 0 || b.Position.Y > float32(rl.GetScreenHeight()) {
			b.Active = false
		}
	}
}

// Draw renders the bullet as a yellow square
func (b *Bullet) Draw() {
	if b.Active {
		rl.DrawRectangle(int32(b.Position.X), int32(b.Position.Y), int32(b.Width), int32(b.Height), rl.Red)
	}
}

// CheckCollision checks if the bullet collides with a given rectangle
func (b *Bullet) CheckCollision(target rl.Rectangle) bool {
	bulletRect := rl.NewRectangle(b.Position.X, b.Position.Y, b.Width, b.Height)
	return rl.CheckCollisionRecs(bulletRect, target)
}

// SpawnBullet spawns a new bullet from the player's position
func SpawnBullet(position rl.Vector2, isLeft bool) *Bullet {
	bulletSpeed := rl.Vector2{X: 1, Y: 0} // Example speed, adjust as needed
	if isLeft {
		bulletSpeed.X = -bulletSpeed.X // Reverse direction if facing left
	}
	newBullet := NewBullet(position, bulletSpeed, 128, 128) // Example width and height, adjust as needed
	return newBullet
}
