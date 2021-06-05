package manip

func SetInt16(d []byte, i int, in int64) {
	for j := 0; j < 2; j++ {
		d[i+j] = byte(in & 255)
		in >>= 8
	}
}

func GetInt16(d []byte, i int) (out int16) {
	var shift uint16
	for j := 0; j < 2; j++ {
		out += int16(d[i+j]) << shift
		shift += 8
	}
	return
}

func GetFloat64(d []byte, i int, byteDepth uint16) float64 {
	switch byteDepth {
	case 1:
		return float64(int8(d[i])) / 128.0
	case 2:
		return float64(GetInt16(d, i)) / 32768.0
	}
	return 0.0
}

func SetInt16_f64(d []byte, i int, in float64) {
	SetInt16(d, i, int64(in*32768))
}

func Round(f float64) int64 {
	if f < 0 {
		return int64(f - .5)
	}
	return int64(f + .5)
}
