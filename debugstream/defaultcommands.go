package debugstream

import (
	"context"
	"io"
	"sync"

	"github.com/oakmound/oak/v3/window"
)

var (
	// DefaultCommands to attach to.
	DefaultCommands *ScopedCommands
	defaultsOnce    sync.Once
)

func checkOrCreateDefaults() {
	defaultsOnce.Do(func() {
		DefaultCommands = NewScopedCommands()
	})
}

// AddCommand to the default command set.
// See ScopedCommands' AddComand.
func AddCommand(c Command) error {
	checkOrCreateDefaults()
	return DefaultCommands.AddCommand(c)
}

// AttachToStream if possible to start consuming the stream
// and executing commands per the stored infomraiton in the ScopeCommands.
func AttachToStream(ctx context.Context, input io.Reader, output io.Writer) {
	checkOrCreateDefaults()
	DefaultCommands.AttachToStream(ctx, input, output)
}

// AddDefaultsForScope for debugging.
func AddDefaultsForScope(scopeID int32, controller interface{}) {
	checkOrCreateDefaults()
	if c, ok := controller.(window.Window); ok {
		DefaultCommands.AddDefaultsForScope(scopeID, c)
	}
}
