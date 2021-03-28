package oak

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
	"github.com/oakmound/oak/v2/timing"
	"github.com/oakmound/shiny/driver"
	"github.com/oakmound/shiny/screen"
)

func (c *Controller) windowController(s screen.Screen, x, y int32, width, height int) (screen.Window, error) {
	return s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title(conf.Title),
		screen.Position(x, y),
		screen.Fullscreen(SetupFullscreen),
		screen.Borderless(SetupBorderless),
		screen.TopMost(SetupTopMost),
	))
}

type Controller struct {
	// TODO: most of these channels should take struct{}s and be closed sometime
	transitionCh chan bool

	// The Scene channel receives a signal
	// when a scene's .loop() function should
	// be called.
	sceneCh chan bool

	// The skip scene channel receives a debug
	// signal to forcibly go to the next
	// scene.
	skipSceneCh chan bool

	// The quit channel receives a signal when
	// the program should stop.
	quitCh chan bool

	// The draw channel receives a signal when
	// drawing should cease (or resume)
	drawCh chan bool

	// The debug reset channel represents
	// when the debug console should forget the
	// commands that have been sent to it.
	debugResetCh chan bool

	// The viewport channel controls when new
	// viewport positions should be drawn
	viewportCh chan intgeom.Point2

	// The viewport shift channel controls when new
	// viewport positions should be shifted to and drawn
	viewportShiftCh chan intgeom.Point2

	debugResetInProgress bool

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

	winBuffer     screen.Image
	screenControl screen.Screen
	windowControl screen.Window

	windowRect     image.Rectangle
	windowUpdateCh chan bool

	// TODO V3: cleanup this interface
	// BackgroundColor is the starting background color for the draw loop. BackgroundImage or SetBackground will overwrite it.
	BackgroundColor image.Image
	// BackgroundImage is the starting custom background for the draw loop. SetBackground will overwrite it.
	BackgroundImage Background
	// DrawTicker is the parallel to LogicTicker to set the draw framerate
	DrawTicker *timing.DynamicTicker

	setBackgroundCh chan Background

	bkgFn func() image.Image

	// SceneMap is a global map of scenes referred to when scenes advance to
	// determine what the next scene should be.
	// It can be replaced or modified so long as these modifications happen
	// during a scene or before the controller has started.
	SceneMap *scene.Map

	// ViewPos represents the point in the world which the viewport is anchored at.
	ViewPos intgeom.Point2
	// ViewPosMutex is used to grant extra saftey in viewpos operations
	viewPosMutex  sync.Mutex
	useViewBounds bool
	viewBounds    intgeom.Rect2

	// ColorPalette is the current color palette oak is set to conform to. Modification of this
	// value directly will not effect oak's palette, use SetPalette instead. If SetPallete is never called,
	// this is the zero value ([]Color of length 0).
	ColorPalette color.Palette

	// UseAspectRatio determines whether new window changes will distort or
	// maintain the relative width to height ratio of the screen buffer.
	UseAspectRatio bool
	aspectRatio    float64

	// Driver is the driver oak will call during initialization
	Driver Driver

	drawLoopPublishDef func(c *Controller, tx screen.Texture)
	drawLoopPublish    func(c *Controller, tx screen.Texture)

	keyState     map[string]bool
	keyDurations map[string]time.Time
	keyLock      sync.RWMutex
	durationLock sync.RWMutex

	startupLoadCh chan bool
	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable

	firstScene string
	// ErrorScene is a scene string that will be entered if the scene handler
	// fails to enter some other scene, for example, because it's name was
	// undefined in the scene map. If the scene map does not have ErrorScene
	// as well, it will fall back to panicking.
	ErrorScene string

	logicHandler event.Handler
	CallerMap    *event.CallerMap

	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
	// TODO: separate initial configuration from controller
	DrawStack        *render.DrawStack
	InitialDrawStack *render.DrawStack

	lastRelativePress mouse.Event

	// LastMouseEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastMouseEvent mouse.Event

	TrackMouseClicks bool

	// LastPress is the last triggered mouse event,
	// where the mouse event was a press.
	// If TrackMouseClicks is set to false then this will not be tracked
	LastMousePress mouse.Event

	FirstSceneInput interface{}

	viewportLocked bool
	commands       map[string]func([]string)

	ControllerID int32
}

