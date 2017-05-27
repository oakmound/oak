package oak

import (
	"fmt"
	"math/rand"
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

const (
	// DefaultSeed is a key int64 sent in to SeedRNG
	// used to indicate that the seed function should just
	// do the default operation for seeding, using the current
	// time.
	DefaultSeed = iota
)

var (
	currentSeed int64
)

// SeedRNG seeds go's random number generator
// and logs the seed set to file.
func SeedRNG(curSeed int64) {
	if currentSeed != 0 && curSeed == DefaultSeed {
		fmt.Println("Oak Seed was already set to :", currentSeed)
		return
	}
	currentSeed = curSeed
	if curSeed == DefaultSeed {
		curSeed = time.Now().UTC().UnixNano()
	}
	rand.Seed(curSeed)

	fmt.Println("\n~~~~~~~~~~~~~~~")
	fmt.Println("Oak Seed:", curSeed)

	// We log here because we want the seed recorded in the
	// logfile for debugging purposes. Maybe a logWrite function
	// would be better.
	dlog.FileWrite("Oak seed:", curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~")
}
