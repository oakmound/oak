package main

import (
	"image/color"
	"time"

	"github.com/200sc/go-dist/floatrange"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
	"github.com/oakmound/oak/v2/timing"
)

var (
	pillarFreq      = floatrange.NewLinear(1, 5)
	gapPosition     = floatrange.NewLinear(10, 370)
	gapSpan         = floatrange.NewLinear(100, 250)
	playerHitPillar bool
	score           int
)

// This const block is used for determining what type
// of entity is colliding with what
const (
	player collision.Label = iota
	pillar
)

func main() {
	oak.Add("bounce", func(string, interface{}) {
		score = 0
		// 1. Make Player
		newFlappy(90, 140)
		// 2. Make Scrolling background
		// 3. Make scrolling repeating pillars
		go func() {
			for {

				select {
				// this uses a signal sent when a scene ends,
				// or when otherwise timing operations need
				// to cease
				case <-timing.ClearDelayCh:
					return
				case <-time.After(time.Duration(pillarFreq.Poll() * float64(time.Second))):
					newPillarPair()
				}
			}
		}()
		// 4. Make Score
		t := render.DefFont().NewIntText(&score, 200, 30)
		render.Draw(t, 0)
	}, func() bool {
		if playerHitPillar {
			playerHitPillar = false
			return false
		}
		return true
	}, func() (string, *scene.Result) {
		return "bounce", nil
	})
	render.SetDrawStack(
		render.NewHeap(false),
		render.NewDrawFPS(),
	)
	oak.Init("bounce")
}

// A Flappy is on a journey to go to the right
type Flappy struct {
	*entities.Interactive
}

// Init satisfies the event.Entity interface
func (f *Flappy) Init() event.CID {
	return event.NextID(f)
}

func newFlappy(x, y float64) *Flappy {
	f := new(Flappy)
	f.Interactive = entities.NewInteractive(x, y, 32, 32, render.NewColorBox(32, 32, color.RGBA{0, 255, 255, 255}), nil, f.Init(), 1)

	f.RSpace.Add(pillar, func(s1, s2 *collision.Space) {
		playerHitPillar = true
	})
	f.RSpace.Space.Label = player
	collision.Add(f.RSpace.Space)

	//f.Vector = f.Vector.Attach(f.R)
	f.R.SetLayer(1)
	render.Draw(f.R, 0)

	f.Bind(func(int, interface{}) int {
		f.ShiftPos(f.Delta.X(), f.Delta.Y())
		f.Add(f.Delta)
		if f.Delta.Y() > 10 {
			f.Delta.SetY(10)
		}
		if f.Delta.Y() < -5 {
			f.Delta.SetY(-5)
		}
		// Gravity
		f.Delta.ShiftY(.15)

		// Todo: attachment as above with f.R for f.Space
		// ShiftPos does this for us right now
		// collision.UpdateSpace(f.X(), f.Y(), f.W, f.H, f.RSpace.Space)
		<-f.RSpace.CallOnHits()
		if f.Y()+f.H > 480 {
			playerHitPillar = true
		}
		if f.Y() < 0 {
			f.SetY(0)
			f.Delta.SetY(0)
		}
		return 0
	}, event.Enter)
	f.Bind(func(int, interface{}) int {
		f.Delta.ShiftY(-4)
		return 0
	}, key.Down+key.W)
	return f
}

// A Pillar blocks flappy from continuing forward
type Pillar struct {
	*entities.Solid
	hasScored bool
}

// Init satisfies the event.Entity interface
func (p *Pillar) Init() event.CID {
	return event.NextID(p)
}

func newPillar(x, y, h float64, isAbove bool) {
	p := new(Pillar)
	p.Solid = entities.NewSolid(x, y, 64, h, render.NewColorBox(64, int(h), color.RGBA{0, 255, 0, 255}), nil, p.Init())
	p.Space.Label = pillar
	collision.Add(p.Space)
	p.Bind(enterPillar, event.Enter)
	p.R.SetLayer(1)
	render.Draw(p.R, 0)
	// Don't score one out of each two pillars
	if isAbove {
		p.hasScored = true
	}
}

func newPillarPair() {
	pos := gapPosition.Poll()
	span := gapSpan.Poll()
	if (pos + span) > 470 {
		span = 470 - pos
	}
	if span < 100 {
		pos = 370
		span = 100
	}
	newPillar(641, 0, pos, true)
	newPillar(641, pos+span, 480-(pos+span), false)
}

func enterPillar(id int, nothing interface{}) int {
	p := event.GetEntity(id).(*Pillar)
	p.ShiftX(-2)
	if p.X()+p.W < 0 {
		p.Destroy()
	}
	if !p.hasScored && p.X()+p.W < 90 {
		p.hasScored = true
		score++
	}
	return 0
}
