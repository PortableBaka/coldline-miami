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

type Collider struct {
	Width, Height float64
	Tag           ColliderTag
}

type Shooter struct {
	Cooldown float64
}
