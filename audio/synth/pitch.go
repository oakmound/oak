package synth

import (
	"sort"

	"github.com/oakmound/oak/v4/audio/pcm"
)

// A Pitch is a frequency value which represents how fast a wave should oscillate to produce a specific tone.
type Pitch uint16

// Pitch frequencies, taken from http://peabody.sapp.org/class/st2/lab/notehz/
// These span octave 0 through octave 8, with sharps suffixed 's' and flats suffixed 'b'
const (
	// 0 is reserved as representing a 'rest' for the purpose of composition
	Rest Pitch = 0
	C0   Pitch = 16
	C0s  Pitch = 17
	D0b  Pitch = 17
	D0   Pitch = 18
	D0s  Pitch = 20
	E0b  Pitch = 20
	E0   Pitch = 21
	F0   Pitch = 22
	F0s  Pitch = 23
	G0b  Pitch = 23
	G0   Pitch = 25
	G0s  Pitch = 26
	A0b  Pitch = 26
	A0   Pitch = 28
	A0s  Pitch = 29
	B0b  Pitch = 29
	B0   Pitch = 31
	C1   Pitch = 33
	C1s  Pitch = 35
	D1b  Pitch = 35
	D1   Pitch = 37
	D1s  Pitch = 39
	E1b  Pitch = 39
	E1   Pitch = 41
	F1   Pitch = 44
	F1s  Pitch = 46
	G1b  Pitch = 46
	G1   Pitch = 49
	G1s  Pitch = 52
	A1b  Pitch = 52
	A1   Pitch = 55
	A1s  Pitch = 58
	B1b  Pitch = 58
	B1   Pitch = 62
	C2   Pitch = 65
	C2s  Pitch = 69
	D2b  Pitch = 69
	D2   Pitch = 73
	D2s  Pitch = 78
	E2b  Pitch = 78
	E2   Pitch = 82
	F2   Pitch = 87
	F2s  Pitch = 93
	G2b  Pitch = 93
	G2   Pitch = 98
	G2s  Pitch = 104
	A2b  Pitch = 104
	A2   Pitch = 110
	A2s  Pitch = 117
	B2b  Pitch = 117
	B2   Pitch = 124
	C3   Pitch = 131
	C3s  Pitch = 139
	D3b  Pitch = 139
	D3   Pitch = 147
	D3s  Pitch = 156
	E3b  Pitch = 156
	E3   Pitch = 165
	F3   Pitch = 175
	F3s  Pitch = 185
	G3b  Pitch = 185
	G3   Pitch = 196
	G3s  Pitch = 208
	A3b  Pitch = 208
	A3   Pitch = 220
	A3s  Pitch = 233
	B3b  Pitch = 233
	B3   Pitch = 247
	C4   Pitch = 262
	C4s  Pitch = 278
	D4b  Pitch = 278
	D4   Pitch = 294
	D4s  Pitch = 311
	E4b  Pitch = 311
	E4   Pitch = 330
	F4   Pitch = 349
	F4s  Pitch = 370
	G4b  Pitch = 370
	G4   Pitch = 392
	G4s  Pitch = 415
	A4b  Pitch = 415
	A4   Pitch = 440
	A4s  Pitch = 466
	B4b  Pitch = 466
	B4   Pitch = 494
	C5   Pitch = 523
	C5s  Pitch = 554
	D5b  Pitch = 554
	D5   Pitch = 587
	D5s  Pitch = 622
	E5b  Pitch = 622
	E5   Pitch = 659
	F5   Pitch = 699
	F5s  Pitch = 740
	G5b  Pitch = 740
	G5   Pitch = 784
	G5s  Pitch = 831
	A5b  Pitch = 831
	A5   Pitch = 880
	A5s  Pitch = 932
	B5b  Pitch = 932
	B5   Pitch = 988
	C6   Pitch = 1047
	C6s  Pitch = 1109
	D6b  Pitch = 1109
	D6   Pitch = 1175
	D6s  Pitch = 1245
	E6b  Pitch = 1245
	E6   Pitch = 1319
	F6   Pitch = 1397
	F6s  Pitch = 1475
	G6b  Pitch = 1475
	G6   Pitch = 1568
	G6s  Pitch = 1661
	A6b  Pitch = 1661
	A6   Pitch = 1760
	A6s  Pitch = 1865
	B6b  Pitch = 1865
	B6   Pitch = 1976
	C7   Pitch = 2093
	C7s  Pitch = 2218
	D7b  Pitch = 2218
	D7   Pitch = 2349
	D7s  Pitch = 2489
	E7b  Pitch = 2489
	E7   Pitch = 2637
	F7   Pitch = 2794
	F7s  Pitch = 2960
	G7b  Pitch = 2960
	G7   Pitch = 3136
	G7s  Pitch = 3322
	A7b  Pitch = 3322
	A7   Pitch = 3520
	A7s  Pitch = 3729
	B7b  Pitch = 3729
	B7   Pitch = 3951
	C8   Pitch = 4186
	C8s  Pitch = 4435
	D8b  Pitch = 4435
	D8   Pitch = 4699
	D8s  Pitch = 4978
	E8b  Pitch = 4978
	E8   Pitch = 5274
	F8   Pitch = 5588
	F8s  Pitch = 5920
	G8b  Pitch = 5920
	G8   Pitch = 6272
	G8s  Pitch = 6645
	A8b  Pitch = 6645
	A8   Pitch = 7040
	A8s  Pitch = 7459
	B8b  Pitch = 7459
	B8   Pitch = 7902
)

