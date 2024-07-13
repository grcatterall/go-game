package characters

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	Texture           rl.Texture2D
	Position          rl.Vector2
	FrameRec          rl.Rectangle
	FrameWidth        float32
	FrameHeight       float32
	FrameSpeed        float32
	FrameCounter      float32
	FramesCount       int32
	CurrentFrame      int32
	Health            int32
	idleTexture       rl.Texture2D
	walkingTexture    rl.Texture2D
	attackingTexture  rl.Texture2D
	Target            *Player
	speed             float32
	detectionDistance float32
	attackRange       float32
	isMoving          bool
	isAttacking       bool
}

func NewEnemy(
	character string,
	health int32,
	position rl.Vector2,
	frameSpeed float32,
	spriteHeight int32,
	spriteWidth int32,
	target *Player) *Enemy {

	idleTexture := rl.LoadTexture(fmt.Sprintf("assets/characters/%s/Idle.png", character))
	walkingTexture := rl.LoadTexture(fmt.Sprintf("assets/characters/%s/Walk.png", character))
	attackingTexture := rl.LoadTexture(fmt.Sprintf("assets/characters/%s/Attack_1.png", character))

	framesCount := idleTexture.Width / spriteWidth // Assuming each frame is 128x128 pixels
	frameWidth := float32(idleTexture.Width) / float32(framesCount)
	frameHeight := float32(idleTexture.Height)

	return &Enemy{
		Texture:  idleTexture,
		Position: position,
		FrameRec: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  frameWidth,
			Height: frameHeight,
		},
		CurrentFrame:      0,
		FrameWidth:        frameWidth,
		FrameHeight:       frameHeight,
		FrameSpeed:        frameSpeed,
		FramesCount:       framesCount,
		idleTexture:       idleTexture,
		walkingTexture:    walkingTexture,
		attackingTexture:  attackingTexture,
		Health:            health,
		Target:            target,
		speed:             0.5,
		detectionDistance: 300,
		attackRange:       30,
	}
}

func (e *Enemy) Update() {
	// Update animation
	e.FrameCounter += e.FrameSpeed
	if e.FrameCounter >= 1 {
		e.FrameCounter = 0
		e.CurrentFrame++
		if e.CurrentFrame >= e.FramesCount {
			e.CurrentFrame = 0
		}
		e.FrameRec.X = float32(e.CurrentFrame) * e.FrameWidth
	}

	distanceToTarget := e.Position.X - e.Target.Position.X

	if distanceToTarget < 0 {
		distanceToTarget *= -1
	}

	if distanceToTarget < e.attackRange {
		e.isAttacking = true
	} else if distanceToTarget <= e.detectionDistance && distanceToTarget > e.attackRange+1 {
		e.isMoving = true
		e.moveToTarget()
	} else {
		e.isMoving = false
		e.isAttacking = false
	}

	if e.isMoving {
		e.Texture = e.walkingTexture
	} else if e.isAttacking {
		e.Texture = e.attackingTexture
	} else {
		e.Texture = e.idleTexture
	}

	e.renderTexture(e.Texture)
}

func (e *Enemy) Draw() {
	drawPosition := rl.Vector2{X: e.Position.X, Y: e.Position.Y}
	if e.Target.Position.X > e.Position.X {
		rl.DrawTextureRec(e.Texture, e.FrameRec, drawPosition, rl.White)
	} else {
		rl.DrawTextureRec(e.Texture, rl.Rectangle{X: e.FrameRec.X + e.FrameRec.Width, Y: e.FrameRec.Y, Width: -e.FrameRec.Width, Height: e.FrameRec.Height}, drawPosition, rl.White)

	}
}

// Unload releases the texture resources
func (e *Enemy) Unload() {
	rl.UnloadTexture(e.Texture)
}

func (e *Enemy) moveToTarget() {
	if e.Target.Position.X < e.Position.X {
		e.Position.X -= e.speed
	} else {
		e.Position.X += e.speed
	}
}

func (e *Enemy) renderTexture(texture rl.Texture2D) {
	e.FrameWidth = float32(texture.Width) / float32(texture.Width/128) // Assuming each frame is 128x128 pixels
	e.FrameHeight = float32(texture.Height)
	e.FramesCount = texture.Width / 128
}
