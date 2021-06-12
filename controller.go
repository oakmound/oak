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

var _ window.Window = &Controller{}

func (c *Controller) windowController(s screen.Screen, x, y int32, width, height int) (screen.Window, error) {
	return s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title(c.config.Title),
		screen.Position(x, y),
		screen.Fullscreen(c.config.Fullscreen),
		screen.Borderless(c.config.Borderless),
		screen.TopMost(c.config.TopMost),
	))
}

type Controller struct {
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
	winBuffer     screen.Image
	screenControl screen.Screen
	windowControl screen.Window

	windowRect image.Rectangle

	// DrawTicker is the parallel to LogicTicker to set the draw framerate
	DrawTicker *time.Ticker

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
	prePublish func(c *Controller, tx screen.Texture)

	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable

	firstScene string
	// ErrorScene is a scene string that will be entered if the scene handler
	// fails to enter some other scene, for example, because it's name was
	// undefined in the scene map. If the scene map does not have ErrorScene
	// as well, it will fall back to panicking.
	ErrorScene string

	logicHandler  event.Handler
	CallerMap     *event.CallerMap
	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
	DrawStack     *render.DrawStack

	// LastMouseEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastMouseEvent    mouse.Event
	lastRelativePress mouse.Event
	// LastPress is the last triggered mouse event,
	// where the mouse event was a press.
	// If TrackMouseClicks is set to false then this will not be tracked
	LastMousePress mouse.Event

	FirstSceneInput interface{}

	commands map[string]func([]string)

	ControllerID int32

	windowTexture screen.Texture

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

func NewController() *Controller {
	c := &Controller{
		State:        key.NewState(),
		transitionCh: make(chan struct{}),
		sceneCh:      make(chan struct{}),
		skipSceneCh:  make(chan string),
		quitCh:       make(chan struct{}),
		drawCh:       make(chan struct{}),
	}

	c.SceneMap = scene.NewMap()
	c.Driver = driver.Main
	c.prePublish = func(*Controller, screen.Texture) {}
	c.bkgFn = func() image.Image {
		return image.Black
	}
	c.startupLoading = true
	c.logicHandler = event.DefaultBus
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
func (c *Controller) Propagate(eventName string, me mouse.Event) {
	c.LastMouseEvent = me
	mouse.LastEvent = me

	hits := c.MouseTree.SearchIntersect(me.ToSpace().Bounds())
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Location.Min.Z() < hits[i].Location.Max.Z()
	})
	for _, sp := range hits {
		<-sp.CID.TriggerBus(eventName, &me, c.logicHandler)
		if me.StopPropagation {
			break
		}
	}
	me.StopPropagation = false

	if c.TrackMouseClicks {
		if eventName == mouse.PressOn+"Relative" {
			c.lastRelativePress = me
		} else if eventName == mouse.PressOn {
			c.LastMousePress = me
		} else if eventName == mouse.ReleaseOn {
			if me.Button == c.LastMousePress.Button {
				pressHits := c.MouseTree.SearchIntersect(c.LastMousePress.ToSpace().Bounds())
				sort.Slice(pressHits, func(i, j int) bool {
					return pressHits[i].Location.Min.Z() < pressHits[i].Location.Max.Z()
				})
				for _, sp1 := range pressHits {
					for _, sp2 := range hits {
						if sp1.CID == sp2.CID {
							c.logicHandler.Trigger(mouse.Click, &me)
							<-sp1.CID.TriggerBus(mouse.ClickOn, &me, c.logicHandler)
							if me.StopPropagation {
								return
							}
						}
					}
				}
			}
		} else if eventName == mouse.ReleaseOn+"Relative" {
			if me.Button == c.lastRelativePress.Button {
				pressHits := c.MouseTree.SearchIntersect(c.lastRelativePress.ToSpace().Bounds())
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

func (c *Controller) Width() int {
	return c.ScreenWidth
}

func (c *Controller) Height() int {
	return c.ScreenHeight
}

func (c *Controller) Viewport() intgeom.Point2 {
	return c.viewPos
}

func (c *Controller) ViewportBounds() intgeom.Rect2 {
	return c.viewBounds
}

func (c *Controller) SetLoadingRenderable(r render.Renderable) {
	c.LoadingR = r
}

func (c *Controller) SetBackground(b Background) {
	c.bkgFn = func() image.Image {
		return b.GetRGBA()
	}
}

func (c *Controller) SetColorBackground(img image.Image) {
	c.bkgFn = func() image.Image {
		return img
	}
}

func (c *Controller) GetBackgroundImage() image.Image {
	return c.bkgFn()
}

// SetLogicHandler swaps the logic system of the engine with some other
// implementation. If this is never called, it will use event.DefaultBus
func (c *Controller) SetLogicHandler(h event.Handler) {
	c.logicHandler = h
}

func (c *Controller) NextScene() {
	c.skipSceneCh <- ""
}

func (c *Controller) GoToScene(nextScene string) {
	c.skipSceneCh <- nextScene
}

func (c *Controller) InFocus() bool {
	return c.inFocus
}

// CollisionTrees helps access the mouse and collision trees from the controller.
// These trees together detail how a controller can drive mouse and entity interactions.
func (c *Controller) CollisionTrees() (mouseTree, collisionTree *collision.Tree) {
	return c.MouseTree, c.CollisionTree
}

func (c *Controller) EventHandler() event.Handler {
	return c.logicHandler
}

// MostRecentInput returns the most recent input type (e.g keyboard/mouse or joystick)
// recognized by the window. This value will only change if the controller's Config is
// set to TrackInputChanges
func (c *Controller) MostRecentInput() InputType {
	return c.mostRecentInput
}

func (c *Controller) exitWithError(err error) {
	c.exitError = err
	c.Quit()
}

func (c *Controller) debugConsole(input io.Reader, output io.Writer) {
	debugstream.AttachToStream(c.ParentContext, input, output)
	debugstream.AddDefaultsForScope(c.ControllerID, c)
}
