package core

import (
	"github.com/google/uuid"
)

type WorldDimensions struct {
	Width  float64
	Height float64
}

type World struct {
	entities        map[uuid.UUID]*Entity
	transforms      map[uuid.UUID]*Transform
	movements       map[uuid.UUID]*Movement
	staticMovements map[uuid.UUID]*StaticMovement
	healths         map[uuid.UUID]*Health
	colliders       map[uuid.UUID]*Collider
	shooters        map[uuid.UUID]*Shooter
	enemyShooters   map[uuid.UUID]*EnemyShooter

	endState   GameEndState
	settings   *Settings
	dimensions *WorldDimensions
}

func NewWorld() *World {
	return &World{
		entities:        make(map[uuid.UUID]*Entity),
		transforms:      make(map[uuid.UUID]*Transform),
		movements:       make(map[uuid.UUID]*Movement),
		staticMovements: make(map[uuid.UUID]*StaticMovement),
		healths:         make(map[uuid.UUID]*Health),
		colliders:       make(map[uuid.UUID]*Collider),
		shooters:        make(map[uuid.UUID]*Shooter),
		enemyShooters:   make(map[uuid.UUID]*EnemyShooter),

		endState: In_Progress,
		settings: &Settings{
			debug: false,
		},
		dimensions: &WorldDimensions{},
	}
}

func (w *World) NewEntity(name string, entityType EntityType) *Entity {
	id := uuid.New()

	entity := Entity{
		ID:         id,
		name:       name,
		entityType: entityType,
	}

	w.entities[id] = &entity

	return &entity
}

func (w *World) AddTransform(entity *Entity, x, y float64) {
	w.transforms[entity.ID] = &Transform{
		X: x,
		Y: y,
	}
}

func (w *World) AddMovement(entity *Entity, speed float64) {
	w.movements[entity.ID] = &Movement{
		Speed: speed,
	}
}

func (w *World) AddStaticMovement(entity *Entity, directionX, directionY, speed float64) {
	w.staticMovements[entity.ID] = &StaticMovement{
		DirectionX: directionX,
		DirectionY: directionY,
		Speed:      speed,
	}
}

func (w *World) AddHealth(entity *Entity, current, max int) {
	w.healths[entity.ID] = &Health{
		Current: current,
		Max:     max,
	}
}

func (w *World) AddCollider(entity *Entity, width, height float64, tag ColliderTag) {
	w.colliders[entity.ID] = &Collider{
		Width:  width,
		Height: height,
		Tag:    tag,
	}
}

func (w *World) AddShooter(entity *Entity, cooldown float64) {
	w.shooters[entity.ID] = &Shooter{
		Cooldown: cooldown,
	}
}

func (w *World) AddEnemyShooter(entity *Entity, cooldown float64) {
	w.enemyShooters[entity.ID] = &EnemyShooter{
		Cooldown: cooldown,
	}
}

func (w *World) KillEntity(id uuid.UUID) (*Entity, bool) {
	entity, exists := w.entities[id]

	if !exists {
		return nil, false
	}

	delete(w.entities, id)
	delete(w.transforms, id)
	delete(w.movements, id)
	delete(w.staticMovements, id)
	delete(w.healths, id)
	delete(w.colliders, id)
	delete(w.shooters, id)
	delete(w.enemyShooters, id)
	return entity, true
}

func (w *World) TurnOnDebug() {
	w.settings.debug = true
}

func (w *World) TurnOffDebug() {
	w.settings.debug = false
}
