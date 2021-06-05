package sequence

type VolumePattern []float64

type HasVolumes interface {
	GetVolumePattern() []float64
	SetVolumePattern([]float64)
}

func (vp *VolumePattern) GetVolumePattern() []float64 {
	return *vp
}

func (vp *VolumePattern) SetVolumePattern(vs []float64) {
	*vp = vs
}

// Volumes sets the generator's Volume pattern
func Volumes(vs ...float64) Option {
	return func(g Generator) {
		if hvs, ok := g.(HasVolumes); ok {
			hvs.SetVolumePattern(vs)
		}
	}
}

// VolumeAt sets the n'th value in the entire play sequence
// to be Volume p. This could involve duplicating a pattern
// until it is long enough to reach n. Meaningless if the
// Volume pattern has not been set yet.
func VolumeAt(v float64, n int) Option {
	return func(g Generator) {
		if hvs, ok := g.(HasVolumes); ok {
			if hl, ok := hvs.(HasLength); ok {
				if hl.GetLength() < n {
					volumes := hvs.GetVolumePattern()
					if len(volumes) == 0 {
						return
					}
					// If the pattern is not long enough, there are two things
					// we could do-- 1. Extend the pattern and replace the
					// individual note, or 2. Replace the note that would be
					// played at n and thus all earlier and later plays within
					// the pattern as well.
					//
					// This uses approach 1.
					for len(volumes) < n {
						volumes = append(volumes, volumes...)
					}
					volumes[n] = v
					hvs.SetVolumePattern(volumes)
				}
			}
		}
	}
}

// VolumePatternAt sets the n'th value in the Volume pattern
// to be Volume p. Meaningless if the Volume pattern has not
// been set yet.
func VolumePatternAt(v float64, n int) Option {
	return func(g Generator) {
		if hvs, ok := g.(HasVolumes); ok {
			volumes := hvs.GetVolumePattern()
			if len(volumes) < n {
				return
			}
			volumes[n] = v
			hvs.SetVolumePattern(volumes)
		}
	}
}
