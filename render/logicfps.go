package render

import (
	"time"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/timing"
)

// LogicFPS is a Stackable that will draw the logical fps onto the screen when a part
// of the draw stack.
type LogicFPS struct {
	event.CallerID
	*Text
	fps       int
	lastTime  time.Time
	Smoothing float64
}

// Init satisfies event.Entity
func (lf *LogicFPS) Init() event.CallerID {
	// TODO: not default caller map
	id := event.DefaultCallerMap.Register(lf)
	lf.CallerID = id
	return id
}

// NewLogicFPS returns a LogicFPS, which will render a counter of how fast it receives event.Enter events.
// If font is not provided, DefaultFont is used. If smoothing is 0, a reasonable default is used.
func NewLogicFPS(smoothing float64, font *Font, x, y float64) *LogicFPS {
	if smoothing == 0.0 {
		smoothing = defaultFpsSmoothing
	}
	if font == nil {
		font = DefaultFont().Copy()
	}
	lf := &LogicFPS{
		Smoothing: smoothing,
		lastTime:  time.Now(),
	}
	lf.Text = font.NewIntText(&lf.fps, x, y)
	lf.Init()
	// TODO: not default bus
	event.Bind(event.DefaultBus, event.Enter, lf.CallerID, logicFPSBind)

	return lf
}

func logicFPSBind(id event.CallerID, _ event.EnterPayload) event.Response {
	// TODO v4: should bindings give you an interface instead of a callerID, so bindings don't need to
	// know what caller map to look up the caller from?
	lf := event.DefaultCallerMap.GetEntity(id).(*LogicFPS)
	t := time.Now()
	lf.fps = int((timing.FPS(lf.lastTime, t) * lf.Smoothing) + (float64(lf.fps) * (1 - lf.Smoothing)))
	lf.lastTime = t
	return 0
}