var (
	nextControllerID = new(int32)
)

func NewController() *Controller {
	c := &Controller{}
	c.transitionCh = make(chan bool)
	c.sceneCh = make(chan bool)
	c.skipSceneCh = make(chan bool)
	c.quitCh = make(chan bool)
	c.drawCh = make(chan bool)
	c.debugResetCh = make(chan bool)
	c.viewportCh = make(chan intgeom.Point2)
	c.viewportShiftCh = make(chan intgeom.Point2)
	c.windowUpdateCh = make(chan bool)
	c.SceneMap = scene.NewMap()
	c.BackgroundColor = image.Black
	c.setBackgroundCh = make(chan Background)
	c.Driver = driver.Main
	c.drawLoopPublishDef = func(c *Controller, tx screen.Texture) {
		tx.Upload(zeroPoint, c.winBuffer, c.winBuffer.Bounds())
		c.windowControl.Scale(c.windowRect, tx, tx.Bounds(), draw.Src)
		c.windowControl.Publish()
	}
	c.drawLoopPublish = c.drawLoopPublishDef
	c.bkgFn = func() image.Image {
		return c.BackgroundColor
	}
	c.keyState = make(map[string]bool)
	c.keyDurations = make(map[string]time.Time)
	c.startupLoadCh = make(chan bool)
	c.logicHandler = event.DefaultBus
	c.MouseTree = mouse.DefTree
	c.CollisionTree = collision.DefTree
	c.CallerMap = event.DefaultCallerMap
	c.DrawStack = render.GlobalDrawStack
	c.InitialDrawStack = render.NewDrawStack(
		render.NewDynamicHeap(),
	)
	c.TrackMouseClicks = true
	c.viewportLocked = true
	c.commands = make(map[string]func([]string))
	c.ControllerID = atomic.AddInt32(nextControllerID, 1)
	return c
}

// Propagate triggers direct mouse events on entities which are clicked
func (c *Controller) Propagate(eventName string, me mouse.Event) {
	c.LastMouseEvent = me

	hits := c.MouseTree.SearchIntersect(me.ToSpace().Bounds())
	for _, sp := range hits {
		sp.CID.Trigger(eventName, me)
	}

	if c.TrackMouseClicks {
		if eventName == mouse.PressOn+"Relative" {
			c.lastRelativePress = me
		} else if eventName == mouse.PressOn {
			c.LastMousePress = me
		} else if eventName == mouse.ReleaseOn {
			if me.Button == c.LastMousePress.Button {
				pressHits := c.MouseTree.SearchIntersect(c.LastMousePress.ToSpace().Bounds())
				for _, sp1 := range pressHits {
					for _, sp2 := range hits {
						if sp1.CID == sp2.CID {
							event.Trigger(mouse.Click, me)
							sp1.CID.Trigger(mouse.ClickOn, me)
						}
					}
				}
			}
		} else if eventName == mouse.ReleaseOn+"Relative" {
			if me.Button == c.lastRelativePress.Button {
				pressHits := c.MouseTree.SearchIntersect(c.lastRelativePress.ToSpace().Bounds())
				for _, sp1 := range pressHits {
					for _, sp2 := range hits {
						if sp1.CID == sp2.CID {
							sp1.CID.Trigger(mouse.ClickOn+"Relative", me)
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
	return c.ViewPos
}

func (c *Controller) SetLoadingRenderable(r render.Renderable) {
	c.LoadingR = r
}

func (c *Controller) GetBackgroundColor() image.Image {
	return c.BackgroundColor
}

func (c *Controller) SetBackgroundColor(img image.Image) {
	c.BackgroundColor = img
}
