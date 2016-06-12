package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"image"
)

// Modeled after Parcycle
type ParticleGenerator struct {
	MaxParticles           int
	X, Y                   int
	Size, SizeRand         int
	LifeSpan, LifeSpanRand int
	// 0 - between quadrant 1 and 4
	// 90 - between quadrant 2 and 1
	Angle, AngleRand           int
	Speed, SpeedRand           float64
	Spread                     int
	Duration                   int
	GravityX, GravityY         float64
	OffsetX, OffsetY           int
	StartColor, StartColorRand image.Color
	EndColor, EndColorRand     image.Color
}

type ParticleSource struct {
	Generator ParticleGenerator
	particles []Particle
}

// A particle is a colored pixel at a given position, moving in a certain direction.
// After a while, it will dissipate.
type Particle struct {
	x, y       int
	velX, velY float64
	color      image.Color
	life       int
}

func (ps *ParticleSource) Init() event.CID {
	return plastic.NextID(ps)
}

func (pg *ParticleGenerator) Generate() (*ParticleSource, event.Binding) {
	// Make a source
	ps := ParticleSource{&pg, make([]Particle)}

	// Bind things to that source:
	cID := ps.Init()
	binding, _ := cID.Bind(rotateParticles, "EnterFrame")

	return &ps, binding
}

func rotateParticles(id int, nothing interface{}) error {
	// Regularly create particles (up until max particles)
	// Regularly destroy old particles
	// Regularly modify particles' colors
	// Regularly apply gravity to particles
	ps_p := plastic.GetEntity(id)

}

func Stop(cID int) {

	cID.Unbind(rotateParticles, "EnterFrame")
	// Unbind things
	// Delete the source
}
