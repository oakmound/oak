package main

import (
	"errors"
	"image/color"
	"log"
	"strconv"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg"
	"github.com/oakmound/oak/v3/alg/range/floatrange"
	"github.com/oakmound/oak/v3/alg/range/intrange"
	"github.com/oakmound/oak/v3/debugstream"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	pt "github.com/oakmound/oak/v3/render/particle"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/shape"
)

var (
	startColor     color.Color
	startColorRand color.Color
	endColor       color.Color
	endColorRand   color.Color
	src            *pt.Source
)

func parseShape(args []string) shape.Shape {
	if len(args) > 0 {
		switch args[0] {
		case "heart":
			return shape.Heart
		case "square":
			return shape.Square
		case "circle":
			return shape.Circle
		case "diamond":
			return shape.Diamond
		case "checkered":
			return shape.Checkered
		case "or":
			return shape.JustIn(shape.OrIn(parseShape(args[1:2]).In, parseShape(args[2:]).In))
		case "and":
			return shape.JustIn(shape.AndIn(parseShape(args[1:2]).In, parseShape(args[2:]).In))
		case "not":
			return shape.JustIn(shape.NotIn(parseShape(args[1:]).In))
		}
	}
	return nil
}

func main() {

	debugstream.AddCommand("followMouse", nil, func(args []string) error {
		event.GlobalBind(event.Enter, func(event.CID, interface{}) int {
			// It'd be interesting to attach to the mouse position
			src.SetPos(float64(mouse.LastEvent.X()), float64(mouse.LastEvent.Y()))
			return 0
		})
		return nil
	})

	debugstream.AddCommand("shape", nil, func(args []string) error {
		if len(args) > 0 {
			sh := parseShape(args)
			if sh != nil {
				src.Generator.(pt.Shapeable).SetShape(sh)
			}
		}
		return nil
	})

	debugstream.AddCommand("size", nil, func(args []string) error {
		f1, f2, two, err := parseInts(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.(pt.Sizeable).SetSize(intrange.NewConstant(f1))
		} else {
			src.Generator.(pt.Sizeable).SetSize(intrange.NewLinear(f1, f2))
		}

		return nil
	})

	debugstream.AddCommand("endsize", nil, func(args []string) error {
		f1, f2, two, err := parseInts(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.(pt.Sizeable).SetEndSize(intrange.NewConstant(f1))
		} else {
			src.Generator.(pt.Sizeable).SetEndSize(intrange.NewLinear(f1, f2))
		}
		return nil
	})

	debugstream.AddCommand("count", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.GetBaseGenerator().NewPerFrame = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().NewPerFrame = floatrange.NewLinear(npf, npf2)
		}
		return nil
	})

	debugstream.AddCommand("life", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.GetBaseGenerator().LifeSpan = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().LifeSpan = floatrange.NewLinear(npf, npf2)
		}
		return nil
	})

	debugstream.AddCommand("rotation", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.GetBaseGenerator().Rotation = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Rotation = floatrange.NewLinear(npf, npf2)
		}
		return nil
	})

	debugstream.AddCommand("angle", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.GetBaseGenerator().Angle = floatrange.NewConstant(npf * alg.DegToRad)
		} else {
			src.Generator.GetBaseGenerator().Angle = floatrange.NewLinear(npf*alg.DegToRad, npf2*alg.DegToRad)
		}
		return nil
	})

	debugstream.AddCommand("speed", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			src.Generator.GetBaseGenerator().Speed = floatrange.NewConstant(npf)
		} else {
			src.Generator.GetBaseGenerator().Speed = floatrange.NewLinear(npf, npf2)
		}
		return nil
	})

	debugstream.AddCommand("spread", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "speeds"}
		}
		src.Generator.GetBaseGenerator().Spread.SetPos(npf, npf2)
		return nil
	})

	debugstream.AddCommand("gravity", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "speeds"}
		}
		src.Generator.GetBaseGenerator().Gravity.SetPos(npf, npf2)

		return nil
	})

	debugstream.AddCommand("speeddecay", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "speeds"}
		}
		src.Generator.GetBaseGenerator().SpeedDecay.SetPos(npf, npf2)

		return nil
	})

	debugstream.AddCommand("pos", nil, func(args []string) error {
		npf, npf2, two, err := parseFloats(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		if !two {
			return oakerr.InsufficientInputs{AtLeast: 2, InputName: "positions"}
		}
		src.Generator.SetPos(npf, npf2)

		return nil
	})

	debugstream.AddCommand("startcolor", nil, func(args []string) error {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		startColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		return nil

	})

	debugstream.AddCommand("startrand", nil, func(args []string) error {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		startColorRand = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetStartColor(startColor, startColorRand)
		return nil

	})

	debugstream.AddCommand("endcolor", nil, func(args []string) error {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		endColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetEndColor(endColor, endColorRand)
		return nil
	})

	debugstream.AddCommand("endrand", nil, func(args []string) error {
		if len(args) < 3 {
			return oakerr.InsufficientInputs{AtLeast: 3, InputName: "colorvalues"}
		}
		r, g, b, a, err := parseRGBA(args)
		if err != nil {
			return oakerr.UnsupportedFormat{Format: err.Error()}
		}
		endColorRand = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		src.Generator.(pt.Colorable).SetEndColor(endColor, endColorRand)
		return nil
	})

	oak.AddScene("demo", scene.Scene{Start: func(*scene.Context) {
		render.Draw(render.NewDrawFPS(0, nil, 10, 10))
		x := 320.0
		y := 240.0
		newPf := floatrange.NewLinear(1, 2)
		life := floatrange.NewLinear(100, 120)
		angle := floatrange.NewLinear(0, 360)
		speed := floatrange.NewLinear(1, 5)
		size := intrange.NewConstant(1)
		layerFn := func(v physics.Vector) int {
			return 1
		}
		startColor = color.RGBA{255, 255, 255, 255}
		startColorRand = color.RGBA{0, 0, 0, 0}
		endColor = color.RGBA{255, 255, 255, 255}
		endColorRand = color.RGBA{0, 0, 0, 0}
		shape := shape.Square
		src = pt.NewColorGenerator(
			pt.Pos(x, y),
			pt.Duration(pt.Inf),
			pt.LifeSpan(life),
			pt.Angle(angle),
			pt.Speed(speed),
			pt.Layer(layerFn),
			pt.Shape(shape),
			pt.Size(size),
			pt.Color(startColor, startColorRand, endColor, endColorRand),
			pt.NewPerFrame(newPf)).Generate(0)
	}})

	render.SetDrawStack(
		render.NewCompositeR(),
	)

	err := oak.Init("demo", oak.FileConfig("oak.config"), func(c oak.Config) (oak.Config, error) {
		c.Debug.Level = "VERBOSE"
		c.DrawFrameRate = 1200
		c.FrameRate = 60
		c.EnableDebugConsole = true
		return c, nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func parseRGBA(args []string) (r, g, b, a int, err error) {
	if len(args) < 4 {
		return
	}
	r, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	g, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	b, err = strconv.Atoi(args[2])
	if err != nil {
		return
	}
	a, err = strconv.Atoi(args[3])
	return
}

func parseFloats(args []string) (f1, f2 float64, two bool, err error) {
	if len(args) < 1 {
		err = errors.New("no args")
		return
	}
	f1, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		return
	}
	if len(args) < 2 {
		return
	}
	f2, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return
	}
	two = true
	return
}

func parseInts(args []string) (i1, i2 int, two bool, err error) {
	if len(args) < 1 {
		err = errors.New("No args")
		return
	}
	i1, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if len(args) < 2 {
		return
	}
	i2, err = strconv.Atoi(args[1])
	if err != nil {
		return
	}
	two = true
	return
}
