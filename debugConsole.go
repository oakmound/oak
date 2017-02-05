package oak

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	// "bitbucket.org/oakmoundstudio/oak/dlog"
	"reflect"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
)

var (
	viewportLocked = true
)

func DebugConsole(resetCh, skipScene chan bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-resetCh: //reset all vars in debug console that save state
			viewportLocked = true
		default:
		}
		for scanner.Scan() {
			select {
			case <-resetCh: //reset all vars in debug console that save state
				viewportLocked = true
			default:
			}
			//Parse the Input
			tokenString := strings.Fields(scanner.Text())
			switch tokenString[0] {
			case "viewport":
				switch tokenString[1] {
				case "unlock":
					if viewportLocked {
						speed := parseTokenAsInt(tokenString, 2, 5)
						viewportLocked = false
						event.GlobalBind(moveViewportBinding(speed), "EnterFrame")
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

			case "fade":
				if len(tokenString) > 1 {
					toFade, ok := render.GetDebugRenderable(tokenString[1])
					fadeVal := parseTokenAsInt(tokenString, 2, 255)
					if ok {
						toFade.(render.Modifiable).Fade(fadeVal)
					} else {
						fmt.Println("Could not fade input")
					}
				} else {
					fmt.Println("Unrecognized length for fade")
				}
			case "skip":
				if len(tokenString) > 1 {
					switch tokenString[1] {
					case "scene":
						skipScene <- true

					default:
						fmt.Println("Bad Skip Input")
					}
				}
			case "print":
				if len(tokenString) > 1 {
					if i, err := strconv.Atoi(tokenString[1]); err == nil {
						if i > 0 && event.HasEntity(i) {
							e := event.GetEntity(i)
							fmt.Println(reflect.TypeOf(e), e)
						} else {
							fmt.Println("No entity ", i)
						}
					} else {
						fmt.Println("Unable to parse", tokenString[1])
					}
				}
			case "mouse":
				if len(tokenString) > 1 {
					switch tokenString[1] {
					case "details":
						event.GlobalBind(mouseDetails, "MouseRelease")
					default:
						fmt.Println("Bad Mouse Input")
					}
				}
			default:
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
	me := mevent.(mouse.MouseEvent)
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
			fmt.Println(reflect.TypeOf(e), e)
		} else {
			fmt.Println("No entity ", i)
		}
	}

	return 0
}
