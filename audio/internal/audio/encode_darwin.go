//+build darwin

package audio

import "errors"

type darwinNopAudio struct {
	Encoding
}

func (dna *darwinNopAudio) Play() <-chan error {
	ch := make(chan error)
	go func() {
		ch <- errors.New("Playback on Darwin is not supported")
	}()
	return ch
}

func (dna *darwinNopAudio) Stop() error {
	return errors.New("Playback on Darwin is not supported")
}

func (dna *darwinNopAudio) SetVolume(int32) error {
	return errors.New("SetVolume on Darwin is not supported")
}

func (dna *darwinNopAudio) Filter(fs ...Filter) (Audio, error) {
	var a Audio = dna
	var err, consErr error
	for _, f := range fs {
		a, err = f.Apply(a)
		if err != nil {
			if consErr == nil {
				consErr = err
			} else {
				consErr = errors.New(err.Error() + ":" + consErr.Error())
			}
		}
	}
	return dna, consErr
}

func (dna *darwinNopAudio) MustFilter(fs ...Filter) Audio {
	a, _ := dna.Filter(fs...)
	return a
}

func EncodeBytes(enc Encoding) (Audio, error) {
	return &darwinNopAudio{enc}, nil
}
