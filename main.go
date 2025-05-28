package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/internal/algorithm/bwt"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/internal/algorithm/mtf"
	"github.com/goferwplynie/kompresja/logger"
)

func main() {
	start := time.Now()
	compress("test.json")
	logger.Cute(time.Since(start))

	logger.Cute("decompression")
	start = time.Now()
	decompress("test.gofr")
	logger.Cute(time.Since(start))
}

func compress(filename string) {
	data, _ := os.ReadFile(filename)

	b, bwtIndex := bwt.Encode(data)
	b = mtf.Encode(b)
	bb, treeBuffer := huffman.Encode(b)
	addMetaData(bb, treeBuffer)
	os.WriteFile(filename+".gofr", b, 0644)

	fmt.Printf("bwtIndex: %v\n", bwtIndex)

	logger.Cute(fmt.Sprintf("original: %dB\ncompressed: %dB", len(data), len(b)))
}

func decompress(filename string) {
	data, _ := os.ReadFile(filename)
	b := huffman.Decode(data)
	b = mtf.Decode(b)
	b = bwt.Decode(b, 340)
	logger.Cute(string(b))
}

func addMetaData(bitBuffer *bitbuffer.BitBuffer, treeBuffer *bitbuffer.BitBuffer) {
	logger.Cute("adding metadata")
	meta := bitbuffer.New()

	padding := (8 - bitBuffer.CurrentBit) % 8
	bitBuffer.Finalize()

	logger.Cute(len(treeBuffer.Bytes))

	treePadding := (8 - treeBuffer.CurrentBit) % 8
	treeBuffer.Finalize()

	treeSize := len(treeBuffer.Bytes)
	treeSizeB := binary.BigEndian.AppendUint16(make([]byte, 0, 2), uint16(treeSize))
	logger.Cute(treeSizeB)

	meta.AddByte(byte(padding))
	meta.AddBytes(treeSizeB)
	meta.AddByte(byte(treePadding))

	meta.MergeRight(treeBuffer)
	bitBuffer.MergeLeft(meta)

	//[padding][tree size][tree padding][tree][data]

	logger.Warn(fmt.Sprintf("padding: %v", padding))
	logger.Warn(fmt.Sprintf("tree padding: %v", treePadding))
	logger.Warn(fmt.Sprintf("tree size: %vB", treeSize))
	logger.Warn(fmt.Sprintf("metadata size: %vB", len(meta.Bytes)))
	logger.Log(fmt.Sprintf("serialized tree: %v", treeBuffer.Bytes))
}
