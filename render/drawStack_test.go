package render

import "testing"

func TestDrawStack(t *testing.T) {
	resetDraw = true
	PreDraw()
	resetDraw = false
	PreDraw()
	SetDrawStack(
		NewHeap(false),
		NewDrawFPS(),
	)
}
