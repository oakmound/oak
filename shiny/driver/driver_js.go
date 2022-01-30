//go:build js && !nooswindow && !windows && !darwin && !linux
// +build js,!nooswindow,!windows,!darwin,!linux,!android

package driver

import (
	"github.com/oakmound/oak/v3/shiny/driver/jsdriver"
	"github.com/oakmound/oak/v3/shiny/screen"
)

func main(f func(screen.Screen)) {
	jsdriver.Main(f)
}

func monitorSize() (int, int) {
	return 0, 0
}
