//go:build js && !nooswindow && !windows && !darwin && !linux
// +build js,!nooswindow,!windows,!darwin,!linux

package driver

import (
	"github.com/oakmound/oak/v4/shiny/driver/jsdriver"
	"github.com/oakmound/oak/v4/shiny/screen"
)

func main(f func(screen.Screen)) {
	jsdriver.Main(f)
}

type Window = jsdriver.Window
