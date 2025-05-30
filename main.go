package main

import (
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/goferwplynie/kompresja/archive"
	//"github.com/goferwplynie/kompresja/compression"
	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/workerpool"
)

const partSize = 100 * 1024 * 1024

func main() {
	start := time.Now()
	compress("../kompresja", "../kompresja.gofr")
	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	//compression.Decompress("test.gofr")
	logger.Cute(time.Since(start))
}

func compress(path string, dest string) {
	os.Create(dest)

	wp := workerpool.New(20)
	wp.Run(dest)

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, _ := d.Info()
			if info.Size() > partSize {
				parts := math.Ceil(float64(info.Size()) / partSize)
				for i := range int(parts) {
					wp.AddTask(archive.NewFile(path, i))
				}
			}
			wp.AddTask(archive.NewFile(path, 0))

		}
		return nil
	})

}
