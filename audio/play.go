//+build windows

package audio

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/physics"
)

// Play plays a sound from an audio
func (a *Audio) Play() error {
	origVolume, _ := a.GetVolume()
	origFrequency, _ := a.GetFrequency()
	origPan, _ := a.GetPan()
	var err error
	ad := a.Data
	if a.F.Volume != 0 {
		err = a.ShiftVolume(a.F.Volume)
		if err != nil {
			return err
		}
	}
	if a.F.Frequency != 0 {
		err = a.ShiftFrequency(a.F.Frequency)
		if err != nil {
			return err
		}
	}
	switch a.F.ForceLoop {
	case ForceLoop:
		ad.SetLooping(true)
	case ForceNoLoop:
		ad.SetLooping(false)
	}

	if a.F.Ears != nil && a.X != nil && a.Y != nil {
		err = a.ScaleVolume(a.F.CalculateVolume(physics.NewVector(*a.X, *a.Y)))
		if err != nil {
			return err
		}
		err = a.SetPan(a.F.CalculatePan(*a.X))
		if err != nil {
			return err
		}
	} else {
		err = a.ShiftPan(a.F.Pan)
		if err != nil {
			return err
		}
	}

	err = ad.Play()
	ad.SetVolume(origVolume)
	ad.SetFrequency(origFrequency)
	ad.SetPan(origPan)
	return err
}

// PlayWav is shorthand for GetWav followed by Play.
func (f *Font) PlayWav(filename string) error {
	ad, err := GetWav(filename)
	if err == nil {
		a := Audio{ad, f, nil, nil}
		err = a.Play()
	} else {
		dlog.Error(err)
	}
	return err
}

// PlayWav with no font calls PlayWav on the default font.
func PlayWav(filename string) error {
	return defFont.PlayWav(filename)
}
