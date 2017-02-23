package alg

import (
	"errors"
	"math/rand"
)

//
// This algorithm reacts very poorly to zero or negative weights.
//
func WeightedChoose(weights []float64, toChoose int) ([]int, error) {
	toChoose_f64 := float64(toChoose)
	lengthWeights := len(weights)
	out := make([]int, toChoose)
	remainingWeights := make([]float64, lengthWeights)
	remainingWeights[lengthWeights-1] = weights[lengthWeights-1]
	for i := lengthWeights - 2; i >= 0; i-- {
		remainingWeights[i] = remainingWeights[i+1] + weights[i]
	}

	for i, v := range remainingWeights {
		if toChoose == 0 {
			return out, nil
		}
		if lengthWeights-i < toChoose {
			//ERROR Out
			return []int{}, errors.New("Tried to choose too many items from a slice")
		}
		if lengthWeights-i == toChoose {
			toChoose--
			toChoose_f64--
			out[toChoose] = i
		} else if rand.Float64() < toChoose_f64/v {
			toChoose--
			toChoose_f64--
			out[toChoose] = i
		}
	}

	return out, nil

}

func UniqueChooseX(weights []float64, x int) []int {
	out := make([]int, x)
	stwh := NewSTWHeap(weights)
	for i := 0; i < x; i++ {
		out[i] = stwh.Pop()
	}
	return out
}

// AKA Roulette search
//
// This algorithm works well, the only issue relative to above
// is that it can choose the same element multiple times, which
// is not always the desired effect.
//
// A version of it could easily be made to only pick each element
// once, however. It would benefit from the linear psuedo-random
// roulette search, where forced increments would happen once an
// index was chosen.
func ChooseX(weights []float64, x int) []int {
	lengthWeights := len(weights)
	out := make([]int, x)
	remainingWeights := make([]float64, lengthWeights)
	remainingWeights[lengthWeights-1] = weights[lengthWeights-1]
	for i := lengthWeights - 2; i >= 0; i-- {
		remainingWeights[i] = remainingWeights[i+1] + weights[i]
	}
	for i := 0; i < x; i++ {
		j := CumWeightedChooseOne(remainingWeights)
		for weights[j] == 0 {
			j = CumWeightedChooseOne(remainingWeights)
		}
		out[i] = j
	}
	return out
}

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

/*
* Inefficient implementation but allows the utility to have allocations based on a map
*@returns the key chosen from the map passed in.
 */
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
