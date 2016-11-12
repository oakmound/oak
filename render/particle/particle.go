// Package particle provides options for generating renderable
// particle sources.
package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"image"
	"image/draw"
	"math"
	"time"
)

type Generator interface {
	GetBaseGenerator() *BaseGenerator
	GenerateParticle(BaseParticle) Particle
	Generate(int) *Source
	GetParticleSize() (float64, float64, bool)
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
}

// A Source is used to store and control a set of particles.
type Source struct {
	render.Layered
	Generator     Generator
	particles     []Particle
	rotateBinding event.Binding
	clearBinding  event.Binding
	cID           event.CID
}

type Particle interface {
	GetBaseParticle() *BaseParticle
	Draw(Generator, draw.Image)
	GetPos() (float64, float64)
	GetSize() (float64, float64)
}

type BaseParticle struct {
	x, y       float64
	velX, velY float64
	life       float64
	totalLife  float64
}

func (ps *Source) Init() event.CID {
	cID := event.NextID(ps)
	binding, _ := cID.Bind(rotateParticles, "EnterFrame")
	ps.rotateBinding = binding
	ps.cID = cID
	if ps.Generator.GetBaseGenerator().Duration != -1 {
		go func(ps_p *Source, duration int) {
			select {
			case <-time.After(time.Duration(duration) * time.Millisecond):
				if ps_p.GetLayer() != -1 {
					ps_p.Stop()
				}
			}
		}(ps, ps.Generator.GetBaseGenerator().Duration)
	}
	return event.NextID(ps)
}

func (ps *Source) Draw(buff draw.Image) {
	for _, p := range ps.particles {
		p.Draw(ps.Generator, buff)
	}
}

// rotateParticles updates particles over time as long
// as a Source is active.
func rotateParticles(id int, nothing interface{}) int {
	ps := event.GetEntity(id).(*Source)
	pg := ps.Generator.GetBaseGenerator()

	newParticles := make([]Particle, 0)

	for _, p := range ps.particles {

		bp := p.GetBaseParticle()

		// Ignore dead particles
		if bp.life > 0 {

			// Move towards doom
			bp.life--

			// Apply rotational acceleration
			if pg.Rotation != 0 || pg.RotationRand != 0 {
				magnitude := math.Abs(bp.velX) + math.Abs(bp.velY)
				angle := math.Atan2(bp.velX, bp.velY)
				angle += pg.Rotation + floatFromSpread(pg.RotationRand)
				bp.velX = math.Sin(angle)
				bp.velY = math.Cos(angle)
				magnitude = magnitude / (math.Abs(bp.velX) + math.Abs(bp.velY))
				bp.velX = bp.velX * magnitude
				bp.velY = bp.velY * magnitude
			}

			if pg.SpeedDecayX != 0 {
				bp.velX *= pg.SpeedDecayX
				if math.Abs(bp.velX) < 0.2 {
					bp.velX = 0.0
				}
			}
			if pg.SpeedDecayY != 0 {
				bp.velY *= pg.SpeedDecayY
				if math.Abs(bp.velY) < 0.2 {
					bp.velY = 0.0
				}
			}

			// Be dragged down by the weight of the soul
			bp.velX += pg.GravityX
			bp.velY += pg.GravityY

			bp.x += bp.velX
			bp.y += bp.velY
			render.SetDirty(bp.x, bp.y)

			newParticles = append(newParticles, p)
		} else {
			if pg.EndFunc != nil {
				pg.EndFunc(p)
			}
		}
	}

	// Regularly create particles (up until max particles)
	newParticleRand := roundFloat(floatFromSpread(pg.NewPerFrameRand))
	newParticleCount := int(pg.NewPerFrame) + newParticleRand
	for i := 0; i < newParticleCount; i++ {

		angle := (pg.Angle + floatFromSpread(pg.AngleRand)) * math.Pi / 180.0
		speed := pg.Speed + floatFromSpread(pg.SpeedRand)
		startLife := pg.LifeSpan + floatFromSpread(pg.LifeSpanRand)

		bp := BaseParticle{
			x:         pg.X + floatFromSpread(pg.SpreadX),
			y:         pg.Y + floatFromSpread(pg.SpreadY),
			velX:      speed * math.Cos(angle) * -1,
			velY:      speed * math.Sin(angle) * -1,
			life:      startLife,
			totalLife: startLife,
		}

		p := ps.Generator.GenerateParticle(bp)

		newParticles = append(newParticles, p)
	}

	ps.particles = newParticles

	return 0
}

