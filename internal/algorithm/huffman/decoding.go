package huffman

import (
	"encoding/binary"
	"fmt"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/models"
)

func Decode(b []byte) []byte {
	var bytes []byte
	var codes = make(map[string]byte)
	metadata, _ := extractData(b)
	tree := rebuildTree(metadata.Tree)
	printTree(tree)

	logger.Cute(codes)

	return bytes
}

func extractData(b []byte) (metadata models.FileMetadata, data *bitbuffer.BitBuffer) {
	byteCount := 0

	metadata.Padding = b[byteCount]
	byteCount += 1

	metadata.TreeSize = binary.BigEndian.Uint16(b[byteCount : byteCount+1])
	byteCount += 2

	metadata.TreePadding = b[byteCount]
	byteCount += 1

	logger.Warn(metadata.TreePadding)

	treeBytes := b[byteCount : byteCount+int(metadata.TreeSize)]
	metadata.Tree = bitbuffer.New(treeBytes)

	byteCount += len(metadata.Tree.Bytes)

	data = bitbuffer.New(b[byteCount:])

	logger.Cute(fmt.Sprintf("padding: %v", int(metadata.Padding)))

	return metadata, data
}

func rebuildTree(bb *bitbuffer.BitBuffer) *huffmantree.Node {
	var node huffmantree.Node

	return &node
}
