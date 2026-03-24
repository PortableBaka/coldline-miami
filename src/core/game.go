package core

type Game struct {
	world *World
}

func NewGame() *Game {
	return &Game{
		world: NewWorld(),
	}
}
