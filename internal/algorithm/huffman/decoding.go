package huffman

import (
	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/bits/bitreader"
	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
	"github.com/goferwplynie/kompresja/models"
)

func Decode(data *bitbuffer.BitBuffer, metadata models.FileMetadata) []byte {
	var bytes []byte
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
	return bytes
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
