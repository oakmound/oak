package oak

import (
	"fmt"
	"math/rand"
	"time"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

const (
	DEFAULT_SEED = iota
)

func SeedRNG(curSeed int64) {

	if curSeed == DEFAULT_SEED {
		curSeed = time.Now().UTC().UnixNano()
	}
	rand.Seed(curSeed)

	fmt.Println("\n~~~~~~~~~~~~~~~")
	fmt.Println("Oak Seed:", curSeed)

	// We log here because we want the seed recorded in the
	// logfile for debugging purposes. Maybe a logWrite function
	// would be better.
	dlog.Info("Oak seed:", curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~")
}
