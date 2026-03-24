package core

import (
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
	rl.SetTargetFPS(60)
	rl.InitWindow(800, 600, "Coldline Miami")

	g.world.TurnOnDebug()

	player := g.world.NewEntity("Player")
	enemy := g.world.NewEntity("Enemy")

	g.world.AddTransform(enemy, 500, 300)
	g.world.AddHealth(enemy, 5, 5)
	g.world.AddCollider(enemy, 20, 20, TagEnemy)

	g.world.AddTransform(player, 300, 300)
	g.world.AddMovement(player, 100)
	g.world.AddHealth(player, 10, 10)
	g.world.AddCollider(player, 20, 20, TagPlayer)
	g.world.AddShooter(player, 0.5)

}

func (g *Game) deinit() {
	rl.CloseWindow()
}

func (g *Game) update(dt float64) {
	MovementSystem(g.world, dt)
	ShootingSystem(g.world, dt)
	StaticMovementSystem(g.world, dt)
	CollisionSystem(g.world)
	HealthSystem(g.world, dt)
	GameOverSystem(g.world)
}

func (g *Game) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	if g.world.endState != In_Progress {
		var msg string
		if g.world.endState == Victory {
			msg = "You Win!"
		} else {
			msg = "Game Over"
		}
		rl.DrawText(msg, 350, 280, 30, rl.Black)
		rl.EndDrawing()
		return
	}

	if g.world.settings.debug {
		rl.DrawText("Debug Mode ON", 10, 10, 20, rl.Red)
		rl.DrawText(strconv.Itoa(int(rl.GetFPS())), 10, 20, 20, rl.Red)
	}

	for _, entity := range g.world.entities {
		transform := g.world.transforms[entity.ID]

		switch entity.name {
		case "Player":
			rl.DrawCircle(int32(transform.X), int32(transform.Y), 10, rl.Blue)
		case "Enemy":
			rl.DrawCircle(int32(transform.X), int32(transform.Y), 10, rl.Red)
		case "Bullet":
			rl.DrawCircle(int32(transform.X), int32(transform.Y), 5, rl.Black)
		}

		if g.world.settings.debug {
			if health, ok := g.world.healths[entity.ID]; ok {
				rl.DrawText("HP: "+strconv.Itoa(health.Current), int32(transform.X-10), int32(transform.Y-20), 10, rl.Red)
			}
		}
	}

	rl.EndDrawing()
}
