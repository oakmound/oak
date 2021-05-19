package oak

import (
	"github.com/oakmound/oak/v3/shiny/screen"
)

// A Driver is a function which can take in our lifecycle function
// and initialize oak with the OS interfaces it needs.
type Driver func(f func(screen.Screen))
