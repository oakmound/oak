package oak

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/oakmound/oak/v2/oakerr"

	"github.com/davecgh/go-spew/spew"
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/render/mod"
)

var (
	viewportLocked = true
	commands       = make(map[string]func([]string))
)

// AddCommand adds a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func AddCommand(s string, fn func([]string)) error {
	if _, ok := commands[s]; ok {
		return oakerr.ExistingElement{
			InputName:   "s",
			InputType:   "string",
			Overwritten: false,
		}
	}
	dlog.Info("Adding command", s)
	commands[s] = fn
	return nil
}

// ForceAddCommand adds or overwrites a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace. If a command is overwritten
// the overwritten command will be returned.
func ForceAddCommand(s string, fn func([]string)) func([]string) {

	existing, overwritten := commands[s]
	if overwritten {
		dlog.Info("Overwriting command", s)
	} else {
		dlog.Info("Adding command", s)
	}
	commands[s] = fn
	return existing
}

// ClearCommand clears an existing debug command by key: <s>
func ClearCommand(s string) {
	delete(commands, s)
}

// ResetCommands will throw out all existing debug commands from the
// debug console.
func ResetCommands() {
	commands = map[string]func([]string){}
}

// GetDebugKeys returns the current debug console commands as a string array
func GetDebugKeys() []string {
	dkeys := make([]string, len(commands))
	i := 0
	for k := range commands {
		dkeys[i] = k
		i++
	}
	return dkeys
}

func debugConsole(resetCh, skipScene chan bool, input io.Reader) {
	scanner := bufio.NewScanner(input)
	spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 2

	// built in commands
	if conf.LoadBuiltinCommands {
		dlog.ErrorCheck(AddCommand("viewport", viewportCommands))
		dlog.ErrorCheck(AddCommand("fade", fadeCommands))
		dlog.ErrorCheck(AddCommand("skip", skipCommands(skipScene)))
		dlog.ErrorCheck(AddCommand("print", printCommands))
		dlog.ErrorCheck(AddCommand("mouse", mouseCommands))
		dlog.ErrorCheck(AddCommand("move", moveWindow))
		dlog.ErrorCheck(AddCommand("fullscreen", fullScreen))
		dlog.ErrorCheck(AddCommand("quit", func([]string) { Quit() }))
	}

	for {
		select {
		case <-resetCh: //reset all vars in debug console that save state
			viewportLocked = true
		default:
		}
		for scanner.Scan() {
			tokenString := strings.Fields(scanner.Text())
			if len(tokenString) == 0 {
				continue
			}
			if fn, ok := commands[tokenString[0]]; ok {
				fn(tokenString[1:])
			} else {
				fmt.Println("Unknown command", tokenString[0])
			}
		}
	}
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

func mouseDetails(nothing int, mevent interface{}) int {
	me := mevent.(mouse.Event)
	x := int(me.X()) + ViewPos.X
	y := int(me.Y()) + ViewPos.Y
	loc := collision.NewUnassignedSpace(float64(x), float64(y), 16, 16)
	results := collision.Hits(loc)
	fmt.Println("Mouse at:", x, y, "rel:", me.X(), me.Y())
	if len(results) == 0 {
		results = mouse.Hits(loc)
	}
	if len(results) > 0 {
		i := int(results[0].CID)
		if i > 0 && event.HasEntity(i) {
			e := event.GetEntity(i)
			spew.Dump(e)
		} else {
			fmt.Println("No entity ", i)
		}
	}

	return 0
}

func viewportCommands(tokenString []string) {
	switch tokenString[0] {
	case "unlock":
		if viewportLocked {
			speed := parseTokenAsInt(tokenString, 1, 5)
			viewportLocked = false
			event.GlobalBind(moveViewportBinding(speed), event.Enter)
		} else {
			fmt.Println("Viewport is already unbound")
		}
	case "lock":
		if viewportLocked {
			fmt.Println("Viewport is already locked")
		} else {
			viewportLocked = true
		}
	default:
		fmt.Println("Unrecognized command for viewport")
	}
}

func fadeCommands(tokenString []string) {
	toFade, ok := render.GetDebugRenderable(tokenString[0])
	fadeVal := parseTokenAsInt(tokenString, 1, 255)
	if ok {
		toFade.(render.Modifiable).Filter(mod.Fade(fadeVal))
	} else {
		fmt.Println("Could not fade input")
	}
}

func skipCommands(skipScene chan bool) func(tokenString []string) {
	return func(tokenString []string) {
		switch tokenString[0] {
		case "scene":
			skipScene <- true
		default:
			fmt.Println("Bad Skip Input")
		}
	}
}

func printCommands(tokenString []string) {
	if i, err := strconv.Atoi(tokenString[0]); err == nil {
		if i > 0 && event.HasEntity(i) {
			e := event.GetEntity(i)
			fmt.Println(reflect.TypeOf(e), e)
		} else {
			fmt.Println("No entity ", i)
		}
	} else {
		fmt.Println("Unable to parse", tokenString[0])
	}
}

func mouseCommands(tokenString []string) {
	switch tokenString[0] {
	case "details":
		event.GlobalBind(mouseDetails, "MouseRelease")
	default:
		fmt.Println("Bad Mouse Input")
	}
}

func moveWindow(in []string) {
	if len(in) < 4 {
		dlog.Error("Insufficient integer arguments for moving window")
		return
	}
	ints := make([]int, 4)
	var err error
	for i := range ints {
		ints[i], err = strconv.Atoi(in[i])
		if err != nil {
			dlog.Error(err)
			return
		}
	}
	err = MoveWindow(ints[0], ints[1], ints[2], ints[3])
	dlog.ErrorCheck(err)
}

func fullScreen(sub []string) {
	on := true
	if len(sub) > 0 {
		if sub[0] == "off" {
			on = false
		}
	}
	err := SetFullScreen(on)
	dlog.ErrorCheck(err)
}
