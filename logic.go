package oak

import "github.com/oakmound/oak/v2/event"

var (
	logicHandler event.Handler = event.DefaultBus
)

// SetLogicHandler swaps the logic system of the engine with some other
// implementation. If this is never called, it will use event.DefaultBus
func SetLogicHandler(h event.Handler) {
	logicHandler = h
}
