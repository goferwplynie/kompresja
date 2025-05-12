package huffman

import (
	"fmt"

	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
	"github.com/goferwplynie/kompresja/logger"
	"github.com/goferwplynie/kompresja/models"
)

func Decode(b []byte) []string {
	var chars []string
	var codes = make(map[string]string)
	metadata, data := extractData(b)
	tree, _ := rebuildTree(metadata.Tree)
	printTree(tree)

	makeCodesReversed(codes, tree)

	logger.Cute(codes)
	currentCode := ""
	for _, v := range data {
		currentCode += string(v)
		val, exists := codes[currentCode]
		if exists {
			chars = append(chars, val)
			currentCode = ""
		}
	}

	return chars
}

func makeCodesReversed(codes map[string]string,
	node *huffmantree.Node, currentCode ...string) {
	code := ""
	if len(currentCode) > 0 {
		code = currentCode[0]
	} else {
		logger.Cute("making codes")

	}

	if node.IsLast() {
		codes[code] = node.Value.Chars
		return
	}

	makeCodesReversed(codes, node.Left, code+"0")
	makeCodesReversed(codes, node.Right, code+"1")
}

func extractData(b []byte) (metadata models.FileMetadata, data string) {
	byteCount := 0

	metadata.Padding = b[byteCount]
	byteCount += 1

	metadata.TreeSize = uint32(b[byteCount])
	byteCount += 1

	metadata.TreePadding = b[byteCount]
	byteCount += 1

	logger.Warn(metadata.TreePadding)

	metadata.Tree = bitutils.BytesToBitString(b[byteCount : byteCount+int(metadata.TreeSize)])
	metadata.Tree = metadata.Tree[:len(metadata.Tree)-int(metadata.TreePadding)]
	byteCount += int(metadata.TreeSize)

	data = bitutils.BytesToBitString(b[byteCount:])
	data = data[:len(data)-int(metadata.Padding)]

	logger.Cute(fmt.Sprintf("padding: %v", int(metadata.Padding)))

	return metadata, data
}

func rebuildTree(bitstring string) (*huffmantree.Node, string) {
	var node huffmantree.Node
	if bitstring[0] == '1' {
		charBits := bitstring[1:9]
		charByte := bitutils.BitStringToBytes(charBits)

		bitstring = bitstring[9:]
		node.Value = huffmantree.NewValue(string(charByte))
		return &node, bitstring
	} else {
		bitstring = bitstring[1:]
		left, bitstring := rebuildTree(bitstring)
		right, bitstring := rebuildTree(bitstring)
		node.Value = huffmantree.NewValue(left.Value.Chars + right.Value.Chars)
		node.AddLeft(left)
		node.AddRight(right)

		return &node, bitstring
	}

}
