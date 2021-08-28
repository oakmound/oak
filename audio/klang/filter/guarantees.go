// Package filter provides various audio filters to be applied to audios through the
// Filter function.
package filter

import (
	"github.com/oakmound/oak/v3/audio/klang"
	"github.com/oakmound/oak/v3/audio/klang/filter/supports"
)

// These declarations guarantee that the filters in this package satisfy the filter interface
var (
	_ klang.Filter = SampleRate(func(*uint32) {})
	_ klang.Filter = Data(func(*[]byte) {})
	_ klang.Filter = Loop(func(*bool) {})
	_ klang.Filter = Encoding(func(supports.Encoding) {})
)
