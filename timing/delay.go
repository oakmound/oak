package timing

import "time"

var (
	ClearDelayCh = make(chan bool)
)

func DoAfter(d time.Duration, f func()) {
	select {
	case <-time.After(d):
		f()
	case <-ClearDelayCh:
	}
}
