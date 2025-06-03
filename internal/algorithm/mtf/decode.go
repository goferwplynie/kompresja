package mtf

func Decode(b []byte) []byte {
	result := make([]byte, 0, len(b))
	available := make([]byte, 0, 256)

	for i := 0; i <= 255; i++ {
		available = append(available, uint8(i))
	}
	for _, v := range b {
		symbol := available[v]
		index := search(symbol, available)
		result = append(result, symbol)
		available = moveToFront(index, available)
	}
	return result
}
