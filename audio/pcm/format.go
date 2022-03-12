package pcm

// Format is a PCM format. Equivalent to klang.Format.
type Format struct {
	SampleRate uint32
	Channels   uint16
	Bits       uint16
}

// PCMFormat returns this format.
func (f Format) PCMFormat() Format {
	return f
}

// The Formatted interface represents types that are aware of a PCM Format they expect or provide.
type Formatted interface {
	// PCMFormat will return the Format used by an encoded audio or expected by an audio consumer.
	// Implementations can embed a Format struct to simplify this.
	PCMFormat() Format
}

// BytesPerSecond returns how many bytes this format would be encoded into per second in an audio stream.
func (f Format) BytesPerSecond() uint32 {
	blockAlign := f.Channels * f.Bits / 8
	return f.SampleRate * uint32(blockAlign)
}
