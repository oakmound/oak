package oak

import "github.com/oakmound/oak/v2/event"

// SetLogicHandler swaps the logic system of the engine with some other
// implementation. If this is never called, it will use event.DefaultBus
func (c *Controller) SetLogicHandler(h event.Handler) {
	c.logicHandler = h
}
