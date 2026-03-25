package core

import rl "github.com/gen2brain/raylib-go/raylib"

func getLogicalMousePosition(dimensions WorldDimensions) rl.Vector2 {
	m := rl.GetMousePosition()
	scaleX := float32(dimensions.Width) / float32(rl.GetScreenWidth())
	scaleY := float32(dimensions.Height) / float32(rl.GetScreenHeight())
	return rl.Vector2{X: m.X * scaleX, Y: m.Y * scaleY}
}
