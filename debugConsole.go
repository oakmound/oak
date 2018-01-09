package oak

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

var (
	viewportLocked  = true
	commands        = make(map[string]func([]string))
	builtinCommands = make(map[string]func([]string))
)

// AddCommand adds a console command to call fn when
// 'c <s> <args>' is input to the console
func AddCommand(s string, fn func([]string)) {
	commands[s] = fn
}

func debugConsole(resetCh, skipScene chan bool, input io.Reader) {
	scanner := bufio.NewScanner(input)
	spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 2

	builtinCommands = map[string]func([]string){
		"viewport": viewportCommands,
		"fade":     fadeCommands,
		"skip":     skipCommands(skipScene),
		"print":    printCommands,
		"mouse":    mouseCommands,
	}

	for {
		select {
		case <-resetCh: //reset all vars in debug console that save state
			viewportLocked = true
		default:
		}
		for scanner.Scan() {
			//Parse the Input
			tokenString := strings.Fields(scanner.Text())
			if len(tokenString) < 2 {
				continue
			}

			// The builtin commands should probably be split off, so that
			// they aren't on by default always. It's worth considering making
			// all commands through the AddCommand function and removing the
			// requirement to precede custom commands with 'c', which would
			// then require that we return an error for overwriting old command
			// names with new commands.
			if tokenString[0] == "c" || tokenString[0] == "cheat" {
				// Requires that cheats are all one word! <-- don't forget
				if fn, ok := commands[tokenString[1]]; ok {
					fn(tokenString[1:])
				} else {
					fmt.Println("Unknown command", tokenString[1])
				}
			} else if fn, ok := builtinCommands[tokenString[0]]; ok {
				fn(tokenString[1:])
			} else {
				fmt.Println("Unrecognized Input")
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
	x := int(me.X) + ViewPos.X
	y := int(me.Y) + ViewPos.Y
	loc := collision.NewUnassignedSpace(float64(x), float64(y), 16, 16)
	results := collision.Hits(loc)
	fmt.Println("Mouse at:", x, y, "rel:", me.X, me.Y)
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
		toFade.(render.Modifiable).Modify(render.Fade(fadeVal))
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
