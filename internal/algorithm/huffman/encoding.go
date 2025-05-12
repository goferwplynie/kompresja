package huffman

import (
	"fmt"
	"sort"

	"encoding/binary"

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

func Encode(bytes []byte) []byte {
	bb := bitbuffer.New()

	nodes := MakeNodes(bytes)
	tree := buildTree(nodes)

	codes := make(map[byte][]bool)

	makeCodes(codes, tree)
	logger.Cute(codes)
	printTree(tree)

	for _, b := range bytes {
		bb.AddBits(codes[b])
	}

	addMetaData(bb, tree)

	return bb.Bytes
}

func addMetaData(bitBuffer *bitbuffer.BitBuffer, tree *huffmantree.Node) {
	logger.Cute("adding metadata")
	meta := bitbuffer.New()

	padding := 8 - bitBuffer.CurrentBit
	bitBuffer.Finalize()

	treeBuffer := bitbuffer.New()

	treeToBits(tree, treeBuffer)

	treePadding := 8 - treeBuffer.CurrentBit
	treeBuffer.Finalize()

	treeSize := len(treeBuffer.Bytes)
	treeSizeB := binary.BigEndian.AppendUint16(make([]byte, 2), uint16(treeSize))

	meta.AddByte(byte(padding))
	meta.AddBytes(treeSizeB)
	meta.AddByte(byte(treePadding))

	meta.MergeRight(treeBuffer)
	bitBuffer.MergeLeft(meta)

	logger.Warn(fmt.Sprintf("padding: %v", padding))
	logger.Warn(fmt.Sprintf("tree padding: %v", treePadding))
	logger.Warn(fmt.Sprintf("tree size: %vB", treeSize))
	logger.Log(fmt.Sprintf("serialized tree: %v", treeBuffer.Bytes))
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

	logger.Cute("calculating frequencies")
	for _, b := range bytes {
		frequencies[b]++
	}

	logger.Cute("creating nodes")
	for key, value := range frequencies {
		nodes = append(nodes, huffmantree.NewNode(
			huffmantree.NewValue([]byte{key}, value)))
	}

	return nodes
}

func sortNodes(nodes []*huffmantree.Node) {
	logger.Log("sorting nodes")
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
	logger.Cute("making tree")
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
	} else {
		logger.Cute("making codes")

	}
	if node.IsLast() {
		codes[node.Value.Bytes[0]] = code
		return
	}

	makeCodes(codes, node.Left, append(code, false))
	makeCodes(codes, node.Right, append(code, true))
}
