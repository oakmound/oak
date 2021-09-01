package debugstream

import (
	"bufio"
	"context"
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
	commands     map[int32]map[string]Command
}

// Command is a local format for performing these debug stream things.
type Command struct {
	Name      string
	ScopeID   int32
	Operation func([]string) string // the actual operation to execute
	Usage     string                // usage string, print when 'help' is called
	Force     bool                  // replace any existing command by this name
}

// NewScopedCommands creates set of standard help functions.
func NewScopedCommands() *ScopedCommands {
	sc := &ScopedCommands{commands: map[int32]map[string]Command{}}
	sc.AddCommand(Command{
		Name:      "help",
		Operation: sc.printHelp,
	})
	sc.AddCommand(Command{
		Name:      "scope",
		Usage:     explainAssumeScope,
		Operation: sc.assumeScope,
	})
	sc.AddCommand(Command{
		Name:      "fade",
		Usage:     explainFade,
		Operation: fadeCommands,
	})
	return sc
}

// AttachToStream and start executing the registered commands on input to said stream.
// Currently a given set of scoped commands may be attached once and only once. It will stop
// parsing commands when the provided context is done.
func (sc *ScopedCommands) AttachToStream(ctx context.Context, input io.Reader, out io.Writer) {
	sc.attachOnce.Do(
		func() {
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
					case <-ctx.Done():
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
								out.Write([]byte(fmt.Sprintf("unknown scopeID %d\n", scopeID)))
								continue
							}
							if len(tokenString) == 1 {
								out.Write([]byte(fmt.Sprintf("only provided scopeID %d without command\n", scopeID)))
								continue
							}

							tokenIDX++
							// see if specified
							potentialOp := strings.ToLower(tokenString[tokenIDX])
							if cmd, ok := sc.commands[scopeID][potentialOp]; ok {
								commandOut := cmd.Operation(tokenString[tokenIDX+1:])
								out.Write([]byte(commandOut))
								continue
							}
						}
						potentialOp := strings.ToLower(tokenString[0])
						// assumedscope
						if cmd, ok := sc.commands[sc.assumedScope][potentialOp]; ok {
							commandOut := cmd.Operation(tokenString[1:])
							out.Write([]byte(commandOut))
							continue
						}

						// fallback to scope 0
						if cmd, ok := sc.commands[0][potentialOp]; ok {
							commandOut := cmd.Operation(tokenString[1:])
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
func (sc *ScopedCommands) AddCommand(c Command) error {

	s := strings.ToLower(c.Name)
	scopeID := c.ScopeID

	sc.Lock()
	defer sc.Unlock()

	if _, ok := sc.commands[scopeID]; !ok {

		sc.commands[scopeID] = map[string]Command{}
		sc.scopes = append(sc.scopes, scopeID)
	}

	if _, ok := sc.commands[scopeID][s]; ok {
		if !c.Force {
			return oakerr.ExistingElement{
				InputName:   c.Name,
				InputType:   "string",
				Overwritten: false,
			}
		}
	}
	sc.commands[scopeID][s] = c
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
	sc.commands = map[int32]map[string]Command{}
}

// ResetCommandsForScope will throw out all existing debug commands from the
// debug console for hte given scope.
func (sc *ScopedCommands) ResetCommandsForScope(scope int32) {
	sc.commands[scope] = map[string]Command{}
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

// CommandsInScope returns the current debug console commands as a string array
func (sc *ScopedCommands) CommandsInScope(scope int32, showUsage bool) []string {

	cmds, ok := sc.commands[scope]
	if !ok {
		return []string{}
	}

	dkeys := make([]string, len(cmds))
	i := 0
	for k, command := range cmds {
		dkeys[i] = k + "\n"
		if showUsage && command.Usage != "" {
			dkeys[i] = fmt.Sprintf("%s: %s\n", k, command.Usage)
		}
		i++
	}
	sort.Strings(dkeys)
	return dkeys
}

// printHelp descriptions.
// Either for everything, a given scopeID, a given command, or a scopeID with a command.
func (sc *ScopedCommands) printHelp(tokenString []string) (out string) {
	scopeID := sc.assumedScope
	commandStr := ""
	var err error
	tknIndex := 0
	if len(tokenString) > 0 {

		// Check for a scope
		scopeID, err = strToInt32(tokenString[0])
		if err == nil {
			tknIndex++
		}
		// check for a command of interest
		if len(tokenString) > tknIndex {
			commandStr = tokenString[tknIndex]
		}
	}

	// error out if the scopeID is invalid for one reason or another
	if _, ok := sc.commands[scopeID]; !ok {
		if scopeID == sc.assumedScope {
			out += fmt.Sprintf("current scope %v is not usable please use 'scope 0' or 'scope\n", scopeID)
		} else {
			out += fmt.Sprintf("inactive scope %d see correct usage by using help without the scope\n", scopeID)
		}
		return
	}

	if scopeID == sc.assumedScope {
		out += "help <scopeID> to see commands linked to a given window\n" +
			fmt.Sprintf("Active Scopes: %v\n", sc.scopes)
	}

	out += fmt.Sprintf("Current Assumed Scope: %v\n", sc.assumedScope) // TODO: if in a verbose mode present usage.

	// give a general overview if a specific command is not specified
	if commandStr == "" {
		out += fmt.Sprintf("General Commands:\n%s%s\n", indent, strings.Join(sc.CommandsInScope(0, true), indent))
		if scopeID != 0 {
			out += fmt.Sprintf("Current Window Commands:\n%s%s\n", indent, strings.Join(sc.CommandsInScope(scopeID, true), indent))
		}
		return
	}

	out += fmt.Sprintf("Registered Instances of %s\n", commandStr)
	// return just the usage for the given command
	for scope, cmdSet := range sc.commands {
		if c, ok := cmdSet[commandStr]; ok {
			out += fmt.Sprintf("%sscope%v %s: %s\n", indent, scopeID, commandStr, c.Usage)
		} else {
			if scope == scopeID {
				out += fmt.Sprintf("%sWarning scope '%v' did not have the specified command %q\n", indent, scopeID, commandStr)
			}
		}

	}

	return
}

const indent = "  "
const explainAssumeScope = "provide a scopeID to use commands without a scopeID prepended"

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
			fmt.Sprintf("you provided %q which errored with %v\n", tokenString[0], err)
		return
	}
	if _, ok := sc.commands[scopeID]; !ok {
		out += fmt.Sprintf("inactive scope %d\n", scopeID)
		return
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
