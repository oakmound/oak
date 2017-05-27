package render

import (
	"image"
	"image/draw"
	"time"

	"bitbucket.org/oakmoundstudio/oak/timing"
)

const (
	// FpsSmoothing is how much of the fps of the upcoming frame is used
	// relative to the previous fps total when calculating the fps at a given
	// frame
	FpsSmoothing = .25
)

// DrawFPS is a draw stack element that will draw the fps onto the screen
type DrawFPS struct {
	fps      int
	lastTime time.Time
	txt      *Text
}

// NewDrawFPS returns a zero-initialized DrawFPS
func NewDrawFPS() *DrawFPS {
	df := new(DrawFPS)
	df.lastTime = time.Now()
	df.fps = 0
	return df
}

// PreDraw does nothing for a drawFPS
func (df *DrawFPS) PreDraw() {
	//NOP
	// This is not done in NewDrawFPS because at the time when
	// NewDrawFPS is called, DefFont() does not exist.
	if df.txt == nil {
		df.txt = DefFont().NewIntText(&df.fps, 10, 20)
	}
}

// Add does nothing for a drawFPS
func (df *DrawFPS) Add(Renderable, int) Renderable {
	//NOP
	return nil
}

// Replace does nothing for a drawFPS
func (df *DrawFPS) Replace(Renderable, Renderable, int) {
	//NOP
}

// Copy does effectively nothing for a drawFPS
func (df *DrawFPS) Copy() Addable {
	return new(DrawFPS)
}

func (df *DrawFPS) draw(world draw.Image, view image.Point, w, h int) {
	df.fps = int((timing.FPS(df.lastTime, time.Now()) * FpsSmoothing) + (float64(df.fps) * (1 - FpsSmoothing)))
	df.txt.Draw(world)
	df.lastTime = time.Now()
}
