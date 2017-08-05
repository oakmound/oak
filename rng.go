package oak

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oakmound/oak/dlog"
)

const (
	// DefaultSeed is a key int64 sent in to SeedRNG
	// used to indicate that the seed function should just
	// do the default operation for seeding, using the current
	// time.
	DefaultSeed int64 = iota
)

var (
	currentSeed int64
)

// SeedRNG seeds go's random number generator
// and logs the seed set to file.
func SeedRNG(inSeed int64) {
	if currentSeed != 0 && inSeed == DefaultSeed {
		fmt.Println("Oak Seed was already set to :", currentSeed)
		return
	}
	currentSeed = inSeed
	if inSeed == DefaultSeed {
		inSeed = time.Now().UTC().UnixNano()
	}
	rand.Seed(inSeed)

	fmt.Println("\n~~~~~~~~~~~~~~~")
	fmt.Println("Oak Seed:", inSeed)

	dlog.FileWrite("Oak seed:", inSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~")
}
