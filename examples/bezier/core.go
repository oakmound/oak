package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/oak/shape"
)

var (
	cmp *render.Composite
)

func main() {

	// c bezier accepts the same inputs as BezierCurve.
	// The indicator will follow the path given here.
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
		bz, err := shape.BezierCurve(floats...)
		if err != nil {
			fmt.Println(err)
		}
		if cmp != nil {
			cmp.UnDraw()
		}
		cmp = bezierDraw(bz)
		render.Draw(cmp, 0)
	})

	oak.Add("bezier", func(string, interface{}) {
		// Stubs
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "bezier", nil
	})
	oak.Init("bezier")
}

func bezierDraw(b shape.Bezier) *render.Composite {
	list := render.NewComposite([]render.Modifiable{})
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
		sp := render.NewColorBox(3, 3, color.RGBA{255, 0, 0, 255})
		sp.SetPos(bzn.X-1, bzn.Y-1)
		list.Append(sp)
	default:
	}
}
