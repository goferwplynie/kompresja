package workerpool

import (
	"fmt"
	"os"

	"github.com/goferwplynie/kompresja/archive"
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
	for i := range wp.maxWorker {
		go func(workerID int) {
			for file := range wp.queuedTaskC {
				b := compression.Compress(file)
				logger.Cute(fmt.Sprintf("file: %v compressed size: %v", file.Path, len(b)))
				f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					logger.Error(err)
				}
				if _, err := f.Write(b); err != nil {
					logger.Error(err)
				}
				if err := f.Close(); err != nil {
					logger.Error(err)
				}

			}
		}(i + 1)
	}
}

func (wp *workerPool) AddTask(file archive.File) {
	wp.queuedTaskC <- file
}
