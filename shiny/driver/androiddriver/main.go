//go:build android
// +build android

package androiddriver

import (
	"sync"

	"github.com/oakmound/oak/v4/shiny/screen"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

func Main(f func(screen.Screen)) {
	app.Main(func(a app.App) {
		s := &Screen{
			app: a,
		}
		screenOnce := sync.Once{}
		for e := range a.Events() {
			switch e := e.(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					s.glctx, _ = e.DrawContext.(gl.Context)
					s.images = glutil.NewImages(s.glctx)
					screenOnce.Do(func() {
						go f(s)
					})
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					s.glctx = nil
					for _, img := range s.activeImages {
						img.Release()
					}
					s.images.Release()
				}
			case size.Event:
				s.lastSz = e
				s.Deque.Send(e)
			case touch.Event:
				// TODO: expose touch events in a way an oak program can
				// differentiate them from clicks
				switch e.Type {
				case touch.TypeBegin:
					s.Deque.Send(mouse.Event{
						X:         e.X,
						Y:         e.Y,
						Button:    mouse.ButtonLeft,
						Direction: mouse.DirPress,
					})
				case touch.TypeEnd:
					s.Deque.Send(mouse.Event{
						X:         e.X,
						Y:         e.Y,
						Button:    mouse.ButtonLeft,
						Direction: mouse.DirRelease,
					})
				case touch.TypeMove:
					s.Deque.Send(mouse.Event{
						X:         e.X,
						Y:         e.Y,
						Button:    mouse.ButtonLeft,
						Direction: mouse.DirNone,
					})
				}
			}
		}
	})
}
