package particle

import (
	"goevo/alg"
	"image/color"
	"math"
)

type ColorGenerator struct {
	BaseGenerator
	StartColor, StartColorRand color.Color
	EndColor, EndColorRand     color.Color
	// The size, in pixel radius, of spawned particles
	Size alg.IntRange
	//
	// Some sort of particle type, for rendering triangles or squares or circles...
	Shape ShapeFunction
}

func NewColorGenerator(options ...func(Generator)) Generator {
	g := new(ColorGenerator)
	g.SetDefaults()

	for _, opt := range options {
		opt(g)
	}

	return g
}

func (cg *ColorGenerator) SetDefaults() {
	cg.BaseGenerator.SetDefaults()
	cg.StartColor = color.RGBA{0, 0, 0, 0}
	cg.StartColorRand = color.RGBA{0, 0, 0, 0}
	cg.EndColor = color.RGBA{0, 0, 0, 0}
	cg.EndColorRand = color.RGBA{0, 0, 0, 0}
	cg.Size = alg.Constant(1)
	cg.Shape = Square
}

func (cg *ColorGenerator) Generate(layer int) *Source {
	// Convert rotation from degrees to radians
	if cg.Rotation != nil {
		cg.Rotation = cg.Rotation.Mult(math.Pi / 180)
	}
	return NewSource(cg)
}

func (cg *ColorGenerator) GenerateParticle(bp BaseParticle) Particle {
	return &ColorParticle{
		BaseParticle: bp,
		startColor:   randColor(cg.StartColor, cg.StartColorRand),
		endColor:     randColor(cg.EndColor, cg.EndColorRand),
		size:         cg.Size.Poll(),
	}
}

func (cg *ColorGenerator) GetBaseGenerator() *BaseGenerator {
	return &cg.BaseGenerator
}

func (cp *ColorGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

// Coloration
//

type Colorable interface {
	SetStartColor(color.Color, color.Color)
	SetEndColor(color.Color, color.Color)
}

func Color(sc, scr, ec, ecr color.Color) func(Generator) {
	return func(g Generator) {
		c := g.(Colorable)
		c.SetStartColor(sc, scr)
		c.SetEndColor(ec, ecr)
	}
}

func (cg *ColorGenerator) SetStartColor(sc, scr color.Color) {
	cg.StartColor = sc
	cg.StartColorRand = scr
}

func (cg *ColorGenerator) SetEndColor(ec, ecr color.Color) {
	cg.EndColor = ec
	cg.EndColorRand = ecr
}

//
// Sizing
//

type Sizeable interface {
	SetSize(i alg.IntRange)
}

func Size(i alg.IntRange) func(Generator) {
	return func(g Generator) {
		g.(Sizeable).SetSize(i)
	}
}

func (cg *ColorGenerator) SetSize(i alg.IntRange) {
	cg.Size = i
}

//
// Shaping
//

type ShapeFunction func(x, y, size int) bool

var (
	Square = func(x, y, size int) bool {
		return true
	}
	Diamond = func(x, y, size int) bool {
		radius := size / 2
		return math.Abs(float64(x-radius))+math.Abs(float64(y-radius)) < float64(radius)
	}
	Circle = func(x, y, size int) bool {
		radius := size / 2
		dx := math.Abs(float64(x - radius))
		dy := math.Abs(float64(y - radius))
		radiusf64 := float64(radius)
		if dx+dy <= radiusf64 {
			return true
		}
		return math.Pow(dx, 2)+math.Pow(dy, 2) < math.Pow(radiusf64, 2)
	}
)

// Shapeable generators can have the Shape option called on them
type Shapeable interface {
	SetShape(ShapeFunction)
}

// Shape is an option to set a generator's shape
func Shape(sf ShapeFunction) func(Generator) {
	return func(g Generator) {
		g.(Shapeable).SetShape(sf)
	}
}

func (cg *ColorGenerator) SetShape(sf ShapeFunction) {
	cg.Shape = sf
}
