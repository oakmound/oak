//+build windows

package audio

import "bitbucket.org/StephenPatrick/go-winaudio/winaudio"

var (
	loadedWavs = make(map[string]AudioData, 0)
	defFont    = &Font{}
)

// We alias the winaudio package's interface here
// so game files don't need to import winaudio
type AudioData winaudio.Audio

func InitWinAudio() {
	err := winaudio.InitWinAudio()
	if err != nil {
		panic(err)
	}
}