var (
	allPitches = []Pitch{
		C0,
		C0s,
		D0,
		D0s,
		E0,
		F0,
		F0s,
		G0,
		G0s,
		A0,
		A0s,
		B0,
		C1,
		C1s,
		D1,
		D1s,
		E1,
		F1,
		F1s,
		G1,
		G1s,
		A1,
		A1s,
		B1,
		C2,
		C2s,
		D2,
		D2s,
		E2,
		F2,
		F2s,
		G2,
		G2s,
		A2,
		A2s,
		B2,
		C3,
		C3s,
		D3,
		D3s,
		E3,
		F3,
		F3s,
		G3,
		G3s,
		A3,
		A3s,
		B3,
		C4,
		C4s,
		D4,
		D4s,
		E4,
		F4,
		F4s,
		G4,
		G4s,
		A4,
		A4s,
		B4,
		C5,
		C5s,
		D5,
		D5s,
		E5,
		F5,
		F5s,
		G5,
		G5s,
		A5,
		A5s,
		B5,
		C6,
		C6s,
		D6,
		D6s,
		E6,
		F6,
		F6s,
		G6,
		G6s,
		A6,
		A6s,
		B6,
		C7,
		C7s,
		D7,
		D7s,
		E7,
		F7,
		F7s,
		G7,
		G7s,
		A7,
		A7s,
		B7,
		C8,
		C8s,
		D8,
		D8s,
		E8,
		F8,
		F8s,
		G8,
		G8s,
		A8,
		A8s,
		B8,
	}

	// Reverse lookup for allPitches
	noteIndices = map[Pitch]int{
		C0:  0,
		C0s: 1,
		D0:  2,
		D0s: 3,
		E0:  4,
		F0:  5,
		F0s: 6,
		G0:  7,
		G0s: 8,
		A0:  9,
		A0s: 10,
		B0:  11,
		C1:  12,
		C1s: 13,
		D1:  14,
		D1s: 15,
		E1:  16,
		F1:  17,
		F1s: 18,
		G1:  19,
		G1s: 20,
		A1:  21,
		A1s: 22,
		B1:  23,
		C2:  24,
		C2s: 25,
		D2:  26,
		D2s: 27,
		E2:  28,
		F2:  29,
		F2s: 30,
		G2:  31,
		G2s: 32,
		A2:  33,
		A2s: 34,
		B2:  35,
		C3:  36,
		C3s: 37,
		D3:  38,
		D3s: 39,
		E3:  40,
		F3:  41,
		F3s: 42,
		G3:  43,
		G3s: 44,
		A3:  45,
		A3s: 46,
		B3:  47,
		C4:  48,
		C4s: 49,
		D4:  50,
		D4s: 51,
		E4:  52,
		F4:  53,
		F4s: 54,
		G4:  55,
		G4s: 56,
		A4:  57,
		A4s: 58,
		B4:  59,
		C5:  60,
		C5s: 61,
		D5:  62,
		D5s: 63,
		E5:  64,
		F5:  65,
		F5s: 66,
		G5:  67,
		G5s: 68,
		A5:  69,
		A5s: 70,
		B5:  71,
		C6:  72,
		C6s: 73,
		D6:  74,
		D6s: 75,
		E6:  76,
		F6:  77,
		F6s: 78,
		G6:  79,
		G6s: 80,
		A6:  81,
		A6s: 82,
		B6:  83,
		C7:  84,
		C7s: 85,
		D7:  86,
		D7s: 87,
		E7:  88,
		F7:  89,
		F7s: 90,
		G7:  91,
		G7s: 92,
		A7:  93,
		A7s: 94,
		B7:  95,
		C8:  96,
		C8s: 97,
		D8:  98,
		D8s: 99,
		E8:  100,
		F8:  101,
		F8s: 102,
		G8:  103,
		G8s: 104,
		A8:  105,
		A8s: 106,
		B8:  107,
	}

	pitchStrings = map[Pitch]string{
		Rest: "Rest",
		C0:   "C0",
		C0s:  "C0#",
		D0:   "D0",
		D0s:  "D0#",
		E0:   "E0",
		F0:   "F0",
		F0s:  "F0#",
		G0:   "G0",
		G0s:  "G0#",
		A0:   "A0",
		A0s:  "A0#",
		B0:   "B0",
		C1:   "C1",
		C1s:  "C1#",
		D1:   "D1",
		D1s:  "D1#",
		E1:   "E1",
		F1:   "F1",
		F1s:  "F1#",
		G1:   "G1",
		G1s:  "G1#",
		A1:   "A1",
		A1s:  "A1#",
		B1:   "B1",
		C2:   "C2",
		C2s:  "C2#",
		D2:   "D2",
		D2s:  "D2#",
		E2:   "E2",
		F2:   "F2",
		F2s:  "F2#",
		G2:   "G2",
		G2s:  "G2#",
		A2:   "A2",
		A2s:  "A2#",
		B2:   "B2",
		C3:   "C3",
		C3s:  "C3#",
		D3:   "D3",
		D3s:  "D3#",
		E3:   "E3",
		F3:   "F3",
		F3s:  "F3#",
		G3:   "G3",
		G3s:  "G3#",
		A3:   "A3",
		A3s:  "A3#",
		B3:   "B3",
		C4:   "C4",
		C4s:  "C4#",
		D4:   "D4",
		D4s:  "D4#",
		E4:   "E4",
		F4:   "F4",
		F4s:  "F4#",
		G4:   "G4",
		G4s:  "G4#",
		A4:   "A4",
		A4s:  "A4#",
		B4:   "B4",
		C5:   "C5",
		C5s:  "C5#",
		D5:   "D5",
		D5s:  "D5#",
		E5:   "E5",
		F5:   "F5",
		F5s:  "F5#",
		G5:   "G5",
		G5s:  "G5#",
		A5:   "A5",
		A5s:  "A5#",
		B5:   "B5",
		C6:   "C6",
		C6s:  "C6#",
		D6:   "D6",
		D6s:  "D6#",
		E6:   "E6",
		F6:   "F6",
		F6s:  "F6#",
		G6:   "G6",
		G6s:  "G6#",
		A6:   "A6",
		A6s:  "A6#",
		B6:   "B6",
		C7:   "C7",
		C7s:  "C7#",
		D7:   "D7",
		D7s:  "D7#",
		E7:   "E7",
		F7:   "F7",
		F7s:  "F7#",
		G7:   "G7",
		G7s:  "G7#",
		A7:   "A7",
		A7s:  "A7#",
		B7:   "B7",
		C8:   "C8",
		C8s:  "C8#",
		D8:   "D8",
		D8s:  "D8#",
		E8:   "E8",
		F8:   "F8",
		F8s:  "F8#",
		G8:   "G8",
		G8s:  "G8#",
		A8:   "A8",
		A8s:  "A8#",
		B8:   "B8",
	}
)

