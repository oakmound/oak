package particle

import (
	"math"
	"time"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

const (
	RECYCLED = -1000
)

// A Source is used to store and control a set of particles.
type Source struct {
	render.Layered
	Generator     Generator
	particles     [BLOCK_SIZE]Particle
	nextPID       int
	recycled      []int
	rotateBinding event.Binding
	clearBinding  event.Binding
	cID           event.CID
	paused        bool
	pIDBlock      int
}

func (ps *Source) Init() event.CID {
	cID := event.NextID(ps)
	cID.Bind(rotateParticles, "EnterFrame")
	ps.cID = cID
	ps.pIDBlock = Allocate(ps.cID)
	if ps.Generator.GetBaseGenerator().Duration != -1 {
		go func(ps_p *Source, duration int) {
			select {
			case <-time.After(time.Duration(duration) * time.Millisecond):
				ps_p.Stop()
			}
		}(ps, ps.Generator.GetBaseGenerator().Duration)
	}
	return event.NextID(ps)
}

func (ps *Source) CycleParticles() {

	pg := ps.Generator.GetBaseGenerator()

	for i := 0; i < ps.nextPID; i++ {
		p := ps.particles[i]
		bp := p.GetBaseParticle()

		// Ignore dead particles
		if bp.life > 0 {

			// Move towards doom
			bp.life--

			// Apply rotational acceleration
			if pg.Rotation != 0 || pg.RotationRand != 0 {
				bp.Vel.Rotate(pg.Rotation, floatFromSpread(pg.RotationRand))
			}

			if pg.SpeedDecayX != 0 {
				bp.Vel.X *= pg.SpeedDecayX
				if math.Abs(bp.Vel.X) < 0.2 {
					bp.Vel.X = 0
				}
			}
			if pg.SpeedDecayY != 0 {
				bp.Vel.Y *= pg.SpeedDecayY
				if math.Abs(bp.Vel.Y) < 0.2 {
					bp.Vel.Y = 0
				}
			}

			bp.Vel.X += pg.GravityX
			bp.Vel.Y += pg.GravityY

			bp.Pos.Add(bp.Vel)
			bp.SetLayer(ps.Layer(bp.Pos))
		} else if bp.life != RECYCLED {
			p.UnDraw()
			if pg.EndFunc != nil {
				pg.EndFunc(p)
			}
			// This relies on life going down by 1 per frame
			bp.life = RECYCLED
			ps.recycled = append(ps.recycled, bp.pID%BLOCK_SIZE)
		}
	}
}

// Shorthand
func (ps *Source) Layer(v *physics.Vector) int {
	return ps.Generator.GetBaseGenerator().LayerFunc(v)
}

func (ps *Source) AddParticles() {
	pg := ps.Generator.GetBaseGenerator()
	// Regularly create particles (up until max particles)
	newParticleRand := roundFloat(floatFromSpread(pg.NewPerFrameRand))
	newParticleCount := int(pg.NewPerFrame) + newParticleRand
	ri := 0
	for ri < len(ps.recycled) && ri < newParticleCount {
		j := ps.recycled[ri]
		bp := ps.particles[j].GetBaseParticle()
		angle := (pg.Angle + floatFromSpread(pg.AngleRand)) * math.Pi / 180.0
		speed := pg.Speed + floatFromSpread(pg.SpeedRand)
		startLife := pg.LifeSpan + floatFromSpread(pg.LifeSpanRand)

		bp.Pos = physics.NewVector(
			pg.X+floatFromSpread(pg.SpreadX),
			pg.Y+floatFromSpread(pg.SpreadY))
		bp.Vel = physics.NewVector(
			speed*math.Cos(angle)*-1,
			speed*math.Sin(angle)*-1)
		bp.life = startLife
		bp.totalLife = startLife
		ps.particles[ps.recycled[ri]] = ps.Generator.GenerateParticle(*bp)

		render.Draw(ps.particles[ps.recycled[ri]], ps.Layer(bp.Pos))

		ri++
	}
	newParticleCount -= len(ps.recycled)
	ps.recycled = ps.recycled[ri:]

	if ps.nextPID >= BLOCK_SIZE {
		return
	}
	for i := 0; i < newParticleCount; i++ {

		angle := (pg.Angle + floatFromSpread(pg.AngleRand)) * math.Pi / 180.0
		speed := pg.Speed + floatFromSpread(pg.SpeedRand)
		startLife := pg.LifeSpan + floatFromSpread(pg.LifeSpanRand)

		bp := BaseParticle{
			Src: ps,
			Pos: physics.NewVector(
				pg.X+floatFromSpread(pg.SpreadX),
				pg.Y+floatFromSpread(pg.SpreadY)),
			Vel: physics.NewVector(
				speed*math.Cos(angle)*-1,
				speed*math.Sin(angle)*-1),
			life:      startLife,
			totalLife: startLife,
			pID:       ps.nextPID + ps.pIDBlock*BLOCK_SIZE,
		}

		p := ps.Generator.GenerateParticle(bp)
		render.Draw(p, ps.Layer(bp.Pos))
		ps.particles[ps.nextPID] = p
		ps.nextPID++
	}
}

// rotateParticles updates particles over time as long
// as a Source is active.
func rotateParticles(id int, nothing interface{}) int {
	ps := event.GetEntity(id).(*Source)
	if !ps.paused {
		ps.CycleParticles()
		ps.AddParticles()
	}
	return 0
}

// clearParticles is used after a Source has been stopped
// to continue moving old particles for as long as they exist.
func clearParticles(id int, nothing interface{}) int {
	ps := event.GetEntity(id).(*Source)
	if !ps.paused {
		if len(ps.particles) > 0 {
			ps.CycleParticles()
		} else {
			event.DestroyEntity(id)
			Deallocate(ps.pIDBlock)
			return event.UNBIND_EVENT
		}
	}
	return 0
}

func clearAtExit(id int, nothing interface{}) int {
	if event.HasEntity(id) {
		t := event.GetEntity(id)
		switch t.(type) {
		case *Source:
			ps := t.(*Source)
			ps.cID.Bind(clearParticles, "ExitFrame")
			return event.UNBIND_EVENT
		}
	}
	return 0
}

// Stop manually stops a Source, if its duration is infinite
// or if it should be stopped before expriring naturally.
func (ps *Source) Stop() {
	ps.cID.Bind(clearAtExit, "EnterFrame")
}

// Pausing a Source just stops the repetition
// of its rotation function, which moves, destroys,
// ages and generates particles. Existing particles will
// stay in place.
func (ps *Source) Pause() {
	ps.paused = true
}

// Unpausing a Source rebinds it's rotate function.
func (ps *Source) UnPause() {
	ps.paused = false
}

// Placeholder
func (ps *Source) String() string {
	return "ParticleSource"
}

func (ps *Source) ShiftX(x float64) {
	ps.Generator.ShiftX(x)
}

func (ps *Source) ShiftY(y float64) {
	ps.Generator.ShiftY(y)
}

func (ps *Source) SetPos(x, y float64) {
	ps.Generator.SetPos(x, y)
}
