package bwt

import (
	"bytes"
	"sort"
)

func Encode(b []byte) ([]byte, int) {
	n := len(b)
	rotations := make([][]byte, 0, n)

	original := make([]byte, n)
	copy(original, b)

	rotation := make([]byte, n)
	copy(rotation, b)
	for range n {
		rotations = append(rotations, append([]byte{}, rotation...))
		rotation = append(rotation[1:], rotation[0])
	}

	sort.Slice(rotations, func(i, j int) bool {
		return bytes.Compare(rotations[i], rotations[j]) < 0
	})

	result := make([]byte, n)
	var bwtIndex int
	for i, rot := range rotations {
		result[i] = rot[n-1]
		if bytes.Equal(rot, original) {
			bwtIndex = i
		}
	}

	return result, bwtIndex
}
