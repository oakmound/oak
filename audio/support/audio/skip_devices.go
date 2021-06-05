package audio

import (
	"os"
)

// SkipDevicesContaining is a environment variable controlled value
// which will cause audio devices containing the given string to be 
// skipped when finding an audio device to play audio through. 
// Currently only supported on linux.
// Todo: find a more elegant fix for bad audio devices being chosen
var SkipDevicesContaining = "HDMI"


func init() {
	skipDevices := os.Getenv("KGS_AUDIO_SKIP_DEVICES") 
	if skipDevices != "" {
		SkipDevicesContaining = skipDevices
	}
}