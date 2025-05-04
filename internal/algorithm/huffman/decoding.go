package huffman

import (
	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
)

func ExtractData(bitstring string) (padding, treeSize, treePadding int,
	tree *huffmantree.Node, data string) {
	paddingInBits := bitstring[0:8]

	return padding, treeSize, treePadding, tree, data
}
