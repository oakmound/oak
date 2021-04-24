package timing

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

const randTestCt = 100

func TestFPS(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())
	t.Run("SecondDiff", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < randTestCt; i++ {
			now := time.Now()
			secs := rand.Intn(100) + 1
			t2 := now.Add(time.Duration(secs) * time.Second)
			got := FPS(now, t2)
			expected := 1 / float64(secs)
			if got != expected {
				t.Fatalf("got fps of %v, expected %v", got, expected)
			}
		}
	})
	t.Run("NoDiff", func(t *testing.T) {
		now := time.Now()
		expected := float64(maximumFPS)
		got := FPS(now, now)
		if got != expected {
			t.Fatalf("got fps of %v, expected %v", got, expected)
		}
	})
}

func TestFPSToNano(t *testing.T) {
	t.Parallel()
	t.Run("NanoPerSecond", func(t *testing.T) {
		t.Parallel()
		fps := nanoPerSecond
		got := FPSToNano(float64(fps))
		expected := int64(1)
		if got != expected {
			t.Fatalf("got nanos of %v, expected %v", got, expected)
		}
	})
	t.Run("0", func(t *testing.T) {
		t.Parallel()
		fps := 0.0
		got := FPSToNano(float64(fps))
		expected := int64(math.MaxInt64)
		if got != expected {
			t.Fatalf("got nanos of %v, expected %v", got, expected)
		}
	})
}

func TestFPSToFrameDelay(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())
	t.Run("1-201", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < randTestCt; i++ {
			rate := rand.Intn(200) + 1
			got := FPSToFrameDelay(rate)
			expected := time.Second / time.Duration(int64(rate))
			if got != expected {
				t.Fatalf("got duration of %v, expected %v", got, expected)
			}
		}
	})
	t.Run("0", func(t *testing.T) {
		t.Parallel()
		got := FPSToFrameDelay(0)
		expected := time.Duration(math.MaxInt64)
		if got != expected {
			t.Fatalf("got duration of %v, expected %v", got, expected)
		}
	})
}
