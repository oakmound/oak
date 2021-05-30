package debugstream

import (
	"fmt"
	"strings"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/window"
)

// AddDefaultsForScope for debugging.
func (sc *ScopedCommands) AddDefaultsForScope(scopeID int32, controller window.Window) {
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "fullscreen", fullScreen(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "mouse", mouseCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "quit", quitCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "skip", skipCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "move", moveWindow(controller)))
}

func moveWindow(w window.Window) func([]string) {
	return func(sub []string) {
		if len(sub) != 2 && len(sub) != 4 {
			fmt.Println("'move' expects 'x y' or 'x y w h'")
			return
		}
		width := parseTokenAsInt(sub, 3, w.Width())
		height := parseTokenAsInt(sub, 4, w.Height())
		v := w.Viewport()
		x := parseTokenAsInt(sub, 0, v.X())
		y := parseTokenAsInt(sub, 1, v.Y())
		w.MoveWindow(x, y, width, height)
	}
}

func fullScreen(w window.Window) func([]string) {
	return func(sub []string) {
		on := true
		if len(sub) > 0 {
			if sub[0] == "off" {
				on = false
			}
		}
		err := w.SetFullScreen(on)
		dlog.ErrorCheck(err)
	}
}

func mouseCommands(w window.Window) func([]string) {
	return func(tokenString []string) {
		if len(tokenString) != 1 {
			fmt.Println("Input must be a single string from the following (\"details\") ")
			return
		}
		switch tokenString[0] {
		case "details":
			w.EventHandler().GlobalBind("MouseRelease", mouseDetails(w))
		default:
			fmt.Println("Bad Mouse Input")
		}

	}

}

func mouseDetails(w window.Window) func(event.CID, interface{}) int {
	return func(nothing event.CID, mevent interface{}) int {
		me := mevent.(mouse.Event)
		viewPos := w.Viewport()
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

func quitCommands(w window.Window) func([]string) {
	return func(tokenString []string) {
		if len(tokenString) > 0 {
			fmt.Println("Quit does not support extra options such as the ones provided: ", tokenString)
		}
		w.Quit()
	}
}

func skipCommands(w window.Window) func([]string) {
	return func(tokenString []string) {

		if len(tokenString) != 1 {
			fmt.Println("Input must be a single string from the following (\"scene\"). ")
			return
		}
		switch tokenString[0] {
		case "scene":
			w.NextScene()
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
