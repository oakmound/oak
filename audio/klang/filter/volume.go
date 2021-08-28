package filter

import (
	"github.com/oakmound/oak/v3/audio/klang/filter/supports"
	"github.com/oakmound/oak/v3/audio/klang/internal/manip"
)

// Volume will magnify the data by mult, increasing or reducing the volume
// of the output sound. For mult <= 1 this should have no unexpected behavior,
// although for mult ~= 1 it might not have any effect. More importantly for
// mult > 1, values may result in the output data clipping over integer overflows,
// which is presumably not desired behavior.
func Volume(mult float64) Encoding {
	return vol(0, 1, mult)
}

// VolumeLeft acts like volume but reduces left channel volume only
func VolumeLeft(mult float64) Encoding {
	return vol(0, 2, mult)
}

// VolumeRight acts like volume but reduces left channel volume only
func VolumeRight(mult float64) Encoding {
	return vol(1, 2, mult)
}

func vol(init, inc int, mult float64) Encoding {
	return mod(init, inc, func(f float64) float64 {
		return f * mult
	})
}

// VolumeBalance will filter audio on two channels such that the left channel
// is (l+r)/2 * lMult, and the right channel is (l+r)/2 * rMult
func VolumeBalance(lMult, rMult float64) Encoding {
	return func(enc supports.Encoding) {
		if *enc.GetChannels() != 2 {
			return
		}
		data := enc.GetData()
		d := *data
		byteDepth := int(*enc.GetBitDepth() / 8)
		switch byteDepth {
		case 2:
			for i := 0; i < len(d); i += (byteDepth * 2) {
				var v int16
				var shift uint16
				for j := 0; j < byteDepth; j++ {
					v += int16(int(d[i+j])+int(d[i+j+byteDepth])) / 2 << shift
					shift += 8
				}
				l := manip.Round(float64(v) * lMult)
				r := manip.Round(float64(v) * rMult)
				for j := 0; j < byteDepth; j++ {
					d[i+j] = byte(l & 255)
					d[i+j+byteDepth] = byte(r & 255)
					l >>= 8
					r >>= 8
				}
			}
		default:
			// log unsupported bit depth
		}
		*data = d
	}
}
