package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"
)

// Modeled after Parcycle
type ParticleGenerator struct {
	NewPerFrame, NewPerFrameRand float64
	X, Y                         float64
	//Size, SizeRand               int
	LifeSpan, LifeSpanRand float64
	// 0 - between quadrant 1 and 4
	// 90 - between quadrant 2 and 1
	Angle, AngleRand           float64
	Speed, SpeedRand           float64
	SpreadX, SpreadY           float64
	Duration                   int
	GravityX, GravityY         float64
	StartColor, StartColorRand color.Color
	EndColor, EndColorRand     color.Color
}

type ParticleSource struct {
	Generator     ParticleGenerator
	particles     []Particle
	rotateBinding event.Binding
	layer         int
	cID           event.CID
}

// A particle is a colored pixel at a given position, moving in a certain direction.
// After a while, it will dissipate.
type Particle struct {
	x, y       float64
	velX, velY float64
	startColor color.Color
	endColor   color.Color
	life       float64
	totalLife  float64
}

func (ps *ParticleSource) Init() event.CID {
	return plastic.NextID(ps)
}

// Todo: add draw priority to call
func (pg *ParticleGenerator) Generate(layer int) *ParticleSource {
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
				Stop(ps_p)
			}
		}(&ps, pg.Duration)
	}
	return &ps
}

func (ps *ParticleSource) Draw(buff screen.Buffer) {
	for _, p := range ps.particles {

		r, g, b, a := p.startColor.RGBA()
		r2, g2, b2, a2 := p.endColor.RGBA()
		progress := p.life / p.totalLife
		c := color.RGBA{
			unit8OnScale(r, r2, progress),
			unit8OnScale(g, g2, progress),
			unit8OnScale(b, b2, progress),
			unit8OnScale(a, a2, progress),
		}

		img := image.NewRGBA(image.Rect(0, 0, 1, 1))

		img.SetRGBA(0, 0, c)

		draw.Draw(buff.RGBA(), buff.Bounds(),
			img, image.Point{int(p.x), int(p.y)}, draw.Over)
	}
}

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
			p.velX -= pg.GravityX
			p.velY -= pg.GravityY
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
			velY:       speed * math.Sin(angle),
			startColor: randColor(pg.StartColor, pg.StartColorRand),
			endColor:   randColor(pg.EndColor, pg.EndColorRand),
			life:       startLife,
			totalLife:  startLife,
		})
	}

	ps.particles = newParticles

	return nil
}

func clearParticles(id int, nothing interface{}) error {

	ps := plastic.GetEntity(id).(*ParticleSource)
	pg := ps.Generator

	if len(ps.particles) > 0 {
		newParticles := make([]Particle, 0)
		for _, p := range ps.particles {

			// Ignore dead particles
			if p.life > 0 {

				p.life--

				p.velX -= pg.GravityX
				p.velY -= pg.GravityY
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

func Stop(ps *ParticleSource) {
	ps.rotateBinding.Unbind()
	ps.rotateBinding, _ = ps.cID.Bind(clearParticles, "EnterFrame")
}

func floatFromSpread(f float64) float64 {
	return (f * 2 * rand.Float64()) - f
}

func roundFloat(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

func randColor(c, ra color.Color) color.Color {
	r, g, b, a := c.RGBA()
	r2, g2, b2, a2 := ra.RGBA()
	return color.RGBA{
		uint8Spread(r, r2),
		uint8Spread(g, g2),
		uint8Spread(b, b2),
		uint8Spread(a, a2),
	}
}

func uint8Spread(n, r uint32) uint8 {
	n = n / 257
	r = r / 257
	return uint8(math.Min(float64(int(n)+roundFloat(floatFromSpread(float64(r)))), 255.0))
}

func unit8OnScale(n, endN uint32, progress float64) uint8 {
	return uint8((float64(n) - float64(n)*(1.0-progress) + float64(endN)*(1.0-progress)) / 257)
}

func (ps *ParticleSource) GetLayer() int {
	return ps.layer
}
func (ps *ParticleSource) SetLayer(l int) {
	ps.layer = l
}
func (ps *ParticleSource) UnDraw() {
	ps.layer = -1
}

func (ps *ParticleSource) GetRGBA() *image.RGBA {
	return nil
}

func (ps *ParticleSource) SetPos(x, y float64) {
	ps.Generator.X = x
	ps.Generator.Y = y
}

func (ps *ParticleSource) Pause() {
	ps.rotateBinding.Unbind()
}

func (ps *ParticleSource) UnPause() {
	binding, _ := ps.cID.Bind(rotateParticles, "EnterFrame")
	ps.rotateBinding = binding
}
