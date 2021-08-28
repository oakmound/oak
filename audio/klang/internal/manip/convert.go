package manip

func BytesToF64(data []byte, channels, bitRate uint16, channel int) []float64 {
	byteDepth := bitRate / 8
	out := make([]float64, (len(data)/int(byteDepth*channels))+1)
	for i := channel * int(byteDepth); i < len(data); i += int(byteDepth * channels) {
		out[i/int(byteDepth*channels)] = GetFloat64(data, i, byteDepth)
	}
	return out
}
