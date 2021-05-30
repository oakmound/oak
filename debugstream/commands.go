package debugstream

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/oakmound/oak/v3/oakerr"
)

var (
	// DefaultCommands to attach to. TODO: init should be lazy.
	DefaultCommands = NewScopedCommands()
)

// ScopedCommands for the debug stream commands.
// Contains a set of scopes that align with oak.Controller.
// Currently can only be attached to a single stream
type ScopedCommands struct {
	sync.Mutex
	attachOnce   sync.Once
	assumedScope int32
	scopes       []int32
	commands     map[int32]map[string]func([]string)
}

// NewScopedCommands creates set of standard help functions.
func NewScopedCommands() *ScopedCommands {
	sc := &ScopedCommands{commands: map[int32]map[string]func([]string){}}
	sc.AddCommand("help", sc.printHelp)
	sc.AddCommand("scope", sc.assumeScope)
	sc.AddCommand("fade", fadeCommands)
	return sc
}

// AttachToStream and start executing the registered commands on input to said stream.
// Currently a given set of scoped commands may be attached once and only once.
func (c *ScopedCommands) AttachToStream(input io.Reader) {
	if c == nil {
		return
	}
	c.attachOnce.Do(
		func() {
			scanner := bufio.NewScanner(input)
			go func() {
				for {
					for scanner.Scan() {

						// TODO: accept interrupts

						tokenString := strings.Fields(scanner.Text())
						if len(tokenString) == 0 {
							continue
						}
						// Attempt to parse the first arg as a scope
						scopeID, err := strToInt32(tokenString[0])

						tokenIDX := 0
						// if there was a scope specified then increment what we care about
						if err == nil {
							_, ok := c.commands[scopeID]
							if !ok {
								fmt.Printf("unknown scopeID %d see correct usage via help\n", scopeID)
								continue
							}
							if len(tokenString) == 1 {
								fmt.Printf("Only provided scopeID %d without see usage via help\n", scopeID)
								continue
							}

							tokenIDX++
							// see if specified
							if fn, ok := c.commands[scopeID][tokenString[tokenIDX]]; ok {
								fn(tokenString[tokenIDX+1:])
								continue
							}
						}

						// assumedscope
						if fn, ok := c.commands[c.assumedScope][tokenString[0]]; ok {
							fn(tokenString[1:])
							continue
						}

						fmt.Println("fall ", tokenString)
						// fallback to scope 0
						if fn, ok := c.commands[0][tokenString[0]]; ok {
							fn(tokenString[1:])
							continue
						}

						fmt.Printf("Unknown command '%s' for scopeID %d see correct usage via help or help %d\n", tokenString[tokenIDX], scopeID, scopeID)

					}
				}
			}()
		})
}

// AddCommand adds a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (c *ScopedCommands) AddCommand(s string, fn func([]string)) error {
	// We tightly link to controllerIDs here for better or for worse and controllerIDs shall always be > 0
	// This means our unscoped commands are safe when set as 0.
	return c.AddScopedCommand(0, s, fn)
}

func AddCommand(s string, fn func([]string)) error {
	return DefaultCommands.AddCommand(s, fn)
}

// AddCommand adds a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (c *ScopedCommands) AddScopedCommand(scopeID int32, s string, fn func([]string)) error {
	return c.addCommand(scopeID, s, fn, false)
}

func AddScopedCommand(scopeID int32, s string, fn func([]string)) error {
	return DefaultCommands.AddScopedCommand(scopeID, s, fn)
}

// addCommand for future executions.
func (c *ScopedCommands) addCommand(scopeID int32, s string, fn func([]string), force bool) error {

	c.Lock()
	defer c.Unlock()

	if _, ok := c.commands[scopeID]; !ok {

		c.commands[scopeID] = map[string]func([]string){}
		c.scopes = append(c.scopes, scopeID)
	}

	if _, ok := c.commands[scopeID][s]; ok {
		if !force {
			return oakerr.ExistingElement{
				InputName:   s,
				InputType:   "string",
				Overwritten: false,
			}
		}
	}
	c.commands[scopeID][s] = fn

	return nil
}

