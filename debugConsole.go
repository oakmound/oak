package oak

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/oakmound/oak/v3/oakerr"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
)

// AddCommand adds a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (c *Controller) AddCommand(s string, fn func([]string)) error {
	return c.addCommand(s, fn, false)
}

// ForceAddCommand adds or overwrites a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (c *Controller) ForceAddCommand(s string, fn func([]string)) {
	c.addCommand(s, fn, true)
}

func (c *Controller) addCommand(s string, fn func([]string), force bool) error {
	if _, ok := c.commands[s]; ok {
		if !force {
			return oakerr.ExistingElement{
				InputName:   "s",
				InputType:   "string",
				Overwritten: false,
			}
		}
	}
	c.commands[s] = fn
	return nil
}

// ClearCommand clears an existing debug command by key: <s>
func (c *Controller) ClearCommand(s string) {
	delete(c.commands, s)
}

// ResetCommands will throw out all existing debug commands from the
// debug console.
func (c *Controller) ResetCommands() {
	c.commands = map[string]func([]string){}
}

// GetDebugKeys returns the current debug console commands as a string array
func (c *Controller) GetDebugKeys() []string {
	dkeys := make([]string, len(c.commands))
	i := 0
	for k := range c.commands {
		dkeys[i] = k
		i++
	}
	return dkeys
}

func (c *Controller) debugConsole(input io.Reader) {
	scanner := bufio.NewScanner(input)

	// built in commands
	if c.config.LoadBuiltinCommands {
		dlog.ErrorCheck(c.AddCommand("fade", c.fadeCommands))
		dlog.ErrorCheck(c.AddCommand("skip", c.skipCommands))
		dlog.ErrorCheck(c.AddCommand("print", c.printCommands))
		dlog.ErrorCheck(c.AddCommand("mouse", c.mouseCommands))
		dlog.ErrorCheck(c.AddCommand("move", c.moveWindow))
		dlog.ErrorCheck(c.AddCommand("fullscreen", c.fullScreen))
		dlog.ErrorCheck(c.AddCommand("help", c.printDebugCommands))
		dlog.ErrorCheck(c.AddCommand("quit", func([]string) { c.Quit() }))
	}

	for {
		for scanner.Scan() {
			tokenString := strings.Fields(scanner.Text())
			if len(tokenString) == 0 {
				continue
			}
			if fn, ok := c.commands[tokenString[0]]; ok {
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

func (c *Controller) mouseDetails(nothing event.CID, mevent interface{}) int {
	me := mevent.(*mouse.Event)
	x := int(me.X()) + c.viewPos[0]
	y := int(me.Y()) + c.viewPos[1]
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

	return 0
}

func (c *Controller) fadeCommands(tokenString []string) {
	if len(tokenString) > 0 {
		fmt.Println("Input must start with the name of the renderable to fade")
		return
	}
	toFade, ok := render.GetDebugRenderable(tokenString[0])
	if ok {
		fadeVal := parseTokenAsInt(tokenString, 1, 255)
		toFade.(render.Modifiable).Filter(mod.Fade(fadeVal))
	} else {
		fmt.Println("Could not fade input")
	}
}

func (c *Controller) skipCommands(tokenString []string) {
	if len(tokenString) != 1 {
		fmt.Println("Input must be a single string from the following (\"scene\"). ")
		return
	}
	switch tokenString[0] {
	case "scene":
		c.NextScene()
	default:
		fmt.Println("Bad Skip Input")
	}
}

func (c *Controller) printCommands(tokenString []string) {
	if len(tokenString) != 1 {
		fmt.Println("Input must be a single number that corresponds to an entity.")
		return
	}
	i, err := strconv.Atoi(tokenString[0])
	if err != nil {
		fmt.Println("Unable to parse", tokenString[0])
		return
	}
	if i > 0 && event.HasEntity(event.CID(i)) {
		e := event.GetEntity(event.CID(i))
		fmt.Println(reflect.TypeOf(e), e)
	} else {
		fmt.Println("No entity ", i)
	}

}

func (c *Controller) mouseCommands(tokenString []string) {
	if len(tokenString) != 1 {
		fmt.Println("Input must be a single string from the following (\"details\") ")
		return
	}
	switch tokenString[0] {
	case "details":
		event.GlobalBind("MouseRelease", c.mouseDetails)
	default:
		fmt.Println("Bad Mouse Input")
	}
}

func (c *Controller) printDebugCommands(tokenString []string) {
	dbgKeys := c.GetDebugKeys()
	sort.Strings(dbgKeys)
	fmt.Printf("Commands: %s\n", strings.Join(dbgKeys, ", "))
}

func (c *Controller) moveWindow(in []string) {
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
	//err = c.MoveWindow(ints[0], ints[1], ints[2], ints[3])
	dlog.ErrorCheck(err)
}

func (c *Controller) fullScreen(sub []string) {
	// on := true
	// if len(sub) > 0 {
	// 	if sub[0] == "off" {
	// 		on = false
	// 	}
	// }
	//err := c.SetFullScreen(on)
}

// RunCommand runs a command added with AddCommand.
// It's intended use is making it easier to
// alias commands/subcommands.
// It returns an error if the command doesn't exist.
func (c *Controller) RunCommand(cmd string, args ...string) error {
	fn, ok := c.commands[cmd]
	if !ok {
		return fmt.Errorf("Unknown command %s", cmd)
	}
	fn(args)
	return nil
}
