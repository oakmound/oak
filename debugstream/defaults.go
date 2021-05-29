package debugstream

import (
	"fmt"
	"strings"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
)

// AddDefaultsForScope for debugging.
// This is a nasty set of reflection all to break this out into a subpackage.
func (sc *ScopedCommands) AddDefaultsForScope(scopeID int32, controller interface{}) {

	if fs, ok := controller.(fullScreenable); ok {
		dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "fullscreen", fullScreen(fs)))
	}

	if md, ok := controller.(mouseDetailer); ok {
		dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "mouse", mouseCommands(md)))
	}

	if cq, ok := controller.(canQuit); ok {
		dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "quit", quitCommands(cq)))
	}

	if hs, ok := controller.(hasScenes); ok {
		dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "skip", skipCommands(hs)))
	}

	// 	dlog.ErrorCheck(c.AddCommand("move", c.moveWindow))

}

type fullScreenable interface {
	SetFullScreen(bool) error
}

func fullScreen(fs fullScreenable) func([]string) {
	return func(sub []string) {
		on := true
		if len(sub) > 0 {
			if sub[0] == "off" {
				on = false
			}
		}
		err := fs.SetFullScreen(on)
		dlog.ErrorCheck(err)
	}
}

type hasCollisionTrees interface {
	CollisionTrees() (mouseTree, collisionTree *collision.Tree)
}
type hasViewport interface {
	Viewport() intgeom.Point2
}
type hasGlobalBind interface {
	GlobalBind(name string, fn event.Bindable)
}

type mouseDetailer interface {
	hasCollisionTrees
	hasViewport
	hasGlobalBind
}

func mouseCommands(md mouseDetailer) func([]string) {
	return func(tokenString []string) {
		if len(tokenString) != 1 {
			fmt.Println("Input must be a single string from the following (\"details\") ")
			return
		}
		switch tokenString[0] {
		case "details":
			// CONSIDER: scoping to the controllers logicHandler
			md.GlobalBind("MouseRelease", mouseDetails(md))
		default:
			fmt.Println("Bad Mouse Input")
		}

	}

}

func mouseDetails(md mouseDetailer) func(event.CID, interface{}) int {
	return func(nothing event.CID, mevent interface{}) int {
		me := mevent.(mouse.Event)
		viewPos := md.Viewport()
		x := int(me.X()) + viewPos[0]
		y := int(me.Y()) + viewPos[1]
		loc := collision.NewUnassignedSpace(float64(x), float64(y), 16, 16)
		results := collision.Hits(loc)
		fmt.Println("Mouse at:", x, y, "rel:", me.X(), me.Y())
		if len(results) == 0 {
			results = mouse.Hits(loc)
		}
		if len(results) > 0 {
			i := int(results[0].CID)
			if i > 0 && event.HasEntity(event.CID(i)) {
				e := event.GetEntity(event.CID(i))
				fmt.Printf("%+v\n", e)
			} else {
				fmt.Println("No entity ", i)
			}
		}

		return event.UnbindSingle
	}
}

type canQuit interface {
	Quit()
}

func quitCommands(quitter canQuit) func([]string) {
	return func(tokenString []string) {
		if len(tokenString) > 0 {
			fmt.Println("Quit does not support extra options such as the ones provided: ", tokenString)
		}
		quitter.Quit()

	}
}

type hasScenes interface {
	NextScene()
}

func skipCommands(scener hasScenes) func([]string) {
	return func(tokenString []string) {

		if len(tokenString) != 1 {
			fmt.Println("Input must be a single string from the following (\"scene\"). ")
			return
		}
		switch tokenString[0] {
		case "scene":
			scener.NextScene()
		default:
			fmt.Println("Bad Skip Input")
		}
	}
}

func fadeCommands(tokenString []string) {
	if len(tokenString) == 0 {
		fmt.Println("Input must start with the name of the renderable to fade")
		return
	}
	toFade, ok := render.GetDebugRenderable(tokenString[0])
	if ok {
		fadeVal := parseTokenAsInt(tokenString, 1, 255)
		toFade.(render.Modifiable).Filter(mod.Fade(fadeVal))
		return
	}

	fmt.Println("Could not fade input", tokenString[0])
	fmt.Printf("Possible inputs are '%s'\n", strings.Join(render.EnumerateDebugRenderableKeys(), ", "))

}
