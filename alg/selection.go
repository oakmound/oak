package alg

import "math/rand"

// UniqueChooseX returns n indices from the input weights at a count
// relative to the weight of each index. This will never return duplicate indices.
// if n > len(weights), it will return -1 after depleting the n elements from
// weights.
func UniqueChooseX(weights []float64, n int) []int {
	out := make([]int, n)
	stwh := newSTWHeap(weights)
	for i := 0; i < n; i++ {
		out[i] = stwh.Pop()
	}
	return out
}

// ChooseX AKA Roulette search
//
// This returns n indices from the input weights at a count
// relative to the weight of each index. It can return the same index
// multiple times.
func ChooseX(weights []float64, n int) []int {
	lengthWeights := len(weights)
	out := make([]int, n)
	remainingWeights := make([]float64, lengthWeights)
	remainingWeights[lengthWeights-1] = weights[lengthWeights-1]
	for i := lengthWeights - 2; i >= 0; i-- {
		remainingWeights[i] = remainingWeights[i+1] + weights[i]
	}
	for i := 0; i < n; i++ {
		out[i] = CumWeightedChooseOne(remainingWeights)
	}
	return out
}

// CumWeightedChooseOne returns a single index from the weights given
// at a rate relative to the magnitude of each weight
func CumWeightedChooseOne(remainingWeights []float64) int {
	totalWeight := remainingWeights[0]
	choice := rand.Float64() * totalWeight
	i := len(remainingWeights) / 2
	start := 0
	end := len(remainingWeights) - 1
	for {
		if remainingWeights[i] < choice {
			if remainingWeights[i-1] < choice {
				end = i
				i = (start + end) / 2
			} else {
				return i - 1
			}
		} else {
			if i != len(remainingWeights)-1 && remainingWeights[i+1] > choice {
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

// CumWeightedFromMap converts the input map into a set where keys are
// indices and values are weights for CumWeightedChooseOne, then returns
// the key for CumWeightedChooseOne of the weights.
func CumWeightedFromMap(weightMap map[int]float64) int {
	keys := make([]int, len(weightMap))
	values := make([]float64, len(weightMap))
	idx := 0
	for key, value := range weightMap {
		keys[idx] = key
		values[idx] = value
		idx++
	}

	return keys[CumWeightedChooseOne(values)]
}

// CumulativeWeights converts a slice of weights into
// a slice of cumulative weights, where each index
// is the sum of all weights up until that index in
// the original slice
func CumulativeWeights(weights []float64) []float64 {
	cum := make([]float64, len(weights))
	cum[0] = weights[0]
	for i := 1; i < len(weights); i++ {
		cum[i] = cum[i-1] + weights[i]
	}
	return cum
}
