package alg

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const testCt = 1000000

func TestUniqueChooseXSuccess(t *testing.T) {
	t.Parallel()

	rand.Seed(int64(time.Now().UTC().Nanosecond()))
	// Assert that we choose everything when n = len(weights)
	weights := []float64{1.0, .9, .8, .7, .6, .5, .4, .3, .2, .1}
	chosenCts := make([]int, len(weights))
	chosen := UniqueChooseX(weights, len(weights))
	for _, c := range chosen {
		chosenCts[c]++
	}
	expected := make([]int, len(weights))
	for i := 0; i < len(expected); i++ {
		expected[i] = 1
	}
	if !reflect.DeepEqual(expected, chosenCts) {
		t.Fatalf("UniqueChooseX did not choose every index once")
	}
}
func TestUniqueChooseXFailure(t *testing.T) {
	t.Parallel()
	weights := []float64{1.0, .9, .8, .7, .6, .5, .4, .3, .2, .1}

	// -1 represents an error from this
	chosen := UniqueChooseX(weights, 20)
	for _, v := range chosen {
		if v == -1 {
			return
		}
	}
	t.Fatalf("expected -1 index in results from 20 choices on 10 inputs")
}

func TestUniqueChooseXSeededFailure(t *testing.T) {
	t.Parallel()
	weights := []float64{1.0, .9, .8, .7, .6, .5, .4, .3, .2, .1}
	chosen := UniqueChooseXSeeded(weights, 20, rand.New(rand.NewSource(0)))
	for _, v := range chosen {
		if v == -1 {
			return
		}
	}
	t.Fatalf("expected -1 index in results from 20 choices on 10 inputs")
}

func TestUniqueChooseXAverage(t *testing.T) {
	t.Parallel()

	rand.Seed(int64(time.Now().UTC().Nanosecond()))
	weights := []float64{1.0, .9, .8, .7, .6, .5, .4, .3, .2, .1}
	chosenCts := make([]int, len(weights))
	for i := 0; i < testCt; i++ {
		chosen := UniqueChooseX(weights, 1)
		for _, c := range chosen {
			chosenCts[c]++
		}
	}
	outWeights := make([]float64, len(weights))
	for i, v := range chosenCts {
		outWeights[i] = float64(v) / float64(testCt)
	}
	for i := 0; i < len(weights)-1; i++ {
		diff := math.Abs(outWeights[i] - outWeights[i+1])
		if (outWeights[i] <= outWeights[i+1]) && diff > .1 {
			t.Fatalf("chooseX chose improbably unlikely choices")
		}
	}
}

func TestChooseXAverage(t *testing.T) {
	t.Parallel()

	rand.Seed(int64(time.Now().UTC().Nanosecond()))
	weights := []float64{1.0, .9, .8, .7, .6, .5, .4, .3, .2, .1}
	chosenCts := make([]int, len(weights))
	for i := 0; i < testCt; i++ {
		chosen := ChooseX(weights, 1)
		for _, c := range chosen {
			chosenCts[c]++
		}
	}
	outWeights := make([]float64, len(weights))
	for i, v := range chosenCts {
		outWeights[i] = float64(v) / float64(testCt)
	}
	for i := 0; i < len(weights)-1; i++ {
		diff := math.Abs(outWeights[i] - outWeights[i+1])
		if (outWeights[i] <= outWeights[i+1]) && diff > .1 {
			t.Fatalf("chooseX chose improbably unlikely choices")
		}
	}
}

func TestChooseXZeroWeights(t *testing.T) {
	weights := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0}
	for i := 0; i < testCt; i++ {
		chosen := ChooseX(weights, 1)
		if len(chosen) != 1 {
			t.Fatalf("chooseX with n = 1 did not return slice of length 1")
		}
		if chosen[0] != 8 {
			t.Fatalf("slice weighted at index 8 did not return index 8")
		}
	}
}

func TestWeightedMapChoice(t *testing.T) {
	t.Parallel()

	m := map[int]float64{
		0: 1.0,
		1: .9,
		2: .8,
		3: .7,
		4: .6,
		5: .5,
		6: .4,
		7: .3,
		8: .2,
		9: .1,
	}
	// This uses the same underlying function as chooseX internally
	chosen := WeightedMapChoice(m)
	if chosen < 0 || chosen > 9 {
		t.Fatalf("WeightedMapChoice returned impossible value")
	}
	chosen = WeightedMapChoiceSeeded(m, rand.New(rand.NewSource(0)))
	if chosen < 0 || chosen > 9 {
		t.Fatalf("WeightedMapChoiceSeeded returned impossible value")
	}
}

func TestWeightedChooseOne(t *testing.T) {
	t.Parallel()

	weights := []float64{0, 0, 0, 0, 0, 0, 1}
	remaining := CumulativeWeights(weights)
	choice := WeightedChooseOne(remaining)
	if choice != 6 {
		t.Fatalf("Expected choice of 6, got %v", choice)
	}
}

type float64er struct {
	val float64
}

func (f *float64er) Float64() float64 {
	return f.val
}

func TestWeightedChooseOneSeeded(t *testing.T) {
	t.Parallel()

	weights := []float64{1, 1, 1, 1, 1, 1, 1}
	remaining := CumulativeWeights(weights)
	choice := WeightedChooseOneSeeded(remaining, &float64er{val: 0})
	if choice != 6 {
		t.Fatalf("Expected choice of 6, got %v", choice)
	}
}
