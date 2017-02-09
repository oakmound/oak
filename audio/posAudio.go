package audio

import (
	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
)

type PosAudio struct {
	Audio
	X, Y *int
}

type PosAudioF struct {
	Audio
	X, Y *float64
}

func (pa *PosAudioF) Play() error {
	if usingEars {
		PlayPositional(pa.Audio, *pa.X, *pa.Y)
	} else {
		return pa.Audio.Play()
	}
	return nil
}

func (pa *PosAudio) Play() error {
	if usingEars {
		return PlayPositional(pa.Audio, float64(*pa.X), float64(*pa.Y))
	} else {
		return pa.Audio.Play()
	}
	return nil
}

func PlayPositional(sound Audio, x, y float64) error {
	volume := CalculateVolume(x, y)
	if volume > winaudio.MIN_VOLUME {
		sound.SetPan(CalculatePan(x))
		sound.SetVolume(volume)
		return sound.Play()
	}
	return nil
}
