package debugtools

import (
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/scene"
)

// DebugMouseRelease will print the position and button pressed of the mouse when the mouse is released, if the given
// key is held down at the time. If 0 is given, it will always be printed
func DebugMouseRelease(ctx *scene.Context, k key.Code) {
	event.GlobalBind(ctx, mouse.Release, func(mev *mouse.Event) event.Response {
		if k == 0 || ctx.KeyState.IsDown(k) {
			dlog.Info(mev)
		}
		return 0
	})
}
