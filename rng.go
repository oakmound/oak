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
	dlog.Info("Oak seed:", curSeed)
	fmt.Println("\n~~~~~~~~~~~~~~~")
}
