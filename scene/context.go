package scene

import (
	"context"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/window"
)

// A Context contains all transient engine components used in a scene, including
// the draw stack, event bus, known event callers, collision trees, keyboard state,
// and a reference to the OS window itself. When a scene ends, modifications made
// to these structures will be reset, excluding window modifications.
// TODO oak v4: consider embedding these system objects on the context to change
// ctx.DrawStack.Draw to ctx.Draw and ctx.EventHandler.Bind to ctx.Bind
type Context struct {
	// This context will be canceled when the scene ends
	context.Context

	PreviousScene string
	SceneInput    interface{}
	Window        window.Window

	DrawStack     *render.DrawStack
	EventHandler  event.Handler
	CallerMap     *event.CallerMap
	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
	KeyState      *key.State
}
