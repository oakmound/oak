// Package supports holds interface types for filter supports
package supports

// Data types support filters that manipulate their raw audio data
type Data interface {
	GetData() *[]byte
}

// Loop types support filters that manipulate whether they loop
type Loop interface {
	GetLoop() *bool
}

// SampleRate types support filters that manipulate their SampleRate
type SampleRate interface {
	GetSampleRate() *uint32
}

// BitDepth types support filters that manipulate bit depth. Probably
// only useful in combination as an encoding
type BitDepth interface {
	GetBitDepth() *uint16
}

// Channels types support filters that manipulate channels. Probably
// only useful in combination as an encoding
type Channels interface {
	GetChannels() *uint16
}

// Encoding types can get any variable on an audio.Encoding. They do
// not just return an audio.Encoding because that would be an import
// loop or another package to avoid said import loop.
type Encoding interface {
	SampleRate
	BitDepth
	Data
	Channels
}
