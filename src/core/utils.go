package core

import rl "github.com/gen2brain/raylib-go/raylib"

func getLogicalMousePosition(dimensions WorldDimensions) rl.Vector2 {
	logW := float32(dimensions.Width)
	logH := float32(dimensions.Height)
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	scale := min(screenW/logW, screenH/logH)
	offsetX := (screenW - logW*scale) / 2
	offsetY := (screenH - logH*scale) / 2
	m := rl.GetMousePosition()
	return rl.Vector2{
		X: (m.X - offsetX) / scale,
		Y: (m.Y - offsetY) / scale,
	}
}
