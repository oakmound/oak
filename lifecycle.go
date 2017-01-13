// Package oak is a game engine...
package oak

import (
	"fmt"
	"image"
	"image/draw"
	"strconv"
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
	pmouse "bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func lifecycleLoop(s screen.Screen) {

	// The event loop requires information about
	// the size of the world and screen that is
	// being dealt with, and so initializes it here.
	//
	// Todo: add world size to config
	worldBuffer, _ = s.NewBuffer(image.Point{4000, 4000})

	winBuffer, _ = s.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	w, err := s.NewWindow(&screen.NewWindowOptions{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
	}
	defer w.Release()

	sscreen = s

	eb = event.GetEventBus()

	// Todo: add frame rate to config
	frameRate := 60
	frameCh := make(chan bool)

	// This goroutine maintains a logical framerate
	go func(frameCh chan bool, frameRate int64) {
		c := time.Tick(time.Second / time.Duration(frameRate))
		for range c {
			frameCh <- true
		}
	}(frameCh, int64(frameRate))

	// Native go event handler
	go func() {
		for {
			e := w.NextEvent()
			switch e := e.(type) {

			// We only currently respond to death lifecycle events.
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					quitCh <- true
					return
				}

			// Send key events
			//
			// Key events have two varieties:
			// The "KeyDown" and "KeyUp" events, which trigger for all keys
			// and specific "KeyDown$key", etc events which trigger only for $key.
			// The specific key that is pressed is passed as the data interface for
			// the former events, but not for the latter.
			case key.Event:
				k := GetKeyBind(e.Code.String()[4:])
				if e.Direction == key.DirPress {
					fmt.Println("--------------------", e.Code.String()[4:], k)
					setDown(k)
					eb.Trigger("KeyDown", k)
					eb.Trigger("KeyDown"+k, nil)
				} else if e.Direction == key.DirRelease {
					setUp(k)
					eb.Trigger("KeyUp", k)
					eb.Trigger("KeyUp"+k, nil)
				}

			// Send mouse events
			//
			// Mouse events are parsed based on their button
			// and direction into an event name and then triggered:
			// 'MousePress', 'MouseRelease', 'MouseScrollDown', 'MouseScrollUp', and 'MouseDrag'
			//
			// The basic event name is meant for entities which
			// want to respond to the mouse event happening -anywhere-.
			//
			// For events which have mouse collision enabled, they'll recieve
			// $eventName+"On" when the event occurs within their collision area.
			//
			// Mouse events all recieve an x, y, and button string.
			case mouse.Event:
				button := pmouse.GetMouseButton(int32(e.Button))
				dlog.Verb("Mouse direction ", e.Direction.String(), " Button ", button)
				mevent := pmouse.MouseEvent{e.X, e.Y, button}
				var eventName string
				if e.Direction == mouse.DirPress {
					setDown(button)
					eventName = "MousePress"
				} else if e.Direction == mouse.DirRelease {
					setUp(button)
					eventName = "MouseRelease"
				} else if e.Button == -2 {
					eventName = "MouseScrollDown"
				} else if e.Button == -1 {
					eventName = "MouseScrollUp"
				} else {
					eventName = "MouseDrag"
				}
				pmouse.LastMouseEvent = mevent
				eb.Trigger(eventName, mevent)
				pmouse.Propagate(eventName+"On", mevent)

			// I don't really know what a paint event is to be honest.
			// We hypothetically don't allow the user to manually resize
			// their window, so we don't do anything special for such events.
			case size.Event, paint.Event:
			case error:
				dlog.Error(e)
			}

			// This is a hardcoded quit function bound to the escape key.
			if IsDown("Escape") {
				if esc {
					dlog.Warn("Quiting oak from holding ESCAPE")
					w.Send(lifecycle.Event{0, 0, nil})
				}
				esc = true
			} else {
				esc = false
			}
		}
	}()

	// This sends a signal to initiate the first scene
	initCh <- true

	// The draw loop
	// Unless told to stop, the draw channel will repeatedly
	// 1. draw black to a temporary buffer
	// 2. run any functions bound to precede drawing.
	// 3. draw all elements onto the temporary buffer.
	// 4. run any functions bound to follow drawing.
	// 5. draw the buffer's data at the viewport's position to the screen.
	// 6. publish the screen to display in window.
	go func() {
		<-drawChannel
		//cb := render.CompositeFilter(render.NewColorBox(4096, 4096, color.RGBA{0, 0, 0, 125}).Sprite)
		lastTime := time.Now()

		text := render.DefFont().NewText("", 10, 20)
		render.StaticDraw(text, 60000)
		for {
			dlog.Verb("Draw Loop")
		drawSelect:
			select {

			case <-drawChannel:
				dlog.Verb("Got something from draw channel")
				for {
					select {
					case <-drawChannel:
						render.StaticDraw(text, 60000)
						break drawSelect
					case viewPoint := <-viewportChannel:
						dlog.Verb("Got something from viewport channel (waiting on draw)")
						updateScreen(viewPoint[0], viewPoint[1])
					}

				}
			case viewPoint := <-viewportChannel:
				dlog.Verb("Got something from viewport channel")
				updateScreen(viewPoint[0], viewPoint[1])
			default:

				eb = event.GetEventBus()

				draw.Draw(worldBuffer.RGBA(), worldBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)

				render.PreDraw()
				render.DrawHeap(worldBuffer, ViewX, ViewY, ScreenWidth, ScreenHeight)
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), worldBuffer.RGBA(), image.Point{ViewX, ViewY}, screen.Src)
				render.DrawStaticHeap(winBuffer)

				w.Upload(image.Point{0, 0}, winBuffer, winBuffer.Bounds())
				w.Publish()

				timeSince := 1000000000.0 / float64(time.Since(lastTime).Nanoseconds())
				text.SetText(strconv.Itoa(int(timeSince)))

				timeSince = 1000000000.0 / float64(time.Since(lastTime).Nanoseconds())
				text.SetText(strconv.Itoa(int(timeSince)))

				lastTime = time.Now()
			}
		}
	}()

	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then runs any functions bound to when a frame ends.
	// It then allows a scene to perform it's loop operation.
	go func() {
		for runEventLoop {
			event.ResolvePending()
		}
	}()
	for {
		for runEventLoop {
			<-frameCh
			<-eb.Trigger("EnterFrame", nil)
			<-eb.Trigger("ExitFrame", nil)
			sceneCh <- true
		}
	}
}

func CurrentScene() string {
	return scene
}

func GetScreen() draw.Image {
	return winBuffer.RGBA()
}

func GetWorld() draw.Image {
	return worldBuffer.RGBA()
}

func SetWorldSize(x, y int) {
	worldBuffer, _ = sscreen.NewBuffer(image.Point{x, y})
}
