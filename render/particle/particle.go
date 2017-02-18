// Package particle provides options for generating renderable
// particle sources.
package particle

import (
	"image"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Particle interface {
	render.Renderable
	GetBaseParticle() *BaseParticle
	GetPos() *physics.Vector
	GetSize() (float64, float64)
	DrawOffsetGen(gen Generator, buff draw.Image, xOff, yOff float64)
	Cycle(gen Generator)
}

type BaseParticle struct {
	render.Layered
	Src       *Source
	Pos       *physics.Vector
	Vel       *physics.Vector
	Life      float64
	totalLife float64
	pID       int
}

// A particle has no concept of an individual
// rgba buffer, and so it returns nothing when its
// rgba buffer is queried. This may change.
func (bp *BaseParticle) GetRGBA() *image.RGBA {
	return nil
}

func (bp *BaseParticle) ShiftX(x float64) {
	bp.Pos.X += x
}

func (bp *BaseParticle) ShiftY(y float64) {
	bp.Pos.Y += y
}

func (bp *BaseParticle) GetX() float64 {
	return bp.Pos.X
}

func (bp *BaseParticle) GetY() float64 {
	return bp.Pos.Y
}
func (bp *BaseParticle) SetPos(x, y float64) {
	bp.Pos.X = x
	bp.Pos.Y = y
}

func (bp *BaseParticle) Cycle(gen Generator){}

func (bp *BaseParticle) String() string {
	return "BaseParticle"
}
