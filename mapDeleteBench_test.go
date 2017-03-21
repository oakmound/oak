package oak

import "testing"

var (
	potStrings = []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
		"L",
		"M",
		"N",
		"O",
		"P",
		"Q",
		"R",
		"S",
		"T",
		"U",
		"V",
		"W",
		"X",
		"Y",
		"Z",
	}
)

func BenchmarkMapDelete(b *testing.B) {
	m := make(map[string]bool)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, s := range potStrings {
			m[s] = true
		}
		for _, s := range potStrings {
			delete(m, s)
		}
	}
}

func BenchmarkMapSetFalse(b *testing.B) {
	m := make(map[string]bool)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, s := range potStrings {
			m[s] = true
		}
		for _, s := range potStrings {
			m[s] = false
		}
	}
}
