package audio

import (
	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
	"bitbucket.org/oakmoundstudio/oak/dlog"
)

type PosAudio struct {
	Audio
	X, Y *int
}

type PosAudioF struct {
	Audio
	X, Y *float64
}

func (pa *PosAudioF) Play() (err error) {
	if usingEars {
		return PlayPositional(pa.Audio, *pa.X, *pa.Y)
	}
	return pa.Audio.Play()
}

func (pa *PosAudio) Play() error {
	if usingEars {
		return PlayPositional(pa.Audio, float64(*pa.X), float64(*pa.Y))
	}
	return pa.Audio.Play()
}

func PlayPositional(sound Audio, x, y float64) (err error) {
	volume := CalculateVolume(x, y)
	if volume > winaudio.MIN_VOLUME {
		err = sound.SetPan(CalculatePan(x))
		if err != nil {
			dlog.Error(err)
		}
		err = sound.SetVolume(volume)
		if err != nil {
			dlog.Error(err)
		}
		err = sound.Play()
	}
	return err
}
