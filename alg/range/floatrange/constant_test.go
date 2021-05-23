package floatrange

import (
	"math/rand"
	"testing"
	"time"
)

func TestConstant(t *testing.T) {
	rand.Seed(time.Now().Unix())
	const testCount = 100
	const maxInt = 100000
	const minInt = -100000
	for i := 0; i < testCount; i++ {
		val := rand.Float64()*(maxInt-minInt) + minInt
		cons := NewConstant(val)
		if cons.Poll() != val {
			t.Fatal("Constant.Poll did not return initialized value")
		}
		magnitude := rand.Float64()
		cons2 := cons.Mult(magnitude)
		if cons2.Poll() != float64(val)*magnitude {
			t.Fatal("Constant.Mult result did not match expected Poll")
		}
		if cons.EnforceRange(rand.Float64()*(maxInt-minInt)+minInt) != val {
			t.Fatal("Constant.EnforceRange did not return initialized value")
		}
		if cons.Percentile(rand.Float64()) != val {
			t.Fatal("Constant.Percentile did not return initialized value")
		}
	}
}
