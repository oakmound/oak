package synth

// Detune detunes between -1.0 and 1.0, 1.0 representing a half step up.
// Q: What is detuning? A: It's taking the pitch of the audio and adjusting it less than
//    a single tone up or down. If you detune too far, you've just made the next pitch,
//    but if you detune a little, you get a resonant sound.
func Detune(percent float64) func(src Source) Source {
	return func(src Source) Source {
		curPitch := src.Pitch
		var nextPitch Pitch
		if percent > 0 {
			nextPitch = curPitch.Up(HalfStep)
		} else {
			nextPitch = curPitch.Down(HalfStep)
		}
		rawDelta := float64(int16(curPitch) - int16(nextPitch))
		delta := rawDelta * percent
		// TODO: does pitch need to be a float?
		src.Pitch = Pitch(float64(curPitch) + delta)
		return src
	}
}
