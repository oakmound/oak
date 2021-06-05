package debugstream

import (
	"math"
	"testing"
)

func Test_jaroDecreased(t *testing.T) {
	type args struct {
		candidate  string
		registered string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want2 float64
	}{
		{"fullmatch", args{"super", "super"}, 1, 1},
		{"partial by paper", args{"CRATE", "TRACE"}, 2.2 / 3.0, (2.2 / 3.0) * 1.2 * 1.05},
		{"nomatch", args{"aaaaa", "super"}, 0, 0},

		{"partialex", args{"afulls", "fullscreen"}, 7.0 / 9.0, 1},
		{"partialex2", args{"scope", "help"}, 3.0/20.0 + 1.0/3.0, (3.0/20.0 + 1.0/3.0) * 1.04},
		{"low", args{"full", "help"}, 0.5, 0.5 * 1.04},

		{"single", args{"f", "fullscreen"}, (2.1) / 3.0, (2.1) / 3.0 * 1.1 * 1.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pseudoJ, psuedoJWithPref := jaroDecreased(tt.args.candidate, tt.args.registered)
			if floatSig(4, pseudoJ) != floatSig(4, tt.want) {
				t.Errorf("jaroDecreased of %s val = %v, wanted %v", tt.name, pseudoJ, tt.want)
			}
			if floatSig(4, psuedoJWithPref) != floatSig(4, tt.want2) {
				t.Errorf("jaroDecreased of %s with boost = %v, wanted %v", tt.name, psuedoJWithPref, tt.want2)
			}
		})
	}
}

func floatSig(bits int, target float64) float64 {
	factor := math.Pow10(bits)
	return math.Floor(target*factor) / factor
}
