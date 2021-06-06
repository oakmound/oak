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
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "fullscreen", explainFullScreen, fullScreen(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "mouse", nil, mouseCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "quit", nil, quitCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "skip", nil, skipCommands(controller)))
	dlog.ErrorCheck(sc.AddScopedCommand(scopeID, "move", nil, moveWindow(controller)))

	if sc.assumedScope != 0 {
		return
	}
	// assume the scope for easy usage here
	sc.assumedScope = scopeID
}

func moveWindow(w window.Window) func([]string) string {
	return func(sub []string) string {
		if len(sub) != 2 && len(sub) != 4 {
			return oakerr.InsufficientInputs{
				AtLeast:   2,
				InputName: "coordinates",
			}.Error()
		}
		width := parseTokenAsInt(sub, 3, w.Width())
		height := parseTokenAsInt(sub, 4, w.Height())
		v := w.Viewport()
		x := parseTokenAsInt(sub, 0, v.X())
		y := parseTokenAsInt(sub, 1, v.Y())
		w.MoveWindow(x, y, width, height)
		return ""
	}
}

func explainFullScreen([]string) string {
	return "specify off 'fullscreen off' to exit fullscreen\n"
}
func fullScreen(w window.Window) func([]string) string {
	return func(sub []string) (out string) {
		on := true
		if len(sub) > 0 {
			if sub[0] == "off" {
				on = false
			}
		}
		err := w.SetFullScreen(on)
		dlog.ErrorCheck(err)
		return
	}
}

func mouseCommands(w window.Window) func([]string) string {
	return func(tokenString []string) string {
		if len(tokenString) != 1 {
			return oakerr.InsufficientInputs{
				AtLeast:   1,
				InputName: "arguments",
			}.Error()
		}
		switch tokenString[0] {
		case "details":
			w.EventHandler().GlobalBind("MouseRelease", mouseDetails(w))
		default:
			return "Bad Mouse Input"
		}

		return ""

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

func quitCommands(w window.Window) func([]string) string {
	return func(tokenString []string) string {
		w.Quit()
		if len(tokenString) > 0 {
			return fmt.Sprintf("Quit does not support extra options such as the ones provided: %v\n ", tokenString) +
				oakerr.InvalidInput{InputName: "any arguments"}.Error()
		}
		return ""

	}
}

func skipCommands(w window.Window) func([]string) string {
	return func(tokenString []string) string {

		if len(tokenString) != 1 {
			return oakerr.InsufficientInputs{
				AtLeast:   1,
				InputName: "arguments",
			}.Error()
		}
		switch tokenString[0] {
		case "scene":
			w.NextScene()
		default:
			return oakerr.NotFound{InputName: tokenString[0]}.Error()

		}
		return ""
	}
}

func fadeCommands(tokenString []string) (out string) {
	if len(tokenString) == 0 {
		return oakerr.InsufficientInputs{
			AtLeast:   1,
			InputName: "arguments",
		}.Error()
	}
	toFade, ok := render.GetDebugRenderable(tokenString[0])
	if ok {
		fadeVal := parseTokenAsInt(tokenString, 1, 255)
		toFade.(render.Modifiable).Filter(mod.Fade(fadeVal))
	}
	out += fmt.Sprintf("Could not fade input %s\n", tokenString[0]) +
		fmt.Sprintf("Possible inputs are '%s'\n", strings.Join(render.EnumerateDebugRenderableKeys(), ", "))

	out += oakerr.NotFound{InputName: tokenString[0]}.Error() + "\n"
	return
}
