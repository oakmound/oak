//+build linux

package klang

import (
	"errors"
	"strings"
	"sync"

	"github.com/yobert/alsa"
)

type alsaAudio struct {
	*Encoding
	*alsa.Device
	playAmount   int
	playProgress int
	stopCh       chan struct{}
	playing      bool
	playCh       chan error
	period       int
}

func (aa *alsaAudio) Play() <-chan error {
	// If currently playing, restart
	if aa.playing {
		aa.playProgress = 0
		return aa.playCh
	}
	aa.playing = true
	aa.playCh = make(chan error)
	go func() {
		for {
			var data []byte
			if len(aa.Encoding.Data)-aa.playProgress <= aa.playAmount {
				data = aa.Encoding.Data[aa.playProgress:]
				if aa.Loop {
					delta := aa.playAmount - (len(aa.Encoding.Data) - aa.playProgress)
					data = append(data, aa.Encoding.Data[:delta]...)
				}
			} else {
				data = aa.Encoding.Data[aa.playProgress : aa.playProgress+aa.playAmount]
			}
			if len(data) != 0 {
				err := aa.Device.Write(data, aa.period)
				if err != nil {
					select {
					case aa.playCh <- err:
					default:
					}
					break
				}
			}
			aa.playProgress += aa.playAmount
			if aa.playProgress > len(aa.Encoding.Data) {
				if aa.Loop {
					aa.playProgress %= len(aa.Encoding.Data)
				} else {
					select {
					case aa.playCh <- nil:
					default:
					}
					break
				}
			}
			select {
			case <-aa.stopCh:
				select {
				case aa.playCh <- nil:
				default:
				}
				break
			default:
			}
		}
		aa.playing = false
		aa.playProgress = 0
	}()
	return aa.playCh
}

func (aa *alsaAudio) Stop() error {
	if aa.playing {
		go func() {
			aa.stopCh <- struct{}{}
		}()
	} else {
		return errors.New("Audio not playing, cannot stop")
	}
	return nil
}

func (aa *alsaAudio) SetVolume(int32) error {
	return errors.New("SetVolume on Linux is not supported")
}

func (aa *alsaAudio) Filter(fs ...Filter) (Audio, error) {
	var a Audio = aa
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
	return aa, consErr
}

// MustFilter acts like Filter, but ignores errors (it does not panic,
// as filter errors are expected to be non-fatal)
func (aa *alsaAudio) MustFilter(fs ...Filter) Audio {
	a, _ := aa.Filter(fs...)
	return a
}

func EncodeBytes(enc Encoding) (Audio, error) {
	handle, err := openDevice()
	if err != nil {
		return nil, err
	}
	// Todo: annotate these errors with more info
	format, err := alsaFormat(enc.Bits)
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateFormat(format)
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateRate(int(enc.SampleRate))
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateChannels(int(enc.Channels))
	if err != nil {
		return nil, err
	}
	// Default value at recommendation of library
	period, err := handle.NegotiatePeriodSize(2048)
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateBufferSize(4096)
	if err != nil {
		return nil, err
	}
	err = handle.Prepare()
	if err != nil {
		return nil, err
	}
	return &alsaAudio{
		playAmount: period * int(enc.Bits) / 4,
		period:     period,
		Encoding:   &enc,
		Device:     handle,
		stopCh:     make(chan struct{}),
	}, nil
}

var (
	// Todo: support more customized audio device usage
	openDeviceLock sync.Mutex
	openedDevice   *alsa.Device
)

func openDevice() (*alsa.Device, error) {
	openDeviceLock.Lock()
	defer openDeviceLock.Unlock()

	if openedDevice != nil {
		return openedDevice, nil
	}
	cards, err := alsa.OpenCards()
	if err != nil {
		return nil, err
	}
	defer alsa.CloseCards(cards)
	for i, c := range cards {
		dvcs, err := c.Devices()
		if err != nil {
			continue
		}
		for _, d := range dvcs {
			if d.Type != alsa.PCM || !d.Play {
				continue
			}
			if strings.Contains(d.Title, SkipDevicesContaining) {
				continue
			}
			d.Close()
			err := d.Open()
			if err != nil {
				continue
			}
			// We've a found a device we can hypothetically use
			cards = append(cards[:i], cards[i+1:]...)
			openedDevice = d
			return d, nil
		}
	}
	return nil, errors.New("No valid device found")
}

func alsaFormat(bits uint16) (alsa.FormatType, error) {
	switch bits {
	case 8:
		return alsa.S8, nil
	case 16:
		return alsa.S16_LE, nil
	}
	return 0, errors.New("Undefined alsa format for encoding bits")
}
