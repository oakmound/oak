package particle

import (
	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
)

type CollisionGenerator struct {
	Gen     Generator
	Fragile bool
	HitMap  map[int]collision.OnHit
}

func NewCollisionGenerator(g Generator, options ...func(*CollisionGenerator)) Generator {
	g2 := new(CollisionGenerator)

	g2.Gen = g

	for _, opt := range options {
		opt(g2)
	}

	return g2
}

func (cg *CollisionGenerator) SetDefault() {
	cg.HitMap = make(map[int]collision.OnHit)
}

func (cg *CollisionGenerator) Generate(layer int) *Source {
	ps := cg.Gen.Generate(layer)
	ps.Generator = cg
	return ps
}

func (cg *CollisionGenerator) GenerateParticle(bp *BaseParticle) Particle {
	p := cg.Gen.GenerateParticle(bp)

	w, h, dynamic := cg.Gen.GetParticleSize()
	if dynamic {
		w, h = p.GetSize()
	}
	pos := p.GetPos()
	return &CollisionParticle{
		p,
		collision.NewReactiveSpace(collision.NewFullSpace(pos.X, pos.Y, w, h, 0, event.CID(bp.pID)), cg.HitMap),
	}
}

func (cg *CollisionGenerator) GetBaseGenerator() *BaseGenerator {
	return cg.Gen.GetBaseGenerator()
}

func (cg *CollisionGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

func (cg *CollisionGenerator) ShiftX(x float64) {
	cg.Gen.ShiftX(x)
}

func (cg *CollisionGenerator) ShiftY(y float64) {
	cg.Gen.ShiftY(y)
}

func (cg *CollisionGenerator) SetPos(x, y float64) {
	cg.Gen.SetPos(x, y)
}

func (cg *CollisionGenerator) GetPos() (float64, float64) {
	return cg.Gen.GetPos()
}

func Fragile(f bool) func(*CollisionGenerator) {
	return func(cg *CollisionGenerator) {
		cg.Fragile = f
	}
}

func HitMap(hm map[int]collision.OnHit) func(*CollisionGenerator) {
	return func(cg *CollisionGenerator) {
		cg.HitMap = hm
	}
}
