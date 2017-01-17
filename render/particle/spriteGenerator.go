package particle

import (
	"math"

	"bitbucket.org/oakmoundstudio/oak/alg"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type SpriteGenerator struct {
	BaseGenerator
	SpriteRotation alg.FloatRange
	Base           *render.Sprite
}

func NewSpriteGenerator(options ...func(Generator)) Generator {
	g := new(SpriteGenerator)
	g.SetDefaults()

	for _, opt := range options {
		opt(g)
	}

	return g
}

func (sg *SpriteGenerator) SetDefaults() {
	sg.BaseGenerator.SetDefaults()
	sg.SpriteRotation = alg.Constantf(0)
}

func (sg *SpriteGenerator) Generate(layer int) *Source {
	// Convert rotation from degrees to radians
	if sg.Rotation != nil {
		sg.Rotation = sg.Rotation.Mult(math.Pi / 180)
	}
	return NewSource(sg)
}

func (sg *SpriteGenerator) GenerateParticle(bp BaseParticle) Particle {
	return &SpriteParticle{
		BaseParticle: bp,
		rotation:     sg.SpriteRotation.Poll(),
	}
}

func (sg *SpriteGenerator) GetBaseGenerator() *BaseGenerator {
	return &sg.BaseGenerator
}

type Sprited interface {
	SetSprite(*render.Sprite)
	SetSpriteRotation(f alg.FloatRange)
}

func Sprite(s *render.Sprite) func(Generator) {
	return func(g Generator) {
		g.(Sprited).SetSprite(s)
	}
}

func (sg *SpriteGenerator) SetSprite(s *render.Sprite) {
	sg.Base = s
}

func SpriteRotation(f alg.FloatRange) func(Generator) {
	return func(g Generator) {
		g.(Sprited).SetSpriteRotation(f)
	}
}

func (sg *SpriteGenerator) SetSpriteRotation(f alg.FloatRange) {
	sg.SpriteRotation = f
}
