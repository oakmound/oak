package main

import (
	"image/color"
	"time"

	"github.com/oakmound/oak/v3/alg/range/floatrange"
	"github.com/oakmound/oak/v3/mouse"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
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
	oak.AddScene("bounce", scene.Scene{Start: func(ctx *scene.Context) {
		render.Draw(render.NewDrawFPS(0.03, nil, 10, 10))

		score = 0
		// 1. Make Player
		newFlappy(ctx, 90, 140)
		// 2. Make scrolling repeating pillars
		var pillarLoop func()
		pillarLoop = func() {
			newPillarPair(ctx)
			ctx.DoAfter(time.Duration(pillarFreq.Poll()*float64(time.Second)), pillarLoop)
		}
		go ctx.DoAfter(time.Duration(pillarFreq.Poll()*float64(time.Second)), pillarLoop)

		// 3. Make Score
		t := render.DefaultFont().NewIntText(&score, 200, 30)
		render.Draw(t, 0)
	}, End: func() (string, *scene.Result) {
		return "bounce", nil
	}})
	render.SetDrawStack(
		render.NewDynamicHeap(),
	)
	oak.Init("bounce")
}

// A Flappy is on a journey to go to the right
type Flappy struct {
	*entities.Interactive
}

// CID returns the event.CallerID so that this can be bound to.
func (flap *Flappy) CID() event.CallerID {
	return flap.CallerID
}

func newFlappy(ctx *scene.Context, x, y float64) *Flappy {
	f := new(Flappy)
	f.Interactive = entities.NewInteractive(x, y, 32, 32, render.NewColorBox(32, 32, color.RGBA{0, 255, 255, 255}), nil, ctx.Register(f), 1)

	f.RSpace.Add(pillar, func(s1, s2 *collision.Space) {
		ctx.Window.NextScene()
	})
	f.RSpace.Space.Label = player
	collision.Add(f.RSpace.Space)

	f.R.SetLayer(1)
	render.Draw(f.R, 0)
	event.Bind(ctx, event.Enter, f, func(f *Flappy, ev event.EnterPayload) event.Response {
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

		<-f.RSpace.CallOnHits()
		if f.Y()+f.H > 480 {
			ctx.Window.NextScene()
		}
		if f.Y() < 0 {
			f.SetY(0)
			f.Delta.SetY(0)
		}
		return 0
	})

	event.Bind(ctx, mouse.Press, f, func(f *Flappy, me *mouse.Event) event.Response {
		f.Delta.ShiftY(-4)
		return 0
	})
	event.Bind(ctx, key.Down(key.W), f, func(f *Flappy, k key.Event) event.Response {
		f.Delta.ShiftY(-4)
		return 0
	})
	return f
}

// A Pillar blocks flappy from continuing forward
type Pillar struct {
	*entities.Solid
	hasScored bool
}

// CID returns the event.CallerID so that this can be bound to.
func (p *Pillar) CID() event.CallerID {
	return p.CallerID
}

func newPillar(ctx *scene.Context, x, y, h float64, isAbove bool) {
	p := new(Pillar)
	p.Solid = entities.NewSolid(x, y, 64, h, render.NewColorBox(64, int(h), color.RGBA{0, 255, 0, 255}), nil, ctx.Register(p))
	p.Space.Label = pillar
	collision.Add(p.Space)
	event.Bind(ctx, event.Enter, p, enterPillar)

	p.R.SetLayer(1)
	render.Draw(p.R, 0)
	// Don't score one out of each two pillars
	if isAbove {
		p.hasScored = true
	}
}

func newPillarPair(ctx *scene.Context) {
	pos := gapPosition.Poll()
	span := gapSpan.Poll()
	if (pos + span) > 470 {
		span = 470 - pos
	}
	if span < 100 {
		pos = 370
		span = 100
	}
	newPillar(ctx, 641, 0, pos, true)
	newPillar(ctx, 641, pos+span, 480-(pos+span), false)
}

func enterPillar(p *Pillar, ev event.EnterPayload) event.Response {
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
