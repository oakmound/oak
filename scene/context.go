package scene

import (
	"context"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/window"
)

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
	// todo: ...
}
