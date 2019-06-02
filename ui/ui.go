package ui

import "image/color"

// Need Func for drawRect
// Need fucn to make main window
// Need way to get keyboard events / keyboard state
// Need func to get mouse events / mouse state

type Renderer interface {
	// DrawFilledRect(x int, y int, h int, w int, c color.RGBA)
	DrawFilledSquare(x, y, size int, c color.RGBA)
	// DrawFilledCircle(x, y, radius int, c color.RGBA)
}

type Renderable interface {
	Render(r Renderer)
}
