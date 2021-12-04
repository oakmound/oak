package oak

import (
	"context"
	"image"
	"io"
	"sort"
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/debugstream"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/shiny/driver"
	"github.com/oakmound/oak/v3/shiny/screen"
	"github.com/oakmound/oak/v3/window"
)

var _ window.Window = &Window{}

func (w *Window) windowController(s screen.Screen, x, y int32, width, height int) (screen.Window, error) {
	return s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title(w.config.Title),
		screen.Position(x, y),
		screen.Fullscreen(w.config.Fullscreen),
		screen.Borderless(w.config.Borderless),
		screen.TopMost(w.config.TopMost),
	))
}

type Window struct {
	key.State

	// TODO: most of these channels are not closed cleanly
	transitionCh chan struct{}

	// The Scene channel receives a signal
	// when a scene's .loop() function should
	// be called.
	sceneCh chan struct{}

	// The skip scene channel receives a debug
	// signal to forcibly go to the next
	// scene.
	skipSceneCh chan string

	// The quit channel receives a signal when
	// oak should stop active workers and return from Init.
	quitCh chan struct{}

	// The draw channel receives a signal when
	// drawing should cease (or resume)
	drawCh chan struct{}

	betweenDrawCh chan func()

	// ScreenWidth is the width of the screen
	ScreenWidth int
	// ScreenHeight is the height of the screen
	ScreenHeight int

	// FrameRate is the current logical frame rate.
	// Changing this won't directly effect frame rate, that
	// requires changing the LogicTicker, but it will take
	// effect next scene
	FrameRate int

	// DrawFrameRate is the equivalent to FrameRate for
	// the rate at which the screen is drawn.
	DrawFrameRate int

	// IdleDrawFrameRate is how often the screen will be redrawn
	// when the window is out of focus.
	IdleDrawFrameRate int

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	winBuffers     [2]screen.Image
	screenControl  screen.Screen
	windowControl  screen.Window
	windowTextures [2]screen.Texture
	bufferIdx      uint8

	windowRect image.Rectangle

	// DrawTicker is the parallel to LogicTicker to set the draw framerate
	DrawTicker *time.Ticker
	// animationFrame is used by the javascript driver instead of DrawTicker
	animationFrame chan struct{}

	bkgFn func() image.Image

	// SceneMap is a global map of scenes referred to when scenes advance to
	// determine what the next scene should be.
	// It can be replaced or modified so long as these modifications happen
	// during a scene or before the controller has started.
	SceneMap *scene.Map

	// viewPos represents the point in the world which the viewport is anchored at.
	viewPos    intgeom.Point2
	viewBounds intgeom.Rect2

	aspectRatio float64

	// Driver is the driver oak will call during initialization
	Driver Driver

	// prePublish is a function called each draw frame prior to
	prePublish func(w *Window, tx screen.Texture)

	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable

	firstScene string
	// ErrorScene is a scene string that will be entered if the scene handler
	// fails to enter some other scene, for example, because it's name was
	// undefined in the scene map. If the scene map does not have ErrorScene
	// as well, it will fall back to panicking.
	ErrorScene string

	eventHandler  event.Handler
	CallerMap     *event.CallerMap
	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
	DrawStack     *render.DrawStack

	// LastMouseEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastMouseEvent         mouse.Event
	LastRelativeMouseEvent mouse.Event
	lastRelativePress      mouse.Event
	// LastPress is the last triggered mouse event,
	// where the mouse event was a press.
	// If TrackMouseClicks is set to false then this will not be tracked
	LastMousePress mouse.Event

	FirstSceneInput interface{}

	commands map[string]func([]string)

	ControllerID int32

	config Config

	mostRecentInput InputType

	exitError     error
	ParentContext context.Context

	TrackMouseClicks bool
	startupLoading   bool
	useViewBounds    bool
	// UseAspectRatio determines whether new window changes will distort or
	// maintain the relative width to height ratio of the screen buffer.
	UseAspectRatio bool

	inFocus bool
}

var (
	nextControllerID = new(int32)
)

// NewWindow creates a window with defauklt settings.
func NewWindow() *Window {
	c := &Window{
		State:         key.NewState(),
		transitionCh:  make(chan struct{}),
		sceneCh:       make(chan struct{}),
		skipSceneCh:   make(chan string),
		quitCh:        make(chan struct{}),
		drawCh:        make(chan struct{}),
		betweenDrawCh: make(chan func()),
	}

	c.SceneMap = scene.NewMap()
	c.Driver = driver.Main
	c.prePublish = func(*Window, screen.Texture) {}
	c.bkgFn = func() image.Image {
		return image.Black
	}
	c.eventHandler = event.DefaultBus
	c.MouseTree = mouse.DefaultTree
	c.CollisionTree = collision.DefaultTree
	c.CallerMap = event.DefaultCallerMap
	c.DrawStack = render.GlobalDrawStack
	c.TrackMouseClicks = true
	c.commands = make(map[string]func([]string))
	c.ControllerID = atomic.AddInt32(nextControllerID, 1)
	c.ParentContext = context.Background()
	return c
}

