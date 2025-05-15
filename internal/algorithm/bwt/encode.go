package bwt

import (
	"bytes"
	"slices"
	"sort"
)

func Encode(b []byte) ([]byte, int) {
	rotations := make([][]byte, 0, len(b))
	rotation := b
	result := make([]byte, 0, len(b))
	var bwtIndex int

	rotations = append(rotations, rotation)

	for range b {
		freeRotation := make([]byte, len(b))
		copy(freeRotation, rotation)
		freeRotation = nextRotation(freeRotation)
		rotations = append(rotations, freeRotation)
		rotation = freeRotation
	}

	sort.Slice(rotations, func(i, j int) bool {
		return rotations[i][0] < rotations[j][0]
	})

	for i, v := range rotations {
		if bytes.Equal(b, v) {
			bwtIndex = i
		}
		result = append(result, v[len(v)-1])
	}

	return result, bwtIndex

}

func nextRotation(b []byte) []byte {
	by := b[0]
	b = slices.Delete(b, 0, 1)
	return append(b, by)
}
