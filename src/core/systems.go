package core

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

func MovementSystem(w *World, dt float64) {
	for id, movement := range w.movements {
		transform, ok := w.transforms[id]

		if !ok {
			continue
		}

		if rl.IsKeyDown(rl.KeyRight) {
			transform.X += movement.Speed * dt
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			transform.X -= movement.Speed * dt
		}
		if rl.IsKeyDown(rl.KeyUp) {
			transform.Y -= movement.Speed * dt
		}
		if rl.IsKeyDown(rl.KeyDown) {
			transform.Y += movement.Speed * dt
		}
	}
}

func HealthSystem(w *World, dt float64) {
	for id, health := range w.healths {
		if health.InvincibleFor > 0 {
			health.InvincibleFor -= dt
		}
		if health.Current <= 0 {
			w.KillEntity(id)
		}
	}
}

func StaticMovementSystem(w *World, dt float64) {
	for id, staticMovement := range w.staticMovements {
		transform, ok := w.transforms[id]

		if !ok {
			continue
		}

		transform.X += staticMovement.DirectionX * staticMovement.Speed * dt
		transform.Y += staticMovement.DirectionY * staticMovement.Speed * dt
	}
}

func ShootingSystem(w *World, dt float64) {
	for id, shooter := range w.shooters {
		if shooter.Cooldown > 0 {
			shooter.Cooldown -= dt
			continue
		}

		if rl.IsKeyDown(rl.KeySpace) {
			playerTransform, ok := w.transforms[id]
			if !ok {
				continue
			}

			mouse := getLogicalMousePosition(*w.dimensions)
			dirX := float64(mouse.X) - playerTransform.X
			dirY := float64(mouse.Y) - playerTransform.Y
			length := math.Sqrt(dirX*dirX + dirY*dirY)
			if length > 0 {
				dirX /= length
				dirY /= length
			}

			bullet := w.NewEntity("Bullet", TypeBullet)
			w.AddTransform(bullet, playerTransform.X, playerTransform.Y)
			w.AddStaticMovement(bullet, dirX, dirY, 400)
			w.AddCollider(bullet, 5, 5, TagBullet)
			shooter.Cooldown = 0.5
		}
	}
}

func GameOverSystem(w *World) {
	playerAlive := false
	enemyAlive := false

	for _, entity := range w.entities {
		if entity.name == "Player" {
			playerAlive = true
		}
		if entity.name == "Enemy" {
			enemyAlive = true
		}
	}

	if !playerAlive {
		w.endState = Defeat
	} else if !enemyAlive {
		w.endState = Victory
	}
}

func CollisionSystem(w *World) {
	type collidable struct {
		id        uuid.UUID
		transform *Transform
		collider  *Collider
	}

	candidates := []collidable{}
	for id, collider := range w.colliders {
		transform, ok := w.transforms[id]

		if !ok {
			continue
		}

		candidates = append(candidates, collidable{id, transform, collider})
	}

	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			a := candidates[i]
			b := candidates[j]

			if aabbOverlap(a.transform, a.collider, b.transform, b.collider) {
				onCollision(w, a.id, b.id)
			}
		}
	}
}

func aabbOverlap(ta *Transform, ca *Collider, tb *Transform, cb *Collider) bool {
	return ta.X < tb.X+cb.Width &&
		ta.X+ca.Width > tb.X &&
		ta.Y < tb.Y+cb.Height &&
		ta.Y+ca.Height > tb.Y
}

func onCollision(w *World, aID, bID uuid.UUID) {
	aCol := w.colliders[aID]
	bCol := w.colliders[bID]

	pair := [2]ColliderTag{aCol.Tag, bCol.Tag}

	switch pair {
	case [2]ColliderTag{TagBullet, TagEnemy}, [2]ColliderTag{TagEnemy, TagBullet}:
		enemyID, bulletID := bID, aID
		if aCol.Tag == TagEnemy {
			enemyID, bulletID = aID, bID
		}
		if health, ok := w.healths[enemyID]; ok {
			health.Current--
		}
		w.KillEntity(bulletID)

	case [2]ColliderTag{TagPlayer, TagEnemy}, [2]ColliderTag{TagEnemy, TagPlayer}:
		playerID := bID
		if aCol.Tag == TagPlayer {
			playerID = aID
		}
		if health, ok := w.healths[playerID]; ok && health.InvincibleFor <= 0 {
			health.Current--
			health.InvincibleFor = PLAYER_INVICIBLE_AFTER_HIT_TIME // 1 second cooldown
		}

	case [2]ColliderTag{TagBullet, TagPlayer}, [2]ColliderTag{TagPlayer, TagBullet}:
		playerID, bulletID := bID, aID
		if aCol.Tag == TagPlayer {
			playerID, bulletID = aID, bID
		}
		if health, ok := w.healths[playerID]; ok && health.InvincibleFor <= 0 {
			health.Current--
			health.InvincibleFor = PLAYER_INVICIBLE_AFTER_HIT_TIME // 1 second cooldown
		}

		w.KillEntity(bulletID)

	case [2]ColliderTag{TagPlayer, TagWall}, [2]ColliderTag{TagWall, TagPlayer},
		[2]ColliderTag{TagEnemy, TagWall}, [2]ColliderTag{TagWall, TagEnemy}:
		// wall blocking — movement resolution not yet implemented
	}
}

func DespawnSystem(w *World) {
	for id := range w.entities {
		transform, ok := w.transforms[id]
		if !ok {
			continue
		}

		if transform.X < -DESPAWN_OFFSET || transform.X > w.dimensions.Width+DESPAWN_OFFSET || transform.Y < -DESPAWN_OFFSET || transform.Y > w.dimensions.Height+DESPAWN_OFFSET {
			w.KillEntity(id)
		}
	}
}

func EnemyShootingSystem(w *World, dt float64) {
	var playerX, playerY float64
	playerFound := false
	for id, entity := range w.entities {
		if entity.entityType == TypePlayer {
			if t, ok := w.transforms[id]; ok {
				playerX, playerY = t.X, t.Y
				playerFound = true
			}
			break
		}
	}
	if !playerFound {
		return
	}

	for id, entity := range w.entities {
		if entity.entityType != TypeEnemy {
			continue
		}
		shooter, ok := w.enemyShooters[id]
		if !ok {
			continue
		}
		if shooter.Cooldown > 0 {
			shooter.Cooldown -= dt
			continue
		}
		enemyTransform, ok := w.transforms[id]
		if !ok {
			continue
		}
		dirX := playerX - enemyTransform.X
		dirY := playerY - enemyTransform.Y
		length := math.Sqrt(dirX*dirX + dirY*dirY)
		if length > 0 {
			dirX /= length
			dirY /= length
		}
		bullet := w.NewEntity("EnemyBullet", TypeBullet)
		spawnOffset := 15.0
		w.AddTransform(bullet, enemyTransform.X+dirX*spawnOffset, enemyTransform.Y+dirY*spawnOffset)
		w.AddStaticMovement(bullet, dirX, dirY, 300)
		w.AddCollider(bullet, 5, 5, TagBullet)
		shooter.Cooldown = 1.0
	}
}

func EnemyMovementSystem(w *World, dt float64) {

}