// Propagate triggers direct mouse events on entities which are clicked
func (w *Window) Propagate(eventName string, me mouse.Event) {
	hits := w.MouseTree.SearchIntersect(me.ToSpace().Bounds())
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Location.Min.Z() < hits[i].Location.Max.Z()
	})
	for _, sp := range hits {
		<-sp.CID.TriggerBus(eventName, &me, w.eventHandler)
		if me.StopPropagation {
			break
		}
	}
	me.StopPropagation = false

	if w.TrackMouseClicks {
		if eventName == mouse.PressOn+"Relative" {
			w.lastRelativePress = me
		} else if eventName == mouse.PressOn {
			w.LastMousePress = me
		} else if eventName == mouse.ReleaseOn {
			if me.Button == w.LastMousePress.Button {
				pressHits := w.MouseTree.SearchIntersect(w.LastMousePress.ToSpace().Bounds())
				sort.Slice(pressHits, func(i, j int) bool {
					return pressHits[i].Location.Min.Z() < pressHits[i].Location.Max.Z()
				})
				for _, sp1 := range pressHits {
					for _, sp2 := range hits {
						if sp1.CID == sp2.CID {
							w.eventHandler.Trigger(mouse.Click, &me)
							<-sp1.CID.TriggerBus(mouse.ClickOn, &me, w.eventHandler)
							if me.StopPropagation {
								return
							}
						}
					}
				}
			}
		} else if eventName == mouse.ReleaseOn+"Relative" {
			if me.Button == w.lastRelativePress.Button {
				pressHits := w.MouseTree.SearchIntersect(w.lastRelativePress.ToSpace().Bounds())
				sort.Slice(pressHits, func(i, j int) bool {
					return pressHits[i].Location.Min.Z() < pressHits[i].Location.Max.Z()
				})
				for _, sp1 := range pressHits {
					for _, sp2 := range hits {
						if sp1.CID == sp2.CID {
							sp1.CID.Trigger(mouse.ClickOn+"Relative", &me)
							if me.StopPropagation {
								return
							}
						}
					}
				}
			}
		}
	}
}

// Width returns the absolute width of the window in pixels.
func (w *Window) Width() int {
	return w.ScreenWidth
}

// Height returns the absolute height of the window in pixels.
func (w *Window) Height() int {
	return w.ScreenHeight
}

// Viewport returns the viewport's position. Its width and height are the window's
// width and height. This position plus width/height cannot exceed ViewportBounds.
func (w *Window) Viewport() intgeom.Point2 {
	return w.viewPos
}

// ViewportBounds returns the boundary of this window's viewport, or the rectangle
// that the viewport is not allowed to exit as it moves around. It often represents
// the total size of the world within a given scene.
func (w *Window) ViewportBounds() intgeom.Rect2 {
	return w.viewBounds
}

// SetLoadingRenderable sets what renderable should display between scenes
// during loading phases.
func (w *Window) SetLoadingRenderable(r render.Renderable) {
	w.LoadingR = r
}

// SetBackground sets this window's background.
func (w *Window) SetBackground(b Background) {
	w.bkgFn = func() image.Image {
		return b.GetRGBA()
	}
}

// SetColorBackground sets this window's background to be a standar image.Image,
// commonly a uniform color.
func (w *Window) SetColorBackground(img image.Image) {
	w.bkgFn = func() image.Image {
		return img
	}
}

// GetBackgroundImage returns the image this window will display as its background
func (w *Window) GetBackgroundImage() image.Image {
	return w.bkgFn()
}

// SetLogicHandler swaps the logic system of the engine with some other
// implementation. If this is never called, it will use event.DefaultBus
func (w *Window) SetLogicHandler(h event.Handler) {
	w.eventHandler = h
}

// NextScene  causes this window to immediately end the current scene.
func (w *Window) NextScene() {
	go func() {
		w.skipSceneCh <- ""
	}()
}

// GoToScene causes this window to skip directly to the given scene.
func (w *Window) GoToScene(nextScene string) {
	go func() {
		w.skipSceneCh <- nextScene
	}()
}

// InFocus returns whether this window is currently in focus.
func (w *Window) InFocus() bool {
	return w.inFocus
}

// CollisionTrees helps access the mouse and collision trees from the controller.
// These trees together detail how a controller can drive mouse and entity interactions.
func (w *Window) CollisionTrees() (mouseTree, collisionTree *collision.Tree) {
	return w.MouseTree, w.CollisionTree
}

// EventHandler returns this window's event handler.
func (w *Window) EventHandler() event.Handler {
	return w.eventHandler
}

// MostRecentInput returns the most recent input type (e.g keyboard/mouse or joystick)
// recognized by the window. This value will only change if the controller's Config is
// set to TrackInputChanges
func (w *Window) MostRecentInput() InputType {
	return w.mostRecentInput
}

func (w *Window) exitWithError(err error) {
	w.exitError = err
	w.Quit()
}

func (w *Window) debugConsole(input io.Reader, output io.Writer) {
	debugstream.AttachToStream(w.ParentContext, input, output)
	debugstream.AddDefaultsForScope(w.ControllerID, w)
}
