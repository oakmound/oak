//go:build linux

package audio

import (
	"fmt"
	"os"

	"github.com/jfreymuth/pulse"
	"github.com/oakmound/oak/v3/audio/pcm"
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
		newWriter = newPulseWriter
	case DriverALSA:
		//???
		newWriter = newALSAWriter
		if skipDevices := os.Getenv("OAK_SKIP_AUDIO_DEVICES"); skipDevices != "" {
			SkipDevicesContaining = skipDevices
		}
	default:
		return oakerr.UnsupportedPlatform{
			Operation: "pcm.Init:" + driver.String(),
		}
	}
	return nil
}

var newWriter = func(f pcm.Format) (pcm.Writer, error) {
	return nil, fmt.Errorf("this package has not been initialized")
}

// TODO: do other drivers need this? Can we pick devices more intelligently?
var SkipDevicesContaining string = "HDMI"
