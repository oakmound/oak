package random

import (
	"math/rand"
	"sync/atomic"
	"time"
)

var seed int64

func init() {
	seed = time.Now().UTC().UnixNano()
}

func Rand() *rand.Rand {
	return rand.New(rand.NewSource(atomic.AddInt64(&seed, 1)))
}
