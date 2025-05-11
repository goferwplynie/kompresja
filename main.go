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
	compress("test", "txt")
	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	decompress("test.gofr")
	logger.Cute(time.Since(start))
}

func compress(filename string, extension string) {
	file, err := fileutils.ReadFile(filename + "." + extension)
	if err != nil {
		logger.Error("cant read file :c")
		logger.Error(err)
	}
	chars := make([]string, 0, len(file))

	for _, b := range file {
		chars = append(chars, string(b))
	}

	encoded := huffman.Encode(chars)
	logger.Log(len(encoded))
	bytes := bitutils.BitStringToBytes(encoded)
	err = fileutils.SaveToFile(filename+".gofr", bytes)
	if err != nil {
		logger.Error("cant save to file :c")
		logger.Error(err)
	}

	logger.Log(encoded)

	logger.Cute(fmt.Sprintf("original: %vB", len(file)))
	logger.Cute(fmt.Sprintf("compressed: %vB", len(encoded)/8))

}

func decompress(filename string) {
	data, err := fileutils.ReadFile(filename)
	if err != nil {
		logger.Error(err)
		logger.Error(":c")
		return
	}

	//using slice for later MTF
	chars := huffman.Decode(data)
	logger.Log(chars)

	str := ""

	for _, v := range chars {
		str += v
	}

	fileutils.SaveToFile("decompressed.txt", []byte(str))
}
