package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goferwplynie/kompresja/internal/algorithm/bwt"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/internal/algorithm/mtf"
	"github.com/goferwplynie/kompresja/logger"
)

func main() {
	start := time.Now()
	compress("test", "json")
	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	decompress("test.gofr")
	logger.Cute(time.Since(start))
}

func compress(filename string, extension string) {
	data, _ := os.ReadFile(filename + "." + extension)

	b, _ := bwt.Encode(data)
	fmt.Printf("b: %v\n", b)
	b = mtf.Encode(b)
	b = huffman.Encode(b)
	os.WriteFile(filename+".gofr", b, 0644)

	logger.Cute(fmt.Sprintf("original: %dB\ncompressed: %dB", len(data), len(b)))
}

func decompress(filename string) {
	data, _ := os.ReadFile(filename)
	b := huffman.Decode(data)
	b = mtf.Decode(b)
	logger.Cute(string(b))
}
