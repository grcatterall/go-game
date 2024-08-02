package characters

import (
	"fmt"

	"github.com/grcatterall/go-game/classes/game_manager"
	"github.com/grcatterall/go-game/classes/helpers"
	"github.com/grcatterall/go-game/classes/objects/weapons"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Position           rl.Vector2
	Velocity           rl.Vector2
	Gravity            float32
	Speed              float32
	IsMoving           bool
	IsRunning          bool
	IsGrounded         bool
	IsLeft             bool
	IsShooting         bool
	IsAttacking        bool
	IdleAnimation      helpers.Animation
	WalkingAnimation   helpers.Animation
	RunningAnimation   helpers.Animation
	ShootingAnimation  helpers.Animation
	AttackingAnimation helpers.Animation
	CurrentAnimation   *helpers.Animation
	Bullets            []*weapons.Bullet
}

var mainSprite = "Soldier_1"

// NewPlayer creates a new player with the given sprite and position
func NewPlayer(position rl.Vector2, frameSpeed float32) *Player {
	return &Player{
		Position:           position,
		Velocity:           rl.Vector2{X: 0, Y: 0},
		Speed:              1.0,
		Gravity:            0.1,
		IdleAnimation:      helpers.LoadAnimation(fmt.Sprintf("assets/characters/%s/Idle.png", mainSprite), frameSpeed, 128),
		WalkingAnimation:   helpers.LoadAnimation(fmt.Sprintf("assets/characters/%s/Walk.png", mainSprite), frameSpeed, 128),
		RunningAnimation:   helpers.LoadAnimation(fmt.Sprintf("assets/characters/%s/Run.png", mainSprite), frameSpeed, 128),
		ShootingAnimation:  helpers.LoadAnimation(fmt.Sprintf("assets/characters/%s/Shot_1.png", mainSprite), 0.1, 128),
		AttackingAnimation: helpers.LoadAnimation(fmt.Sprintf("assets/characters/%s/Attack.png", mainSprite), 0.1, 128),
		CurrentAnimation:   nil,
		Bullets:            []*weapons.Bullet{},
	}
}

// Update updates the player animation and movement
func (p *Player) Update(tileMap *game_manager.TileMap) {
	p.Velocity.Y += p.Gravity

	p.updateAnimation()
	p.updateMovement()
	p.updateActions()
	p.selectCurrentAnimation()
	p.updateFrameRec()

	// Update bullets
	for _, bullet := range p.Bullets {
		bullet.Update()
	}

	// Remove inactive bullets
	p.Bullets = filterActiveBullets(p.Bullets)

	// Apply velocity to position
	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
}

// updateAnimation updates the animation frames
func (p *Player) updateAnimation() {
	if p.CurrentAnimation == nil {
		return
	}

	p.CurrentAnimation.FrameCounter += p.CurrentAnimation.FrameSpeed
	if p.CurrentAnimation.FrameCounter >= 1 {
		p.CurrentAnimation.FrameCounter = 0
		p.CurrentAnimation.CurrentFrame++
		if p.CurrentAnimation.CurrentFrame >= p.CurrentAnimation.Frames {
			p.CurrentAnimation.CurrentFrame = 0
			if p.IsShooting || p.IsAttacking {
				p.IsShooting = false
				p.IsAttacking = false
			}
		}
		p.CurrentAnimation.FrameRec.X = float32(p.CurrentAnimation.CurrentFrame) * p.CurrentAnimation.FrameRec.Width
	}
}

// updateMovement updates the player's position based on input
func (p *Player) updateMovement() {
	if rl.IsKeyDown(rl.KeySpace) && p.IsGrounded {
		p.Velocity.Y = -3
		p.IsGrounded = false
	}

	if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyA)) && !p.IsShooting {
		p.IsMoving = true
		p.IsRunning = false

		var speed = p.Speed

		if rl.IsKeyDown(rl.KeyLeftShift) {
			p.IsRunning = true
			speed *= 2
		}

		if rl.IsKeyDown(rl.KeyD) {
			p.IsLeft = false
			p.Velocity.X = speed
		}
		if rl.IsKeyDown(rl.KeyA) {
			p.IsLeft = true
			p.Velocity.X = -speed
		}
	} else {
		p.IsMoving = false
		p.Velocity.X = 0
	}
}

