package alg

import "math/rand"

// A Float64Generator must be able to generate a float64. This is generally used to implement randomness
// See rand.Float64 for an example
type Float64Generator interface {
	Float64() float64
}

// UniqueChooseX returns n indices from the input weights at a count
// relative to the weight of each index. This will never return duplicate indices.
// if n > len(weights), it will return -1 after depleting the n elements from
// weights.
func UniqueChooseX(weights []float64, n int) []int {
	return uniqueChooseX(weights, n, rand.Float64)
}

// UniqueChooseXSeeded returns n indices from the input weights at a count
// relative to the weight of each index. This will never return duplicate indices.
// if n > len(weights), it will return -1 after depleting the n elements from
// weights. If you do not want to use your own rand
// use UniqueChooseX.
func UniqueChooseXSeeded(weights []float64, n int, rng Float64Generator) []int {
	return uniqueChooseX(weights, n, rng.Float64)
}
func uniqueChooseX(weights []float64, n int, rngFxn func() float64) []int {
	out := make([]int, n)
	stwh := newSTWHeap(weights)
	for i := 0; i < n; i++ {
		out[i] = stwh.Pop(rngFxn())
	}
	return out
}

// ChooseX - also known as Roulette Search.
// This returns n indices from the input weights at a count
// relative to the weight of each index. It can return the same index
// multiple times.
func ChooseX(weights []float64, n int) []int {
	out := make([]int, n)
	cumulative := CumulativeWeights(weights)
	for i := 0; i < n; i++ {
		out[i] = WeightedChooseOne(cumulative)
	}
	return out
}

// WeightedMapChoice converts the input map into a set where keys are
// indices and values are weights for WeightedChooseOne, then returns
// the key for WeightedChooseOne of the weights.
func WeightedMapChoice(weightMap map[int]float64) int {
	return weightedMapChoice(weightMap, rand.Float64)
}

// WeightedMapChoiceSeeded converts the input map into a set where keys are
// indices and values are weights for WeightedChooseOne, then returns
// the key for WeightedChooseOne of the weight. If you do not want to use your own rand
// use WeightedMapChoice.
func WeightedMapChoiceSeeded(weightMap map[int]float64, rng Float64Generator) int {
	return weightedMapChoice(weightMap, rng.Float64)
}
func weightedMapChoice(weightMap map[int]float64, rngFxn func() float64) int {
	keys := make([]int, len(weightMap))
	values := make([]float64, len(weightMap))
	idx := 0
	for key, value := range weightMap {
		keys[idx] = key
		values[idx] = value
		idx++
	}
	values = CumulativeWeights(values)
	return keys[weightedChooseOne(values, rngFxn)]
}

// CumulativeWeights converts a slice of weights into
// a slice of cumulative weights, where each index
// is the sum of all weights following that index in
// the original slice
func CumulativeWeights(weights []float64) []float64 {
	cumulative := make([]float64, len(weights))
	cumulative[len(weights)-1] = weights[len(weights)-1]
	for i := len(weights) - 2; i >= 0; i-- {
		cumulative[i] = cumulative[i+1] + weights[i]
	}
	return cumulative
}

// WeightedChooseOne returns a single index from the weights given
// at a rate relative to the magnitude of each weight. It expects
// the input to be in the form of CumulativeWeights, cumulative with
// the total at index 0.
func WeightedChooseOne(cumulative []float64) int {
	return weightedChooseOne(cumulative, rand.Float64)
}

// WeightedChooseOneSeeded returns a single index from the weights given
// at a rate relative to the magnitude of each weight. It expects
// the input to be in the form of CumulativeWeights, cumulative with
// the total at index 0. If you do not want to use your own rand
// use WeightedChooseOne.
func WeightedChooseOneSeeded(cumulative []float64, rng Float64Generator) int {
	return weightedChooseOne(cumulative, rng.Float64)
}

// weightedChooseOne returns a single index from the weights given
// at a rate relative to the magnitude of each weight. It expects
// the input to be in the form of CumulativeWeights, cumulative with
// the total at index 0.
func weightedChooseOne(cumulative []float64, rngFxn func() float64) int {
	totalWeight := cumulative[0]
	choice := rngFxn() * totalWeight
	i := len(cumulative) / 2
	start := 0
	end := len(cumulative) - 1
	for {
		if cumulative[i] < choice {
			if cumulative[i-1] < choice {
				end = i
				i = (start + end) / 2
			} else {
				return i - 1
			}
		} else {
			if i != len(cumulative)-1 && cumulative[i+1] > choice {
				start = i

				i = (start + end) / 2
				if (start+end)%2 == 1 {
					i++
				}
			} else {
				return i
			}
		}
	}
}
