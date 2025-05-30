package compression

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/internal/algorithm/bwt"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/internal/algorithm/mtf"
	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/models"
)

func Decompress(filename string) {
	by, _ := os.ReadFile(filename)
	fmt.Printf("by: %v\n", by)
	meta, data := extractData(by)
	b := huffman.Decode(data, meta)
	b = mtf.Decode(b)
	b = bwt.Decode(b, int(meta.BwtIndex))
	logger.Cute(string(b))
}

func extractData(b []byte) (metadata models.FileMetadata, data *bitbuffer.BitBuffer) {
	byteCount := 0
	logger.Cute(b)

	metadata.Padding = b[byteCount]
	byteCount += 1
	fmt.Printf("metadata.Padding: %v\n", metadata.Padding)

	metadata.TreeSize = binary.BigEndian.Uint16(b[byteCount : byteCount+2])
	byteCount += 2
	fmt.Printf("metadata.TreeSize: %v\n", metadata.TreeSize)

	metadata.TreePadding = b[byteCount]
	byteCount += 1
	fmt.Printf("metadata.TreePadding: %v\n", metadata.TreePadding)

	treeBytes := b[byteCount : byteCount+int(metadata.TreeSize)]
	metadata.Tree = bitbuffer.New(treeBytes)
	fmt.Printf("metadata.Tree: %v\n", metadata.Tree)

	byteCount += len(metadata.Tree.Bytes)

	metadata.BwtIndex = binary.BigEndian.Uint32(b[byteCount : byteCount+4])
	logger.Cute(b[byteCount : byteCount+4])
	byteCount += 4
	logger.Cute(metadata.BwtIndex)

	data = bitbuffer.New(b[byteCount:])

	logger.Cute(fmt.Sprintf("padding: %v", int(metadata.Padding)))
	logger.Cute(fmt.Sprintf("tree padding: %v", int(metadata.TreePadding)))
	logger.Cute(fmt.Sprintf("tree size: %v", int(metadata.TreeSize)))
	logger.Cute(fmt.Sprintf("tree: %v", metadata.Tree.Bytes))

	return metadata, data
}
