package filter

// these fourier functions did not work for me.
// In case I can fix them, I leave them here.
// Credit Arnaud Gatouillat <fu [AT] iq12 [DOT] com>

// fourier1 has a bad name
// fourier1 is a helper function that does some kind of fourier transform math
// What are nn and isign?
// func fourier1(data []float64, nn, isign int) {
// 	n := nn << 1
// 	j := 1
// 	for i := 1; i < n; i += 2 {
// 		if j > i {
// 			data[j], data[i] = data[i], data[j]
// 			data[j+1], data[i+1] = data[i+1], data[j+1]
// 		}
// 		m := n >> 1
// 		for m >= 2 && j > m {
// 			j -= m
// 			m >>= 1
// 		}
// 		j += m
// 	}
// 	mmax := 2
// 	for n > mmax {
// 		stp := 2 * mmax
// 		theta := math.Pi * 2 / float64(isign*mmax)
// 		wpr, wpi := wprWpi(theta)
// 		wr := 1.0
// 		wi := 0.0
// 		for m := 1; m < mmax; m += 2 {
// 			for i := m; i <= n; i += stp {
// 				tr := wr*data[j] - wi*data[j+1]
// 				ti := wr*data[j+1] - wi*data[i]
// 				data[j] = data[i] - tr
// 				data[j+1] = data[i+1] - ti
// 				data[i] += tr
// 				data[i+1] += ti
// 			}
// 			wt := wr
// 			wr = wr*wpr - wi*wpi + wr
// 			wi = wi*wpr + wt*wpi + wi
// 		}
// 		mmax = stp
// 	}
// }

// func RealFourierTransform(data []float64, n, isign int) {
// 	theta := math.Pi / float64(n)
// 	var c2 float64
// 	if isign == 1 {
// 		c2 = -.5
// 		fourier1(data, n, 1)
// 	} else {
// 		c2 = .5
// 		theta *= -1
// 	}
// 	wpr, wpi := wprWpi(theta)
// 	wr := 1.0 + wpr
// 	wi := wpi
// 	// Wow what a great name for this variable
// 	n2p3 := 2*n + 3
// 	for i := 2; i <= n/2; i++ {
// 		i1 := i + i - 1
// 		i2 := i1 + 1
// 		i3 := n2p3 - i2
// 		i4 := i3 + 1
// 		h1r := .5 * (data[i1] + data[i3])
// 		h1i := .5 * (data[i2] - data[i4])
// 		h2r := -c2 * (data[i2] + data[i4])
// 		h2i := c2 * (data[i1] - data[i3])
// 		data[i1] = h1r + wr*h2r - wi*h2i
// 		data[i2] = h1i + wr*h2i + wi*h2r
// 		data[i3] = h1r - wr*h2r + wi*h2i
// 		data[i4] = -h1i + wr*h2i + wi*h2r
// 		wt := wr
// 		wr = wr*wpr - wi*wpi + wr
// 		wi = wi*wpr + wt*wpi + wi
// 	}
// 	if isign == 1 {
// 		data[1], data[2] = (data[1] + data[2]), (data[1] - data[2])
// 	} else {
// 		data[1], data[2] = .5*(data[1]+data[2]), .5*(data[1]-data[2])
// 		fourier1(data, n, -1)
// 	}
// }

// func wprWpi(theta float64) (float64, float64) {
// 	w := math.Sin(0.5 * theta)
// 	wpr := -2 * math.Pow(w, 2)
// 	wpi := math.Sin(theta)
// 	return wpr, wpi
// }
