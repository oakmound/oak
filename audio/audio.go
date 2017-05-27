//+build windows

package audio

import "bitbucket.org/StephenPatrick/go-winaudio/winaudio"

// Audio is a struct of some audio data and the variables
// required to filter it through a sound font.
type Audio struct {
	Data
	F    *Font
	X, Y *float64
}

// SetVolume is a wrapper around Data.SetVolume
// to protect against invalid inputs
// this will eventually change to our own range that isn't silly,
// but right now the range is -10000 to 0 because thats what windows
// uses.
func (a *Audio) SetVolume(vol int32) error {
	if vol < winaudio.MIN_VOLUME {
		vol = winaudio.MIN_VOLUME
	} else if vol > winaudio.MAX_VOLUME {
		vol = winaudio.MAX_VOLUME
	}
	return a.Data.SetVolume(vol)
}

// SetPan is a wrapper around Data.SetPan
// to protect against invalid inputs with
// a slightly less silly range of -10000 to 10000
func (a *Audio) SetPan(pan int32) error {
	if pan < winaudio.LEFT_PAN {
		pan = winaudio.LEFT_PAN
	} else if pan > winaudio.RIGHT_PAN {
		pan = winaudio.RIGHT_PAN
	}
	return a.Data.SetPan(pan)
}

// SetFrequency is a wrapper around Data.SetFrequency
// to protect against invalid inputs.
// The range is 100 to 10000, but modifying frequency by
// anything more than a little bit will usually result in
// something inaudible.
func (a *Audio) SetFrequency(freq uint32) error {
	if freq < winaudio.MIN_FREQUENCY {
		freq = winaudio.MIN_FREQUENCY
	} else if freq > winaudio.MAX_FREQUENCY {
		freq = winaudio.MAX_FREQUENCY
	}
	return a.Data.SetFrequency(freq)
}

// ScaleVolume changes an audio's volume to be scale * the existing volume,
// ignoring the silly -10000 to 0 scale and assuming a positive scale.
// I.E. this works like you would expect: -5000 volume * scale .5 will give
// -7500, or half as loud.
func (a *Audio) ScaleVolume(scale float64) error {
	vol, err := a.GetVolume()
	if err != nil {
		return err
	}
	// Todo: magic numbers
	vol += 10000
	vol = int32(float64(vol) * scale)
	vol -= 10000
	return a.SetVolume(vol)
}

// ShiftVolume increases volume by v.
func (a *Audio) ShiftVolume(v float64) error {
	vol, err := a.GetVolume()
	if err != nil {
		return err
	}
	vol += int32(v)
	return a.SetVolume(vol)
}

// ShiftPan increases pan by p.
func (a *Audio) ShiftPan(p float64) error {
	pan, err := a.GetPan()
	if err != nil {
		return err
	}
	pan += int32(p)
	return a.SetPan(pan)
}

// ShiftFrequency increases frequency by f.
// For some reason, GetFrequency appears to be wrong sometimes,
// and this might have unexpected behavior of making the audio inaudible.
func (a *Audio) ShiftFrequency(f float64) error {
	freq, err := a.GetFrequency()
	if err != nil {
		return err
	}
	freq += uint32(f)
	return a.SetFrequency(freq)
}
