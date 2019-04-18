package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/event"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/scene"
)

type renderer struct {
	event.CID
	joy *joystick.Joystick
	rs  map[string]render.Modifiable
}

func (r *renderer) Init() event.CID {
	r.CID = event.NextID(r)
	return r.CID
}

var initialOffsets = map[string]floatgeom.Point2{
	"Outline":   floatgeom.Point2{0, 0},
	"LtStick":   floatgeom.Point2{50, 75},
	"RtStick":   floatgeom.Point2{210, 115},
	"Up":        floatgeom.Point2{90, 115},
	"Left":      floatgeom.Point2{80, 125},
	"Right":     floatgeom.Point2{100, 125},
	"Down":      floatgeom.Point2{90, 135},
	"Back":      floatgeom.Point2{130, 85},
	"Start":     floatgeom.Point2{190, 85},
	"X":         floatgeom.Point2{240, 85},
	"Y":         floatgeom.Point2{250, 75},
	"A":         floatgeom.Point2{250, 95},
	"B":         floatgeom.Point2{260, 85},
	"LB":        floatgeom.Point2{60, 40},
	"RB":        floatgeom.Point2{260, 40},
	"LtTrigger": floatgeom.Point2{40, 30},
	"RtTrigger": floatgeom.Point2{240, 30},
}

func newRenderer(joy *joystick.Joystick) {
	outline, err := render.LoadSprite("", "controllerOutline.png")
	if err != nil {
		dlog.Error(err)
		return
	}
	rend := &renderer{
		joy: joy,
		rs:  make(map[string]render.Modifiable),
	}
	rend.Init()
	rend.rs["Outline"] = outline
	rend.rs["LtStick"] = render.NewCircle(color.RGBA{255, 255, 255, 255}, 15, 12)
	rend.rs["RtStick"] = render.NewCircle(color.RGBA{255, 255, 255, 255}, 15, 12)
	rend.rs["Up"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["Down"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["Left"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["Right"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["Back"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["Start"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["X"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["Y"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["A"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["B"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["LB"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["RB"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	// Draw the triggers behind the outline to simulate pressing down
	rend.rs["LtTrigger"] = render.NewColorBox(40, 30, color.RGBA{255, 255, 255, 255})
	rend.rs["RtTrigger"] = render.NewColorBox(40, 30, color.RGBA{255, 255, 255, 255})

	for k, r := range rend.rs {
		offset := initialOffsets[k]
		r.SetPos(offset.X(), offset.Y())
		switch joy.ID() {
		case 0:
		case 1:
			r.ShiftX(float64(oak.ScreenWidth / 2))
		case 2:
			r.ShiftY(float64(oak.ScreenHeight / 2))
		case 3:
			r.ShiftX(float64(oak.ScreenWidth / 2))
			r.ShiftY(float64(oak.ScreenHeight / 2))
		}
		if k == "LtTrigger" || k == "RtTrigger" {
			render.Draw(r, 0)
		} else if k == "Outline" {
			render.Draw(r, 1)
		} else {
			render.Draw(r, 2)
		}
	}

	joy.Handler = rend
	joy.Listen(nil)

	rend.Bind(func(id int, _ interface{}) int {
		rend, ok := event.GetEntity(id).(*renderer)
		if !ok {
			return 0
		}
		fmt.Println("Got disconnected event")
		for _, r := range rend.rs {
			r.Undraw()
		}
		return 0
	}, joystick.Disconnected)
}

func main() {
	oak.Add("viz", func(string, interface{}) {
		// Todo: Handle disconnection
		joystick.Init()
		go func() {
			jCh, cancel := joystick.WaitForJoysticks(1 * time.Second)
			defer cancel()
			for {
				select {
				case joy := <-jCh:
					newRenderer(joy)
				}
			}
		}()
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "viz", nil
	})
	oak.SetupConfig.Assets.ImagePath = "."
	oak.SetupConfig.Assets.AssetPath = "."
	oak.Init("viz")
}
