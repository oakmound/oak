package filter

import (
	"math"

	"github.com/oakmound/oak/v3/audio/support/audio/filter/supports"
	"github.com/oakmound/oak/v3/audio/support/audio/manip"
)

/*****************************************************************************
* HOME URL: http://blogs.zynaptiq.com/bernsee
* KNOWN BUGS: none
*
* SYNOPSIS: Routine for doing pitch shifting while maintaining
* duration using the Short Time Fourier Transform.
*
* DESCRIPTION: The routine takes a pitchShift factor value which is between 0.5
* (one octave down) and 2. (one octave up). A value of exactly 1 does not change
* the pitch. numSampsToProcess tells the routine how many samples in indata[0...
* numSampsToProcess-1] should be pitch shifted and moved to outdata[0 ...
* numSampsToProcess-1]. The two buffers can be identical (ie. it can process the
* data in-place). fftFrameSize defines the FFT frame size used for the
* processing. Typical values are 1024, 2048 and 4096. It may be any value <=
* MAX_FRAME_LENGTH but it MUST be a power of 2. osamp is the STFT
* oversampling factor which also determines the overlap between adjacent STFT
* frames. It should at least be 4 for moderate scaling ratios. A value of 32 is
* recommended for best quality. sampleRate takes the sample rate for the signal
* in unit Hz, ie. 44100 for 44.1 kHz audio. The data passed to the routine in
* indata[] should be in the range [-1.0, 1.0), which is also the output range
* for the data, make sure you scale the data accordingly (for 16bit signed integers
* you would have to divide (and multiply) by 32768).
*
* COPYRIGHT 1999-2015 Stephan M. Bernsee <s.bernsee [AT] zynaptiq [DOT] com>
*
* 						The Wide Open License (WOL)
*
* Permission to use, copy, modify, distribute and sell this software and its
* documentation for any purpose is hereby granted without fee, provided that
* the above copyright notice and this license appear in all source copies.
* THIS SOFTWARE IS PROVIDED "AS IS" WITHOUT EXPRESS OR IMPLIED WARRANTY OF
* ANY KIND. See http://www.dspguru.com/wol.htm for more information.
*
*****************************************************************************/
// As is standard with translations of this code to other languages,
// Go translation copyright Patrick Stephen 2017
// To be clear, the PitchShift function + FFT is what had to be translated

// A PitchShifter has an encoding function that will shift
// a pitch up to an octave up or down (0.5 -> octave down, 2.0 -> octave up)
// these are for lower-level use, and a similar type that takes in steps to
// shift by (and eventually pitches to set to) will follow.
type PitchShifter interface {
	PitchShift(float64) Encoding
}

// FFTShifter holds buffers and settings for performing a pitch shift on PCM audio
type FFTShifter struct {
	fftFrameSize                      int
	oversampling                      int
	step                              int
	latency                           int
	stack, frame                      []float64
	workBuffer                        []float64
	magnitudes, frequencies           []float64
	synthMagnitudes, synthFrequencies []float64
	lastPhase, sumPhase               []float64
	outAcc                            []float64
	expected                          float64
	window, windowFactors             []float64
}

// These are built in shifters with some common inputs
var (
	LowQualityShifter, _  = NewFFTShifter(1024, 8)
	HighQualityShifter, _ = NewFFTShifter(1024, 32)
)

