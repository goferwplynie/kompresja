package main

import (
	"time"

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

}

func decompress(filename string) {

}