func (p Pitch) String() string {
	return pitchStrings[p]
}

var accidentals = map[Pitch]struct{}{
	C0s: {},
	D0s: {},
	F0s: {},
	G0s: {},
	A0s: {},
	C1s: {},
	D1s: {},
	F1s: {},
	G1s: {},
	A1s: {},
	C2s: {},
	D2s: {},
	F2s: {},
	G2s: {},
	A2s: {},
	C3s: {},
	D3s: {},
	F3s: {},
	G3s: {},
	A3s: {},
	C4s: {},
	D4s: {},
	F4s: {},
	G4s: {},
	A4s: {},
	C5s: {},
	D5s: {},
	F5s: {},
	G5s: {},
	A5s: {},
	C6s: {},
	D6s: {},
	F6s: {},
	G6s: {},
	A6s: {},
	C7s: {},
	D7s: {},
	F7s: {},
	G7s: {},
	A7s: {},
	C8s: {},
	D8s: {},
	F8s: {},
	G8s: {},
	A8s: {},
}

// A Step is an index offset on a pitch
// to raise or lower it to a relative new pitch
type Step int

// Step values
const (
	HalfStep  Step = 1
	WholeStep Step = 2
	Octave    Step = 12
)

// Up raises a pitch s steps
func (p Pitch) Up(s Step) Pitch {
	i := noteIndices[p]
	if i+int(s) >= len(allPitches) {
		return allPitches[len(allPitches)-1]
	}
	return allPitches[i+int(s)]
}

