package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goferwplynie/kompresja/bits/bitutils"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/logger"
)

func main() {
	start := time.Now()
	compress("test", "txt")
	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	decompress("test.gofr")
	logger.Cute(time.Since(start))
}

func compress(filename string, extension string) {
	file, err := os.ReadFile(filename + "." + extension)
	if err != nil {
		logger.Error("cant read file :c")
		logger.Error(err)
	}
	text := string(file)
	chars := make([]string, 0, len(file))

	for _, v := range text {
		chars = append(chars, string(v))
	}

	encoded := huffman.Encode(chars)
	logger.Log(len(encoded))
	bytes := bitutils.BitStringToBytes(encoded)
	err = os.WriteFile(filename+".gofr", bytes, 0644)
	if err != nil {
		logger.Error("cant save to file :c")
		logger.Error(err)
	}

	logger.Log(encoded)

	logger.Cute(fmt.Sprintf("original: %vB", len(file)))
	logger.Cute(fmt.Sprintf("compressed: %vB", len(encoded)/8))

}

func decompress(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error(err)
		logger.Error(":c")
		return
	}

	//using slice for later MTF
	chars := huffman.Decode(data)
	logger.Log(chars)

	str := strings.Join(chars, "")

	os.WriteFile("decompressed.txt", []byte(str), 0644)
}
