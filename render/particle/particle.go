// Package particle provides options for generating renderable
// particle sources.
package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"
)

var (
	Square = func(x, y, size int) bool {
		return true
	}
	Diamond = func(x, y, size int) bool {
		radius := size / 2
		return math.Abs(float64(x-radius))+math.Abs(float64(y-radius)) < float64(radius)
	}
)

// ParticleGenerator represents the various options
// one needs to or may provide in order to generate a
// ParticleSource.
// Modeled after Parcycle
type ParticleGenerator struct {
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
	GravityX, GravityY         float64
	StartColor, StartColorRand color.Color
	EndColor, EndColorRand     color.Color
	// The size, in pixel radius, of spawned particles
	Size, SizeRand int
	//
	// Future potential options:
	// Some sort of particle type, for rendering triangles or squares or circles...
	Shape ShapeFunction
}

// A ParticleSource is used to store and control a set of particles.
type ParticleSource struct {
	render.Layered
	Generator     ParticleGenerator
	particles     []Particle
	rotateBinding event.Binding
	cID           event.CID
}

type ShapeFunction func(x, y, size int) bool

// A particle is a colored pixel at a given position, moving in a certain direction.
// After a while, it will dissipate.
type Particle struct {
	x, y       float64
	velX, velY float64
	startColor color.Color
	endColor   color.Color
	life       float64
	totalLife  float64
	size       int
}

func (ps *ParticleSource) Init() event.CID {
	return plastic.NextID(ps)
}

// Generate takes a generator and converts it into a source,
// drawing particles and binding functions for particle generation
// and rotation.
func (pg *ParticleGenerator) Generate(layer int) *ParticleSource {

	// Convert rotation from degrees to radians
	pg.Rotation = pg.Rotation / 180 * math.Pi
	pg.RotationRand = pg.Rotation / 180 * math.Pi

	// Make a source
	ps := ParticleSource{
		Generator: *pg,
		particles: make([]Particle, 0),
	}

	// Bind things to that source:
	cID := ps.Init()
	binding, _ := cID.Bind(rotateParticles, "EnterFrame")
	ps.rotateBinding = binding
	ps.cID = cID
	render.Draw(&ps, layer)
	if pg.Duration != -1 {
		go func(ps_p *ParticleSource, duration int) {
			select {
			case <-time.After(time.Duration(duration) * time.Millisecond):
				ps_p.Stop()
			}
		}(&ps, pg.Duration)
	}
	return &ps
}

func (ps *ParticleSource) Draw(buff draw.Image) {
	for _, p := range ps.particles {

		r, g, b, a := p.startColor.RGBA()
		r2, g2, b2, a2 := p.endColor.RGBA()
		progress := p.life / p.totalLife
		c := color.RGBA64{
			uint16OnScale(r, r2, progress),
			uint16OnScale(g, g2, progress),
			uint16OnScale(b, b2, progress),
			uint16OnScale(a, a2, progress),
		}

		img := image.NewRGBA64(image.Rect(0, 0, p.size, p.size))

		for i := 0; i < p.size; i++ {
			for j := 0; j < p.size; j++ {
				if ps.Generator.Shape(i, j, p.size) {
					img.SetRGBA64(i, j, c)
				}
			}
		}

		halfSize := float64(p.size / 2)

		render.ShinyDraw(buff, img, int(p.x-halfSize), int(p.y-halfSize))
	}
}

