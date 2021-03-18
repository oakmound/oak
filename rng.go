package oak

import (
	"math/rand"
	"time"
)

// seedRNG seeds math/rand with time.Now, useful for minimal examples
// that would tend to forget to do this. TODO v3: add a way to disable this being called
func seedRNG() {
	rand.Seed(time.Now().UTC().UnixNano())
}
