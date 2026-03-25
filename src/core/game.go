package core

import rl "github.com/gen2brain/raylib-go/raylib"

type Game struct {
	world        *World
	renderTarget rl.RenderTexture2D
}

func NewGame() *Game {
	return &Game{
		world: NewWorld(),
	}
}
