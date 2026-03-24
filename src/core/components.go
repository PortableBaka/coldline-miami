package core

type Transform struct {
	X, Y float64
}

type Movement struct {
	Speed float64
}

type StaticMovement struct {
	DirectionX, DirectionY float64
	Speed                  float64
}

type Sprite struct {
	ImagePath string
}

type Health struct {
	Current       int
	Max           int
	InvincibleFor float64
}

type ColliderTag string

const (
	TagPlayer ColliderTag = "player"
	TagEnemy  ColliderTag = "enemy"
	TagWall   ColliderTag = "wall"
	TagBullet ColliderTag = "bullet"
)

type Collider struct {
	Width, Height float64
	Tag           ColliderTag
}

type Shooter struct {
	Cooldown float64
}

type GameState string

const (
	StatePlaying  GameState = "playing"
	StatePaused   GameState = "paused"
	StateGameOver GameState = "game_over"
)

type GameEndState string

const (
	Victory     GameEndState = "victory"
	Defeat      GameEndState = "defeat"
	In_Progress GameEndState = "in_progress"
)
