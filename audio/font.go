package audio

type LoopType int

const (
	AUDIO_DEFINED LoopType = iota
	FORCE_LOOP
	FORCE_NO_LOOP
)

type Font struct {
	*Ears
	Volume    float64
	Pan       float64
	Frequency float64
	ForceLoop LoopType
}
