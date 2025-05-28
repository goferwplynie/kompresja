package huffman

import (
	"encoding/binary"
	"fmt"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/bits/bitreader"
	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/models"
)

func Decode(b []byte) []byte {
	var bytes []byte
	metadata, data := extractData(b)
	root := rebuildTree(bitreader.New(metadata.Tree.Bytes, len(metadata.Tree.Bytes)*8))
	//printTree(root)

	br := bitreader.New(data.Bytes, len(data.Bytes)*8-int(metadata.Padding))
	currentNode := root

	for {
		bit, err := br.Next()
		if err != nil {
			break
		}
		if bit {
			currentNode = currentNode.Right
		} else {
			currentNode = currentNode.Left
		}

		if currentNode.IsLast() {
			bytes = append(bytes, currentNode.Value.Bytes[0])
			currentNode = root
		}

	}
	fmt.Printf("bytes: %v\n", bytes)

	return bytes
}

func extractData(b []byte) (metadata models.FileMetadata, data *bitbuffer.BitBuffer) {
	byteCount := 0
	logger.Cute(b)

	metadata.Padding = b[byteCount]
	byteCount += 1

	metadata.TreeSize = binary.BigEndian.Uint16(b[byteCount : byteCount+2])
	byteCount += 2

	metadata.TreePadding = b[byteCount]
	byteCount += 1

	treeBytes := b[byteCount : byteCount+int(metadata.TreeSize)]
	metadata.Tree = bitbuffer.New(treeBytes)

	byteCount += len(metadata.Tree.Bytes)

	data = bitbuffer.New(b[byteCount:])

	logger.Cute(fmt.Sprintf("padding: %v", int(metadata.Padding)))
	logger.Cute(fmt.Sprintf("tree padding: %v", int(metadata.TreePadding)))
	logger.Cute(fmt.Sprintf("tree size: %v", int(metadata.TreeSize)))
	logger.Cute(fmt.Sprintf("tree: %v", metadata.Tree.Bytes))

	return metadata, data
}

func rebuildTree(br *bitreader.BitReader) *huffmantree.Node {
	bit, _ := br.Next()

	if bit {
		b, _ := br.ReadByte()
		return huffmantree.NewNode(huffmantree.NewValue([]byte{b}, 0))
	}

	left := rebuildTree(br)
	right := rebuildTree(br)

	node := huffmantree.NewNode(huffmantree.NewValue([]byte{}))
	node.AddLeft(left)
	node.AddRight(right)

	return node
}
