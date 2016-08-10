package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"image"
	"image/draw"
	"math"
	"time"
)

type SpriteGenerator struct {
	NewPerFrame, NewPerFrameRand       float64
	X, Y                               float64
	LifeSpan, LifeSpanRand             float64
	Angle, AngleRand                   float64
	Speed, SpeedRand                   float64
	SpreadX, SpreadY                   float64
	Duration                           int
	Rotation, RotationRand             float64
	SpriteRotation, SpriteRotationRand float64
	GravityX, GravityY                 float64
	Base                               *render.Sprite
}

type SpriteSource struct {
	render.Layered
	Generator     SpriteGenerator
	particles     []SpriteParticle
	rotateBinding event.Binding
	clearBinding  event.Binding
	cID           event.CID
}

type SpriteParticle struct {
	x, y       float64
	velX, velY float64
	rotation   float64
	life       float64
}

func (ss *SpriteSource) Init() event.CID {
	return event.NextID(ss)
}

// Generate takes a generator and converts it into a source,
// drawing particles and binding functions for particle generation
// and rotation.
func (sg *SpriteGenerator) Generate(layer int) *SpriteSource {

	// Convert rotation from degrees to radians
	sg.Rotation = sg.Rotation / 180 * math.Pi
	sg.RotationRand = sg.Rotation / 180 * math.Pi

	// Make a source
	ss := SpriteSource{
		Generator: *sg,
		particles: make([]SpriteParticle, 0),
	}

	// Bind things to that source:
	cID := ss.Init()
	binding, _ := cID.Bind(rotateSprites, "EnterFrame")
	ss.rotateBinding = binding
	ss.cID = cID
	render.Draw(&ss, layer)
	if sg.Duration != -1 {
		go func(ss *SpriteSource, duration int) {
			select {
			case <-time.After(time.Duration(duration) * time.Millisecond):
				if ss.GetLayer() != -1 {
					ss.Stop()
				}
			}
		}(&ss, sg.Duration)
	}
	return &ss
}

func (ss *SpriteSource) Draw(buff draw.Image) {
	base := ss.Generator.Base
	for _, s := range ss.particles {
		img := base.Copy()
		img.Rotate(int(s.rotation))
		render.ShinyDraw(buff, img.GetRGBA(), int(s.x), int(s.y))
	}
}

func rotateSprites(id int, nothing interface{}) int {
	ss := event.GetEntity(id).(*SpriteSource)
	sg := ss.Generator

	newParticles := make([]SpriteParticle, 0)

	for _, s := range ss.particles {

		// Ignore dead particles
		if s.life > 0 {

			// Move towards doom
			s.life--

			// Be dragged down by the weight of the soul
			s.velX += sg.GravityX
			s.velY += sg.GravityY

			// Apply rotational acceleration
			if sg.Rotation != 0 && sg.RotationRand != 0 {
				magnitude := math.Abs(s.velX) + math.Abs(s.velY)
				angle := math.Atan2(s.velX, s.velY)
				angle += sg.Rotation + floatFromSpread(sg.RotationRand)
				s.velX = math.Sin(angle)
				s.velY = math.Cos(angle)
				magnitude = magnitude / (math.Abs(s.velX) + math.Abs(s.velY))
				s.velX = s.velX * magnitude
				s.velY = s.velY * magnitude
			}
			s.rotation += s.rotation

			s.x += s.velX
			s.y += s.velY

			newParticles = append(newParticles, s)
		}
	}

	// Regularly create particles (up until max particles)
	newParticleRand := roundFloat(floatFromSpread(sg.NewPerFrameRand))
	newParticleCount := int(sg.NewPerFrame) + newParticleRand
	for i := 0; i < newParticleCount; i++ {

		angle := (sg.Angle + floatFromSpread(sg.AngleRand)) * math.Pi / 180.0
		speed := sg.Speed + floatFromSpread(sg.SpeedRand)
		startLife := sg.LifeSpan + floatFromSpread(sg.LifeSpanRand)
		rotation := sg.SpriteRotation + floatFromSpread(sg.SpriteRotationRand)

		newParticles = append(newParticles, SpriteParticle{
			x:        sg.X + floatFromSpread(sg.SpreadX),
			y:        sg.Y + floatFromSpread(sg.SpreadY),
			velX:     speed * math.Cos(angle) * -1,
			velY:     speed * math.Sin(angle) * -1,
			life:     startLife,
			rotation: rotation,
		})
	}

	ss.particles = newParticles

	return 0
}

func clearSprites(id int, nothing interface{}) int {
	ss := event.GetEntity(id).(*SpriteSource)
	sg := ss.Generator

	if len(ss.particles) > 0 {
		newParticles := make([]SpriteParticle, 0)
		for _, s := range ss.particles {

			// Ignore dead particles
			if s.life > 0 {

				// Move towards doom
				s.life--

				// Be dragged down by the weight of the soul
				s.velX += sg.GravityX
				s.velY += sg.GravityY

				// Apply rotational acceleration
				if sg.Rotation != 0 && sg.RotationRand != 0 {
					magnitude := math.Abs(s.velX) + math.Abs(s.velY)
					angle := math.Atan2(s.velX, s.velY)
					angle += sg.Rotation + floatFromSpread(sg.RotationRand)
					s.velX = math.Sin(angle)
					s.velY = math.Cos(angle)
					magnitude = magnitude / (math.Abs(s.velX) + math.Abs(s.velY))
					s.velX = s.velX * magnitude
					s.velY = s.velY * magnitude
				}
				s.rotation += s.rotation

				s.x += s.velX
				s.y += s.velY

				newParticles = append(newParticles, s)
			}
		}
		ss.particles = newParticles
	} else {
		ss.UnDraw()
		ss.rotateBinding.Unbind()
	}
	return 0
}

func clearSpritesAtExit(id int, nothing interface{}) int {
	ss := event.GetEntity(id).(*SpriteSource)
	ss.clearBinding.Unbind()
	ss.rotateBinding.Unbind()
	ss.rotateBinding, _ = ss.cID.Bind(clearSprites, "EnterFrame")
	return 0
}

func (ss *SpriteSource) Stop() {
	ss.clearBinding, _ = ss.cID.Bind(clearSpritesAtExit, "ExitFrame")
}

// A particle source has no concept of an individual
// rgba buffer, and so it returns nothing when its
// rgba buffer is queried. This may change.
func (ss *SpriteSource) GetRGBA() *image.RGBA {
	return nil
}

func (ss *SpriteSource) ShiftX(x float64) {
	ss.Generator.X += x
}

func (ss *SpriteSource) ShiftY(y float64) {
	ss.Generator.Y += y
}

func (ss *SpriteSource) GetX() float64 {
	return ss.Generator.X
}

func (ss *SpriteSource) GetY() float64 {
	return ss.Generator.Y
}
func (ss *SpriteSource) SetPos(x, y float64) {
	ss.Generator.X = x
	ss.Generator.Y = y
}

func (ss *SpriteSource) Pause() {
	ss.rotateBinding.Unbind()
}

func (ss *SpriteSource) UnPause() {
	binding, _ := ss.cID.Bind(rotateSprites, "EnterFrame")
	ss.rotateBinding = binding
}
