package core

import (
	cfg "coldline-miami/src/config"
	"fmt"
	"log"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) Loop() {
	g.init()
	defer g.deinit()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		g.update(float64(dt))
		g.render()
	}
}

func (g *Game) init() {
	config, err := cfg.LoadConfig()
	if err != nil {
		log.Printf("failed to load config: %v, using defaults", err)
		config = cfg.DefaultConfig()
	}

	g.world.dimensions.Width = float64(config.LogicalWidth)
	g.world.dimensions.Height = float64(config.LogicalHeight)

	rl.SetTargetFPS(config.FPSLimit)
	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, config.Title)
	g.renderTarget = rl.LoadRenderTexture(config.LogicalWidth, config.LogicalHeight)

	fmt.Println(config)

	if config.Debug {
		g.world.TurnOnDebug()
	}

	player := g.world.NewEntity("Player", TypePlayer)
	enemy := g.world.NewEntity("Enemy", TypeEnemy)

	g.world.AddTransform(enemy, 500, 300)
	g.world.AddHealth(enemy, 5, 5)
	g.world.AddCollider(enemy, 20, 20, TagEnemy)
	g.world.AddEnemyShooter(enemy, 1.0)

	g.world.AddTransform(player, 300, 300)
	g.world.AddMovement(player, 100)
	g.world.AddHealth(player, 10, 10)
	g.world.AddCollider(player, 20, 20, TagPlayer)
	g.world.AddShooter(player, 0.5)

}

func (g *Game) deinit() {
	rl.UnloadRenderTexture(g.renderTarget)
	rl.CloseWindow()
}

func (g *Game) update(dt float64) {
	DespawnSystem(g.world)
	MovementSystem(g.world, dt)
	ShootingSystem(g.world, dt)
	EnemyShootingSystem(g.world, dt)
	StaticMovementSystem(g.world, dt)
	CollisionSystem(g.world)
	HealthSystem(g.world, dt)
	GameOverSystem(g.world)
}

func (g *Game) render() {
	rl.BeginTextureMode(g.renderTarget)
	rl.ClearBackground(rl.RayWhite)

	renderDebugInfo(g.world)
	renderGameState(g.world)

	rl.EndTextureMode()

	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	logW := float32(g.renderTarget.Texture.Width)
	logH := float32(g.renderTarget.Texture.Height)
	scale := min(screenW/logW, screenH/logH)
	src := rl.Rectangle{X: 0, Y: 0, Width: logW, Height: -logH}
	dst := rl.Rectangle{
		X:      (screenW - logW*scale) / 2,
		Y:      (screenH - logH*scale) / 2,
		Width:  logW * scale,
		Height: logH * scale,
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(g.renderTarget.Texture, src, dst, rl.Vector2{}, 0, rl.White)
	rl.EndDrawing()
}

func renderDebugInfo(w *World) {
	if w.settings.debug {
		rl.DrawText("Debug Mode ON", 10, 10, 20, rl.Red)
		rl.DrawText("FPS: "+strconv.Itoa(int(rl.GetFPS())), 10, 30, 20, rl.Red)

		for _, entity := range w.entities {
			transform := w.transforms[entity.ID]

			if health, ok := w.healths[entity.ID]; ok {
				rl.DrawText("HP: "+strconv.Itoa(health.Current), int32(transform.X-10), int32(transform.Y-20), 10, rl.Red)
			}
		}
	}
}

func renderGameState(w *World) {
	if w.endState != In_Progress {
		var msg string
		if w.endState == Victory {
			msg = "You Win!"
		} else {
			msg = "Game Over"
		}
		rl.DrawText(msg, 350, 280, 30, rl.Black)
	} else {
		for _, entity := range w.entities {
			transform := w.transforms[entity.ID]

			switch entity.name {
			case "Player":
				rl.DrawCircle(int32(transform.X), int32(transform.Y), 10, rl.Blue)
			case "Enemy":
				rl.DrawCircle(int32(transform.X), int32(transform.Y), 10, rl.Red)
			case "Bullet":
				rl.DrawCircle(int32(transform.X), int32(transform.Y), 5, rl.Black)
			}
		}
	}
}