// ClearCommand clears an existing debug command for scope with key: <s>
func (c *ScopedCommands) ClearCommand(scopeID int32, s string) {
	_, ok := c.commands[scopeID]
	if !ok {
		return
	}
	delete(c.commands[scopeID], s)
}

// ResetCommands will throw out all existing debug commands from the
// debug console.
func (c *ScopedCommands) ResetCommands() {
	c.commands = map[int32]map[string]func([]string){}
}

// ResetCommandsForScope will throw out all existing debug commands from the
// debug console for hte given scope.
func (c *ScopedCommands) ResetCommandsForScope(scope int32) {
	c.commands[scope] = map[string]func([]string){}
}

// RemoveScope from the command set.
// Usually done on the close of a scope.
func (c *ScopedCommands) RemoveScope(scope int32) {
	delete(c.commands, scope)
	for i := 0; i < len(c.scopes); i++ {
		if c.scopes[i] == scope {
			c.scopes = append(c.scopes[:i], c.scopes[i+1:]...)
			return
		}
	}
}

// GetDebugKeys returns the current debug console commands as a string array
func (c *ScopedCommands) CommandsInScope(scope int32) []string {

	cmds, ok := c.commands[scope]
	if !ok {
		return []string{}
	}

	dkeys := make([]string, len(cmds))
	i := 0
	for k := range cmds {
		dkeys[i] = k
		i++
	}
	return dkeys
}

func (c *ScopedCommands) printHelp(tokenString []string) {
	if len(tokenString) == 0 {
		fmt.Println("help <scopeID> to see commands linked to a window or help 0 to see general commands")
		fmt.Printf("Active Scopes: %v\n", c.scopes)
		return
	}

	scopeID, err := strToInt32(tokenString[0])
	if err != nil {
		fmt.Println("help <scopeID> expects a valid int scope")
		fmt.Printf("you provided %s which errored with %v \n", tokenString[0], err)
		return
	}
	if _, ok := c.commands[scopeID]; !ok {
		fmt.Printf("inactive scope %d see correct usage via help\n", scopeID)
	}
	fmt.Println("Current Assumed Scope:", c.assumedScope)
	fmt.Printf("Commands: %s\n", strings.Join(c.CommandsInScope(scopeID), ","))
}

// assumeScope of the given windowID if possible
// This allows for easier usage of windows when multiple windows exist.
func (c *ScopedCommands) assumeScope(tokenString []string) {
	if len(tokenString) == 0 {
		fmt.Println("assume scope requires a scopeID or -")
		fmt.Printf("Active Scopes: %v\n", c.scopes)
		return
	}

	scopeID, err := strToInt32(tokenString[0])
	if err != nil {
		fmt.Println("assume scope <scopeID> expects a valid int32 scope")
		fmt.Printf("you provided %s which errored with %v \n", tokenString[0], err)
		return
	}
	if _, ok := c.commands[scopeID]; !ok {
		fmt.Printf("inactive scope %d see correct usage via help\n", scopeID)
	}
	c.assumedScope = scopeID
	fmt.Println("assumed scope ", scopeID)
}

// strToInt32 helps align to the window scope handles
func strToInt32(potentialInt string) (int32, error) {
	i64, err := strconv.ParseInt(potentialInt, 10, 32)
	return int32(i64), err
}

// parseTokenAsInt is a convience function for parsing from our string slice of strings
func parseTokenAsInt(tokenString []string, arrIndex int, defaultVal int) int {
	if len(tokenString) > arrIndex {
		tmp, err := strconv.Atoi(tokenString[arrIndex])
		if err == nil {
			return tmp
		}
	}
	return defaultVal
}
