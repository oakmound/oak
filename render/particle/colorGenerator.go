package particle

import (
	"image/color"
	"math"

	"github.com/200sc/go-dist/intrange"
)

// A ColorGenerator generates ColorParticles
type ColorGenerator struct {
	BaseGenerator
	StartColor, StartColorRand color.Color
	EndColor, EndColorRand     color.Color
	// The size, in pixel radius, of spawned particles
	Size intrange.Range
	//
	// Some sort of particle type, for rendering triangles or squares or circles...
	Shape ShapeFunction
}

// NewColorGenerator returns a new color generator
func NewColorGenerator(options ...func(Generator)) Generator {
	g := new(ColorGenerator)
	g.setDefaults()

	for _, opt := range options {
		opt(g)
	}

	return g
}

func (cg *ColorGenerator) setDefaults() {
	cg.BaseGenerator.setDefaults()
	cg.StartColor = color.RGBA{0, 0, 0, 0}
	cg.StartColorRand = color.RGBA{0, 0, 0, 0}
	cg.EndColor = color.RGBA{0, 0, 0, 0}
	cg.EndColorRand = color.RGBA{0, 0, 0, 0}
	cg.Size = intrange.Constant(1)
	cg.Shape = Square
}

// Generate creates a source using this generator
func (cg *ColorGenerator) Generate(layer int) *Source {
	// Convert rotation from degrees to radians
	if cg.Rotation != nil {
		cg.Rotation = cg.Rotation.Mult(math.Pi / 180)
	}
	return NewSource(cg, layer)
}

// GenerateParticle creates a particle from a generator
func (cg *ColorGenerator) GenerateParticle(bp *baseParticle) Particle {
	return &ColorParticle{
		baseParticle: bp,
		startColor:   randColor(cg.StartColor, cg.StartColorRand),
		endColor:     randColor(cg.EndColor, cg.EndColorRand),
		size:         cg.Size.Poll(),
	}
}

// GetParticleSize on a color generator returns that the particles
// are per-particle specificially sized
func (cg *ColorGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

// Coloration
//

// SetStartColor lets cg have its color be set
func (cg *ColorGenerator) SetStartColor(sc, scr color.Color) {
	cg.StartColor = sc
	cg.StartColorRand = scr
}

// SetEndColor lets cg have its end color be set
func (cg *ColorGenerator) SetEndColor(ec, ecr color.Color) {
	cg.EndColor = ec
	cg.EndColorRand = ecr
}

//
// Sizing
//

// A Sizeable is a generator that can have some size set to it
type Sizeable interface {
	SetSize(i intrange.Range)
}

// Size is an option to set a Sizeable size
func Size(i intrange.Range) func(Generator) {
	return func(g Generator) {
		g.(Sizeable).SetSize(i)
	}
}

// SetSize satisfies Sizeable
func (cg *ColorGenerator) SetSize(i intrange.Range) {
	cg.Size = i
}

//
// Shaping
//

// A ShapeFunction takes in a total size of a shape and coordinates within the
// shape, and reports whether that coordinate pair lies in the shape.
type ShapeFunction func(x, y, size int) bool

// ShapeFunction block
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

// SetShape satisfies Shapeable
func (cg *ColorGenerator) SetShape(sf ShapeFunction) {
	cg.Shape = sf
}