// clearParticles is used after a Source has been stopped
// to continue moving old particles for as long as they exist.
func clearParticles(id int, nothing interface{}) int {
	ps := event.GetEntity(id).(*Source)
	pg := ps.Generator.GetBaseGenerator()

	if len(ps.particles) > 0 {
		newParticles := make([]Particle, 0)
		for _, p := range ps.particles {

			bp := p.GetBaseParticle()

			// Ignore dead particles
			if bp.life > 0 {

				// Move towards doom
				bp.life--

				// Apply rotational acceleration
				if pg.Rotation != 0 || pg.RotationRand != 0 {
					magnitude := math.Abs(bp.velX) + math.Abs(bp.velY)
					angle := math.Atan2(bp.velX, bp.velY)
					angle += pg.Rotation + floatFromSpread(pg.RotationRand)
					bp.velX = math.Sin(angle)
					bp.velY = math.Cos(angle)
					magnitude = magnitude / (math.Abs(bp.velX) + math.Abs(bp.velY))
					bp.velX = bp.velX * magnitude
					bp.velY = bp.velY * magnitude
				}

				if pg.SpeedDecayX != 0 {
					bp.velX *= pg.SpeedDecayX
					if math.Abs(bp.velX) < 0.2 {
						bp.velX = 0.0
					}
				}
				if pg.SpeedDecayY != 0 {
					bp.velY *= pg.SpeedDecayY
					if math.Abs(bp.velY) < 0.2 {
						bp.velY = 0.0
					}
				}

				// Be dragged down by the weight of the soul
				bp.velX += pg.GravityX
				bp.velY += pg.GravityY

				bp.x += bp.velX
				bp.y += bp.velY
				render.SetDirty(bp.x, bp.y)

				newParticles = append(newParticles, p)
			} else {
				if pg.EndFunc != nil {
					pg.EndFunc(p)
				}
			}
		}
		ps.particles = newParticles
	} else {
		ps.UnDraw()
		ps.rotateBinding.Unbind()
		event.DestroyEntity(id)
	}
	return 0
}

func clearAtExit(id int, nothing interface{}) int {
	if event.HasEntity(id) {
		ps := event.GetEntity(id).(*Source)
		ps.clearBinding.Unbind()
		ps.rotateBinding.Unbind()
		ps.rotateBinding, _ = ps.cID.Bind(clearParticles, "EnterFrame")
	}
	return 0
}

// Stop manually stops a Source, if its duration is infinite
// or if it should be stopped before expriring naturally.
func (ps *Source) Stop() {
	ps.clearBinding, _ = ps.cID.Bind(clearAtExit, "ExitFrame")
}

// A particle source has no concept of an individual
// rgba buffer, and so it returns nothing when its
// rgba buffer is queried. This may change.
func (ps *Source) GetRGBA() *image.RGBA {
	return nil
}

func (ps *Source) ShiftX(x float64) {
	ps.Generator.GetBaseGenerator().X += x
}

func (ps *Source) ShiftY(y float64) {
	ps.Generator.GetBaseGenerator().Y += y
}

func (ps *Source) GetX() float64 {
	return ps.Generator.GetBaseGenerator().X
}

func (ps *Source) GetY() float64 {
	return ps.Generator.GetBaseGenerator().Y
}
func (ps *Source) SetPos(x, y float64) {
	ps.Generator.GetBaseGenerator().X = x
	ps.Generator.GetBaseGenerator().Y = y
}

// Pausing a Source just stops the repetition
// of its rotation function, which moves, destroys,
// ages and generates particles. Existing particles will
// stay in place.
func (ps *Source) Pause() {
	ps.rotateBinding.Unbind()
}

// Unpausing a Source rebinds it's rotate function.
func (ps *Source) UnPause() {
	binding, _ := ps.cID.Bind(rotateParticles, "EnterFrame")
	ps.rotateBinding = binding
}

// Placeholder
func (ps *Source) String() string {
	return "ParticleSource"
}
