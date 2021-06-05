package filter

import (
	"fmt"

	"github.com/oakmound/oak/v3/audio/support/audio/filter/supports"
)

// Speed modifies the filtered audio by a speed ratio, changing its sample rate
// in the process while maintaining pitch.
func Speed(ratio float64, pitchShifter PitchShifter) Encoding {
	return func(senc supports.Encoding) {
		r := ratio
		fmt.Println(ratio)
		for r < .5 {
			r *= 2
			pitchShifter.PitchShift(.5)(senc)
		}
		for r > 2.0 {
			r /= 2
			pitchShifter.PitchShift(2.0)(senc)
		}
		pitchShifter.PitchShift(1 / r)(senc)
		ModSampleRate(ratio)(senc.GetSampleRate())
	}
}
