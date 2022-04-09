// Package dls contains data structures for DLS (.dls) file types.
package dls

import "github.com/oakmound/oak/v3/audio/format/riff"

// The DLS is the major struct we care about in this package
// DLS files contain instrument and wave sample information, and
// a bunch of other things users probably don't care about.
type DLS struct {
	Dlid      ID     `riff:"dlid"`
	Colh      uint32 `riff:"colh"`
	Vers      int64  `riff:"vers"`
	Lins      []Ins  `riff:"lins"`
	Ptbl      []byte `riff:"ptbl"` //PoolTable
	Wvpl      []Wave `riff:"wvpl"`
	riff.INFO `riff:"INFO"`
}

// PoolTable is a goofy name for a thing that redirects references
// between instruments and waves, I think.
type PoolTable struct {
	CbSize uint32
	CCues  uint32
	// CCues size
	PoolCues []uint32
}

// An ID is a unique identifer for a dls file or instrument (or wave).
// This could just be written as a complex128.
type ID struct {
	UlData1 uint32
	UlData2 uint16
	UlData3 uint16
	AbData4 [8]byte
}

// Wave is the underlying struct you'd also find in WAV files. It stores raw
// audio information and headers describing how to play that information. The
// DLS Wave struct can also have a DLSID. Todo: Consider moving this out of this
// file entirely and into the WAV package, the downside of which would be that
// if a user wanted access to a DLSID it would no longer be there to get.
type Wave struct {
	Dlid ID        `riff:"dlid"`
	GUID []byte    `riff:"guid"`
	Wavu []byte    `riff:"wavu"`
	Fmt  PCMFormat `riff:"fmt "`
	Wavh []byte    `riff:"wavh"`
	Smpl []byte    `riff:"smpl"`
	Wsmp []byte    `riff:"wsmp"`
	// Data is the stuff you actually care about
	Data      []byte `riff:"data"`
	riff.INFO `riff:"INFO"`
}

// PCMFormat is a wave format that just know how many bits per sample a wave
// takes, beyond common format fields
// Really, there are two formats a Wave can take, this is just the one we hope
// to see. Todo: Fix that
type PCMFormat struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	WhoKnows      uint16 // Test files show there are buffer bytes here?
}

// An Ins holds instrument data
type Ins struct {
	Dlid      ID        `riff:"dlid"`
	Insh      InsHeader `riff:"insh"`
	Lrgn      []Rgn     `riff:"lrgn"`
	Lart      Art       `riff:"lart"`
	riff.INFO `riff:"INFO"`
}

// InsHeader stores header information for an instrument, notably the number
// of regions and the internal instrument bank and number
type InsHeader struct {
	CRegions uint32
	Locale   MIDILOCALE
}

// MIDILOCALE stores two of the fields in an instrument header
type MIDILOCALE struct {
	UlBank       uint32
	UlInstrument uint32
}

// An Art is something we need to look more into
type Art struct {
	// Todo: art1 doesn't fit our unmarshaler's expectations, because
	// it's basically its own type of subchunk with two sizes following
	// 'art1' then a number of structs based on the second size
	Art1 []byte `riff:"art1"`
	// This []byte is really:
	// cbSize uint32
	// cConnectionBlocks uint32
	// ConnectionBlocks []ConnectionBlock
	// Also Art2, which is equivalent to Art1
}

// Rgn is a region linking instruments to waves
type Rgn struct {
	Rgnh RgnHeader `riff:"rgnh"`
	Wsmp []byte    `riff:"wsmp"`
	Wlnk WaveLink  `riff:"wlnk"`
	Lart Art       `riff:"lart"`
}

// An RgnHeader stores header information for regions, notably for one the valid range
// of notes an instrument should be applied to
// Todo: figure out how to distinguish between rgnhs with and without ulLayer
type RgnHeader struct {
	RangeKey      RGNRANGE
	RangeVelocity RGNRANGE
	FusOptions    uint16
	UsKeyGroup    uint16
	UsLayer       uint16 // This field is optional
}

// An RGNRANGE just stores a low and a high value
type RGNRANGE struct {
	UsLow  uint16
	UsHigh uint16
}

// WaveLink stores things I don't know about
type WaveLink struct {
	FusOptions   uint16
	UsPhaseGroup uint16
	UlChannel    uint32
	UlTableIndex uint32
}

// WaveSample also stores things I don't know about
type WaveSample struct {
	CbSize      uint32
	UsUnityNote uint16
	SFineTune   int16
	LGain       int32
	FulOptions  uint32
	// As for art, WaveSampleLoop is CSampleLoops long
	CSampleLoops   uint32
	WaveSampleLoop []WaveSampleLoop
}

// WaveSampleLoop also stores things I don't know about
type WaveSampleLoop struct {
	CbSize       uint32
	UlLoopType   uint32
	UlLoopStart  uint32
	UlLoopLength uint32
}
