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

// ChooseX - also known as Roulette Search.
// This returns n indices from the input weights at a count
// relative to the weight of each index. It can return the same index
// multiple times.
func ChooseX(weights []float64, n int) []int {
	out := make([]int, n)
	remainingWeights := RemainingWeights(weights)
	for i := 0; i < n; i++ {
		out[i] = WeightedChooseOne(remainingWeights)
	}
	return out
}

// WeightedChooseOne returns a single index from the weights given
// at a rate relative to the magnitude of each weight. It expects
// the input to be in the form of RemainingWeights, cumulative with
// the total at index 0.
func WeightedChooseOne(remainingWeights []float64) int {
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

// WeightedMapChoice converts the input map into a set where keys are
// indices and values are weights for WeightedChooseOne, then returns
// the key for WeightedChooseOne of the weights.
func WeightedMapChoice(weightMap map[int]float64) int {
	keys := make([]int, len(weightMap))
	values := make([]float64, len(weightMap))
	idx := 0
	for key, value := range weightMap {
		keys[idx] = key
		values[idx] = value
		idx++
	}

	return keys[WeightedChooseOne(values)]
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

// RemainingWeights is equivalent to CumulativeWeights
// with the slice reversed, where the zeroth element
// will contain the total weight.
func RemainingWeights(weights []float64) []float64 {
	remainingWeights := make([]float64, len(weights))
	remainingWeights[len(weights)-1] = weights[len(weights)-1]
	for i := len(weights) - 2; i >= 0; i-- {
		remainingWeights[i] = remainingWeights[i+1] + weights[i]
	}
	return remainingWeights
}
