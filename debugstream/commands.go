package debugstream

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/oakmound/oak/v3/oakerr"
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
func (sc *ScopedCommands) AttachToStream(input io.Reader) {
	if sc == nil {
		return
	}
	sc.attachOnce.Do(
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
							_, ok := sc.commands[scopeID]
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
							if fn, ok := sc.commands[scopeID][tokenString[tokenIDX]]; ok {
								fn(tokenString[tokenIDX+1:])
								continue
							}
						}

						// assumedscope
						if fn, ok := sc.commands[sc.assumedScope][tokenString[0]]; ok {
							fn(tokenString[1:])
							continue
						}

						// fallback to scope 0
						if fn, ok := sc.commands[0][tokenString[0]]; ok {
							fn(tokenString[1:])
							continue
						}

						fmt.Printf("Unknown command '%s' for scopeID %d see correct usage via help or help %d\n", tokenString[tokenIDX], scopeID, scopeID)
						suggestions := sc.suggestForCandidate(4, tokenString[tokenIDX])
						if len(suggestions) > 0 {
							fmt.Println("Did you mean one of the following?")
							for _, s := range suggestions {
								fmt.Println(indent, s)
							}
						}
					}
				}
			}()
		})
}

// AddCommand adds a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (sc *ScopedCommands) AddCommand(s string, fn func([]string)) error {
	// We tightly link to controllerIDs here for better or for worse and controllerIDs shall always be > 0
	// This means our unscoped commands are safe when set as 0.
	return sc.AddScopedCommand(0, s, fn)
}

// AddScopedCommand adds a console command for a given window scope to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (sc *ScopedCommands) AddScopedCommand(scopeID int32, s string, fn func([]string)) error {
	return sc.addCommand(scopeID, s, fn, false)
}

// addCommand for future executions.
func (sc *ScopedCommands) addCommand(scopeID int32, s string, fn func([]string), force bool) error {

	sc.Lock()
	defer sc.Unlock()

	if _, ok := sc.commands[scopeID]; !ok {

		sc.commands[scopeID] = map[string]func([]string){}
		sc.scopes = append(sc.scopes, scopeID)
	}

	if _, ok := sc.commands[scopeID][s]; ok {
		if !force {
			return oakerr.ExistingElement{
				InputName:   s,
				InputType:   "string",
				Overwritten: false,
			}
		}
	}
	sc.commands[scopeID][s] = fn

	return nil
}

// ClearCommand clears an existing debug command for scope with key: <s>
func (sc *ScopedCommands) ClearCommand(scopeID int32, s string) {
	_, ok := sc.commands[scopeID]
	if !ok {
		return
	}
	delete(sc.commands[scopeID], s)
}

// ResetCommands will throw out all existing debug commands from the
// debug console.
func (sc *ScopedCommands) ResetCommands() {
	sc.commands = map[int32]map[string]func([]string){}
}

// ResetCommandsForScope will throw out all existing debug commands from the
// debug console for hte given scope.
func (sc *ScopedCommands) ResetCommandsForScope(scope int32) {
	sc.commands[scope] = map[string]func([]string){}
}

// RemoveScope from the command set.
// Usually done on the close of a scope.
func (sc *ScopedCommands) RemoveScope(scope int32) {
	delete(sc.commands, scope)
	for i := 0; i < len(sc.scopes); i++ {
		if sc.scopes[i] == scope {
			sc.scopes = append(sc.scopes[:i], sc.scopes[i+1:]...)
			return
		}
	}
}

// GetDebugKeys returns the current debug console commands as a string array
func (sc *ScopedCommands) CommandsInScope(scope int32) []string {

	cmds, ok := sc.commands[scope]
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

func (sc *ScopedCommands) printHelp(tokenString []string) {

	scopeID := sc.assumedScope
	var err error
	if len(tokenString) != 0 {
		scopeID, err = strToInt32(tokenString[0])
		if err != nil {
			fmt.Println("help <scopeID> expects a valid int scope")
			fmt.Printf("you provided %s which errored with %v \n", tokenString[0], err)
			fmt.Println("try using help without arguments for an overview")
			return
		}
	}

	fmt.Println("help <scopeID> to see commands linked to a given window")
	fmt.Printf("Active Scopes: %v\n", sc.scopes)
	if _, ok := sc.commands[scopeID]; !ok {
		fmt.Printf("inactive scope %d see correct usage via help\n", scopeID)
	}

	fmt.Println("Current Assumed Scope:", sc.assumedScope)
	// TODO: if in a verbose mode present usage.

	fmt.Printf("General Commands:\n%s%s\n", indent, strings.Join(sc.CommandsInScope(0), "\n"+indent))
	fmt.Printf("Current Window Commands:\n%s%s\n\n", indent, strings.Join(sc.CommandsInScope(scopeID), "\n"+indent))
}

const indent = "  "

// assumeScope of the given windowID if possible
// This allows for easier usage of windows when multiple windows exist.
func (sc *ScopedCommands) assumeScope(tokenString []string) {
	if len(tokenString) == 0 {
		fmt.Println("assume scope requires a scopeID or -")
		fmt.Printf("Active Scopes: %v\n", sc.scopes)
		return
	}

	scopeID, err := strToInt32(tokenString[0])
	if err != nil {
		fmt.Println("assume scope <scopeID> expects a valid int32 scope")
		fmt.Printf("you provided %s which errored with %v \n", tokenString[0], err)
		return
	}
	if _, ok := sc.commands[scopeID]; !ok {
		fmt.Printf("inactive scope %d see correct usage via help\n", scopeID)
	}
	sc.assumedScope = scopeID
	fmt.Println("assumed scope ", scopeID)
}

func (sc *ScopedCommands) suggestForCandidate(maxSuggestions int, candidate string) (suggestions []string) {

	possibilities := []candidateStore{}
	scopes := []int32{sc.assumedScope}
	if sc.assumedScope != 0 {
		scopes = append(scopes, 0)
	}
	for _, s := range scopes {
		for c := range sc.commands[s] {
			_, val := jaroDecreased(candidate, c)
			if val > suggestionCuttOff {
				possibilities = append(possibilities, candidateStore{c, val})
			}
		}
	}

	sort.Slice(possibilities, func(i, j int) bool {
		return possibilities[i].value > possibilities[j].value
	})

	maxS := maxSuggestions
	if len(possibilities) <= maxS {
		maxS = len(possibilities)
	}
	for i := 0; i < maxS; i++ {
		suggestions = append(suggestions, possibilities[i].name)
	}

	return suggestions
}

const suggestionCuttOff = 0.4

type candidateStore struct {
	name  string
	value float64
}

func strToInt32(potentialInt string) (int32, error) {
	i64, err := strconv.ParseInt(potentialInt, 10, 32)
	return int32(i64), err
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
