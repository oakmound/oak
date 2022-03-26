package main

import (
	"fmt"
	"time"

	"github.com/oakmound/oak/v3/debugtools/inputviz"
	"github.com/oakmound/oak/v3/render"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/event"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/joystick"
	"github.com/oakmound/oak/v3/scene"
)

func main() {
	oak.AddScene("viz", scene.Scene{Start: func(ctx *scene.Context) {
		joystick.Init()
		latestInput := new(string)
		*latestInput = "Latest Input: Keyboard+Mouse"
		ctx.DrawStack.Draw(render.NewStrPtrText(latestInput, 10, 460), 4)
		ctx.DrawStack.Draw(render.NewText("Space to Vibrate", 10, 440), 4)
		ctx.Handler.GlobalBind(event.InputChange, func(_ event.CallerID, payload interface{}) int {
			input := payload.(oak.InputType)
			switch input {
			case oak.InputJoystick:
				*latestInput = "Latest Input: Joystick"
			case oak.InputKeyboardMouse:
				*latestInput = "Latest Input: Keyboard+Mouse"
			}
			return 0
		})
		go func() {
			rWidth := float64(ctx.Window.Width()) / 2
			rHeight := float64(ctx.Window.Height()) / 2
			jCh, cancel := joystick.WaitForJoysticks(1 * time.Second)
			defer cancel()
			for joy := range jCh {
				fmt.Println("new joystick", joy.ID())
				var x, y float64
				switch joy.ID() {
				case 0:
					// 0,0
				case 1:
					x = rWidth
				case 2:
					y = rHeight
				case 3:
					x = rWidth
					y = rHeight
				}
				jrend := inputviz.Joystick{
					Rect:          floatgeom.NewRect2WH(x, y, rWidth, rHeight),
					StickDeadzone: 4000,
					BaseLayer:     -1,
				}
				err := jrend.RenderAndListen(ctx, joy, 1)
				if err != nil {
					fmt.Println("renderer:", err)
				}
			}
		}()
	}})
	oak.Init("viz", func(c oak.Config) (oak.Config, error) {
		c.TrackInputChanges = true
		return c, nil
	})
}
