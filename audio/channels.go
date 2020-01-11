package audio

import (
	"time"

	"github.com/200sc/go-dist/intrange"
	"github.com/200sc/klangsynthese/font"
	"github.com/oakmound/oak/v2/timing"
)

// DefActiveChannel acts like GetActiveChannel when fed DefFont
func DefActiveChannel(freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return GetActiveChannel(DefFont, freq, fileNames...)
}

// GetActiveChannel returns a channel that will block until its frequency
// rotates around. This means that continually sending on ChannelSignal will
// probably cause the game to freeze or substantially slow down. For this reason
// ActiveWavChannels are meant to be used for cases where the user knows they will
// not be sending on the ActiveWavChannel more often than the frequency they send
// in.
// Audio channels serve one purpose: handling audio effects
// which come in at very high or unpredictable frequencies
// while limiting the number of concurrent ongoing audio effects
// from any one source. All channels will only play once per a given
// frequency range.
func GetActiveChannel(f *font.Font, freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getActiveChannel(f, freq, timing.ClearDelayCh, fileNames...)
}

func getActiveChannel(f *font.Font, freq intrange.Range, quitCh chan bool,
	fileNames ...string) (chan ChannelSignal, error) {

	datas, err := GetSounds(fileNames...)
	if err != nil {
		return nil, err
	}

	sounds := make([]*Audio, len(datas))
	for i, d := range datas {
		sounds[i] = New(f, d)
	}

	soundCh := make(chan ChannelSignal)
	go func() {
		// Todo: When a scene ends, we need to clear all of these goroutines out
		for {
			delay := time.Duration(freq.Poll())
			select {
			case <-quitCh:
				return
			case <-time.After(delay * time.Millisecond):
			}
			// Every once in a while, after some delay,
			// we play an audio that slipped through the
			// above routine.
			select {
			case <-quitCh:
				return
			case signal := <-soundCh:
				sound := sounds[signal.GetIndex()]
				usePos, x, y := signal.GetPos()
				if usePos {
					sound.X = &x
					sound.Y = &y
				}
				sound.Play()
			}
		}
	}()
	return soundCh, nil
}

// DefChannel acts like GetChannel when given DefFont
func DefChannel(freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getChannel(DefFont, freq, timing.ClearDelayCh, fileNames...)
}

// GetChannel channels will attempt to steal most sends sent to the output
// audio channel. This will allow a game to constantly send on a channel and
// obtain an output rate of near the sent in frequency instead of locking
// or requiring buffered channel usage.
//
// An important example case-- walking around
// When a character walks, they have some frequency step speed and some
// set of potential fileName sounds that play, and the usage of a channel
// here will let the EnterFrame code which detects the walking status to
// send on the walking audio channel constantly without worrying about
// triggering too many sounds.
func GetChannel(f *font.Font, freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getChannel(f, freq, timing.ClearDelayCh, fileNames...)
}

func getChannel(f *font.Font, freq intrange.Range, quitCh chan bool, fileNames ...string) (chan ChannelSignal, error) {
	soundCh, err := getActiveChannel(f, freq, quitCh, fileNames...)
	if err != nil {
		return nil, err
	}

	// This routine serves to steal almost every
	// attempt to play audio
	go func() {
		for {
			select {
			case <-quitCh:
				return
			case <-soundCh:
			}
		}
	}()
	return soundCh, nil
}
