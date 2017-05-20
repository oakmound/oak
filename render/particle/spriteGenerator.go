package particle

import (
	"math"

	"github.com/200sc/go-dist/floatrange"

	"bitbucket.org/oakmoundstudio/oak/render"
)

type SpriteGenerator struct {
	BaseGenerator
	SpriteRotation floatrange.Range
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
	sg.SpriteRotation = floatrange.Constant(0)
}

func (sg *SpriteGenerator) Generate(layer int) *Source {
	// Convert rotation from degrees to radians
	if sg.Rotation != nil {
		sg.Rotation = sg.Rotation.Mult(math.Pi / 180)
	}
	return NewSource(sg, layer)
}

func (sg *SpriteGenerator) GenerateParticle(bp *BaseParticle) Particle {
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
	SetSpriteRotation(f floatrange.Range)
}

func Sprite(s *render.Sprite) func(Generator) {
	return func(g Generator) {
		g.(Sprited).SetSprite(s)
	}
}

func (sg *SpriteGenerator) SetSprite(s *render.Sprite) {
	sg.Base = s
}

func SpriteRotation(f floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.(Sprited).SetSpriteRotation(f)
	}
}

func (sg *SpriteGenerator) SetSpriteRotation(f floatrange.Range) {
	sg.SpriteRotation = f
}
