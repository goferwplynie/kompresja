package bwt

import "sort"

func Decode(bwt []byte, index int) []byte {
	n := len(bwt)
	result := make([]byte, n)

	count := make(map[byte]int)
	for _, b := range bwt {
		count[b]++
	}

	first := make(map[byte]int)
	sorted := make([]byte, 0, len(count))
	for b := range count {
		sorted = append(sorted, b)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	pos := 0
	for _, b := range sorted {
		first[b] = pos
		pos += count[b]
	}

	rank := make([]int, n)
	seen := make(map[byte]int)
	for i := range n {
		rank[i] = seen[bwt[i]]
		seen[bwt[i]]++
	}

	ptr := index
	for i := n - 1; i >= 0; i-- {
		result[i] = bwt[ptr]
		ptr = first[bwt[ptr]] + rank[ptr]
	}

	return result
}
