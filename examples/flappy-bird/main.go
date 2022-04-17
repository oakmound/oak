package main

import (
	"image/color"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/span"
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
	score int
)

// label pillars with a known constant, so when we hit them, we can restart the scene
const (
	pillar collision.Label = iota
)

func main() {
	oak.AddScene("flappy", scene.Scene{Start: func(ctx *scene.Context) {
		render.Draw(render.NewDrawFPS(0, nil, 10, 10), 2, 0)
		render.Draw(render.NewLogicFPS(0, nil, 10, 20), 2, 0)

		score = 0
		// 1. Make Player
		newFlappy(ctx, 90, 140)
		// 2. Make scrolling repeating pillars
		pillarFreq := span.NewLinear(1.0, 5.0)
		var pillarLoop func()
		pillarLoop = func() {
			newPillarPair(ctx)
			ctx.DoAfter(time.Duration(pillarFreq.Poll()*float64(time.Second)), pillarLoop)
		}
		go ctx.DoAfter(time.Duration(pillarFreq.Poll()*float64(time.Second)), pillarLoop)

		// 3. Make Score
		t := render.DefaultFont().NewIntText(&score, 200, 30)
		render.Draw(t, 0)
	}})
	oak.Init("flappy")
}

func newFlappy(ctx *scene.Context, x, y float64) {
	f := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 32, 32)),
		entities.WithColor(color.RGBA{0, 255, 255, 255}),
		entities.WithDrawLayers([]int{0, 1}),
	)

	event.Bind(ctx, event.Enter, f, func(f *entities.Entity, ev event.EnterPayload) event.Response {
		f.ShiftDelta()
		if f.Delta.Y() > 10 {
			f.Delta[1] = 10
		}
		if f.Delta.Y() < -5 {
			f.Delta[1] = -5
		}
		// Gravity
		f.Delta[1] += .15

		if collision.HitLabel(f.Space, pillar) != nil {
			ctx.Window.NextScene()
		}

		if f.Bottom() > 480 {
			ctx.Window.NextScene()
		}
		if f.Y() < 0 {
			f.ShiftY(-f.Y())
			f.Delta[1] = 0
		}
		return 0
	})
	event.Bind(ctx, mouse.Press, f, func(f *entities.Entity, _ *mouse.Event) event.Response {
		f.Delta[1] -= 4
		return 0
	})
	event.Bind(ctx, key.Down(key.W), f, func(f *entities.Entity, _ key.Event) event.Response {
		f.Delta[1] -= 4
		return 0
	})
}

var (
	gapPosition = span.NewLinear(10.0, 370.0)
	gapSpan     = span.NewLinear(100.0, 250.0)
)

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

func newPillar(ctx *scene.Context, x, y, h float64, isAbove bool) {
	p := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 64, h)),
		entities.WithColor(color.RGBA{0, 255, 0, 255}),
		entities.WithLabel(pillar),
		entities.WithDrawLayers([]int{0, 1}),
	)
	event.Bind(ctx, event.Enter, p, enterPillar(isAbove))
}

func enterPillar(isAbove bool) func(p *entities.Entity, ev event.EnterPayload) event.Response {
	return func(p *entities.Entity, ev event.EnterPayload) event.Response {
		p.ShiftX(-2)
		if p.X()+p.W() < 0 {
			// don't score one out of each two pillars
			if isAbove {
				score++
			}
			p.Destroy()
		}
		return 0
	}
}
