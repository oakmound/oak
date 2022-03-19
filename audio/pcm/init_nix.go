//go:build linux || darwin

package pcm

import (
	"github.com/jfreymuth/pulse"
	"github.com/oakmound/oak/v3/oakerr"
)

func initOS(driver Driver) error {
	switch driver {
	case DriverDefault:
		fallthrough
	case DriverPulse:
		// Sanity check that pulse is installed and a sink is defined
		client, err := pulse.NewClient()
		if err != nil {
			// osx: brew install pulseaudio
			// linux: sudo apt install pulseaudio
			return oakerr.UnsupportedPlatform{
				Operation: "pcm.Init:" + driver.String(),
			}
		}
		defer client.Close()
		_, err = client.DefaultSink()
		if err != nil {
			return err
		}
	default:
		return oakerr.UnsupportedPlatform{
			Operation: "pcm.Init:" + driver.String(),
		}
	}
	return nil
}
