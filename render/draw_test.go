package render

import "image/color"

func DrawExample() {
	// We haven't modified the draw stack, so it contains a single draw heap.
	// Draw a Color Box
	Draw(NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}), 3)
	// Draw a Gradient Box above that color box
	Draw(NewHorizontalGradientBox(5, 5, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}), 4)
}
