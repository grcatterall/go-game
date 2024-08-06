package main

import (
	"github.com/grcatterall/go-game/classes/characters"
	"github.com/grcatterall/go-game/classes/game_manager"
	"github.com/grcatterall/go-game/classes/game_manager/levels"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize the window
	screenWidth := int32(960)
	screenHeight := int32(512)
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - sprite animation")

	// Create a new player
	player := characters.NewPlayer(rl.Vector2{X: float32(screenWidth)/4 - 128, Y: float32(screenHeight)/2 - 30}, 0.2)

	enemy := characters.NewEnemy(
		"Raider_1",
		5,
		rl.Vector2{X: float32(screenWidth) - 128, Y: float32(screenHeight)/2 - 64},
		0.2,
		128,
		128,
		player,
	)

	var tileTextures = map[int]rl.Texture2D{
		1:  game_manager.LoadTile("assets/world/1 Tiles/Tile_01.png"),
		2:  game_manager.LoadTile("assets/world/1 Tiles/Tile_02.png"),
		3:  game_manager.LoadTile("assets/world/1 Tiles/Tile_40.png"),
		4:  game_manager.LoadTile("assets/world/1 Tiles/Tile_39.png"),
		5:  game_manager.LoadTile("assets/world/1 Tiles/Tile_42.png"),
		6:  game_manager.LoadTile("assets/world/1 Tiles/Tile_43.png"),
		7:  game_manager.LoadTile("assets/world/1 Tiles/Tile_37.png"),
		8:  game_manager.LoadTile("assets/world/1 Tiles/Tile_38.png"),
		9:  game_manager.LoadTile("assets/world/1 Tiles/Tile_31.png"),
		10: game_manager.LoadTile("assets/world/1 Tiles/Tile_30.png"),
		11: game_manager.LoadTile("assets/world/1 Tiles/Tile_29.png"),
		12: game_manager.LoadTile("assets/world/1 Tiles/Tile_28.png"),
	}

	tileMap := game_manager.LoadLevel(levels.GetLevel(1), tileTextures)

	// Load parallax background layers
	layerFiles := []string{
		"assets/world/2 Background/Day/1.png",
		"assets/world/2 Background/Day/2.png",
		"assets/world/2 Background/Day/3.png",
		"assets/world/2 Background/Day/4.png",
	}
	speeds := []float32{
		player.Speed,
		0.2,
		0.4,
		0.8,
	}

	parallaxBackground := game_manager.NewParallaxBackground(layerFiles, speeds)

	camera := rl.NewCamera2D(player.Position, rl.Vector2{X: 0, Y: 0}, 0, 1)

	// Set the target frames per second
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		camera.Target = player.Position

		// Start drawing
		rl.BeginDrawing()

		// Clear the background
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(camera)

		parallaxBackground.Update(player.Position.X)
		parallaxBackground.Draw()
		rl.DrawText("< a d > - shift sprint - left click shoot - right click melee", 100, 30, 20, rl.Black)

		// Update
		player.Update(tileMap)

		player.CheckCollisions(tileMap, camera)

		enemy.Update()

		tileMap.Draw()

		// Draw the player
		player.Draw()

		enemy.Draw()

		rl.EndMode2D()

		// End drawing
		rl.EndDrawing()
	}

	// Unload player texture
	player.Unload()

	enemy.Unload()

	// Close the window
	rl.CloseWindow()
}
