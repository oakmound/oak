package particle

import (
	"bitbucket.org/oakmoundstudio/oak/collision"
	// "bitbucket.org/oakmoundstudio/oak/dlog"
	"image/draw"
)

var (
	generated, destroyed int
)

type CollisionGenerator struct {
	Gen    Generator
	HitMap map[int]collision.OnHit
}

func (cg *CollisionGenerator) Generate(layer int) *Source {
	ps := cg.Gen.Generate(layer)
	ps.Generator = cg
	return ps
}

func (cg *CollisionGenerator) GenerateParticle(bp BaseParticle) Particle {
	generated++
	p := cg.Gen.GenerateParticle(bp)

	w, h, dynamic := cg.Gen.GetParticleSize()
	if dynamic {
		w, h = p.GetSize()
	}
	x, y := p.GetPos()
	return &CollisionParticle{
		p,
		collision.NewReactiveSpace(collision.NewUnassignedSpace(x, y, w, h), cg.HitMap),
	}
}

func (cg *CollisionGenerator) GetBaseGenerator() *BaseGenerator {
	return cg.Gen.GetBaseGenerator()
}

func (cg *CollisionGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

type CollisionParticle struct {
	p Particle
	s *collision.ReactiveSpace
}

func (cp *CollisionParticle) DrawOffset(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*CollisionGenerator)
	x, y := cp.p.GetPos()
	cp.s.Space().Location = collision.NewRect(x, y, cp.s.Space().GetW(), cp.s.Space().GetH())
	cp.p.DrawOffset(gen.Gen, buff, xOff, yOff)
	hitFlag := <-cp.s.CallOnHits()
	if hitFlag {
		cp.p.GetBaseParticle().life = 0
	}
}

func (cp *CollisionParticle) Draw(generator Generator, buff draw.Image) {
	gen := generator.(*CollisionGenerator)
	x, y := cp.p.GetPos()
	cp.s.Space().Location = collision.NewRect(x, y, cp.s.Space().GetW(), cp.s.Space().GetH())
	cp.p.Draw(gen.Gen, buff)
	hitFlag := <-cp.s.CallOnHits()
	if hitFlag {
		cp.p.GetBaseParticle().life = 0
	}
}

func (cp *CollisionParticle) GetBaseParticle() *BaseParticle {
	return cp.p.GetBaseParticle()
}

func (cp *CollisionParticle) GetPos() (float64, float64) {
	return cp.p.GetPos()
}
func (cp *CollisionParticle) GetSize() (float64, float64) {
	return cp.s.Space().GetW(), cp.s.Space().GetH()
}