// updateActions handles shooting and attacking actions
func (p *Player) updateActions() {
	if rl.IsMouseButtonDown(0) && !p.IsShooting {
		p.IsShooting = true
		p.ShootingAnimation.CurrentFrame = 0
		p.ShootingAnimation.FrameCounter = 0
		p.Bullets = append(p.Bullets, weapons.SpawnBullet(rl.Vector2{X: p.Position.X + 1, Y: p.Position.Y}, p.IsLeft))

	}

	if rl.IsMouseButtonDown(1) && !p.IsAttacking {
		p.IsAttacking = true
		p.AttackingAnimation.CurrentFrame = 0
		p.AttackingAnimation.FrameCounter = 0
	}
}

// selectCurrentAnimation selects the appropriate animation based on player state
func (p *Player) selectCurrentAnimation() {
	switch {
	case p.IsShooting:
		p.CurrentAnimation = &p.ShootingAnimation
	case p.IsAttacking:
		p.CurrentAnimation = &p.AttackingAnimation
	case p.IsMoving:
		if p.IsRunning {
			p.CurrentAnimation = &p.RunningAnimation
		} else {
			p.CurrentAnimation = &p.WalkingAnimation
		}
	default:
		p.CurrentAnimation = &p.IdleAnimation
	}
}

// updateFrameRec updates the frame rectangle for the current animation
func (p *Player) updateFrameRec() {
	if p.CurrentAnimation != nil {
		p.CurrentAnimation.FrameRec.Width = float32(p.CurrentAnimation.Texture.Width) / float32(p.CurrentAnimation.Frames)
		p.CurrentAnimation.FrameRec.Height = float32(p.CurrentAnimation.Texture.Height)
	}
}

// Draw renders the player sprite on the screen
func (p *Player) Draw() {
	if p.CurrentAnimation == nil {
		return
	}

	flipX := true
	if p.IsLeft {
		flipX = false
	}

	drawPosition := rl.Vector2{X: p.Position.X, Y: p.Position.Y}

	if flipX {
		rl.DrawTextureRec(p.CurrentAnimation.Texture, p.CurrentAnimation.FrameRec, drawPosition, rl.White)
	} else {
		rl.DrawTextureRec(p.CurrentAnimation.Texture, rl.Rectangle{X: p.CurrentAnimation.FrameRec.X + p.CurrentAnimation.FrameRec.Width, Y: p.CurrentAnimation.FrameRec.Y, Width: -p.CurrentAnimation.FrameRec.Width, Height: p.CurrentAnimation.FrameRec.Height}, drawPosition, rl.White)
	}
}

// Unload releases the texture resources
func (p *Player) Unload() {
	rl.UnloadTexture(p.IdleAnimation.Texture)
	rl.UnloadTexture(p.WalkingAnimation.Texture)
	rl.UnloadTexture(p.RunningAnimation.Texture)
	rl.UnloadTexture(p.ShootingAnimation.Texture)
	rl.UnloadTexture(p.AttackingAnimation.Texture)
}

// filterActiveBullets removes inactive bullets from the list
func filterActiveBullets(bullets []*weapons.Bullet) []*weapons.Bullet {
	activeBullets := []*weapons.Bullet{}
	for _, bullet := range bullets {
		if bullet.Active {
			activeBullets = append(activeBullets, bullet)
		}
	}
	return activeBullets
}

