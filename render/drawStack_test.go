package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawStack(t *testing.T) {
	resetDraw = true
	PreDraw()
	resetDraw = false
	PreDraw()
	assert.Equal(t, 1, len(GlobalDrawStack.as))
	SetDrawStack(
		NewHeap(false),
		NewDrawFPS(),
	)
	assert.Equal(t, 2, len(GlobalDrawStack.as))
	GlobalDrawStack.Pop()
	PreDraw()
	assert.Equal(t, 1, len(GlobalDrawStack.as))
}
