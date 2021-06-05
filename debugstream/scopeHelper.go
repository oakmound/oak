package debugstream

import (
	"fmt"
	"strings"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/window"
)

// AddDefaultsForScope for debugging.
func (sc *ScopedCommands) AddDefaultsForScope(scopeID int32, controller window.Window) {
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "fullscreen", nil, fullScreen(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "mouse", nil, mouseCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "quit", nil, quitCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "skip", nil, skipCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "move", nil, moveWindow(controller)))
}

func moveWindow(w window.Window) func([]string) error {
	return func(sub []string) error {
		if len(sub) != 2 && len(sub) != 4 {
			fmt.Println("'move' expects 'x y' or 'x y w h'")
			return oakerr.InsufficientInputs{
				AtLeast:   2,
				InputName: "coordinates",
			}
		}
		width := parseTokenAsInt(sub, 3, w.Width())
		height := parseTokenAsInt(sub, 4, w.Height())
		v := w.Viewport()
		x := parseTokenAsInt(sub, 0, v.X())
		y := parseTokenAsInt(sub, 1, v.Y())
		w.MoveWindow(x, y, width, height)
		return nil
	}
}

func fullScreen(w window.Window) func([]string) error {
	return func(sub []string) error {
		on := true
		if len(sub) > 0 {
			if sub[0] == "off" {
				on = false
			}
		}
		err := w.SetFullScreen(on)
		dlog.ErrorCheck(err)
		return nil
	}
}

func mouseCommands(w window.Window) func([]string) error {
	return func(tokenString []string) error {
		if len(tokenString) != 1 {
			fmt.Println("Input must be a single string from the following (\"details\") ")
			return oakerr.InsufficientInputs{
				AtLeast:   1,
				InputName: "arguments",
			}
		}
		switch tokenString[0] {
		case "details":
			w.EventHandler().GlobalBind("MouseRelease", mouseDetails(w))
		default:
			fmt.Println("Bad Mouse Input")
		}

		return nil

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

func quitCommands(w window.Window) func([]string) error {
	return func(tokenString []string) error {
		w.Quit()
		if len(tokenString) > 0 {
			fmt.Println("Quit does not support extra options such as the ones provided: ", tokenString)
			return oakerr.InvalidInput{InputName: "any arguments"}
		}
		return nil

	}
}

func skipCommands(w window.Window) func([]string) error {
	return func(tokenString []string) error {

		if len(tokenString) != 1 {
			fmt.Println("Input must be a single string from the following (\"scene\"). ")
			return oakerr.InsufficientInputs{
				AtLeast:   1,
				InputName: "arguments",
			}
		}
		switch tokenString[0] {
		case "scene":
			w.NextScene()
		default:
			fmt.Println("Bad Skip Input")
			return oakerr.NotFound{InputName: tokenString[0]}

		}
		return nil
	}
}

func fadeCommands(tokenString []string) error {
	if len(tokenString) == 0 {
		fmt.Println("Input must start with the name of the renderable to fade")
		return oakerr.InsufficientInputs{
			AtLeast:   1,
			InputName: "arguments",
		}
	}
	toFade, ok := render.GetDebugRenderable(tokenString[0])
	if ok {
		fadeVal := parseTokenAsInt(tokenString, 1, 255)
		toFade.(render.Modifiable).Filter(mod.Fade(fadeVal))
	}

	fmt.Println("Could not fade input", tokenString[0])
	fmt.Printf("Possible inputs are '%s'\n", strings.Join(render.EnumerateDebugRenderableKeys(), ", "))
	return oakerr.NotFound{InputName: tokenString[0]}

}
