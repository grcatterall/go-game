package main

import (
	"example/go-game/classes/characters"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize the window
	screenWidth := int32(800)
	screenHeight := int32(450)
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - sprite animation")

	// Create a new player
	player := characters.NewPlayer(rl.Vector2{X: float32(screenWidth)/4 - 128, Y: float32(screenHeight)/2 - 64}, 0.2)

	enemy := characters.NewEnemy(
		"Raider_1",
		5,
		rl.Vector2{X: float32(screenWidth) - 128, Y: float32(screenHeight)/2 - 64},
		0.2,
		128,
		128,
		player,
	)

	// Set the target frames per second
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Update
		player.Update()

		enemy.Update()

		// Start drawing
		rl.BeginDrawing()

		// Clear the background
		rl.ClearBackground(rl.RayWhite)

		// Draw the player
		player.Draw()

		enemy.Draw()

		// End drawing
		rl.EndDrawing()
	}

	// Unload player texture
	player.Unload()

	enemy.Unload()

	// Close the window
	rl.CloseWindow()
}
