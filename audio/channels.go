//+build windows

package audio

import (
	"math/rand"
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

// Audio channels serve one purpose: handling audio effects
// which come in at very high or unpredictable frequencies
// while limiting the number of concurrent ongoing audio effects
// from any one source. All channels will only play once per a given
// frequency range, calculated on each iteration as a range on frequency
// and a random value in addition to frequency.

func GetActiveWavChannel(frequency, freqRand int, fileNames ...string) (chan ChannelSignal, error) {
	return defFont.GetActiveWavChannel(frequency, freqRand, fileNames...)
}

func (f *Font) GetActiveWavChannel(frequency, freqRand int, fileNames ...string) (chan ChannelSignal, error) {

	datas, err := GetSounds(fileNames...)
	if err != nil {
		return nil, err
	}

	sounds := make([]*Audio, len(datas))
	for i, d := range datas {
		sounds[i] = &Audio{d, f, nil, nil}
	}

	soundCh := make(chan ChannelSignal)
	go func() {
		for {
			delay := time.Duration(rand.Intn(freqRand) + frequency)
			<-time.After(delay * time.Millisecond)
			// Every once in a while, after some delay,
			// we play an audio that slipped through the
			// above routine.
			signal := <-soundCh
			sound := sounds[signal.GetIndex()]
			usePos, x, y := signal.GetPos()
			if usePos {
				sound.X = &x
				sound.Y = &y
			}
			err := sound.Play()
			if err != nil {
				dlog.Error(err)
			}
		}
	}()
	return soundCh, nil
}

// Non-Active channels will attempt to steal most sends sent to the output
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

func GetWavChannel(frequency, freqRand int, fileNames ...string) (chan ChannelSignal, error) {
	return defFont.GetWavChannel(frequency, freqRand, fileNames...)
}

func (f *Font) GetWavChannel(frequency, freqRand int, fileNames ...string) (chan ChannelSignal, error) {

	soundCh, err := f.GetActiveWavChannel(frequency, freqRand, fileNames...)
	if err != nil {
		return nil, err
	}
	// This routine serves to steal almost every
	// attempt to play audio
	go func() {
		for {
			<-soundCh
		}
	}()
	return soundCh, nil
}
