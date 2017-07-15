package timing

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFPS(t *testing.T) {
	now := time.Now()
	t2 := now.Add(1 * time.Second)
	assert.Equal(t, FPS(now, t2), 1.0)
	fps := math.Pow(10, 9)
	assert.Equal(t, FPSToNano(fps), int64(1))
	// I hate this test variety, because it just mirrors the implementation
	rate := 60
	assert.Equal(t, FPSToDuration(rate), time.Second/time.Duration(int64(rate)))
}
