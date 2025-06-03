package main

import (
	"encoding/binary"
	"flag"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/goferwplynie/kompresja/archive"
	"github.com/goferwplynie/kompresja/compression"
	"github.com/goferwplynie/kompresja/models"

	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/workerpool"
)

const partSize = 100 * 1024 * 1024

var regexFilter *regexp.Regexp

func main() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "compress":
			start := time.Now()
			compress(args)
			logger.Cute(time.Since(start))
		case "decompress":
			start := time.Now()
			decompress(args)
			logger.Cute(time.Since(start))
		}
	}
}

func compress(args []string) {
	if len(args) < 4 {
		logger.Error("not enough arguments provided")
		logger.Cute("Expected: go run main.go compress <path to compress> <compressed name>")
		return
	}
	path := args[2]
	path, _ = filepath.Abs(path)
	dest := args[3]

	flagSet := flag.NewFlagSet("compress", flag.ExitOnError)
	ignoreHidden := flagSet.Bool("ih", false, "ignore hidden")
	regex := flagSet.String("i", "", "ignore by regex")

	flagSet.Parse(args[4:])

	if *regex != "" {
		var err error
		regexFilter, err = regexp.Compile(*regex)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	os.Create(dest)

	wp := workerpool.New(25)
	go wp.Run(dest)

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error(err)
			return err
		}
		if *ignoreHidden && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if matchRegex(d.Name()) {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			info, _ := d.Info()
			logger.Cute(path)
			if info.Size() > partSize {
				parts := math.Ceil(float64(info.Size()) / partSize)
				for i := range int(parts) {
					wp.AddTask(archive.NewFile(path, i))
				}
			} else {
				wp.AddTask(archive.NewFile(path, 0))
			}

		}
		return nil
	})
	wp.Close()
	wp.Wait()

}

func matchRegex(name string) bool {
	if regexFilter != nil {
		return (*regexFilter).MatchString(name)
	}
	return true
}

func decompress(args []string) {
	if len(args) < 2 {
		logger.Error("not enough arguments provided")
		logger.Cute("Expected: go run main.go decompress <compressed file>")
		return
	}

	filename := args[2]

	var err error
	compressedF, err := os.Open(filename)

	if err != nil {
		logger.Error(err)
		return
	}
	compressedFInfo, err := compressedF.Stat()
	if err != nil {
		logger.Error(err)
		return
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
			return
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
