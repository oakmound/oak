package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/oakmound/oak/v2/key"

	"github.com/oakmound/oak/v2/render/mod"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/render"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/event"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/joystick"
	"github.com/oakmound/oak/v2/scene"
)

type renderer struct {
	event.CID
	joy          *joystick.Joystick
	rs           map[string]render.Modifiable
	lastState    *joystick.State
	triggerY     float64
	lStickCenter floatgeom.Point2
	rStickCenter floatgeom.Point2
}

// Init sets CID on a renderable
func (r *renderer) Init() event.CID {
	r.CID = event.NextID(r)
	return r.CID
}

var initialOffsets = map[string]floatgeom.Point2{
	"Outline":       {0, 0},
	"LtStick":       {50, 75},
	"RtStick":       {210, 115},
	"Up":            {90, 115},
	"Left":          {80, 125},
	"Right":         {100, 125},
	"Down":          {90, 135},
	"Back":          {130, 85},
	"Start":         {190, 85},
	"X":             {240, 85},
	"Y":             {250, 75},
	"A":             {250, 95},
	"B":             {260, 85},
	"LeftShoulder":  {60, 40},
	"RightShoulder": {260, 40},
	"LtTrigger":     {40, 6},
	"RtTrigger":     {240, 6},
}

func newRenderer(joy *joystick.Joystick) {
	outline, err := render.LoadSprite("", "controllerOutline.png")
	if err != nil {
		dlog.Error(err)
		return
	}
	rend := &renderer{
		joy:       joy,
		rs:        make(map[string]render.Modifiable),
		lastState: &joystick.State{},
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
	rend.rs["LeftShoulder"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	rend.rs["RightShoulder"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
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
			rend.triggerY = r.Y()
		} else if k == "Outline" {
			render.Draw(r, 1)
		} else {
			render.Draw(r, 2)
		}
	}

	rend.lStickCenter = floatgeom.Point2{rend.rs["LtStick"].X(), rend.rs["LtStick"].Y()}
	rend.rStickCenter = floatgeom.Point2{rend.rs["RtStick"].X(), rend.rs["RtStick"].Y()}

	joy.Handler = rend
	joy.Listen(nil)

	bts := []string{
		"X",
		"A",
		"Y",
		"B",
		"Up",
		"Down",
		"Left",
		"Right",
		"Back",
		"Start",
		"LeftShoulder",
		"RightShoulder",
	}

	rend.Bind(func(id int, _ interface{}) int {
		rend, ok := event.GetEntity(id).(*renderer)
		if !ok {
			return 0
		}
		for _, r := range rend.rs {
			r.Undraw()
		}
		return 0
	}, joystick.Disconnected)

	rend.Bind(func(id int, _ interface{}) int {
		rend, ok := event.GetEntity(id).(*renderer)
		if !ok {
			return 0
		}
		fmt.Println("Triggering vibration")
		rend.joy.Vibrate(math.MaxUint16, math.MaxUint16)
		go func() {
			time.Sleep(1 * time.Second)
			rend.joy.Vibrate(0, 0)
		}()
		return 0
	}, key.Down+key.Spacebar)

	rend.Bind(func(id int, state interface{}) int {
		rend, ok := event.GetEntity(id).(*renderer)
		if !ok {
			return 0
		}
		st, ok := state.(*joystick.State)
		if !ok {
			return 0
		}
		for _, b := range bts {
			r := rend.rs[b]
			if st.Buttons[b] && !rend.lastState.Buttons[b] {
				r.Filter(mod.Brighten(-40))
			} else if !st.Buttons[b] && rend.lastState.Buttons[b] {
				r.(*render.Reverting).Revert(1)
			}
		}
		rend.lastState = st

		tgr := "LtTrigger"
		x := rend.rs[tgr].X()
		rend.rs[tgr].SetPos(x, rend.triggerY+float64(st.TriggerL/16))

		tgr = "RtTrigger"
		x = rend.rs[tgr].X()
		rend.rs[tgr].SetPos(x, rend.triggerY+float64(st.TriggerR/16))

		pos := rend.lStickCenter
		pos = pos.Add(floatgeom.Point2{
			float64(st.StickLX / 2048),
			-float64(st.StickLY / 2048),
		})
		rend.rs["LtStick"].SetPos(pos.X(), pos.Y())

		pos = rend.rStickCenter
		pos = pos.Add(floatgeom.Point2{
			float64(st.StickRX / 2048),
			-float64(st.StickRY / 2048),
		})
		rend.rs["RtStick"].SetPos(pos.X(), pos.Y())

		return 0
	}, joystick.Change)
}

func main() {
	oak.Add("viz", func(string, interface{}) {
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
