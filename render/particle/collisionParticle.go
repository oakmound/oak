package particle

import (
	"image/draw"

	"github.com/oakmound/oak/collision"
)

// A CollisionParticle is a wrapper around other particles that also
// has a collision space and can functionally react with the environment
// on collision
type CollisionParticle struct {
	Particle
	s *collision.ReactiveSpace
}

// Draw redirects to DrawOffset
func (cp *CollisionParticle) Draw(buff draw.Image) {
	cp.DrawOffset(buff, 0, 0)
}

// DrawOffset redirects to DrawOffsetGen
func (cp *CollisionParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	cp.DrawOffsetGen(cp.Particle.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (cp *CollisionParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*CollisionGenerator)
	cp.Particle.DrawOffsetGen(gen.Generator, buff, xOff, yOff)
}

// Cycle updates the collisiion particles variables once per rotation
func (cp *CollisionParticle) Cycle(generator Generator) {
	gen := generator.(*CollisionGenerator)
	pos := cp.Particle.GetPos()
	cp.s.Space.Location = collision.NewRect(pos.X(), pos.Y(), cp.s.GetW(), cp.s.GetH())

	hitFlag := <-cp.s.CallOnHits()
	if gen.Fragile && hitFlag {
		cp.Particle.GetBaseParticle().Life = 0
	}
}

// GetDims returns the dimensions of the space of the particle
func (cp *CollisionParticle) GetDims() (int, int) {
	return int(cp.s.GetW()), int(cp.s.GetH())
}

// String returns the type as string
func (cp *CollisionParticle) String() string {
	return "CollisionParticle"
}
