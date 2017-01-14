package particle

import "bitbucket.org/oakmoundstudio/oak/physics"

type Generator interface {
	GetBaseGenerator() *BaseGenerator
	GenerateParticle(BaseParticle) Particle
	Generate(int) *Source
	GetParticleSize() (float64, float64, bool)
	ShiftX(float64)
	ShiftY(float64)
	SetPos(float64, float64)
	GetPos() (float64, float64)
}

// Represents the various options
// one needs to or may provide in order to generate a
// Source.
// Modeled after Parcycle
type BaseGenerator struct {
	// This float is currently forced to an integer
	// at new particle rotation. This should be changed
	// to something along the lines of 'new per 30 frames',
	// or allow low fractional values to be meaningful,
	// so that more fine-tuned particle generation speeds are possible.
	NewPerFrame, NewPerFrameRand float64
	X, Y                         float64
	// The number of frames each particle should persist
	// before being removed.
	LifeSpan, LifeSpanRand float64
	// 0 - between quadrant 1 and 4
	// 90 - between quadrant 2 and 1
	Angle, AngleRand float64
	Speed, SpeedRand float64
	SpreadX, SpreadY float64
	// Duration in milliseconds for the particle source.
	// After this many milliseconds have passed, it will
	// stop sending out new particles. Old particles will
	// not be removed until their individual lifespans run
	// out.
	// A duration of -1 represents never stopping.
	Duration int
	// Rotational acceleration, to change angle over time
	Rotation, RotationRand float64
	// Gravity X and Gravity Y represent particle acceleration per frame.
	GravityX, GravityY       float64
	SpeedDecayX, SpeedDecayY float64
	EndFunc                  func(Particle)
	LayerFunc                func(*physics.Vector) int
}

func (bg *BaseGenerator) ShiftX(x float64) {
	bg.X += x
}

func (bg *BaseGenerator) ShiftY(y float64) {
	bg.Y += y
}

func (bg *BaseGenerator) SetPos(x, y float64) {
	bg.X = x
	bg.Y = y
}

func (bg *BaseGenerator) GetPos() (float64, float64) {
	return bg.X, bg.Y
}
