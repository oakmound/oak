package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/oak/shape"
)

var (
	cmp *render.Composite
)

func renderCurve(floats []float64) {
	bz, err := shape.BezierCurve(floats...)
	if err != nil {
		fmt.Println(err)
	}
	if cmp != nil {
		cmp.Undraw()
	}
	cmp = bezierDraw(bz)
	render.Draw(cmp, 0)
}

func main() {

	// c bezier X Y X Y X Y ...
	// for defining custom points without using the mouse.
	// does not interact with the mouse points tracked through left clicks.
	oak.AddCommand("bezier", func(tokens []string) {
		if len(tokens) < 4 {
			return
		}
		tokens = tokens[1:]
		var err error
		floats := make([]float64, len(tokens))
		for i, s := range tokens {
			floats[i], err = strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		renderCurve(floats)
	})

	oak.Add("bezier", func(string, interface{}) {
		mouseFloats := []float64{}
		event.GlobalBind(func(_ int, mouseEvent interface{}) int {
			me := mouseEvent.(mouse.Event)
			// Left click to add a point to the curve
			if me.Button == "LeftMouse" {
				mouseFloats = append(mouseFloats, float64(me.X()), float64(me.Y()))
				renderCurve(mouseFloats)
				// Perform any other click to reset the drawn curve
			} else {
				mouseFloats = []float64{}
				cmp.Undraw()
			}
			return 0
		}, "MousePress")
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "bezier", nil
	})
	oak.Init("bezier")
}

func bezierDraw(b shape.Bezier) *render.Composite {
	list := render.NewComposite()
	bezierDrawRec(b, list, 255)
	return list
}

func bezierDrawRec(b shape.Bezier, list *render.Composite, alpha uint8) {
	switch bzn := b.(type) {
	case shape.BezierNode:
		sp := render.BezierLine(b, color.RGBA{alpha, 0, 0, alpha})
		list.Append(sp)

		bezierDrawRec(bzn.Left, list, uint8(float64(alpha)*.5))
		bezierDrawRec(bzn.Right, list, uint8(float64(alpha)*.5))
	case shape.BezierPoint:
		sp := render.NewColorBox(5, 5, color.RGBA{255, 255, 255, 255})
		sp.SetPos(bzn.X-2, bzn.Y-2)
		list.Append(sp)
	}
}

// Todo: could add a little animation that follows each of the bezier curves
// around as progress increases
