package audio

import "bitbucket.org/StephenPatrick/go-winaudio/winaudio"

type Audio struct {
	AudioData
	F    *Font
	X, Y *float64
}

func (a *Audio) SetVolume(vol int32) error {
	if vol < winaudio.MIN_VOLUME {
		vol = winaudio.MIN_VOLUME
	} else if vol > winaudio.MAX_VOLUME {
		vol = winaudio.MAX_VOLUME
	}
	return a.AudioData.SetVolume(vol)
}

func (a *Audio) SetPan(pan int32) error {
	if pan < winaudio.LEFT_PAN {
		pan = winaudio.LEFT_PAN
	} else if pan > winaudio.RIGHT_PAN {
		pan = winaudio.RIGHT_PAN
	}
	return a.AudioData.SetPan(pan)
}

func (a *Audio) SetFrequency(freq uint32) error {
	if freq < winaudio.MIN_FREQUENCY {
		freq = winaudio.MIN_FREQUENCY
	} else if freq > winaudio.MAX_FREQUENCY {
		freq = winaudio.MAX_FREQUENCY
	}
	return a.AudioData.SetFrequency(freq)
}

func (a *Audio) ScaleVolume(scale float64) error {
	vol, err := a.GetVolume()
	if err != nil {
		return err
	}
	// Todo: magic numbers
	//fmt.Println("Scaling", vol, scale, int32(scale))
	vol += 10000
	vol = int32(float64(vol) * scale)
	vol -= 10000
	return a.SetVolume(vol)
}

func (a *Audio) ShiftVolume(v float64) error {
	vol, err := a.GetVolume()
	if err != nil {
		return err
	}
	vol += int32(v)
	return a.SetVolume(vol)
}

func (a *Audio) ShiftPan(p float64) error {
	pan, err := a.GetPan()
	if err != nil {
		return err
	}
	pan += int32(p)
	return a.SetPan(pan)
}

func (a *Audio) ShiftFrequency(f float64) error {
	freq, err := a.GetFrequency()
	if err != nil {
		return err
	}
	freq += uint32(f)
	return a.SetFrequency(freq)
}
