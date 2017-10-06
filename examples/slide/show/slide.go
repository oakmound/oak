package show

import (
	"fmt"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
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

func AddSlides(slides ...Slide) {
	for i, sl := range slides {
		i := i
		sl := sl
		fmt.Println("slide" + strconv.Itoa(i))
		oak.AddScene("slide"+strconv.Itoa(i), scene.Scene{
			Start: func(string, interface{}) { sl.Init() },
			Loop:  sl.Continue,
			End: func() (string, *scene.Result) {
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

	// Todo: customizable end slide
	oak.AddScene("slide"+strconv.Itoa(len(slides)),
		scene.Scene{
			Start: func(string, interface{}) {
				render.Draw(
					Express.NewStrText(
						"Spacebar to restart show ...",
						float64(oak.ScreenWidth/2),
						float64(oak.ScreenHeight-50),
					),
				)
				event.GlobalBind(func(int, interface{}) int {
					reset = true
					return 0
				}, "KeyDownSpacebar")
			},
			Loop: func() bool {
				return !reset
			},
			End: func() (string, *scene.Result) {
				reset = false
				return "slide0", nil
			},
		},
	)
}

func Start() {
	oak.Init("slide0")
}
