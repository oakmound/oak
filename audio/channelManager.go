package audio

import (
	"github.com/200sc/go-dist/intrange"
	"github.com/200sc/klangsynthese/font"
)

// A ChannelManager can create audio channels that won't be stopped at scene end,
// but can be stopped at any time by calling Close on the manager.
type ChannelManager struct {
	quitCh chan bool
	Font   *font.Font
}

// NewChannelManager creates a channel manager whose Def functions will use the
// given font.
func NewChannelManager(f *font.Font) *ChannelManager {
	return &ChannelManager{
		quitCh: make(chan bool),
		Font:   f,
	}
}

// DefChannel creates an audio channel using the manager's Font.
func (cm *ChannelManager) DefChannel(freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getChannel(cm.Font, freq, cm.quitCh, fileNames...)
}

// DefActiveChannel creates an active channel using the manager's font.
func (cm *ChannelManager) DefActiveChannel(freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getActiveChannel(cm.Font, freq, cm.quitCh, fileNames...)
}

// GetChannel creates a channel using the given font.
func (cm *ChannelManager) GetChannel(f *font.Font, freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getChannel(f, freq, cm.quitCh, fileNames...)
}

// GetActiveChannel creates an active channel using the given font.
func (cm *ChannelManager) GetActiveChannel(f *font.Font, freq intrange.Range, fileNames ...string) (chan ChannelSignal, error) {
	return getActiveChannel(f, freq, cm.quitCh, fileNames...)
}

// Close closes the manager's internal channel handling audio channels. This will
// prevent further audio from being played via any of those channels, and their spawned
// routines will return. As Close does close a channel, it should not be called multiple
// times.
func (cm *ChannelManager) Close() {
	close(cm.quitCh)
}
