package debugstream

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/debugtools"
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
	dlog.ErrorCheck(sc.AddCommand(Command{ScopeID: scopeID, Name: "fullscreen", Usage: explainFullScreen, Operation: fullScreen(controller)}))
	dlog.ErrorCheck(sc.AddCommand(Command{ScopeID: scopeID, Name: "mouse-details", Usage: explainMouseDetails, Operation: mouseCommands(controller)}))
	dlog.ErrorCheck(sc.AddCommand(Command{ScopeID: scopeID, Name: "quit", Usage: explainQuit, Operation: quitCommands(controller)}))
	dlog.ErrorCheck(sc.AddCommand(Command{ScopeID: scopeID, Name: "skip-scene", Operation: skipCommands(controller)}))
	dlog.ErrorCheck(sc.AddCommand(Command{ScopeID: scopeID, Name: "move", Operation: moveWindow(controller)}))

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
		width := parseTokenAsInt(sub, 2, w.Width())
		height := parseTokenAsInt(sub, 3, w.Height())
		v := w.Viewport()
		x := parseTokenAsInt(sub, 0, v.X())
		y := parseTokenAsInt(sub, 1, v.Y())
		w.MoveWindow(x, y, width, height)
		return ""
	}
}

const explainFullScreen = "specify off 'fullscreen off' to exit fullscreen"

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

const explainMouseDetails = "the mext mouse click on the given window will print the cursor's location"

func mouseCommands(w window.Window) func([]string) string {
	return func(tokenString []string) string {
		event.GlobalBind(w.EventHandler(), mouse.Release, mouseDetails(w))
		return ""
	}
}

func mouseDetails(w window.Window) func(*mouse.Event) event.Response {
	return func(me *mouse.Event) event.Response {
		viewPos := w.Viewport()
		x := int(me.X()) + viewPos[0]
		y := int(me.Y()) + viewPos[1]
		loc := collision.NewUnassignedSpace(float64(x), float64(y), 16, 16)
		results := collision.Hits(loc)
		fmt.Println("Mouse at:", x, y, "rel:", me.X(), me.Y())
		if len(results) == 0 {
			results = mouse.Hits(loc)
		}
		cm := w.GetCallerMap()

		if len(results) > 0 {
			i := results[0].CID
			if i > 0 && cm.HasEntity(i) {
				e := cm.GetEntity(i)
				fmt.Printf("%+v\n", e)
			} else {
				fmt.Println("No entity ", i)
			}
		}

		return event.UnbindThis
	}
}

const explainQuit = "close the given window"

func quitCommands(w window.Window) func([]string) string {
	return func(tokenString []string) string {
		w.Quit()
		return ""

	}
}

func skipCommands(w window.Window) func([]string) string {
	return func(tokenString []string) string {
		w.NextScene()
		return ""
	}
}

const explainFade = "fade the specified renderable by the given int if given. Renderable must be registered in debugtools"

func fadeCommands(tokenString []string) (out string) {
	if len(tokenString) == 0 {
		return oakerr.InsufficientInputs{
			AtLeast:   1,
			InputName: "arguments",
		}.Error() + "\n"
	}
	toFade, ok := debugtools.GetDebugRenderable(tokenString[0])
	if ok {
		fadeVal := parseTokenAsInt(tokenString, 1, 255)
		toFade.(render.Modifiable).Filter(mod.Fade(fadeVal))
		return
	}
	out += fmt.Sprintf("Could not fade input %s\n", tokenString[0]) +
		fmt.Sprintf("Possible inputs are '%s'\n", strings.Join(debugtools.EnumerateDebugRenderableKeys(), ", "))
	return
}

func parseTokenAsInt(tokenString []string, arrIndex int, defaultVal int) int {
	if len(tokenString) > arrIndex {
		tmp, err := strconv.Atoi(tokenString[arrIndex])
		if err == nil {
			return tmp
		}
	}
	return defaultVal
}
