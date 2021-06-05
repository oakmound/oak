package audio

import "time"

// Encoding contains all information required to convert raw data
// (currently assumed PCM data but that may/will change) into playable Audio
type Encoding struct {
	// Consider: non []byte data?
	// Consider: should Data be a type just like Format and CanLoop?
	Data []byte
	Format
	CanLoop
}

// Copy returns an audio encoded from this encoding.
// Consider: Copy might be tied to HasEncoding
func (enc *Encoding) Copy() (Audio, error) {
	return EncodeBytes(*enc.copy())
}

// MustCopy acts like Copy, but will panic if err != nil
func (enc *Encoding) MustCopy() Audio {
	a, err := EncodeBytes(*enc.copy())
	if err != nil {
		panic(err)
	}
	return a
}

// GetData satisfies filter.SupportsData
func (enc *Encoding) GetData() *[]byte {
	return &enc.Data
}

// PlayLength returns how long this encoding will play its data for
func (enc *Encoding) PlayLength() time.Duration {
	return time.Duration(
		1000000000*float64(len(enc.Data))/
			float64(enc.SampleRate)/
			float64(enc.Channels)/
			float64(enc.Bits/8)) * time.Nanosecond
}

// copy for an encoding just copies the encoding data,
// it does not return an audio.
func (enc *Encoding) copy() *Encoding {
	newEnc := new(Encoding)
	newEnc.Format = enc.Format
	newEnc.CanLoop = enc.CanLoop
	newEnc.Data = make([]byte, len(enc.Data))
	copy(newEnc.Data, enc.Data)
	return newEnc
}
