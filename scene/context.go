package scene

import (
	"context"

	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
)

// A Context contains all transient engine components used in a scene, including
// the draw stack, event bus, known event callers, collision trees, keyboard state,
// and a reference to the OS window itself. When a scene ends, modifications made
// to these structures will be reset, excluding window modifications.
type Context struct {
	// This context will be canceled when the scene ends
	context.Context

	PreviousScene string
	SceneInput    interface{}
	Window        Window

	*event.CallerMap
	event.Handler
	*render.DrawStack
	*key.State

	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
}

// DoEachFrame is a helper method to call a function on each frame for the duration of this scene.
func (ctx *Context) DoEachFrame(f func()) {
	event.GlobalBind(ctx, event.Enter, func(_ event.EnterPayload) event.Response {
		f()
		return 0
	})
}
