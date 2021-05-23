package intrange

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewLinear_Constant(t *testing.T) {
	linear := NewLinear(1, 1)
	if _, ok := linear.(constant); !ok {
		t.Fatalf("NewLinear with no variance did not create constant")
	}
}

func TestNewSpread_Constant(t *testing.T) {
	linear := NewSpread(1, 0)
	if _, ok := linear.(constant); !ok {
		t.Fatalf("NewSpread with no spread did not create constant")
	}
}

func TestNewSpread(t *testing.T) {
	linear := NewSpread(10, -10).(linear)
	if linear.flipped {
		t.Fatalf("new spread should not produce flipped linear range")
	}
}

func TestLinear(t *testing.T) {
	rand.Seed(time.Now().Unix())
	const testCount = 100
	const maxInt = 100000
	const minInt = -100000
	for i := 0; i < testCount; i++ {
		min := rand.Intn(maxInt-minInt) + minInt
		max := rand.Intn(maxInt-minInt) + minInt
		linear := NewLinear(min, max)
		flipped := false
		if max < min {
			min, max = max, min
			flipped = true
		}
		poll := linear.Poll()
		if poll < min || poll > max {
			t.Fatal("Linear.Poll did not return a value in its range")
		}
		magnitude := rand.Float64()
		linear2 := linear.Mult(magnitude)
		poll2 := linear2.Poll()
		if poll2 < int(float64(min)*magnitude) || poll2 > int(float64(max)*magnitude) {
			t.Fatal("Linear.Mult result did not match expected Poll")
		}
		underMin := (rand.Intn(maxInt-minInt) + minInt) - (maxInt - minInt)
		if linear.EnforceRange(underMin) != min {
			t.Fatal("Linear.EnforceRange under min did not return min")
		}
		overMax := (rand.Intn(maxInt-minInt) + minInt) + (maxInt - minInt)
		if linear.EnforceRange(overMax) != max {
			t.Fatal("Linear.EnforceRange over max did not return max")
		}
		within := rand.Intn(max-min) + min
		if linear.EnforceRange(within) != within {
			t.Fatal("Linear.EnforceRange within range did not return input")
		}
		percent := rand.Float64()
		if !flipped {
			if linear.Percentile(percent) != min+int(float64((max-min))*percent) {
				t.Fatal("Linear.Percentile did not return percentile value")
			}
		} else {
			if linear.Percentile(percent) != max+int(float64((min-max))*percent) {
				t.Fatal("flipped Linear.Percentile did not return percentile value")
			}
		}
	}
}
