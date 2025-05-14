package mtf

import (
	"slices"
)

func Encode(b []byte) []byte {
	available := make([]byte, 0, 256)
	result := make([]byte, 0, len(b))

	for i := 0; i <= 255; i++ {
		available = append(available, uint8(i))
	}
	for _, v := range b {
		index := search(v, available)
		result = append(result, byte(index))
		available = moveToFront(index, available)
	}

	return result
}

func search(b byte, available []byte) int {
	for i, v := range available {
		if v == b {
			return i
		}
	}
	return 0
}

func moveToFront(index int, available []byte) []byte {
	b := available[index]
	available = slices.Delete(available, index, index+1)
	return append([]byte{b}, available...)
}
