package inputviz

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/png"
	"math"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/joystick"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/oakmound/oak/v4/scene"
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

	// BaseLayer is the base layer to render the resulting renderables to
	// if -1, it will render only to the layer provided to RenderAndListen.
	BaseLayer int

	ctx *scene.Context
	event.CallerID
	joy          *joystick.Joystick
	rs           map[string]render.Modifiable
	lastState    *joystick.State
	triggerY     float64
	lStickCenter floatgeom.Point2
	rStickCenter floatgeom.Point2
	cancel       func()

	bindings []event.Binding
}

func (j *Joystick) CID() event.CallerID {
	return j.CallerID
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
	j.joy = joy
	j.rs = make(map[string]render.Modifiable)
	j.lastState = &joystick.State{}
	j.ctx = ctx
	j.CallerID = ctx.CallerMap.Register(j)
	j.rs["Outline"] = outline
	j.rs["LtStick"] = render.NewCircle(color.RGBA{255, 255, 255, 255}, 15, 12)
	j.rs["RtStick"] = render.NewCircle(color.RGBA{255, 255, 255, 255}, 15, 12)
	j.rs[string(joystick.InputUp)] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["Down"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["Left"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["Right"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["Back"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["Start"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["X"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["Y"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["A"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["B"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["LeftShoulder"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	j.rs["RightShoulder"] = render.NewReverting(render.NewColorBox(10, 10, color.RGBA{255, 255, 255, 255}))
	// Draw the triggers behind the outline to simulate pressing down
	j.rs["LtTrigger"] = render.NewColorBox(40, 30, color.RGBA{255, 255, 255, 255})
	j.rs["RtTrigger"] = render.NewColorBox(40, 30, color.RGBA{255, 255, 255, 255})

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

	for k, r := range j.rs {
		offset := offsets[k]
		r.SetPos(offset.X(), offset.Y())
		r.ShiftX(j.Rect.Min.X())
		r.ShiftY(j.Rect.Min.Y())
		var l int
		if k == "LtTrigger" || k == "RtTrigger" {
			l = layer
			j.triggerY = r.Y()
		} else if k == "Outline" {
			l = layer + 1
		} else {
			l = layer + 2
		}
		if j.BaseLayer == -1 {
			ctx.DrawStack.Draw(r, l)
		} else {
			ctx.DrawStack.Draw(r, j.BaseLayer, l)
		}
	}

	j.lStickCenter = floatgeom.Point2{j.rs["LtStick"].X(), j.rs["LtStick"].Y()}
	j.rStickCenter = floatgeom.Point2{j.rs["RtStick"].X(), j.rs["RtStick"].Y()}

	opts := &joystick.ListenOptions{
		JoystickChanges: true,
		StickChanges:    true,
		StickDeadzoneLX: j.StickDeadzone,
		StickDeadzoneLY: j.StickDeadzone,
		StickDeadzoneRX: j.StickDeadzone,
		StickDeadzoneRY: j.StickDeadzone,
	}
	j.cancel = joy.Listen(opts)

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

	b1 := event.Bind(ctx, joystick.Disconnected, j, func(rend *Joystick, _ uint32) event.Response {
		j.Destroy()
		return 0
	})

	// TODO: it is bad that you need to import two 'key' packages
	b2 := event.Bind(ctx, key.Down(key.Spacebar), j, func(j *Joystick, _ key.Event) event.Response {
		j.joy.Vibrate(math.MaxUint16, math.MaxUint16)
		go func() {
			time.Sleep(1 * time.Second)
			j.joy.Vibrate(0, 0)
		}()
		return 0
	})

	b3 := event.Bind(ctx, joystick.Change, j, func(j *Joystick, st *joystick.State) event.Response {
		for _, inputB := range bts {
			b := string(inputB)
			r := j.rs[b]
			if st.Buttons[b] && !j.lastState.Buttons[b] {
				r.Filter(mod.Brighten(-40))
			} else if !st.Buttons[b] && j.lastState.Buttons[b] {
				r.(*render.Reverting).Revert(1)
			}
		}
		j.lastState = st

		tgr := "LtTrigger"
		x := j.rs[tgr].X()
		j.rs[tgr].SetPos(x, j.triggerY+float64(st.TriggerL/16))

		tgr = "RtTrigger"
		x = j.rs[tgr].X()
		j.rs[tgr].SetPos(x, j.triggerY+float64(st.TriggerR/16))
		return 0
	})

	b4 := event.Bind(ctx, joystick.LtStickChange, j, func(j *Joystick, st *joystick.State) event.Response {
		pos := j.lStickCenter
		pos = pos.Add(floatgeom.Point2{
			float64(st.StickLX / 2048),
			-float64(st.StickLY / 2048),
		})
		j.rs["LtStick"].SetPos(pos.X(), pos.Y())
		return 0
	})

	b5 := event.Bind(ctx, joystick.RtStickChange, j, func(j *Joystick, st *joystick.State) event.Response {
		pos := j.rStickCenter
		pos = pos.Add(floatgeom.Point2{
			float64(st.StickRX / 2048),
			-float64(st.StickRY / 2048),
		})
		j.rs["RtStick"].SetPos(pos.X(), pos.Y())
		return 0
	})
	j.bindings = []event.Binding{b1, b2, b3, b4, b5}
	return nil
}

func (j *Joystick) Destroy() {
	for _, b := range j.bindings {
		b.Unbind()
	}
	for _, r := range j.rs {
		r.Undraw()
	}
	j.cancel()
}
