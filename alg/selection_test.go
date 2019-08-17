package alg

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testCt = 1000000

func TestUniqueChooseX(t *testing.T) {
	rand.Seed(int64(time.Now().UTC().Nanosecond()))
	// Assert that we choose everything when n = len(weights)
	weights := []float64{1.0, .9, .8, .7, .6, .5, .4, .3, .2, .1}
	chosenCts := make([]int, len(weights))
	chosen := UniqueChooseX(weights, len(weights))
	for _, c := range chosen {
		chosenCts[c]++
	}
	for _, v := range chosenCts {
		assert.Equal(t, 1, v)
	}
	// That about does it for uniqueness testing
	// Failure testing
	// -1 represents an error from this
	chosen = UniqueChooseX(weights, 20)
	assert.Contains(t, chosen, -1)
	chosen = UniqueChooseXSeeded(weights, 20, rand.New(rand.NewSource(0)))
	assert.Contains(t, chosen, -1)
	//
	chosenCts = make([]int, len(weights))
	for i := 0; i < testCt; i++ {
		chosen = UniqueChooseX(weights, 1)
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
		assert.True(t, (outWeights[i] > outWeights[i+1]) || diff < .1)
	}
}

func TestChooseX(t *testing.T) {
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
		assert.True(t, (outWeights[i] > outWeights[i+1]) || diff < .1)
	}
	// Zero weight testing
	weights = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0}
	for i := 0; i < testCt; i++ {
		chosen := ChooseX(weights, 1)
		assert.Equal(t, 8, chosen[0])
	}
}

func TestWeightedMapChoice(t *testing.T) {
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
	assert.True(t, chosen < 10)
	chosen = WeightedMapChoiceSeeded(m, rand.New(rand.NewSource(0)))
	assert.True(t, chosen < 10)
}

func TestCumulativeWeights(t *testing.T) {
	weights := []float64{1, 2, 3, 4, 5, 6, 7}
	cum := CumulativeWeights(weights)
	assert.Equal(t, []float64{1, 3, 6, 10, 15, 21, 28}, cum)
}
