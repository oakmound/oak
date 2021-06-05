package filter

import (
	"github.com/oakmound/oak/v3/audio/support/audio"
	"github.com/oakmound/oak/v3/audio/support/audio/filter/supports"
)

// Data filters are functions on []byte types
type Data func(*[]byte)

// Apply checks that the given audio supports Data, filters if it
// can, then returns
func (df Data) Apply(a audio.Audio) (audio.Audio, error) {
	if sd, ok := a.(supports.Data); ok {
		df(sd.GetData())
		return a, nil
	}
	return a, supports.NewUnsupported([]string{"Data"})
}
