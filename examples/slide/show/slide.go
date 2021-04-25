package show

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

type Slide interface {
	Init()
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
	oak.AddCommand("slide", func(args []string) {
		fmt.Println(args)
		if len(args) < 2 {
			return
		}
		v := args[1]
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(i)
		if i < 0 {
			skipTo = "0"
		} else if i <= max {
			skipTo = v
		} else {
			skipTo = strconv.Itoa(max)
		}
		skip = true
	})
}

func Start(width, height int, slides ...Slide) {
	for i, sl := range slides {
		i := i
		sl := sl
		oak.AddScene("slide"+strconv.Itoa(i), scene.Scene{
			Start: func(*scene.Context) { sl.Init() },
			Loop: func() bool {
				cont := sl.Continue() && !skip
				// This should be disable-able
				if !cont {
					oak.SetLoadingRenderable(render.NewSprite(0, 0, oak.ScreenShot()))
				}
				return cont
			},
			End: func() (string, *scene.Result) {

				if skip {
					skip = false
					return "slide" + skipTo, slideResult(sl)
				}
				if sl.Prev() {
					if i > 0 {
						return "slide" + strconv.Itoa(i-1), slideResult(sl)
					}
					return "slide0", slideResult(sl)
				}
				return "slide" + strconv.Itoa(i+1), slideResult(sl)
			},
		})
	}

	reset := false

	var oldBackground image.Image

	// Todo: customizable end slide
	oak.AddScene("slide"+strconv.Itoa(len(slides)),
		scene.Scene{
			Start: func(ctx *scene.Context) {
				oldBackground = oak.GetBackgroundImage()
				oak.SetColorBackground(image.NewUniform(color.RGBA{0, 0, 0, 255}))
				render.Draw(
					Express.NewStrText(
						"Spacebar to restart show ...",
						float64(ctx.Window.Width()/2),
						float64(ctx.Window.Height()-50),
					),
				)
				event.GlobalBind("KeyDownSpacebar", func(event.CID, interface{}) int {
					reset = true
					return 0
				})
			},
			Loop: func() bool {
				return !reset
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
		c.Screen = oak.Screen{
			Width:  width,
			Height: height,
		}
		c.FrameRate = 30
		c.DrawFrameRate = 30
		return c, nil
	})
}
