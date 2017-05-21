// Package particle provides options for generating renderable
// particle sources.
package particle

import (
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// A Particle is a renderable that is spawned by a generator, usually very fast,
// usually very small, for visual effects
type Particle interface {
	render.Renderable
	GetBaseParticle() *baseParticle
	GetPos() physics.Vector
	DrawOffsetGen(gen Generator, buff draw.Image, xOff, yOff float64)
	Cycle(gen Generator)
}

type baseParticle struct {
	render.Layered
	Src       *Source
	Pos       physics.Vector
	Vel       physics.Vector
	Life      float64
	totalLife float64
	pID       int
}

func (bp *baseParticle) GetBaseParticle() *baseParticle {
	return bp
}

func (bp *baseParticle) ShiftX(x float64) {
	bp.Pos.X += x
}

func (bp *baseParticle) ShiftY(y float64) {
	bp.Pos.Y += y
}

func (bp *baseParticle) GetX() float64 {
	return bp.Pos.X
}

func (bp *baseParticle) GetY() float64 {
	return bp.Pos.Y
}
func (bp *baseParticle) SetPos(x, y float64) {
	bp.Pos.X = x
	bp.Pos.Y = y
}

func (bp *baseParticle) GetPos() physics.Vector {
	return bp.Pos
}

func (bp *baseParticle) GetDims() (int, int) {
	return 0, 0
}

func (bp *baseParticle) Cycle(gen Generator) {}

func (bp *baseParticle) String() string {
	return "BaseParticle"
}
