package algorithms

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
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

func ChooseX(weights []float64, times int) []int {
	lengthWeights := len(weights)
	out := make([]int, times)
	remainingWeights := make([]float64, lengthWeights)
	remainingWeights[lengthWeights-1] = weights[lengthWeights-1]
	for i := lengthWeights - 2; i >= 0; i-- {
		remainingWeights[i] = remainingWeights[i+1] + weights[i]
	}
	for i := 0; i < times; i++ {
		out[i] = CumWeightedChooseOne(remainingWeights)
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
		dlog.Error("Idx", i)
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
