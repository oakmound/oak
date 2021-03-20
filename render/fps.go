package render

import (
	"image/draw"
	"time"

	"github.com/oakmound/oak/v2/timing"
)

const (
	defaultFpsSmoothing = .25
)

// DrawFPS is a Renderable that will display how fast it is rendered.
type DrawFPS struct {
	*Text
	fps       int
	lastTime  time.Time
	Smoothing float64
}

// NewDrawFPS returns a DrawFPS, which will render a counter of how fast it is being drawn.
// If font is not provided, DefFont is used. If smoothing is
func NewDrawFPS(smoothing float64, font *Font, x, y float64) *DrawFPS {
	if smoothing == 0.0 {
		smoothing = defaultFpsSmoothing
	}
	if font == nil {
		font = DefFont().Copy()
	}
	df := &DrawFPS{
		Smoothing: smoothing,
		lastTime:  time.Now(),
	}
	df.Text = font.NewIntText(&df.fps, x, y)

	return df
}

func (df *DrawFPS) Draw(buff draw.Image, xOff, yOff float64) {
	t := time.Now()
	df.fps = int((timing.FPS(df.lastTime, t) * df.Smoothing) + (float64(df.fps) * (1 - df.Smoothing)))
	df.lastTime = t
	df.Text.Draw(buff, xOff, yOff)
}
