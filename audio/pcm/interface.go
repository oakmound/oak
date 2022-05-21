// Package pcm provides a interface for interacting with PCM audio streams
package pcm

import (
	"fmt"
	"io"
)

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
	// WritePCM expects PCM bytes matching this Writer's format.
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
	// TODO: Do we need LE vs BE, float vs int representation?
	Bits uint16
}

// PCMFormat returns this format.
func (f Format) PCMFormat() Format {
	return f
}

// BytesPerSecond returns how many bytes this format would be encoded into per second in an audio stream.
func (f Format) BytesPerSecond() uint32 {
	return f.SampleRate * uint32(f.SampleSize())
}

func (f Format) SampleSize() int {
	return int(f.Channels) * int(f.Bits/8)
}

// ReadFloat reads a single sample from an audio stream, respecting bits and channels:
// f.Bits / 8 bytes * f.Channels bytes will be read from b, and this count will be returned as 'read'.
// the length of values will be equal to f.Channels, if no error is returned. If an error is returned,
// it will be io.ErrUnexpectedEOF or ErrUnsupportedBits
func (f Format) SampleFloat(b []byte) (values []float64, read int, err error) {
	values = make([]float64, 0, f.Channels)
	read = f.SampleSize()
	if len(b) < read {
		return nil, 0, io.ErrUnexpectedEOF
	}
	_ = b[read-1]
	switch f.Bits {
	case 8:
		for i := 0; i < int(f.Channels); i++ {
			v := int8(b[i])
			values = append(values, float64(v))
		}
	case 16:
		for i := 0; i < int(f.Channels)*2; i += 2 {
			v := int16(b[i]) +
				int16(b[i+1])<<8
			values = append(values, float64(v))
		}
	case 32:
		for i := 0; i < int(f.Channels)*4; i += 4 {
			v := int32(b[i]) +
				int32(b[i+1])<<8 +
				int32(b[i+2])<<16 +
				int32(b[i+3])<<24
			values = append(values, float64(v))
		}
	default:
		return nil, read, ErrUnsupportedBits
	}
	return
}

// ErrUnsupportedBits represents that the Bits value for a Format was not supported for some operation.
var ErrUnsupportedBits = fmt.Errorf("unsupported bits in pcm format")
