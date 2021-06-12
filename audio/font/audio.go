package font

import audio "github.com/oakmound/oak/v3/audio/klang"

// Audio is an ease-of-use wrapper around an audio
// with an attached font, so that the audio can be played
// with .Play() but can take in the remotely variable
// font filter options.
//
// Note that it is a conscious choice for both Font and
// Audio to have a Filter(...Filter) function, so that when
// a FontAudio is in use the user needs to specify which
// element they want to apply a filter on. The alternative would
// be to have two similarly named functions, and its believed
// that fa.Font.Filter(...) and fa.Audio.Filter(...) is
// more or less equivalent to whatever those names would be.
type Audio struct {
	*Font
	audio.FullAudio
	toStop audio.Audio
}

// NewAudio returns a *FontAudio.
// For preparation against API changes, using NewAudio over Audio{}
// is recommended.
func NewAudio(f *Font, a audio.FullAudio) *Audio {
	return &Audio{f, a, nil}
}

// Play is equivalent to Audio.Font.Play(a.Audio)
func (ad *Audio) Play() <-chan error {
	a2, err := ad.FullAudio.Copy()
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	_, err = a2.Filter(ad.Font.Filters...)
	if err != nil {
		ch := make(chan error)
		go func() {
			ch <- err
		}()
		return ch
	}
	ad.toStop = a2
	return a2.Play()
}

// Stop stops a font.Audio's playback
func (ad *Audio) Stop() error {
	if ad.toStop != nil {
		return ad.toStop.Stop()
	}
	return nil
}
