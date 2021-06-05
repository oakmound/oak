package filter

import "github.com/oakmound/oak/v3/audio/support/audio/filter/supports"

// LeftPan filters audio to only play on the left speaker
func LeftPan() Encoding {
	return func(enc supports.Encoding) {
		data := enc.GetData()
		// Right/Left only makes sense for 2 channel
		if *enc.GetChannels() != 2 {
			return
		}
		// Zero out one channel
		swtch := int((*enc.GetBitDepth()) / 8)
		d := *data
		for i := 0; i < len(d); i += (2 * swtch) {
			for j := 0; j < swtch; j++ {
				d[i+j] = byte((int(d[i+j]) + int(d[i+j+swtch])) / 2)
				d[i+j+swtch] = 0
			}
		}
		*data = d
	}
}

// RightPan filters audio to only play on the right speaker
func RightPan() Encoding {
	return func(enc supports.Encoding) {
		data := enc.GetData()
		// Right/Left only makes sense for 2 channel
		if *enc.GetChannels() != 2 {
			return
		}
		// Zero out one channel
		swtch := int((*enc.GetBitDepth()) / 8)
		d := *data
		for i := 0; i < len(d); i += (2 * swtch) {
			for j := 0; j < swtch; j++ {
				d[i+j+swtch] = byte((int(d[i+j]) + int(d[i+j+swtch])) / 2)
				d[i+j] = 0
			}
		}
		*data = d
	}
}

// Pan takes -1 <= f <= 1.
// An f of -1 represents a full pan to the left, a pan of 1 represents
// a full pan to the right.
func Pan(f float64) Encoding {
	// Todo: test this is accurate
	if f > 0 {
		return VolumeBalance(1-f, 1)
	} else if f < 0 {
		return VolumeBalance(1, 1-(-1*f))
	} else {
		return func(enc supports.Encoding) {
			data := enc.GetData()
			// Right/Left only makes sense for 2 channel
			if *enc.GetChannels() != 2 {
				return
			}
			// Zero out one channel
			swtch := int((*enc.GetBitDepth()) / 8)
			d := *data
			for i := 0; i < len(d); i += (2 * swtch) {
				for j := 0; j < swtch; j++ {
					v := byte((int(d[i+j]) + int(d[i+j+swtch])) / 2)
					d[i+j+swtch] = v
					d[i+j] = v
				}
			}
			*data = d
		}
	}
}
