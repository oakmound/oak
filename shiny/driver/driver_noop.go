//go:build nooswindow
// +build nooswindow

package driver

import (
	"github.com/oakmound/oak/v4/shiny/driver/noop"
	"github.com/oakmound/oak/v4/shiny/screen"
)

func main(f func(screen.Screen)) {
	noop.Main(f)
}

type Window = noop.Window
