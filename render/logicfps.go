package render

import (
	"image"
	"image/draw"
	"time"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/timing"
)

// LogicFPS is a draw stack element that will draw the logical fps onto the screen
type LogicFPS struct {
	event.CID
	fps      int
	lastTime time.Time
	txt      *Text
}

// Init satisfies event.Entity
func (lf *LogicFPS) Init() event.CID {
	id := event.NextID(lf)
	lf.CID = id
	return id
}

// NewLogicFPS returns a zero-initialized LogicFPS
func NewLogicFPS() *LogicFPS {
	lf := new(LogicFPS)
	lf.lastTime = time.Now()
	lf.fps = 0
	return lf
}

// PreDraw does nothing for a drawFPS
func (lf *LogicFPS) PreDraw() {
	//NOP
	// This is not done in NewDrawFPS because at the time when
	// NewDrawFPS is called, DefFont() does not exist.
	if lf.txt == nil {
		lf.Init()
		lf.Bind(logicFPSBind, "EnterFrame")
		lf.txt = DefFont().NewIntText(&lf.fps, 10, 30)
	}
}

// Add does nothing for a drawFPS
func (lf *LogicFPS) Add(Renderable, int) Renderable {
	//NOP
	return nil
}

// Replace does nothing for a drawFPS
func (lf *LogicFPS) Replace(Renderable, Renderable, int) {
	//NOP
}

// Copy does effectively nothing for a drawFPS
func (lf *LogicFPS) Copy() Addable {
	return new(LogicFPS)
}

func (lf *LogicFPS) draw(world draw.Image, view image.Point, w, h int) {
	lf.txt.Draw(world)
}

func logicFPSBind(id int, nothing interface{}) int {
	lf := event.GetEntity(id).(*LogicFPS)
	lf.fps = int((timing.FPS(lf.lastTime, time.Now()) * FpsSmoothing) + (float64(lf.fps) * (1 - FpsSmoothing)))
	lf.lastTime = time.Now()
	return 0
}
