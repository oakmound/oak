package audio

import (
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/physics"
)

func (a *Audio) Play() error {
	var err error
	ad := a.AudioData
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
	case FORCE_LOOP:
		ad.SetLooping(true)
	case FORCE_NO_LOOP:
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

	return ad.Play()
}

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

func PlayWav(filename string) error {
	return defFont.PlayWav(filename)
}
