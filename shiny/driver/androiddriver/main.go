package androiddriver

import (
	"sync"

	"github.com/oakmound/oak/v3/shiny/screen"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

func Main(f func(screen.Screen)) {
	app.Main(func(a app.App) {
		var sz size.Event
		s := &screenImpl{}
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
				sz = e
				//fmt.Println("sending size event", e.HeightPx, e.WidthPx)
				s.Deque.Send(e)
			case paint.Event:
				// drop system driven paint events; oak dictates frame rate
				if s.glctx == nil || e.External {
					continue
				}
				s.glctx.ClearColor(0, 0, 0, 1)
				s.glctx.Clear(gl.COLOR_BUFFER_BIT)
				if s.activeImage != nil && !s.activeImage.dead {
					s.activeImage.img.Upload()
					s.activeImage.img.Draw(sz, geom.Point{}, geom.Point{X: sz.WidthPt}, geom.Point{Y: sz.HeightPt}, sz.Bounds())
				}
				s.Publish()
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				// for now, make this a left click
				// have to worry about this later
				//fmt.Println("sending left click at ", e.X, e.Y)
				s.Deque.Send(mouse.Event{
					X:         e.X,
					Y:         e.Y,
					Button:    mouse.ButtonLeft,
					Direction: mouse.DirPress,
				})
			}
		}
	})
}
