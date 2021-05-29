package oak

import (
	"io"

	"github.com/oakmound/oak/v3/debugstream"
)

func (c *Controller) debugConsole(input io.Reader) {
	debugstream.DefaultCommands.AttachToStream(input)
	debugstream.DefaultCommands.AddDefaultsForScope(c.ControllerID, c)
}
