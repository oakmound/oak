package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var (
	bz          render.Bezier
	progressInc = 0.01
	progress    float64
)

func main() {

	// c bezier accepts the same inputs as BezierCurve.
	// The indicator will follow the path given here.
	oak.AddCommand("bezier", func(tokens []string) {
		if len(tokens) < 2 {
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
		bz, err = render.BezierCurve(floats...)
		if err != nil {
			fmt.Println(err)
		}
	})

	// This command changes how fast the curve will advance.
	// Todo: add utilities to shorthand commands for modifying
	// float and int pointers.
	oak.AddCommand("speed", func(tokens []string) {
		if len(tokens) < 2 {
			return
		}
		tokens = tokens[1:]
		f, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		progressInc = f
	})

	oak.AddScene("bezier", func(string, interface{}) {

		// Use a color box to indicate where we are on the curve.
		cb := render.NewColorBox(3, 3, color.RGBA{255, 0, 0, 255})
		cb.SetPos(320, 240)
		render.Draw(cb, 0)

		// Every frame, move the bezier indicator along the curve.
		event.GlobalBind(func(int, interface{}) int {
			if bz == nil {
				return 0
			}
			cb.SetPos(bz.Pos(progress))
			progress += progressInc
			if progress > 1.0 {
				progress = 0.0
			}
			return 0
		}, "EnterFrame")

		// Stubs
	}, func() bool {
		return true
	}, func() (string, *oak.SceneResult) {
		return "bezier", nil
	})
	oak.Init("bezier")
}
