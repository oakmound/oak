package oak

import (
	"fmt"
	"math/rand"
	"time"
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
	fmt.Println("\n~~~~~~~~~~~~~~~")
}
