package workerpool

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/goferwplynie/kompresja/archive"
	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/compression"
	"github.com/goferwplynie/kompresja/logger"
)

type WorkerPool interface {
	Run()
	AddTask(file archive.File)
}

type workerPool struct {
	maxWorker   int
	queuedTaskC chan archive.File
}

func New(maxWorkers int) *workerPool {
	return &workerPool{
		maxWorker:   maxWorkers,
		queuedTaskC: make(chan archive.File),
	}
}

func (wp *workerPool) Run(dest string) {
	var wg sync.WaitGroup
	for i := range wp.maxWorker {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for file := range wp.queuedTaskC {
				b := compression.Compress(file)
				bb := bitbuffer.New(b)
				addMeta(bb, file)
				logger.Cute(fmt.Sprintf("file: %v compressed size: %v", file.Path, len(bb.Bytes)))
				f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					logger.Error(err)
				}
				if _, err := f.Write(bb.Bytes); err != nil {
					logger.Error(err)
				}
				if err := f.Close(); err != nil {
					logger.Error(err)
				}
			}
		}(i + 1)
	}
	wg.Wait()
}

func (wp *workerPool) AddTask(file archive.File) {
	wp.queuedTaskC <- file
}

func (wp *workerPool) Close() {
	close(wp.queuedTaskC)
}

func addMeta(bb *bitbuffer.BitBuffer, file archive.File) {
	meta := bitbuffer.New()

	fp := []byte(file.Path)

	fpLen := binary.BigEndian.AppendUint16(make([]byte, 0, 2), uint16(len(fp)))
	size := binary.BigEndian.AppendUint64(make([]byte, 0, 8), uint64(len(bb.Bytes)))
	part := binary.BigEndian.AppendUint32(make([]byte, 0, 4), uint32(file.Part))

	meta.AddBytes(fpLen)
	meta.AddBytes(fp)
	meta.AddBytes(size)
	meta.AddBytes(part)

	bb.MergeLeft(meta)
}
