package particle

import (
	"image"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	// "bitbucket.org/oakmoundstudio/oak/dlog"
	"image/draw"
)

var (
	generated, destroyed int
)

type CollisionGenerator struct {
	Gen     Generator
	Fragile bool
	HitMap  map[int]collision.OnHit
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

type CollisionParticle struct {
	P Particle
	s *collision.ReactiveSpace
}

func (cp *CollisionParticle) Draw(buff draw.Image) {
	cp.DrawOffset(buff, 0, 0)
}
func (cp *CollisionParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	cp.DrawOffsetGen(cp.P.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}
func (cp *CollisionParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*CollisionGenerator)
	pos := cp.P.GetPos()
	cp.s.Space().Location = collision.NewRect(pos.X, pos.Y, cp.s.Space().GetW(), cp.s.Space().GetH())
	cp.P.DrawOffsetGen(gen.Gen, buff, xOff, yOff)
	hitFlag := <-cp.s.CallOnHits()
	if gen.Fragile && hitFlag {
		cp.P.GetBaseParticle().Life = 0
	}
}

func (cp *CollisionParticle) GetBaseParticle() *BaseParticle {
	return cp.P.GetBaseParticle()
}

func (cp *CollisionParticle) GetPos() *physics.Vector {
	return cp.P.GetPos()
}
func (cp *CollisionParticle) GetSize() (float64, float64) {
	return cp.s.Space().GetW(), cp.s.Space().GetH()
}
func (cp *CollisionParticle) ShiftX(x float64) {
	cp.P.ShiftX(x)
}

func (cp *CollisionParticle) ShiftY(y float64) {
	cp.P.ShiftY(y)
}

func (cp *CollisionParticle) GetX() float64 {
	return cp.P.GetX()
}

func (cp *CollisionParticle) GetY() float64 {
	return cp.P.GetY()
}
func (cp *CollisionParticle) SetPos(x, y float64) {
	cp.P.SetPos(x, y)
}
func (cp *CollisionParticle) GetLayer() int {
	return cp.P.GetLayer()
}

func (cp *CollisionParticle) SetLayer(l int) {
	cp.P.SetLayer(l)
}

func (cp *CollisionParticle) UnDraw() {
	cp.P.UnDraw()
}

func (cp *CollisionParticle) String() string {
	return "CollisionParticle"
}

func (cp *CollisionParticle) GetRGBA() *image.RGBA {
	return cp.P.GetRGBA()
}
