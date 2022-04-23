package particle

import (
	"math"

	"github.com/oakmound/oak/v4/alg/span"
	"github.com/oakmound/oak/v4/physics"
)

var (
	// Inf represents Infinite duration
	Inf = span.NewConstant(math.MaxInt32)
)

// A Generator holds settings for generating particles
type Generator interface {
	GetBaseGenerator() *BaseGenerator
	GenerateParticle(*baseParticle) Particle
	Generate(int) *Source
	GetParticleSize() (float64, float64, bool)
	ShiftX(float64)
	ShiftY(float64)
	SetPos(float64, float64)
	GetPos() (float64, float64)
}

// A BaseGenerator fulfills all basic requirements to generate
// particles
// Modeled after Parcycle
type BaseGenerator struct {
	physics.Vector
	// This float is currently forced to an integer
	// at new particle rotation. This should be changed
	// to something along the lines of 'new per 30 frames',
	// or allow low fractional values to be meaningful,
	// so that more fine-tuned particle generation speeds are possible.
	NewPerFrame span.Span[float64]
	// The number of frames each particle should persist
	// before being removed.
	LifeSpan span.Span[float64]
	// 0 - between quadrant 1 and 4
	// 90 - between quadrant 2 and 1
	Angle  span.Span[float64]
	Speed  span.Span[float64]
	Spread physics.Vector
	// Duration in milliseconds for the particle source.
	// After this many milliseconds have passed, it will
	// stop sending out new particles. Old particles will
	// not be removed until their individual lifespans run
	// out.
	// A duration of -1 represents never stopping.
	Duration span.Span[int]
	// Rotational acceleration, to change angle over time
	Rotation span.Span[float64]
	// Gravity X() and Gravity Y() represent particle acceleration per frame.
	Gravity       physics.Vector
	SpeedDecay    physics.Vector
	EndFunc       func(Particle)
	LayerFunc     func(physics.Vector) int
	ParticleLimit int
}

// GetBaseGenerator returns this
func (bg *BaseGenerator) GetBaseGenerator() *BaseGenerator {
	return bg
}

func (bg *BaseGenerator) setDefaults() {
	*bg = BaseGenerator{
		Vector:      physics.NewVector(0, 0),
		NewPerFrame: span.NewConstant(1.0),
		LifeSpan:    span.NewConstant(60.0),
		Angle:       span.NewConstant(0.0),
		Speed:       span.NewConstant(1.0),
		Spread:      physics.NewVector(0, 0),
		Duration:    Inf,
		Rotation:    nil,
		Gravity:     physics.NewVector(0, 0),
		SpeedDecay:  physics.NewVector(0, 0),
		EndFunc:     nil,
		LayerFunc:   func(physics.Vector) int { return 1 },
	}
}

// ShiftX moves a base generator by an x value
func (bg *BaseGenerator) ShiftX(x float64) {
	bg.Vector = bg.Vector.ShiftX(x)
}

// ShiftY moves a base generator by a y value
func (bg *BaseGenerator) ShiftY(y float64) {
	bg.Vector = bg.Vector.ShiftY(y)
}

// SetPos sets the position of a base generator
func (bg *BaseGenerator) SetPos(x, y float64) {
	bg.Vector = bg.Vector.SetPos(x, y)
}
