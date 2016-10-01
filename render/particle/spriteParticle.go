package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"image/draw"
	"math"
)

type SpriteGenerator struct {
	BaseGenerator
	SpriteRotation, SpriteRotationRand float64
	Base                               *render.Sprite
}

func (sg *SpriteGenerator) Generate(layer int) *Source {

	// Convert rotation from degrees to radians
	sg.Rotation = sg.Rotation / 180 * math.Pi
	sg.RotationRand = sg.Rotation / 180 * math.Pi

	// Make a source
	ps := Source{
		Generator: sg,
		particles: make([]Particle, 0),
	}

	// Bind things to that source:
	ps.Init()
	render.Draw(&ps, layer)

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

func (sp *SpriteParticle) Draw(generator Generator, buff draw.Image) {
	sp.rotation += sp.rotation
	rgba := generator.(*SpriteGenerator).Base.Copy().Rotate(int(sp.rotation)).GetRGBA()
	render.ShinyDraw(buff, rgba, int(sp.x), int(sp.y))
}

func (sp *SpriteParticle) GetBaseParticle() *BaseParticle {
	return &sp.BaseParticle
}

func (sg *SpriteGenerator) GetParticleSize() (float64, float64, bool) {

	bounds := sg.Base.GetRGBA().Rect.Max

	return float64(bounds.X), float64(bounds.Y), false
}

func (sp *SpriteParticle) GetPos() (float64, float64) {
	return sp.x, sp.y
}

func (sp *SpriteParticle) GetSize() (float64, float64) {

	return 0, 0
}
