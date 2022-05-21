//go:build linux
// +build linux

package audio

import (
	"errors"
	"strings"
	"sync"

	"github.com/oakmound/alsa"
	"github.com/oakmound/oak/v4/audio/pcm"
)

func newALSAWriter(f pcm.Format) (pcm.Writer, error) {
	handle, err := openDevice()
	if err != nil {
		return nil, err
	}
	// Todo: annotate these errors with more info
	format, err := alsaFormat(f.Bits)
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateFormat(format)
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateRate(int(f.SampleRate))
	if err != nil {
		return nil, err
	}
	_, err = handle.NegotiateChannels(int(f.Channels))
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
	return &alsaWriter{
		Format: f,
		period: period,
		Device: handle,
	}, nil
}

type alsaWriter struct {
	sync.Mutex
	pcm.Format
	*alsa.Device
	playing bool
	period  int
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
		devices, err := c.Devices()
		if err != nil {
			continue
		}
		for _, d := range devices {
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
			// don't close this card
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
	case 32:
		return alsa.S32_LE, nil
	}
	return 0, errors.New("Undefined alsa format for encoding bits")
}

func (aw *alsaWriter) Close() error {
	aw.Lock()
	defer aw.Unlock()
	var err error
	if aw.playing {
		aw.playing = false
	}
	return err
}

func (aw *alsaWriter) WritePCM(data []byte) (n int, err error) {
	aw.Lock()
	defer aw.Unlock()
	err = aw.Device.Write(data, aw.period)
	if err != nil {
		return 0, err
	}
	if !aw.playing {
		aw.playing = true
	}
	return len(data), err
}
