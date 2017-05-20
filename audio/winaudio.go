//+build windows

package audio

import "bitbucket.org/StephenPatrick/go-winaudio/winaudio"

var (
	loadedWavs = make(map[string]Data, 0)
	defFont    = &Font{}
)

// Data represents the underlying struct we
// wrap around to control audio signals.
// We alias the winaudio package's interface here
// so game files don't need to import winaudio
type Data winaudio.Audio

// InitWinAudio wraps around winaudio.InitWinAudio
func InitWinAudio() {
	err := winaudio.InitWinAudio()
	if err != nil {
		panic(err)
	}
}
