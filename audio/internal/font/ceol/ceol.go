// Package ceol provides functionality to handle .ceol files and .ceol encoded data (Bosca Ceoil files)
package ceol

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/oakmound/oak/v3/audio/internal/sequence"
	"github.com/oakmound/oak/v3/audio/internal/synth"
)

// Raw Ceol types, holds all information in ceol file

// Ceol represents a complete .ceol file
type Ceol struct {
	Version       int
	Swing         int
	Effect        int
	EffectValue   int
	Bpm           int
	PatternLength int
	BarLength     int
	Instruments   []Instrument
	Patterns      []Pattern
	LoopStart     int
	LoopEnd       int
	Arrangement   [][8]int
}

// Instrument represents a single entry in a .ceol's instrument block
type Instrument struct {
	Index        int
	IsDrumkit    int
	Palette      int
	LPFCutoff    int
	LPFResonance int
	Volume       int
}

// Pattern represents a single entry in a .ceol's pattern block
type Pattern struct {
	Key        int
	Scale      int
	Instrument int
	Palette    int
	Notes      []Note
	Filters    []Filter
}

// Note represents a single entry in a .ceol's pattern's note block
type Note struct {
	PitchIndex int // C4 = 60
	Length     int
	Offset     int
}

// Filter represents a single entry in a .ceol's pattern's filter block
type Filter struct {
	Volume       int
	LPFCutoff    int
	LPFResonance int
}

// ChordPattern converts a Ceol's patterns and arrangement into a playable chord
// pattern for sequences
func (c Ceol) ChordPattern() sequence.ChordPattern {
	chp := sequence.ChordPattern{}
	chp.Pitches = make([][]synth.Pitch, c.PatternLength*len(c.Arrangement))
	chp.Holds = make([][]time.Duration, c.PatternLength*len(c.Arrangement))
	for i, m := range c.Arrangement {
		for _, p := range m {
			if p != -1 {
				for _, n := range c.Patterns[p].Notes {
					chp.Pitches[n.Offset+i*c.PatternLength] =
						append(chp.Pitches[n.Offset+i*c.PatternLength], synth.NoteFromIndex(n.PitchIndex))
					chp.Holds[n.Offset+i*c.PatternLength] =
						append(chp.Holds[n.Offset+i*c.PatternLength], DurationFromQuarters(c.Bpm, n.Length))
				}
			}
		}
	}
	return chp
}

// DurationFromQuarters should not be here, should be in a package
// managing bpm and time
// Duration from quarters expects four quarters to occur per beat,
// (direct complaints at terry cavanagh), and returns a time.Duration
// for n quarters in the given bpm.
func DurationFromQuarters(bpm, quarters int) time.Duration {
	beatTime := time.Duration(60000/bpm) * time.Millisecond
	quarterTime := beatTime / 4
	return quarterTime * time.Duration(quarters)
}

// Open returns a Ceol from an io.Reader
func Open(r io.Reader) (Ceol, error) {
	c := Ceol{}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return c, err
	}
	s := string(b)
	in := strings.Split(s, ",")
	ints := make([]int, len(in))
	for i := 0; i < len(in)-1; i++ {
		ints[i], err = strconv.Atoi(in[i])
		if err != nil {
			return c, err
		}
	}
	i := 0
	c.Version = ints[i]
	i++
	c.Swing = ints[i]
	i++
	c.Effect = ints[i]
	i++
	c.EffectValue = ints[i]
	i++
	c.Bpm = ints[i]
	i++
	c.PatternLength = ints[i]
	i++
	c.BarLength = ints[i]
	i++
	nInstruments := ints[i]
	i++
	c.Instruments = make([]Instrument, nInstruments)
	for j := 0; j < nInstruments; j++ {
		c.Instruments[j].Index = ints[i]
		i++
		c.Instruments[j].IsDrumkit = ints[i]
		i++
		c.Instruments[j].Palette = ints[i]
		i++
		c.Instruments[j].LPFCutoff = ints[i]
		i++
		c.Instruments[j].LPFResonance = ints[i]
		i++
		c.Instruments[j].Volume = ints[i]
		i++
	}
	nPatterns := ints[i]
	i++
	c.Patterns = make([]Pattern, nPatterns)
	for j := 0; j < nPatterns; j++ {
		c.Patterns[j].Key = ints[i]
		i++
		c.Patterns[j].Scale = ints[i]
		i++
		c.Patterns[j].Instrument = ints[i]
		i++
		c.Patterns[j].Palette = ints[i]
		i++
		nNotes := ints[i]
		i++
		c.Patterns[j].Notes = make([]Note, nNotes)
		for k := 0; k < nNotes; k++ {
			c.Patterns[j].Notes[k].PitchIndex = ints[i]
			i++
			c.Patterns[j].Notes[k].Length = ints[i]
			i++
			c.Patterns[j].Notes[k].Offset = ints[i]
			i++
			i++ // Dummy value here
		}
		hasFilter := ints[i]
		i++
		var nFilters int
		if hasFilter == 1 {
			nFilters = ints[i]
			i++
		}
		c.Patterns[j].Filters = make([]Filter, nFilters)
		for k := 0; k < nFilters; k++ {
			c.Patterns[j].Filters[k].Volume = ints[i]
			i++
			c.Patterns[j].Filters[k].LPFCutoff = ints[i]
			i++
			c.Patterns[j].Filters[k].LPFResonance = ints[i]
			i++
		}
	}
	songLength := ints[i]
	i++
	c.LoopStart = ints[i]
	i++
	c.LoopEnd = ints[i]
	i++
	c.Arrangement = make([][8]int, songLength)
	for j := 0; j < songLength; j++ {
		for k := 0; k < 8; k++ {
			c.Arrangement[j][k] = ints[i]
			i++
		}
	}
	return c, nil
}
