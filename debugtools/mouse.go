package debugtools

import (
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/scene"
)

// DebugMouseRelease will print the position and button pressed of the mouse when the mouse is released, if the given
// key is held down at the time. If no key is given, it will always be printed
func DebugMouseRelease(ctx *scene.Context, k string) {
	event.GlobalBind(mouse.Release, func(_ event.CID, ev interface{}) int {
		mev, _ := ev.(*mouse.Event)
		if k == "" || ctx.KeyState.IsDown(k) {
			dlog.Info(mev)
		}
		return 0
	})
}