// NewFFTShifter returns a pitch shifter that uses fast fourier transforms
func NewFFTShifter(fftFrameSize int, oversampling int) (PitchShifter, error) {
	// Todo: check that the frame size and oversampling rate make sense
	ps := FFTShifter{}
	ps.fftFrameSize = fftFrameSize
	ps.oversampling = oversampling
	ps.step = fftFrameSize / oversampling
	ps.latency = fftFrameSize - ps.step
	ps.stack = make([]float64, fftFrameSize)
	ps.workBuffer = make([]float64, 2*fftFrameSize)
	ps.magnitudes = make([]float64, fftFrameSize)
	ps.frequencies = make([]float64, fftFrameSize)
	ps.synthMagnitudes = make([]float64, fftFrameSize)
	ps.synthFrequencies = make([]float64, fftFrameSize)
	ps.lastPhase = make([]float64, fftFrameSize/2+1)
	ps.sumPhase = make([]float64, fftFrameSize/2+1)
	ps.outAcc = make([]float64, 2*fftFrameSize)

	ps.expected = 2 * math.Pi * float64(ps.step) / float64(fftFrameSize)

	ps.window = make([]float64, fftFrameSize)
	ps.windowFactors = make([]float64, fftFrameSize)
	t := 0.0
	for i := 0; i < fftFrameSize; i++ {
		w := -0.5*math.Cos(t) + .5
		ps.window[i] = w
		ps.windowFactors[i] = w * (2.0 / float64(fftFrameSize*oversampling))
		t += (math.Pi * 2) / float64(fftFrameSize)
	}

	ps.frame = make([]float64, fftFrameSize)
	return ps, nil
}

// PitchShift modifies filtered audio by the input float, between 0.5 and 2.0,
// each end of the spectrum representing octave down and up respectively
func (ps FFTShifter) PitchShift(shiftBy float64) Encoding {
	return func(senc supports.Encoding) {
		data := *senc.GetData()
		bitDepth := *senc.GetBitDepth()
		byteDepth := bitDepth / 8
		sampleRate := *senc.GetSampleRate()
		channels := *senc.GetChannels()

		// Jeeez
		out := make([]byte, len(data))
		copy(out, data)

		freqPerBin := float64(sampleRate) / float64(ps.fftFrameSize)
		frameIndex := ps.latency

		// End jeeeez

		// for each channel individually
		for c := 0; c < int(channels); c++ {
			// convert this to a channel-specific float64 buffer
			f64in := manip.BytesToF64(data, channels, bitDepth, c)
			f64out := f64in

			for i := 0; i < len(f64in); i++ {
				// Get a frame
				ps.frame[frameIndex] = f64in[i]
				// Bug here for early i values: they'll all be 0!
				f64out[i] = ps.stack[frameIndex-ps.latency]
				frameIndex++

				// A full frame has been obtained
				if frameIndex >= ps.fftFrameSize {
					frameIndex = ps.latency

					// Windowing
					for k := 0; k < ps.fftFrameSize; k++ {
						ps.workBuffer[2*k] = ps.frame[k] * ps.window[k]
						ps.workBuffer[(2*k)+1] = 0
					}

					ShortTimeFourierTransform(ps.workBuffer, ps.fftFrameSize, -1)

					// Analysis
					for k := 0; k <= ps.fftFrameSize/2; k++ {
						real := ps.workBuffer[2*k]
						imag := ps.workBuffer[(2*k)+1]

						magn := 2 * math.Sqrt(real*real+imag*imag)
						ps.magnitudes[k] = magn

						phase := math.Atan2(imag, real)

						diff := phase - ps.lastPhase[k]
						ps.lastPhase[k] = phase

						diff -= float64(k) * ps.expected

						deltaPhase := int(diff * (1 / math.Pi))
						if deltaPhase >= 0 {
							deltaPhase += deltaPhase & 1
						} else {
							deltaPhase -= deltaPhase & 1
						}

						diff -= math.Pi * float64(deltaPhase)
						diff *= float64(ps.oversampling) / (math.Pi * 2)
						diff = (float64(k) + diff) * freqPerBin

						ps.frequencies[k] = diff
					}

					// Processing
					for k := 0; k < ps.fftFrameSize; k++ {
						ps.synthMagnitudes[k] = 0
						ps.synthFrequencies[k] = 0
					}

					for k := 0; k < ps.fftFrameSize/2; k++ {
						l := int(float64(k) * shiftBy)
						if l < ps.fftFrameSize/2 {
							ps.synthMagnitudes[l] += ps.magnitudes[k]
							ps.synthFrequencies[l] = ps.frequencies[k] * shiftBy
						}
					}

					// Synthesis
					for k := 0; k <= ps.fftFrameSize/2; k++ {
						magn := ps.synthMagnitudes[k]
						tmp := ps.synthFrequencies[k]
						tmp -= float64(k) * freqPerBin
						tmp /= freqPerBin
						tmp *= 2 * math.Pi / float64(ps.oversampling)
						tmp += float64(k) * ps.expected
						ps.sumPhase[k] += tmp

						ps.workBuffer[2*k] = magn * math.Cos(ps.sumPhase[k])
						ps.workBuffer[(2*k)+1] = magn * math.Sin(ps.sumPhase[k])
					}

					// Remove negative frequencies
					// I don't get how we know these ones are negative
					// also this looks like it's going to overflow the slice
					for k := ps.fftFrameSize + 2; k < 2*ps.fftFrameSize; k++ {
						ps.workBuffer[k] = 0.0
					}

					ShortTimeFourierTransform(ps.workBuffer, ps.fftFrameSize, 1)

					// Windowing
					for k := 0; k < ps.fftFrameSize; k++ {
						ps.outAcc[k] += ps.windowFactors[k] * ps.workBuffer[2*k]
					}
					for k := 0; k < ps.step; k++ {
						ps.stack[k] = ps.outAcc[k]
					}

					// Shift accumulator, shift frame
					for k := 0; k < ps.fftFrameSize; k++ {
						ps.outAcc[k] = ps.outAcc[k+ps.step]
					}

					for k := 0; k < ps.latency; k++ {
						ps.frame[k] = ps.frame[k+ps.step]
					}
				}
			}
			// remap this f64in to the output
			for i := c * int(byteDepth); i < len(data); i += int(byteDepth * 2) {
				manip.SetInt16_f64(out, i, f64in[i/int(byteDepth*2)])
			}
		}
		datap := senc.GetData()
		*datap = out
	}
}

