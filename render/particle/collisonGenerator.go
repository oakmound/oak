package particle

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
)

// A CollisionGenerator generates collision particles
type CollisionGenerator struct {
	Generator
	Fragile bool
	HitMap  map[collision.Label]collision.OnHit
}

// NewCollisionGenerator creates a new collision generator
func NewCollisionGenerator(g Generator, options ...func(*CollisionGenerator)) Generator {
	g2 := new(CollisionGenerator)
	g2.setDefault()

	g2.Generator = g

	for _, opt := range options {
		opt(g2)
	}

	return g2
}

func (cg *CollisionGenerator) setDefault() {
	cg.HitMap = make(map[collision.Label]collision.OnHit)
}

// Generate creates a source using this generator
func (cg *CollisionGenerator) Generate(layer int) *Source {
	ps := cg.Generator.Generate(layer)
	ps.Generator = cg
	return ps
}

// GenerateParticle creates a particle from a generator
func (cg *CollisionGenerator) GenerateParticle(bp *baseParticle) Particle {
	p := cg.Generator.GenerateParticle(bp)

	w, h, dynamic := cg.Generator.GetParticleSize()
	if dynamic {
		iw, ih := p.GetDims()
		w, h = float64(iw), float64(ih)
	}
	pos := p.GetPos()
	return &CollisionParticle{
		p,
		collision.NewReactiveSpace(collision.NewFullSpace(pos.X(), pos.Y(), w, h, 0, event.CID(bp.pID)), cg.HitMap),
	}
}

// GetParticleSize on a CollisionGenerator tells the caller that the particle size
// is per-particle specific
func (cg *CollisionGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

// Fragile sets whether the particles from this collisionGenerator are destroyed
// on contact
func Fragile(f bool) func(*CollisionGenerator) {
	return func(cg *CollisionGenerator) {
		cg.Fragile = f
	}
}

// HitMap sets functions to be called when particles from this generator collide
// with other spaces
func HitMap(hm map[collision.Label]collision.OnHit) func(*CollisionGenerator) {
	return func(cg *CollisionGenerator) {
		cg.HitMap = hm
	}
}
