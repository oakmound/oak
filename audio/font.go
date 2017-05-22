//+build windows

package audio

// LoopType is an enum for how a font should loop
// audios using it
type LoopType int

const (
	// AudioDefined refers to a given audio's loop variable
	AudioDefined LoopType = iota
	// ForceLoop loops all audios using a font
	ForceLoop
	// ForceNoLoop will not loop any audio using a font
	ForceNoLoop
)

// A Font is scale settings and other variables that act
// as a filter for all audio effects using it. All audios
// are assigned a font, although there exists a default
// font for audio effects that want no filtering.
type Font struct {
	*Ears
	Volume    float64
	Pan       float64
	Frequency float64
	ForceLoop LoopType
}
