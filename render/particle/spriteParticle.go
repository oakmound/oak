package particle

import (
	"image/draw"
	"math"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type SpriteGenerator struct {
	BaseGenerator
	SpriteRotation, SpriteRotationRand float64
	Base                               *render.Sprite
}

func (sg *SpriteGenerator) Generate(layer int) *Source {

	// Convert rotation from degrees to radians
	sg.Rotation = sg.Rotation / 180 * math.Pi
	sg.RotationRand = sg.RotationRand / 180 * math.Pi

	// Make a source
	ps := Source{
		Generator: sg,
	}

	// Bind things to that source:
	ps.Init()

	return &ps
}

func (sg *SpriteGenerator) GenerateParticle(bp BaseParticle) Particle {
	return &SpriteParticle{
		BaseParticle: bp,
		rotation:     sg.SpriteRotation + floatFromSpread(sg.SpriteRotationRand),
	}
}

func (sg *SpriteGenerator) GetBaseGenerator() *BaseGenerator {
	return &sg.BaseGenerator
}

type SpriteParticle struct {
	BaseParticle
	rotation float64
}

func (sp *SpriteParticle) Draw(buff draw.Image) {
	sp.DrawOffset(buff, 0, 0)
}

func (sp *SpriteParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	sp.DrawOffsetGen(sp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

func (sp *SpriteParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {

	sp.rotation += sp.rotation
	gen := generator.(*SpriteGenerator)
	rgba := gen.Base.Copy().Rotate(int(sp.rotation)).GetRGBA()
	render.ShinyDraw(buff, rgba, int(sp.Pos.X+xOff), int(sp.Pos.Y+yOff))
}

func (sp *SpriteParticle) GetBaseParticle() *BaseParticle {
	return &sp.BaseParticle
}

func (sg *SpriteGenerator) GetParticleSize() (float64, float64, bool) {

	bounds := sg.Base.GetRGBA().Rect.Max

	return float64(bounds.X), float64(bounds.Y), false
}

func (sp *SpriteParticle) GetPos() *physics.Vector {
	return sp.Pos
}

func (sp *SpriteParticle) GetSize() (float64, float64) {

	return 0, 0
}