func (p *Player) CheckCollisions(tileMap *game_manager.TileMap, camera rl.Camera2D) {
	// Adjust the player position by the camera offset
	adjustedPosition := rl.Vector2{
		X: p.Position.X - camera.Target.X + camera.Offset.X,
		Y: p.Position.Y - camera.Target.Y + camera.Offset.Y,
	}

	// Calculate the body and feet rectangles relative to the adjusted player position
	feetRect := rl.NewRectangle(adjustedPosition.X+47, adjustedPosition.Y+p.CurrentAnimation.FrameRec.Height*0.9, p.CurrentAnimation.FrameRec.Width-90, p.CurrentAnimation.FrameRec.Height*0.2)
	bodyRect := rl.NewRectangle(adjustedPosition.X+40, adjustedPosition.Y+64, p.CurrentAnimation.FrameRec.Width-75, p.CurrentAnimation.FrameRec.Height*0.9-64)

	// Draw the rectangles for debugging
	rl.DrawRectangleLines(int32(bodyRect.X), int32(bodyRect.Y), int32(bodyRect.Width), int32(bodyRect.Height), rl.Green)
	rl.DrawRectangleLines(int32(feetRect.X), int32(feetRect.Y), int32(feetRect.Width), int32(feetRect.Height), rl.Green)

	// Calculate the range of tiles to check
	tileWidth := 32  // Assuming each tile is 32 pixels wide
	tileHeight := 32 // Assuming each tile is 32 pixels high

	startX := int(p.Position.X) / tileWidth
	endX := (int(p.Position.X) + int(p.CurrentAnimation.FrameRec.Width)) / tileWidth
	startY := int(p.Position.Y) / tileHeight
	endY := (int(p.Position.Y) + int(p.CurrentAnimation.FrameRec.Height)) / tileHeight

	p.IsGrounded = false

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if y >= 0 && y < len(tileMap.Tiles) && x >= 0 && x < len(tileMap.Tiles[y]) {
				tile := tileMap.Tiles[y][x]
				if tile != nil {
					// Adjust the tile position by the camera offset
					adjustedTilePosition := rl.Vector2{
						X: tile.Position.X - camera.Target.X + camera.Offset.X,
						Y: tile.Position.Y - camera.Target.Y + camera.Offset.Y,
					}
					tileRect := rl.NewRectangle(adjustedTilePosition.X, adjustedTilePosition.Y, float32(tileWidth), float32(tileHeight))
					highlightColor := rl.Green // Default color for checked tiles

					// Check horizontal collisions with the body rectangle
					if rl.CheckCollisionRecs(bodyRect, tileRect) {
						highlightColor = rl.Red // Color for tiles with collision

						// Handle horizontal collisions
						if p.Velocity.X > 0 { // Moving right
							p.Position.X = tile.Position.X - p.CurrentAnimation.FrameRec.Width + 35
							p.Velocity.X = 0
						} else if p.Velocity.X < 0 { // Moving left
							p.Position.X = tile.Position.X + tileRect.Width - 41
							p.Velocity.X = 0
						}
					}

					// Draw a rectangle around the tile being checked
					rl.DrawRectangleLines(int32(tileRect.X), int32(tileRect.Y), int32(tileRect.Width), int32(tileRect.Height), highlightColor)
				}
			}
		}
	}

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if y >= 0 && y < len(tileMap.Tiles) && x >= 0 && x < len(tileMap.Tiles[y]) {
				tile := tileMap.Tiles[y][x]
				if tile != nil {
					// Adjust the tile position by the camera offset
					adjustedTilePosition := rl.Vector2{
						X: tile.Position.X - camera.Target.X + camera.Offset.X,
						Y: tile.Position.Y - camera.Target.Y + camera.Offset.Y,
					}
					tileRect := rl.NewRectangle(adjustedTilePosition.X, adjustedTilePosition.Y, float32(tileWidth), float32(tileHeight))
					highlightColor := rl.Green // Default color for checked tiles

					// Check vertical collisions with the feet rectangle
					if rl.CheckCollisionRecs(feetRect, tileRect) {
						highlightColor = rl.Red // Color for tiles with collision

						// Handle vertical collisions
						if p.Velocity.Y > 0 { // Falling down
							p.Position.Y = tile.Position.Y - p.CurrentAnimation.FrameRec.Height
							p.Velocity.Y = 0
							p.IsGrounded = true
						}
					}

					// Draw a rectangle around the tile being checked
					rl.DrawRectangleLines(int32(tileRect.X), int32(tileRect.Y), int32(tileRect.Width), int32(tileRect.Height), highlightColor)
				}
			}
		}
	}
}