// rotateParticles updates particles over time as long
// as a ParticleSource is active.
func rotateParticles(id int, nothing interface{}) error {

	ps := plastic.GetEntity(id).(*ParticleSource)
	pg := ps.Generator

	newParticles := make([]Particle, 0)

	for _, p := range ps.particles {

		// Ignore dead particles
		if p.life > 0 {

			// Move towards doom
			p.life--

			// Be dragged down by the weight of the soul
			p.velX += pg.GravityX
			p.velY += pg.GravityY

			// Apply rotational acceleration
			if pg.Rotation != 0 && pg.RotationRand != 0 {
				magnitude := math.Abs(p.velX) + math.Abs(p.velY)
				angle := math.Atan2(p.velX, p.velY)
				angle += pg.Rotation + floatFromSpread(pg.RotationRand)
				p.velX = math.Sin(angle)
				p.velY = math.Cos(angle)
				magnitude = magnitude / (math.Abs(p.velX) + math.Abs(p.velY))
				p.velX = p.velX * magnitude
				p.velY = p.velY * magnitude
			}

			p.x += p.velX
			p.y += p.velY

			newParticles = append(newParticles, p)
		}
	}

	// Regularly create particles (up until max particles)
	newParticleRand := roundFloat(floatFromSpread(pg.NewPerFrameRand))
	newParticleCount := int(pg.NewPerFrame) + newParticleRand
	for i := 0; i < newParticleCount; i++ {

		angle := (pg.Angle + floatFromSpread(pg.AngleRand)) * math.Pi / 180.0
		speed := pg.Speed + floatFromSpread(pg.SpeedRand)
		startLife := pg.LifeSpan + floatFromSpread(pg.LifeSpanRand)

		newParticles = append(newParticles, Particle{
			x:          pg.X + floatFromSpread(pg.SpreadX),
			y:          pg.Y + floatFromSpread(pg.SpreadY),
			velX:       speed * math.Cos(angle) * -1,
			velY:       speed * math.Sin(angle) * -1,
			startColor: randColor(pg.StartColor, pg.StartColorRand),
			endColor:   randColor(pg.EndColor, pg.EndColorRand),
			life:       startLife,
			totalLife:  startLife,
			size:       pg.Size + intFromSpread(pg.SizeRand),
		})
	}

	ps.particles = newParticles

	return nil
}

// clearParticles is used after a ParticleSource has been stopped
// to continue moving old particles for as long as they exist.
func clearParticles(id int, nothing interface{}) error {

	ps := plastic.GetEntity(id).(*ParticleSource)
	pg := ps.Generator

	if len(ps.particles) > 0 {
		newParticles := make([]Particle, 0)
		for _, p := range ps.particles {

			// Ignore dead particles
			if p.life > 0 {

				p.life--

				p.velX += pg.GravityX
				p.velY += pg.GravityY

				// Apply rotational acceleration
				if pg.Rotation != 0 && pg.RotationRand != 0 {
					magnitude := math.Abs(p.velX) + math.Abs(p.velY)
					angle := math.Atan2(p.velX, p.velY)
					angle += pg.Rotation + floatFromSpread(pg.RotationRand)
					p.velX = math.Sin(angle)
					p.velY = math.Cos(angle)
					magnitude = magnitude / (math.Abs(p.velX) + math.Abs(p.velY))
					p.velX = p.velX * magnitude
					p.velY = p.velY * magnitude
				}

				p.x += p.velX
				p.y += p.velY

				newParticles = append(newParticles, p)
			}
		}
		ps.particles = newParticles
	} else {
		ps.UnDraw()
		ps.rotateBinding.Unbind()
	}
	return nil
}

// Stop manually stops a ParticleSource, if its duration is infinite
// or if it should be stopped before expriring naturally.
func (ps *ParticleSource) Stop() {
	ps.rotateBinding.Unbind()
	ps.rotateBinding, _ = ps.cID.Bind(clearParticles, "EnterFrame")
}

// A particle source has no concept of an individual
// rgba buffer, and so it returns nothing when its
// rgba buffer is queried. This may change.
func (ps *ParticleSource) GetRGBA() *image.RGBA {
	return nil
}

func (ps *ParticleSource) ShiftX(x float64) {
	ps.Generator.X += x
}

func (ps *ParticleSource) ShiftY(y float64) {
	ps.Generator.Y += y
}

func (ps *ParticleSource) SetPos(x, y float64) {
	ps.Generator.X = x
	ps.Generator.Y = y
}

// Pausing a ParticleSource just stops the repetition
// of its rotation function, which moves, destroys,
// ages and generates particles. Existing particles will
// stay in place.
func (ps *ParticleSource) Pause() {
	ps.rotateBinding.Unbind()
}

// Unpausing a ParticleSource rebinds it's rotate function.
func (ps *ParticleSource) UnPause() {
	binding, _ := ps.cID.Bind(rotateParticles, "EnterFrame")
	ps.rotateBinding = binding
}
