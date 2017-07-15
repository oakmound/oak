package main

import (
	"image/color"
	"time"

	"github.com/200sc/go-dist/floatrange"

	"bitbucket.org/oakmoundstudio/oak"
	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/entities"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
	"bitbucket.org/oakmoundstudio/oak/timing"
)

var (
	pillarFreq      = floatrange.NewLinear(1, 5)
	gapPosition     = floatrange.NewLinear(10, 370)
	gapSpan         = floatrange.NewLinear(100, 250)
	playerHitPillar bool
	score           int
)

const (
	PLAYER collision.Label = iota
	PILLAR
)

func main() {
	oak.AddScene("bounce", func(string, interface{}) {
		score = 0
		// 1. Make Player
		NewFlappy(90, 140)
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
					NewPillarPair()
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
	}, func() (string, *oak.SceneResult) {
		return "bounce", nil
	})
	render.SetDrawStack(
		render.NewHeap(false),
		render.NewDrawFPS(),
	)
	oak.Init("bounce")
}

type Flappy struct {
	entities.Interactive
}

func (f *Flappy) Init() event.CID {
	f.CID = event.NextID(f)
	return f.CID
}

func NewFlappy(x, y float64) *Flappy {
	f := new(Flappy)
	f.Init()
	f.Interactive = entities.NewInteractive(x, y, 32, 32, render.NewColorBox(32, 32, color.RGBA{0, 255, 255, 255}), f.CID, 1)

	// Can use a const block for collision labels, here we
	// don't because we only need two.
	f.RSpace.Add(PILLAR, func(s1, s2 *collision.Space) {
		playerHitPillar = true
	})
	f.RSpace.Space.Label = PLAYER
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
	}, "EnterFrame")
	f.Bind(func(int, interface{}) int {
		f.Delta.ShiftY(-4)
		return 0
	}, "KeyDownW")
	return f
}

type Pillar struct {
	entities.Solid
	hasScored bool
}

func (p *Pillar) Init() event.CID {
	p.CID = event.NextID(p)
	return p.CID
}

func NewPillar(x, y, h float64, isAbove bool) {
	pillar := new(Pillar)
	pillar.Init()
	pillar.Solid = entities.NewSolid(x, y, 64, h, render.NewColorBox(64, int(h), color.RGBA{0, 255, 0, 255}), pillar.CID)
	pillar.Space.Label = PILLAR
	collision.Add(pillar.Space)
	pillar.Bind(enterPillar, "EnterFrame")
	pillar.R.SetLayer(1)
	render.Draw(pillar.R, 0)
}

func NewPillarPair() {
	pos := gapPosition.Poll()
	span := gapSpan.Poll()
	if (pos + span) > 470 {
		span = 470 - pos
	}
	if span < 100 {
		pos = 370
		span = 100
	}
	NewPillar(641, 0, pos, true)
	NewPillar(641, pos+span, 480-(pos+span), false)
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
