package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// This example demonstrates the use of the Propagated boolean on
// mouse event payloads to prevent mouse interactions from falling
// through to lower UI elements after interacting with a higher layer

func main() {
	oak.AddScene("click-propagation", scene.Scene{
		Start: func(ctx *scene.Context) {
			z := 0
			y := 400.0
			for x := 20.0; x < 400; x += 20 {
				z++
				y -= 20
				newHoverButton(x, y, 35, 35, color.RGBA{200, 200, 200, 200}, z)
			}
		},
	})
	oak.Init("click-propagation")
}

type hoverButton struct {
	id event.CID

	mouse.CollisionPhase
	*changingColorBox
}

func (hb *hoverButton) Init() event.CID {
	hb.id = event.NextID(hb)
	return hb.id
}

func newHoverButton(x, y, w, h float64, clr color.RGBA, layer int) {
	hb := &hoverButton{}
	hb.Init()
	hb.changingColorBox = newChangingColorBox(x, y, int(w), int(h), clr)

	sp := collision.NewSpace(x, y, w, h, hb.id)
	sp.SetZLayer(float64(layer))

	mouse.Add(sp)
	mouse.PhaseCollision(sp)

	render.Draw(hb.changingColorBox, 0, layer)
	hb.id.Bind(mouse.ClickOn, func(c event.CID, i interface{}) int {
		hb := event.GetEntity(c).(*hoverButton)
		me := i.(*mouse.Event)
		fmt.Println(c, me.Point2)
		hb.changingColorBox.c = color.RGBA{128, 128, 128, 128}
		me.StopPropagation = true
		return 0
	})
	hb.id.Bind(mouse.Start, func(c event.CID, i interface{}) int {
		fmt.Println("start")
		hb := event.GetEntity(c).(*hoverButton)
		me := i.(*mouse.Event)
		hb.changingColorBox.c = color.RGBA{50, 50, 50, 50}
		me.StopPropagation = true
		return 0
	})
	hb.id.Bind(mouse.Stop, func(c event.CID, i interface{}) int {
		fmt.Println("stop")
		hb := event.GetEntity(c).(*hoverButton)
		me := i.(*mouse.Event)
		hb.changingColorBox.c = clr
		me.StopPropagation = true
		return 0
	})
}

type changingColorBox struct {
	render.LayeredPoint
	c    color.RGBA
	w, h int
}

func newChangingColorBox(x, y float64, w, h int, c color.RGBA) *changingColorBox {
	return &changingColorBox{
		LayeredPoint: render.NewLayeredPoint(x, y, 0),
		c:            c,
		w:            w,
		h:            h,
	}
}

func (ccb *changingColorBox) Draw(buff draw.Image, xOff, yOff float64) {
	x := int(ccb.X() + xOff)
	y := int(ccb.Y() + yOff)
	rect := image.Rect(x, y, ccb.w+x, ccb.h+y)
	draw.Draw(buff, rect, image.NewUniform(ccb.c), image.Point{int(ccb.X() + xOff), int(ccb.Y() + yOff)}, draw.Over)
}

func (ccb *changingColorBox) GetDims() (int, int) {
	return ccb.w, ccb.h
}
