package oak

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRNG(t *testing.T) {
	SeedRNG(DefaultSeed)
	assert.Equal(t, currentSeed, DefaultSeed)
	tme := time.Now().UTC().UnixNano()
	SeedRNG(tme)
	assert.Equal(t, tme, currentSeed)
	SeedRNG(DefaultSeed)
	assert.Equal(t, tme, currentSeed)
}
