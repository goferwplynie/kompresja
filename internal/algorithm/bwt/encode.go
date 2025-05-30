package bwt

import (
	"sort"
)

func Encode(data []byte) ([]byte, int) {
	n := len(data)
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	sort.Slice(indices, func(i, j int) bool {
		a, b := indices[i], indices[j]

		for k := range n {
			ac := data[(a+k)%n]
			bc := data[(b+k)%n]
			if ac != bc {
				return ac < bc
			}
		}
		return false
	})

	result := make([]byte, n)
	var bwtIndex int
	for i, idx := range indices {
		if idx == 0 {
			bwtIndex = i
		}
		result[i] = data[(idx+n-1)%n]
	}

	return result, bwtIndex
}
