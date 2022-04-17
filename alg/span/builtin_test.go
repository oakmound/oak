package span

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewLinear_Constant(t *testing.T) {
	linear := NewLinear(1, 1)
	if _, ok := linear.(constant[int]); !ok {
		t.Fatalf("NewLinear with no variance did not create constant")
	}
}

func TestNewSpread_Constant(t *testing.T) {
	linear := NewSpread(1, 0)
	if _, ok := linear.(constant[int]); !ok {
		t.Fatalf("NewSpread with no spread did not create constant")
	}
}

func TestNewSpread(t *testing.T) {
	linear := NewSpread[float32](10, -10).(linear[float32])
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
		linear2 := linear.MulSpan(magnitude)
		poll2 := linear2.Poll()
		if poll2 < int(float64(min)*magnitude) || poll2 > int(float64(max)*magnitude) {
			t.Fatal("Linear.Mult result did not match expected Poll")
		}
		underMin := (rand.Intn(maxInt-minInt) + minInt) - (maxInt - minInt)
		if linear.Clamp(underMin) != min {
			t.Fatal("Linear.EnforceRange under min did not return min")
		}
		overMax := (rand.Intn(maxInt-minInt) + minInt) + (maxInt - minInt)
		if linear.Clamp(overMax) != max {
			t.Fatal("Linear.EnforceRange over max did not return max")
		}
		within := rand.Intn(max-min) + min
		if linear.Clamp(within) != within {
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

func TestConstant(t *testing.T) {
	rand.Seed(time.Now().Unix())
	const testCount = 100
	const maxInt = 100000
	const minInt = -100000
	for i := 0; i < testCount; i++ {
		val := rand.Intn(maxInt-minInt) + minInt
		cons := NewConstant(val)
		if cons.Poll() != val {
			t.Fatal("Constant.Poll did not return initialized value")
		}
		magnitude := rand.Float64()
		cons2 := cons.MulSpan(magnitude)
		if cons2.Poll() != int(float64(val)*magnitude) {
			t.Fatal("Constant.Mult result did not match expected Poll")
		}
		if cons.Clamp(rand.Intn(maxInt)) != val {
			t.Fatal("Constant.EnforceRange did not return initialized value")
		}
		if cons.Percentile(rand.Float64()) != val {
			t.Fatal("Constant.Percentile did not return initialized value")
		}
	}
}
