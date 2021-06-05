package riff

// INFO is a common RIFF component. Most of these fields will be absent on
// any given INFO struct. Todo: consider if these should be given names
// that are informative instead of representative of their structural tag
type INFO struct {
	// Arhcival Location
	IARL string `riff:"IARL"`
	// Arist
	IART string `riff:"IART"`
	// Commissioned By
	ICMS string `riff:"ICMS"`
	// Comments
	ICMT string `riff:"ICMT"`
	// Copyright
	ICOP string `riff:"ICOP"`
	// Creation Date
	ICRD string `riff:"ICRD"`
	// Engineer
	IENG string `riff:"IENG"`
	// Genre
	IGNR string `riff:"IGNR"`
	// Keywords
	IKEY string `riff:"IKEY"`
	// Medium
	IMED string `riff:"IMED"`
	// Name
	INAM string `riff:"INAM"`
	// Product
	IPRD string `riff:"IPRD"`
	// Subject
	ISBJ string `riff:"ISBJ"`
	// Software
	ISFT string `riff:"ISFT"`
	// Source
	ISRC string `riff:"ISRC"`
	// Source Form
	ISRF string `riff:"ISRF"`
	// Technician
	ITCH string `riff:"ITCH"`
}
