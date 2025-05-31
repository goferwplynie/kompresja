package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/goferwplynie/kompresja/archive"
	"github.com/goferwplynie/kompresja/compression"
	"github.com/goferwplynie/kompresja/models"

	//"github.com/goferwplynie/kompresja/compression"
	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/workerpool"
)

const partSize = 100 * 1024 * 1024

func main() {
	start := time.Now()
	//compress("../kompresja2", "../kompresja2.gofr")
	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	//compression.Decompress("test.gofr")
	decompress("../kompresja2.gofr")
	logger.Cute(time.Since(start))
}

func compress(path string, dest string) {
	os.Create(dest)

	wp := workerpool.New(20)
	go wp.Run(dest)

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
	wp.Close()

}

func decompress(filename string) {
	var err error
	compressedF, err := os.Open(filename)

	if err != nil {
		logger.Error(err)
	}
	compressedFInfo, err := compressedF.Stat()
	if err != nil {
		logger.Error(err)
	}

	defer compressedF.Close()
	files := make(map[string][]models.ArchiveMetaData)
	offset := uint(0)
	for true {
		fpLenBytes := make([]byte, 2)
		_, err = compressedF.ReadAt(fpLenBytes, int64(offset))
		if err != nil && err != io.EOF {
			logger.Error(err)
			break
		}
		fpLen := binary.BigEndian.Uint16(fpLenBytes)

		offset += 2

		fpBytes := make([]byte, fpLen)
		_, err = compressedF.ReadAt(fpBytes, int64(offset))
		if err != nil && err != io.EOF {
			logger.Error(err)
			break
		}
		fp := string(fpBytes)

		offset += uint(len(fpBytes))

		fileSizeBytes := make([]byte, 8)
		_, err = compressedF.ReadAt(fileSizeBytes, int64(offset))
		if err != nil && err != io.EOF {
			logger.Error(err)
			break
		}
		fileSize := binary.BigEndian.Uint64(fileSizeBytes)

		offset += 8

		partBytes := make([]byte, 4)
		_, err = compressedF.ReadAt(partBytes, int64(offset))
		if err != nil && err != io.EOF {
			logger.Error(err)
			break
		}
		part := binary.BigEndian.Uint32(partBytes)

		offset += 4

		files[fp] = append(files[fp], models.ArchiveMetaData{
			PathLen:  fpLen,
			Path:     fp,
			FileSize: fileSize,
			Part:     part,
			Start:    offset,
		})
		offset += uint(fileSize)
		fmt.Printf("offset: %v\n", offset)
		if offset >= uint(compressedFInfo.Size()) {
			break
		}
	}
	logger.Cute(files)

	for path, parts := range files {
		dir := filepath.Dir(path)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error(err)
		}

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Error(err)
		}
		for _, p := range parts {
			compressed := make([]byte, p.FileSize)
			compressedF.ReadAt(compressed, int64(p.Start))
			b := compression.Decompress(compressed)
			if _, err := f.Write(b); err != nil {
				logger.Error(err)
			}
		}
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}
}
