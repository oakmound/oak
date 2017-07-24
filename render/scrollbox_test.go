package render

import (
	"image/color"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScrollBoxFuncs(t *testing.T) {
	// Scrollbox is another type that takes in w/h args that should
	// return a potential error
	// Alternatively: negative w/h should result in the w/h being multiplied
	// by -1, and the whole renderable should shift -w, -h, but that would be
	// meaningless for types that start at 0,0, because the common use case
	// would be to set their position to some value afterword anyway.
	sb := NewScrollBox(
		// This could follow the same recommended syntax as sequence
		[]Renderable{
			NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}),
		},
		50,
		50,
		20,
		20,
	)
	rgba := sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{255, 255, 255, 255})

	// First shift check
	time.Sleep(50 * time.Millisecond)
	sb.update()
	rgba = sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{0, 0, 0, 0})

	// Reappear / cycling check
	time.Sleep(50 * 19 * time.Millisecond)
	sb.update()
	rgba = sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{255, 255, 255, 255})

	// AddRenderable
	sb.AddRenderable(NewColorBox(5, 5, color.RGBA{0, 0, 255, 255}))
	rgba = sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{0, 0, 255, 255})

	// High scroll rate / Pausing
	sb.SetScrollRate(10000, 10000)
	sb.Pause()
	sb.update()
	rgba = sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{0, 0, 255, 255})
	sb.Unpause()
	time.Sleep(100 * time.Millisecond)
	sb.update()
	rgba = sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{0, 0, 255, 255})

	// Positive scrollRate/ reappearPos
	assert.Nil(t, sb.SetReappearPos(-10, -10))
	assert.NotNil(t, sb.SetReappearPos(10, -10))
	assert.NotNil(t, sb.SetReappearPos(-10, 10))

	// Zero scroll rate
	sb.SetScrollRate(0, 0)
	time.Sleep(100 * time.Millisecond)
	sb.update()
	rgba = sb.GetRGBA()
	assert.Equal(t, rgba.At(0, 0), color.RGBA{0, 0, 255, 255})

	// Negative scrollRate / reappearPos
	sb.SetScrollRate(-100, -100)
	assert.Nil(t, sb.SetReappearPos(10, 10))
	assert.NotNil(t, sb.SetReappearPos(10, -10))
	assert.NotNil(t, sb.SetReappearPos(-10, 10))
}
