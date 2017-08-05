// Package particle provides options for generating renderable
// particle sources.
package particle

import (
	"image/draw"

	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
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
	render.LayeredPoint
	Src       *Source
	Vel       physics.Vector
	Life      float64
	totalLife float64
	pID       int
}

func (bp *baseParticle) GetLayer() int {
	if bp == nil {
		return render.Undraw
	}
	return bp.LayeredPoint.GetLayer()
}

func (bp *baseParticle) GetBaseParticle() *baseParticle {
	return bp
}

func (bp *baseParticle) GetPos() physics.Vector {
	return bp.Vector
}

func (bp *baseParticle) GetDims() (int, int) {
	return 0, 0
}

func (bp *baseParticle) Cycle(gen Generator) {}

func (bp *baseParticle) String() string {
	return "BaseParticle"
}
