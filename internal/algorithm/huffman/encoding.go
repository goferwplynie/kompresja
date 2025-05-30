package huffman

import (
	"sort"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
	"github.com/goferwplynie/kompresja/logger"
)

func printTree(tree *huffmantree.Node) {
	if tree == nil {
		return
	}
	logger.Cute(tree.Value)
	printTree(tree.Left)
	printTree(tree.Right)
}

func Encode(bytes []byte) (*bitbuffer.BitBuffer, *bitbuffer.BitBuffer) {
	bb := bitbuffer.New()

	nodes := MakeNodes(bytes)
	tree := buildTree(nodes)

	codes := make(map[byte][]bool)

	makeCodes(codes, tree)
	//printTree(tree)

	for _, b := range bytes {
		bb.AddBits(codes[b])
	}

	treeBuffer := bitbuffer.New()
	treeToBits(tree, treeBuffer)

	return bb, treeBuffer
}

func treeToBits(node *huffmantree.Node, bitBuffer *bitbuffer.BitBuffer) {
	if node.IsLast() {
		Byte := byte(node.Value.Bytes[0])
		bitBuffer.AddBit(true)
		bitBuffer.AddByte(Byte)
	} else {
		bitBuffer.AddBit(false)
		treeToBits(node.Left, bitBuffer)
		treeToBits(node.Right, bitBuffer)
	}
}

func MakeNodes(bytes []byte) []*huffmantree.Node {
	var nodes = make([]*huffmantree.Node, 0)
	var frequencies = make(map[byte]int)

	for _, b := range bytes {
		frequencies[b]++
	}

	for key, value := range frequencies {
		nodes = append(nodes, huffmantree.NewNode(
			huffmantree.NewValue([]byte{key}, value)))
	}

	return nodes
}

func sortNodes(nodes []*huffmantree.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Value.Freq < nodes[j].Value.Freq
	})
}

func mergeValues(value1, value2 huffmantree.NodeValue) huffmantree.NodeValue {
	return huffmantree.NodeValue{
		Bytes: append(value1.Bytes, value2.Bytes...),
		Freq:  value1.Freq + value2.Freq,
	}
}

func buildTree(nodes []*huffmantree.Node) *huffmantree.Node {
	for len(nodes) > 1 {
		sortNodes(nodes)

		l := nodes[0]
		nodes = nodes[1:]
		r := nodes[0]
		nodes = nodes[1:]

		newValue := mergeValues(l.Value, r.Value)
		newNode := huffmantree.NewNode(newValue)

		newNode.AddLeft(l)
		newNode.AddRight(r)

		nodes = append(nodes, newNode)
	}

	return nodes[0]
}

func makeCodes(codes map[byte][]bool,
	node *huffmantree.Node, currentCode ...[]bool) {
	code := make([]bool, 0)
	if len(currentCode) > 0 {
		code = currentCode[0]
	}

	if node.IsLast() {
		codes[node.Value.Bytes[0]] = code
		return
	}

	leftCode := make([]bool, len(code)+1)
	copy(leftCode, code)
	leftCode[len(code)] = false

	rightCode := make([]bool, len(code)+1)
	copy(rightCode, code)
	rightCode[len(code)] = true

	makeCodes(codes, node.Left, leftCode)
	makeCodes(codes, node.Right, rightCode)

}
