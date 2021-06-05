// Package filter provides various audio filters to be applied to audios through the
// Filter() function
package filter

import (
	"github.com/oakmound/oak/v3/audio/support/audio"
	"github.com/oakmound/oak/v3/audio/support/audio/filter/supports"
)

// These declarations guarantee that the filters in this package satisfy the filter interface
var (
	_ audio.Filter = SampleRate(func(*uint32) {})
	_ audio.Filter = Data(func(*[]byte) {})
	_ audio.Filter = Loop(func(*bool) {})
	_ audio.Filter = Encoding(func(supports.Encoding) {})
)
