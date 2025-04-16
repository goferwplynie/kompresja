package huffman

import "fmt"

func Encode(str string) {

	fmt.Println(calculateFreq(str))

}

func calculateFreq(str string) map[string]int {
	var freq = make(map[string]int)

	for _, r := range str {
		freq[string(r)]++
	}

	return freq
}
