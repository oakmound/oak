package audio

// Format stores the variables which are presumably
// constant for any given type of audio (wav / mp3 / flac ...)
type Format struct {
	SampleRate uint32
	Channels   uint16
	Bits       uint16
}

// GetSampleRate satisfies supports.SampleRate
func (f *Format) GetSampleRate() *uint32 {
	return &f.SampleRate
}

// GetChannels satisfies supports.Channels
func (f *Format) GetChannels() *uint16 {
	return &f.Channels
}

// GetBitDepth satisfied supports.BitDepth
func (f *Format) GetBitDepth() *uint16 {
	return &f.Bits
}

// Wave takes in raw bytes and encodes them according to this format
func (f *Format) Wave(b []byte) (Audio, error) {
	return EncodeBytes(Encoding{b, *f, CanLoop{}})
}
