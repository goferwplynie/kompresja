package main

import (
	"fmt"
	"time"

	"github.com/goferwplynie/kompresja/bits/bitutils"
	fileutils "github.com/goferwplynie/kompresja/fileUtils"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/logger"
)

func main() {
	start := time.Now()
	chars := make([]string, 0)
	str := "A_DEAD_DAD_CEDED_A_BAD_BABE_A_BEADED_ABACA_BED\n"

	for _, r := range str {
		chars = append(chars, string(r))
	}
	encoded := huffman.Encode(chars)
	logger.Log(len(encoded))
	bytes := bitutils.BitStringToBytes(encoded)
	err := fileutils.SaveToFile("gofr.gofr", bytes)
	if err != nil {
		logger.Error("cant save to file :c")
		logger.Error(err)
	}

	logger.Log(encoded)

	logger.Cute(fmt.Sprintf("original: %vB", len(str)))
	logger.Cute(fmt.Sprintf("compressed: %vB", len(encoded)/8))

	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	bytes, err = fileutils.ReadFile("gofr.gofr")
	if err != nil {
		logger.Error("cant read file :c")
		logger.Error(err)
	}
	bitString := bitutils.BytesToBitString(bytes)
}
