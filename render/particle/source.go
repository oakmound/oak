package particle

import (
	"math"
	"time"

	"github.com/200sc/go-dist/intrange"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
)

const (
	// recycled refers to the life value given to particles ready to be reused
	recycled = -1000
	//IgnoreEnd refers to the life value given to particles that want to skip their source's end function.
	IgnoreEnd = recycled / 2
)

// A Source is used to store and control a set of particles.
type Source struct {
	render.Layered
	Generator  Generator
	particles  [blockSize]Particle
	nextPID    int
	recycled   []int
	CID        event.CID
	pIDBlock   int
	stackLevel int
	paused     bool
	EndFunc    func()
}

// NewSource creates a new source
func NewSource(g Generator, stackLevel int) *Source {
	ps := new(Source)
	ps.Generator = g
	ps.stackLevel = stackLevel
	ps.Init()
	return ps
}

// Init allows a source to be considered as an entity, and initializes it
func (ps *Source) Init() event.CID {
	CID := event.NextID(ps)
	CID.Bind(rotateParticles, event.Enter)
	ps.CID = CID
	ps.pIDBlock = Allocate(ps.CID)
	if ps.Generator.GetBaseGenerator().Duration != Inf {
		go func(ps_p *Source, duration intrange.Range) {
			timing.DoAfter(time.Duration(duration.Poll())*time.Millisecond, func() {
				ps_p.Stop()
			})
		}(ps, ps.Generator.GetBaseGenerator().Duration)
	}
	return CID
}

func (ps *Source) cycleParticles() bool {
	pg := ps.Generator.GetBaseGenerator()
	cycled := false
	for i := 0; i < ps.nextPID; i++ {
		p := ps.particles[i]
		bp := p.GetBaseParticle()

		// Ignore dead particles
		if bp.Life > 0 {
			cycled = true
			bp.Life--
			// Apply rotational acceleration
			if pg.Rotation != nil {
				bp.Vel = bp.Vel.Rotate(pg.Rotation.Poll())
			}

			if pg.SpeedDecay.X() != 0 {
				bp.Vel = bp.Vel.SetX(bp.Vel.X() * pg.SpeedDecay.X())
				if math.Abs(bp.Vel.X()) < 0.2 {
					bp.Vel = bp.Vel.SetX(0)
				}
			}
			if pg.SpeedDecay.Y() != 0 {
				bp.Vel = bp.Vel.SetY(bp.Vel.Y() * pg.SpeedDecay.Y())
				if math.Abs(bp.Vel.Y()) < 0.2 {
					bp.Vel = bp.Vel.SetY(0)
				}
			}

			bp.Vel.Add(pg.Gravity)
			bp.Add(bp.Vel)
			bp.SetLayer(ps.Layer(bp.GetPos()))
			p.Cycle(ps.Generator)

		} else if bp.Life != recycled {
			p.UnDraw()
			cycled = true
			if pg.EndFunc != nil && bp.Life > IgnoreEnd {
				pg.EndFunc(p)
			}
			// This relies on Life going down by 1 per frame
			bp.Life = recycled
			ps.recycled = append(ps.recycled, bp.pID%blockSize)
		}

	}
	return cycled
}

// Layer is shorthand for getting the base generator behind a source's layer
func (ps *Source) Layer(v physics.Vector) int {
	return ps.Generator.GetBaseGenerator().LayerFunc(v)
}

