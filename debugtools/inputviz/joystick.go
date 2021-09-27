package inputviz

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/joystick"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/scene"
)

//go:embed controllerOutline.png
var controllerOutline []byte

var pngOutline image.Image

func init() {
	var err error
	pngOutline, err = png.Decode(bytes.NewBuffer(controllerOutline))
	if err != nil {
		dlog.Error("failed to decode background data: %w", err)
	}
}

// Joystick visualizes the inputs sent to a controller
type Joystick struct {
	// Rect is the rect this joystick should be drawn to.
	// Defaults to (0,0)->(320,240)
	Rect floatgeom.Rect2

	// StickDeadzone is the lowest value of stick movements that should
	// be rendered.
	StickDeadzone int16

	ctx *scene.Context
	event.CID
	joy          *joystick.Joystick
	rs           map[string]render.Modifiable
	lastState    *joystick.State
	triggerY     float64
	lStickCenter floatgeom.Point2
	rStickCenter floatgeom.Point2
	cancel       func()
}

func (r *Joystick) Init() event.CID {
	r.CID = event.NextID(r)
	return r.CID
}

func (j *Joystick) RenderAndListen(ctx *scene.Context, joy *joystick.Joystick, layer int) error {
	bounds := pngOutline.Bounds()
	rgba := image.NewRGBA(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba.Set(x, y, color.RGBAModel.Convert(pngOutline.At(x, y)))
		}
	}

	outline := render.NewSprite(0, 0, rgba)
	rend := &Joystick{
		joy:       joy,
		rs:        make(map[string]render.Modifiable),
		lastState: &joystick.State{},
		ctx:       ctx,
	}
	rend.Init()
	rend.rs["Outline"] = outline
	rend.rs["LtStick"] = render.NewCircle(color.RGBA{255, 255, 255, 255}, 15, 12)
	rend.rs["RtStick"] = render.NewCircle(color.RGBA{255, 255, 255, 255}, 15, 12)
	rend.rs[string(joystick.InputUp)] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
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

	var offsets = map[string]floatgeom.Point2{
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
	const defaultW = 320
	const defaultH = 240
	if j.Rect == (floatgeom.Rect2{}) {
		j.Rect = floatgeom.Rect2{
			Max: floatgeom.Point2{defaultW, defaultH},
		}
	} else {
		// adjust all offsets
		for k, v := range offsets {
			offsets[k] = floatgeom.Point2{
				v.X() * (j.Rect.W() / defaultW),
				v.Y() * (j.Rect.H() / defaultH),
			}
		}
	}

	for k, r := range rend.rs {
		offset := offsets[k]
		r.SetPos(offset.X(), offset.Y())
		r.ShiftX(j.Rect.Min.X())
		r.ShiftY(j.Rect.Min.Y())
		if k == "LtTrigger" || k == "RtTrigger" {
			ctx.DrawStack.Draw(r, layer)
			rend.triggerY = r.Y()
		} else if k == "Outline" {
			ctx.DrawStack.Draw(r, layer+1)
		} else {
			ctx.DrawStack.Draw(r, layer+2)
		}
	}

	rend.lStickCenter = floatgeom.Point2{rend.rs["LtStick"].X(), rend.rs["LtStick"].Y()}
	rend.rStickCenter = floatgeom.Point2{rend.rs["RtStick"].X(), rend.rs["RtStick"].Y()}

	joy.Handler = rend
	opts := &joystick.ListenOptions{
		JoystickChanges: true,
		StickChanges:    true,
		StickDeadzoneLX: j.StickDeadzone,
		StickDeadzoneLY: j.StickDeadzone,
		StickDeadzoneRX: j.StickDeadzone,
		StickDeadzoneRY: j.StickDeadzone,
	}
	rend.cancel = joy.Listen(opts)

	bts := []joystick.Input{
		joystick.InputA,
		joystick.InputB,
		joystick.InputX,
		joystick.InputY,
		joystick.InputUp,
		joystick.InputDown,
		joystick.InputLeft,
		joystick.InputRight,
		joystick.InputBack,
		joystick.InputStart,
		joystick.InputLeftShoulder,
		joystick.InputRightShoulder,
	}

	rend.CheckedIDBind(joystick.Disconnected, func(rend *Joystick, _ uint32) {
		fmt.Println("destroying")
		rend.Destroy()
	})

	rend.CheckedBind(key.Down+key.Spacebar, func(rend *Joystick, st *joystick.State) {
		rend.joy.Vibrate(math.MaxUint16, math.MaxUint16)
		go func() {
			time.Sleep(1 * time.Second)
			rend.joy.Vibrate(0, 0)
		}()
	})

	rend.CheckedBind(joystick.Change, func(rend *Joystick, st *joystick.State) {
		for _, inputB := range bts {
			b := string(inputB)
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
	})

	rend.CheckedBind(joystick.LtStickChange, func(rend *Joystick, st *joystick.State) {
		pos := rend.lStickCenter
		pos = pos.Add(floatgeom.Point2{
			float64(st.StickLX / 2048),
			-float64(st.StickLY / 2048),
		})
		rend.rs["LtStick"].SetPos(pos.X(), pos.Y())
	})

	rend.CheckedBind(joystick.RtStickChange, func(rend *Joystick, st *joystick.State) {
		pos := rend.rStickCenter
		pos = pos.Add(floatgeom.Point2{
			float64(st.StickRX / 2048),
			-float64(st.StickRY / 2048),
		})
		rend.rs["RtStick"].SetPos(pos.X(), pos.Y())
	})
	return nil
}

func (j *Joystick) CheckedIDBind(ev string, f func(*Joystick, uint32)) {
	j.Bind(ev, func(id event.CID, jid interface{}) int {
		joy, ok := event.GetEntity(id).(*Joystick)
		if !ok {
			fmt.Println("checked bind failed")
			return 0
		}
		n, ok := jid.(uint32)
		if !ok {
			return 0
		}
		f(joy, n)
		return 0
	})
}

func (j *Joystick) CheckedBind(ev string, f func(*Joystick, *joystick.State)) {
	j.Bind(ev, func(id event.CID, state interface{}) int {
		joy, ok := event.GetEntity(id).(*Joystick)
		if !ok {
			fmt.Println("checked bind failed")
			return 0
		}
		st, ok := state.(*joystick.State)
		if !ok {
			return 0
		}
		f(joy, st)
		return 0
	})
}

func (j *Joystick) Destroy() {
	j.UnbindAll()
	for _, r := range j.rs {
		r.Undraw()
	}
	j.cancel()
}
