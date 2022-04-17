package show

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/debugstream"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

type Slide interface {
	Init(*scene.Context)
	Continue() bool
	Prev() bool
	Transition() scene.Transition
}

func slideResult(sl Slide) *scene.Result {
	return &scene.Result{
		Transition: sl.Transition(),
	}
}

var (
	skip   bool
	skipTo string
)

func AddNumberShortcuts(max int) {
	debugstream.AddCommand(debugstream.Command{Name: "slide", Operation: func(args []string) string {
		if len(args) < 2 {
			return ""
		}
		v := args[1]
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		if i < 0 {
			skipTo = "0"
		} else if i <= max {
			skipTo = v
		} else {
			skipTo = strconv.Itoa(max)
		}
		skip = true
		return ""
	}})
}

func Start(width, height int, slides ...Slide) {
	for i, sl := range slides {
		i := i
		sl := sl
		oak.AddScene("slide"+strconv.Itoa(i), scene.Scene{
			Start: func(ctx *scene.Context) {

				sl.Init(ctx)
				event.GlobalBind(ctx, event.Enter, func(event.EnterPayload) event.Response {
					cont := sl.Continue() && !skip
					oak.SetLoadingRenderable(render.NewSprite(0, 0, oak.ScreenShot()))
					if !cont {
						ctx.Window.NextScene()
						return event.ResponseUnbindThisBinding
					}
					return 0
				})

			},
			End: func() (string, *scene.Result) {
				fmt.Println("ending")
				if skip {
					skip = false
					return "slide" + skipTo, slideResult(sl)
				}
				if sl.Prev() {
					fmt.Println("Prev slide requested from", i)
					if i > 0 {
						return "slide" + strconv.Itoa(i-1), slideResult(sl)
					}
					return "slide0", slideResult(sl)
				}
				fmt.Println("new slide", strconv.Itoa(i+1))
				return "slide" + strconv.Itoa(i+1), slideResult(sl)
			},
		})
	}

	reset := false

	var oldBackground image.Image

	oak.AddScene("slide"+strconv.Itoa(len(slides)),
		scene.Scene{
			Start: func(ctx *scene.Context) {
				oldBackground = ctx.Window.(*oak.Window).GetBackgroundImage()
				oak.SetColorBackground(image.NewUniform(color.RGBA{0, 0, 0, 255}))
				wbds := ctx.Window.Bounds()
				render.Draw(
					Express.NewText(
						"Spacebar to restart show ...",
						float64(wbds.X()/2),
						float64(wbds.Y()-50),
					),
				)
				event.GlobalBind(ctx, key.Down(key.Spacebar), func(key.Event) event.Response {
					reset = true
					return 0
				})

				event.GlobalBind(ctx, event.Enter, func(event.EnterPayload) event.Response {
					if !reset {
						ctx.Window.NextScene()
						return event.ResponseUnbindThisBinding
					}
					return 0
				})
			},
			End: func() (string, *scene.Result) {
				oak.SetColorBackground(oldBackground)
				reset = false
				skip = false
				return "slide0", nil
			},
		},
	)
	oak.Init("slide0", func(c oak.Config) (oak.Config, error) {
		c.Screen.Width = width
		c.Screen.Height = height
		c.FrameRate = 30
		c.DrawFrameRate = 30
		c.EnableDebugConsole = true
		return c, nil
	})
}
