package pcm

import "io"

var _ Reader = &IOReader{}

// A Reader mimics io.Reader for pcm data.
type Reader interface {
	Formatted
	ReadPCM(b []byte) (n int, err error)
}

// An IOReader converts an io.Reader into a pcm.Reader
type IOReader struct {
	Format
	io.Reader
}

func (ior *IOReader) ReadPCM(p []byte) (n int, err error) {
	return ior.Read(p)
}

// A Writer can have PCM formatted audio data written to it. It mimics io.Writer.
type Writer interface {
	io.Closer
	Formatted
	// WritePCM expects PCM bytes matching the format this speaker was initialized with.
	// WritePCM will block until all of the bytes are consumed.
	WritePCM([]byte) (n int, err error)
}

// The Formatted interface represents types that are aware of a PCM Format they expect or provide.
type Formatted interface {
	// PCMFormat will return the Format used by an encoded audio or expected by an audio consumer.
	// Implementations can embed a Format struct to simplify this.
	PCMFormat() Format
}

// Format is a PCM format; it defines how binary audio data should be converted into real audio.
type Format struct {
	// SampleRate defines how many times per second a consumer should read a single value. An example
	// of a common value for this is 44100 or 44.1khz.
	SampleRate uint32
	// Channels defines how many concurrent audio channels are present in audio data. Common values are
	// 1 for mono and 2 for stereo.
	Channels uint16
	// Bits determines how many bits a single sample value takes up. 8, 16, and 32 are common values.
	Bits uint16
}

// PCMFormat returns this format.
func (f Format) PCMFormat() Format {
	return f
}

// BytesPerSecond returns how many bytes this format would be encoded into per second in an audio stream.
func (f Format) BytesPerSecond() uint32 {
	blockAlign := f.Channels * f.Bits / 8
	return f.SampleRate * uint32(blockAlign)
}
