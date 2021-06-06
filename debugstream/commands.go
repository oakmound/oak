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
	commands     map[int32]map[string]command
	closeChan    chan struct{}
}

// command is a local format for performing these debug stream things.
type command struct {
	name      string
	operation func([]string) string // the actual operation to execute
	usage     func([]string) string // optional args to give scoped usage
}

func newCommand(name string, usage func([]string) string, operation func([]string) string) command {
	return command{name, operation, usage}
}

// NewScopedCommands creates set of standard help functions.
func NewScopedCommands() *ScopedCommands {
	sc := &ScopedCommands{commands: map[int32]map[string]command{}}
	sc.AddCommand("help", nil, sc.printHelp)
	sc.AddCommand("scope", sc.explainAssumeScope, sc.assumeScope)
	sc.AddCommand("fade", nil, fadeCommands)
	return sc
}

func (sc *ScopedCommands) UnAttachFromStream() {
	if sc.closeChan == nil {
		return
	}
	close(sc.closeChan)
}

// AttachToStream and start executing the registered commands on input to said stream.
// Currently a given set of scoped commands may be attached once and only once.
func (sc *ScopedCommands) AttachToStream(input io.Reader, out io.Writer) {
	if sc == nil {
		return
	}
	sc.attachOnce.Do(

		func() {
			sc.closeChan = make(chan struct{})
			textIn := make(chan string)
			go func(textBuffer chan string, in io.Reader) {
				scanner := bufio.NewScanner(in)
				for {

					for scanner.Scan() {
						textBuffer <- scanner.Text()
					}
				}
			}(textIn, input)

			go func() {

				for {

					select {
					case <-sc.closeChan:
						out.Write([]byte("stopping debugstream"))
						return

					case scanText := <-textIn:

						// TODO: accept interrupts

						tokenString := strings.Fields(scanText)
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
								out.Write([]byte(fmt.Sprintf("unknown scopeID %d see correct usage via help\n", scopeID)))
								continue
							}
							if len(tokenString) == 1 {
								out.Write([]byte(fmt.Sprintf("Only provided scopeID %d without see usage via help\n", scopeID)))
								continue
							}

							tokenIDX++
							// see if specified
							potentialOp := strings.ToLower(tokenString[tokenIDX])
							if cmd, ok := sc.commands[scopeID][potentialOp]; ok {
								commandOut := cmd.operation(tokenString[tokenIDX+1:])
								out.Write([]byte(commandOut))
								continue
							}
						}
						potentialOp := strings.ToLower(tokenString[0])
						// assumedscope
						if cmd, ok := sc.commands[sc.assumedScope][potentialOp]; ok {
							commandOut := cmd.operation(tokenString[1:])
							out.Write([]byte(commandOut))
							continue
						}

						// fallback to scope 0
						if cmd, ok := sc.commands[0][potentialOp]; ok {
							commandOut := cmd.operation(tokenString[1:])
							out.Write([]byte(commandOut))
							continue
						}

						out.Write([]byte(fmt.Sprintf("Unknown command '%s' for scopeID %d see correct usage via help or help %d\n", tokenString[tokenIDX], scopeID, scopeID)))
						suggestions := sc.suggestForCandidate(4, tokenString[tokenIDX])
						if len(suggestions) > 0 {
							out.Write([]byte("Did you mean one of the following?\n"))
							for _, s := range suggestions {
								out.Write([]byte(indent + s + "\n"))
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
func (sc *ScopedCommands) AddCommand(s string, usageFn func([]string) string, fn func([]string) string) error {
	// We tightly link to controllerIDs here for better or for worse and controllerIDs shall always be > 0
	// This means our unscoped commands are safe when set as 0.
	return sc.AddScopedCommand(0, s, usageFn, fn)
}

// AddScopedCommand adds a console command to call fn when
// '<s> <args>' is input to the console. fn will be called
// with args split on whitespace.
func (sc *ScopedCommands) AddScopedCommand(scopeID int32, s string, usageFn func([]string) string, fn func([]string) string) error {
	return sc.addCommand(scopeID, s, usageFn, fn, false)
}

// addCommand for future executions.
func (sc *ScopedCommands) addCommand(scopeID int32, sname string, usageFn func([]string) string, fn func([]string) string, force bool) error {

	s := strings.ToLower(sname)

	sc.Lock()
	defer sc.Unlock()

	if _, ok := sc.commands[scopeID]; !ok {

		sc.commands[scopeID] = map[string]command{}
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
	sc.commands[scopeID][s] = newCommand(s, usageFn, fn)

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
	sc.commands = map[int32]map[string]command{}
}

// ResetCommandsForScope will throw out all existing debug commands from the
// debug console for hte given scope.
func (sc *ScopedCommands) ResetCommandsForScope(scope int32) {
	sc.commands[scope] = map[string]command{}
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
func (sc *ScopedCommands) CommandsInScope(scope int32, showUsage bool) []string {

	cmds, ok := sc.commands[scope]
	if !ok {
		return []string{}
	}

	dkeys := make([]string, len(cmds))
	i := 0
	empty := []string{}
	for k, command := range cmds {
		dkeys[i] = k + "\n"
		if showUsage && command.usage != nil {
			dkeys[i] = fmt.Sprintf("%s: %s", k, command.usage(empty))
		}
		i++
	}
	return dkeys
}

func (sc *ScopedCommands) printHelp(tokenString []string) (out string) {
	scopeID := sc.assumedScope
	var err error
	if len(tokenString) != 0 {
		scopeID, err = strToInt32(tokenString[0])
		if err != nil {
			out += "help <scopeID> expects a valid int scope\n" +
				fmt.Sprintf("you provided %s which errored with %v \n", tokenString[0], err) +
				"try using help without arguments for an overview\n"
			return
		}
	}

	out += "help <scopeID> to see commands linked to a given window\n" +
		fmt.Sprintf("Active Scopes: %v\n", sc.scopes)
	if _, ok := sc.commands[scopeID]; !ok {
		out += fmt.Sprintf("inactive scope %d see correct usage via help\n", scopeID)
		return
	}

	out += fmt.Sprintf("Current Assumed Scope: %v\n", sc.assumedScope)
	// TODO: if in a verbose mode present usage.
	out += fmt.Sprintf("General Commands:\n%s%s\n", indent, strings.Join(sc.CommandsInScope(0, true), indent))
	if scopeID != 0 {
		out += fmt.Sprintf("Current Window Commands:\n%s%s\n", indent, strings.Join(sc.CommandsInScope(scopeID, true), indent))
	}
	return
}

const indent = "  "

func (sc *ScopedCommands) explainAssumeScope([]string) string {
	return fmt.Sprintf("provide a scopeID to use commands without a scopeID prepended. Current Scopes are: %v\n", sc.scopes)
}

// assumeScope of the given windowID if possible
// This allows for easier usage of windows when multiple windows exist.
func (sc *ScopedCommands) assumeScope(tokenString []string) (out string) {
	if len(tokenString) == 0 {
		out += "assume scope requires a scopeID\n" +
			fmt.Sprintf("Active Scopes: %v\n", sc.scopes)
		return
	}

	scopeID, err := strToInt32(tokenString[0])
	if err != nil {
		out += "assume scope <scopeID> expects a valid int32 scope\n" +
			fmt.Sprintf("you provided %s which errored with %v \n", tokenString[0], err)
		return
	}
	if _, ok := sc.commands[scopeID]; !ok {
		out += fmt.Sprintf("inactive scope %d see correct usage via help\n", scopeID)
	}
	sc.assumedScope = scopeID
	out += fmt.Sprintf("assumed scope %v\n", scopeID)
	return
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