// ShortTimeFourierTransform : FFT routine, (C)1996 S.M.Bernsee. Sign = -1 is FFT, 1 is iFFT (inverse)
// Fills fftBuffer[0...2*fftFrameSize-1] with the Fourier transform of the
// time domain data in fftBuffer[0...2*fftFrameSize-1]. The FFT array takes
// and returns the cosine and sine parts in an interleaved manner, ie.
// fftBuffer[0] = cosPart[0], fftBuffer[1] = sinPart[0], asf. fftFrameSize
// must be a power of 2. It expects a complex input signal (see footnote 2),
// ie. when working with 'common' audio signals our input signal has to be
// passed as {in[0],0.,in[1],0.,in[2],0.,...} asf. In that case, the transform
// of the frequencies of interest is in fftBuffer[0...fftFrameSize].
func ShortTimeFourierTransform(data []float64, fftFrameSize, sign int) {
	for i := 2; i < 2*(fftFrameSize-2); i += 2 {
		j := 0
		for bitm := 2; bitm < 2*fftFrameSize; bitm <<= 1 {
			if (i & bitm) != 0 {
				j++
			}
			j <<= 1
		}
		if i < j {
			data[j], data[i] = data[i], data[j]
			data[j+1], data[i+1] = data[i+1], data[j+1]
		}
	}
	max := int(math.Log(float64(fftFrameSize))/math.Log(2) + .5)
	le := 2
	for k := 0; k < max; k++ {
		le <<= 1
		le2 := le >> 1
		ur := 1.0
		ui := 0.0
		arg := math.Pi / float64(le2>>1)
		wr := math.Cos(arg)
		wi := float64(sign) * math.Sin(arg)
		for j := 0; j < le2; j += 2 {
			for i := j; i < 2*fftFrameSize; i += le {
				tr := data[i+le2]*ur - data[i+le2+1]*ui
				ti := data[i+le2]*ui + data[i+le2+1]*ur
				data[i+le2] = data[i] - tr
				data[i+le2+1] = data[i+1] - ti
				data[i] += tr
				data[i+1] += ti
			}
			tmp := ur*wr - ui*wi
			ui = ur*wi + ui*wr
			ur = tmp
		}
	}
}
