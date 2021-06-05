package filter

import (
	"github.com/oakmound/oak/v3/audio/internal/audio"
	"github.com/oakmound/oak/v3/audio/internal/audio/filter/supports"
	"github.com/oakmound/oak/v3/audio/internal/audio/manip"
)

// Encoding filters are functions on any combination of the values
// in an audio.Encoding
type Encoding func(supports.Encoding)

// Apply checks that the given audio supports Encoding, filters if it
// can, then returns
func (enc Encoding) Apply(a audio.Audio) (audio.Audio, error) {
	if senc, ok := a.(supports.Encoding); ok {
		enc(senc)
		return a, nil
	}
	return a, supports.NewUnsupported([]string{"Encoding"})
}

// AssertStereo does nothing to audio that has two channels, but will convert
// mono audio to two-channeled audio with the same data on both channels
func AssertStereo() Encoding {
	return func(enc supports.Encoding) {
		chs := enc.GetChannels()
		if *chs > 1 {
			// We can't really do this for non-mono audio
			return
		}
		*chs = 2
		data := enc.GetData()
		d := *data
		newData := make([]byte, len(d)*2)
		byteDepth := int(*enc.GetBitDepth() / 8)
		for i := 0; i < len(d); i += 2 {
			for j := 0; j < byteDepth; j++ {
				newData[i*2+j] = d[i+j]
				newData[i*2+j+byteDepth] = d[i+j]
			}
		}
		*data = newData
	}
}

func mod(init, inc int, modFn func(float64) float64) Encoding {
	return func(enc supports.Encoding) {
		data := enc.GetData()
		d := *data
		byteDepth := int(*enc.GetBitDepth() / 8)
		switch byteDepth {
		case 2:
			for i := byteDepth * init; i < len(d); i += byteDepth * inc {
				manip.SetInt16(d, i, manip.Round(modFn(float64(manip.GetInt16(d, i)))))
			}
		default:
			// log unsupported byte depth
		}
		*data = d
	}
}
