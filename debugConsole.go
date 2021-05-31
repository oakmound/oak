package oak

import (
	"io"

	"github.com/oakmound/oak/v3/debugstream"
)

func (c *Controller) debugConsole(input io.Reader) {
	debugstream.AttachToStream(input)
	debugstream.AddDefaultsForScope(c.ControllerID, c)

}
