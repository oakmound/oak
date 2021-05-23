// +build noop

package driver

import (
	"github.com/oakmound/oak/v3/shiny/driver/noop"
	"github.com/oakmound/oak/v3/shiny/screen"
)

func main(f func(screen.Screen)) {
	noop.Main(f)
}

func monitorSize() (int, int) {
	return 0, 0
}
