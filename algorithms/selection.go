package algorithms

import (
	"errors"
	"math/rand"
)

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
