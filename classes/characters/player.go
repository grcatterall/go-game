package characters

import (
	"fmt"

	"example/go-game/classes/helpers"
	"example/go-game/classes/objects/weapons"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Position           rl.Vector2
	Speed              float32
	IsMoving           bool
	IsRunning          bool
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
		Speed:              1.0,
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
func (p *Player) Update() {
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
			p.Position.X += speed
		}
		if rl.IsKeyDown(rl.KeyA) {
			p.IsLeft = true
			p.Position.X -= speed
		}
	} else {
		p.IsMoving = false
	}
}

// updateActions handles shooting and attacking actions
func (p *Player) updateActions() {
	if rl.IsMouseButtonDown(0) && !p.IsShooting {
		p.IsShooting = true
		p.ShootingAnimation.CurrentFrame = 0
		p.ShootingAnimation.FrameCounter = 0
		p.Bullets = append(p.Bullets, weapons.SpawnBullet(rl.Vector2{X: p.Position.X + 1, Y: p.Position.Y / 2}, p.IsLeft))

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
