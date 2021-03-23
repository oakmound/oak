package scene

import (
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/window"
)

type Context struct {
	PreviousScene string
	SceneInput    interface{}
	Window        window.Window

	DrawStack     *render.DrawStack
	EventHandler  event.Handler
	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
	// todo: ...
}