// Down lowers a pitch s steps
func (p Pitch) Down(s Step) Pitch {
	i := noteIndices[p]
	if i-int(s) < 0 {
		return allPitches[0]
	}
	return allPitches[i-int(s)]
}

// IsAccidental reports true if this pitch is represented with a single sharp or a flat, usually.
func (p Pitch) IsAccidental() bool {
	_, ok := accidentals[p]
	return ok
}

type PitchDetector struct {
	pcm.Reader

	format pcm.Format

	// DetectedPitches and DetectedRawPitches store the calculated pitch values as this reader parses data. The length
	// of these slices will be equal to this reader's format's channel count. Consumers should not modify these slices.
	DetectedPitches    []Pitch
	DetectedRawPitches []float64

	indices     []int
	lastValues  []float64
	crossedZero []bool
}

func NewPitchDetector(r pcm.Reader) *PitchDetector {
	return &PitchDetector{
		Reader:             r,
		format:             r.PCMFormat(),
		DetectedPitches:    make([]Pitch, r.PCMFormat().Channels),
		DetectedRawPitches: make([]float64, r.PCMFormat().Channels),
		indices:            make([]int, r.PCMFormat().Channels),
		lastValues:         make([]float64, r.PCMFormat().Channels),
		crossedZero:        make([]bool, r.PCMFormat().Channels),
	}
}

func (pd *PitchDetector) ReadPCM(b []byte) (n int, err error) {
	n, err = pd.Reader.ReadPCM(b)
	if err != nil {
		return n, err
	}
	var read int
	sampleSize := pd.format.SampleSize()
	for len(b[read:]) > sampleSize {
		vals, valReadBytes, err := pd.format.SampleFloat(b[read:])
		if err != nil {
			break
		}
		read += valReadBytes
		for i, val := range vals {
			pd.indices[i]++
			if pd.lastValues[i] < 0 && val > 0 || val < 0 && pd.lastValues[i] > 0 {
				// we've crossed zero
				if !pd.crossedZero[i] {
					pd.crossedZero[i] = true
				} else {
					// assuming this is pitched audio (if it isn't we can't give a correct answer),
					// pd.index is now the number of samples since the last time this audio
					// stream crossed zero. The second last time this audio stream crossed zero defines how
					// frequently this audio is cycling-- the speed the audio cycles at defines the pitch
					// of the audio in hertz; our pitch constants above are also defined in hertz.
					periodLength := pd.indices[i] * 2
					samplesPerSecond := pd.format.SampleRate
					periodHz := 1 / (float64(periodLength) / float64(samplesPerSecond))
					pd.DetectedRawPitches[i] = periodHz
					pd.DetectedPitches[i] = Pitch(periodHz).Round()
				}
				pd.indices[i] = 0
			}
			pd.lastValues[i] = val
		}
	}
	return
}

// Round rounds a pitch value to the closest predefined pitch value in hertz:
//	func main() {
//		hz := synth.Pitch(1024)
// 		hz2 := hz.Round()
// 		fmt.Println(hz2, int(hz2))) // "C6", 1047
//	}
//
func (p Pitch) Round() Pitch {
	// binary search
	i := sort.Search(len(allPitches)-1, func(i int) bool {
		return p < allPitches[i]
	})
	// adjust for near matches
	// we know hz < allPitches[i]
	if i == 0 {
		return allPitches[i]
	}
	if p-allPitches[i-1] < allPitches[i]-p {
		return allPitches[i-1]
	}
	return allPitches[i]
}
