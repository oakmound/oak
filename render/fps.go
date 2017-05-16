package render

import (
	"image"
	"image/draw"
	"time"

	"bitbucket.org/oakmoundstudio/oak/timing"
)

const (
	FPSSMOOTHING = .25
)

type DrawFPS struct {
	fps      int
	lastTime time.Time
	txt      *Text
}

func NewDrawFPS() *DrawFPS {
	df := new(DrawFPS)
	df.lastTime = time.Now()
	df.fps = 0
	return df
}

func (df *DrawFPS) PreDraw() {
	//NOP
	if df.txt == nil {
		df.txt = DefFont().NewIntText(&df.fps, 10, 20)
	}
}

func (df *DrawFPS) Add(Renderable, int) Renderable {
	//NOP
	return nil
}

func (df *DrawFPS) Copy() Addable {
	return new(DrawFPS)
}

func (df *DrawFPS) draw(world draw.Image, view image.Point, w, h int) {
	df.fps = int((timing.FPS(df.lastTime, time.Now()) * FPSSMOOTHING) + (float64(df.fps) * (1 - FPSSMOOTHING)))
	df.txt.DrawOffset(world, float64(view.X), float64(view.Y))
	df.lastTime = time.Now()
}
