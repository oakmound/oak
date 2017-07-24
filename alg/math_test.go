package alg

import "testing"

func TestRoundF64(t *testing.T) {

	inputs := []float64{
		0.1,
		0.7,
		-1.2,
		-1.7,
	}
	outputs := []int{
		0,
		1,
		-1,
		-2,
	}

	for i, in := range inputs {
		if RoundF64(in) != outputs[i] {
			t.Fail()
		}
	}
}

func TestF64Eq(t *testing.T) {
	inputs := [][]float64{
		{0.000000001, 0},
	}
	outputs := []bool{
		true,
	}
	for i, in := range inputs {
		if F64eq(in[0], in[1]) != outputs[i] {
			t.Fail()
		}
	}
}
