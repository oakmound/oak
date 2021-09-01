package render

import (
	"image/draw"
	"time"

	"github.com/oakmound/oak/v3/timing"
)

const (
	// FPS smoothing will include a portion of the last fps count in new fps counts,
	// turning what could be a spiky display to a more smooth gradient. This can hide
	// issues-- if the fps being reported is spiking between 0 and 1000, fps smoothing
	// can make it look like a fps of 100.
	defaultFpsSmoothing = .25
)

// DrawFPS is a Renderable that will display how fast it is rendered. If it is a part of
// a dynamically ordered stackable, like a Heap, how fast it will be rendered can change
// in each iteration of the stackable's draw. For this reason, its recommended to isolate
// a DrawFPS to its own stack layer or layer within a heap.
type DrawFPS struct {
	*Text
	fps       int
	lastTime  time.Time
	Smoothing float64
}

// NewDrawFPS returns a DrawFPS, which will render a counter of how fast it is being drawn.
// If font is not provided, DefaultFont is used. If smoothing is 0, a reasonable
// default will be used.
func NewDrawFPS(smoothing float64, font *Font, x, y float64) *DrawFPS {
	if smoothing == 0.0 {
		smoothing = defaultFpsSmoothing
	}
	if font == nil {
		font = DefaultFont().Copy()
	}
	df := &DrawFPS{
		Smoothing: smoothing,
		lastTime:  time.Now(),
	}
	df.Text = font.NewIntText(&df.fps, x, y)

	return df
}

// Draw renders a DrawFPS to a buffer.
func (df *DrawFPS) Draw(buff draw.Image, xOff, yOff float64) {
	t := time.Now()
	df.fps = int((timing.FPS(df.lastTime, t) * df.Smoothing) + (float64(df.fps) * (1 - df.Smoothing)))
	df.lastTime = t
	df.Text.Draw(buff, xOff, yOff)
}