func (ps *Source) addParticles() {
	pg := ps.Generator.GetBaseGenerator()
	// Regularly create particles (up until max particles)
	newParticleCount := int(pg.NewPerFrame.Poll())
	ri := 0
	for ri < len(ps.recycled) && ri < newParticleCount {

		j := ps.recycled[ri]
		bp := ps.particles[j].GetBaseParticle()
		angle := pg.Angle.Poll() * math.Pi / 180.0
		speed := pg.Speed.Poll()
		startLife := pg.LifeSpan.Poll()

		bp.LayeredPoint = render.NewLayeredPoint(
			pg.X()+floatFromSpread(pg.Spread.X()),
			pg.Y()+floatFromSpread(pg.Spread.Y()),
			0)
		bp.Vel = physics.NewVector(
			speed*math.Cos(angle)*-1,
			speed*math.Sin(angle)*-1)
		bp.Life = startLife
		bp.totalLife = startLife
		ps.particles[ps.recycled[ri]] = ps.Generator.GenerateParticle(bp)

		ps.particles[ps.recycled[ri]].SetLayer(ps.Layer(bp.GetPos()))
		render.Draw(ps.particles[ps.recycled[ri]], ps.stackLevel)
		ri++
	}
	newParticleCount -= len(ps.recycled)
	if ri > 0 {
		ps.recycled = ps.recycled[ri:]
	}

	if ps.nextPID >= blockSize {
		return
	}
	for i := 0; i < newParticleCount; i++ {
		angle := pg.Angle.Poll() * math.Pi / 180.0
		speed := pg.Speed.Poll()
		startLife := pg.LifeSpan.Poll()

		bp := &baseParticle{
			LayeredPoint: render.NewLayeredPoint(
				pg.X()+floatFromSpread(pg.Spread.X()),
				pg.Y()+floatFromSpread(pg.Spread.Y()),
				0,
			),
			Src: ps,
			Vel: physics.NewVector(
				speed*math.Cos(angle)*-1,
				speed*math.Sin(angle)*-1),
			Life:      startLife,
			totalLife: startLife,
			pID:       ps.nextPID + ps.pIDBlock*blockSize,
		}

		p := ps.Generator.GenerateParticle(bp)
		ps.particles[ps.nextPID] = p
		ps.nextPID++
		if ps.nextPID >= blockSize {
			return
		}
		p.SetLayer(ps.Layer(bp.GetPos()))
		render.Draw(p, ps.stackLevel)
	}

}

// rotateParticles updates particles over time as long
// as a Source is active.
func rotateParticles(id int, nothing interface{}) int {
	ps := event.GetEntity(id).(*Source)
	if !ps.paused {
		ps.cycleParticles()
		ps.addParticles()
	}
	return 0
}

// clearParticles is used after a Source has been stopped
// to continue moving old particles for as long as they exist.
func clearParticles(id int, nothing interface{}) int {
	if ps, ok := event.GetEntity(id).(*Source); ok {
		if !ps.paused {
			if ps.cycleParticles() {
			} else {
				if ps.EndFunc != nil {
					ps.EndFunc()
				}
				event.DestroyEntity(id)
				Deallocate(ps.pIDBlock)
				return event.UnbindEvent
			}
		}

		return 0
	}
	return event.UnbindEvent
}

// Stop manually stops a Source, if its duration is infinite
// or if it should be stopped before expriring naturally.
func (ps *Source) Stop() {
	if ps == nil {
		return
	}
	ps.CID.UnbindAllAndRebind([]event.Bindable{clearParticles}, []string{event.Enter})
}

// Pause on a Source just stops the repetition
// of its rotation function, which moves, destroys,
// ages and generates particles. Existing particles will
// stay in place.
func (ps *Source) Pause() {
	ps.paused = true
}

// UnPause on a source a Source rebinds it's rotate function.
func (ps *Source) UnPause() {
	ps.paused = false
}

// Placeholder
func (ps *Source) String() string {
	return "ParticleSource"
}

// ShiftX shift's a source's underlying generator
func (ps *Source) ShiftX(x float64) {
	ps.Generator.ShiftX(x)
}

// ShiftY shift's a source's underlying generator (todo: consider if this shoud be composed)
func (ps *Source) ShiftY(y float64) {
	ps.Generator.ShiftY(y)
}

// SetPos sets a source's underlying generator
func (ps *Source) SetPos(x, y float64) {
	ps.Generator.SetPos(x, y)
}
