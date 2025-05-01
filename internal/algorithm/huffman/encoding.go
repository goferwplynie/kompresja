package huffman

import (
	"sort"

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

func Encode(str string) {
	nodes := MakeNodes(str)
	tree := buildTree(nodes)
	codes := make(map[string]string)
	makeCodes(codes, tree)
	logger.Cute(codes)
	printTree(tree)
}

func MakeNodes(str string) []*huffmantree.Node {
	var nodes = make([]*huffmantree.Node, 0)
	var frequencies = make(map[string]int)

	logger.Cute("calculating frequencies")
	for _, r := range str {
		frequencies[string(r)]++
	}

	logger.Cute("creating nodes")
	for key, value := range frequencies {
		nodes = append(nodes, huffmantree.NewNode(
			huffmantree.NewValue(key, value)))
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
		Chars: value1.Chars + value2.Chars,
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

func makeCodes(codes map[string]string,
	node *huffmantree.Node, currentCode ...string) {
	code := ""
	if len(currentCode) > 0 {
		code = currentCode[0]
		logger.Cute("making codes")
	}
	if node.IsLast() {
		codes[node.Value.Chars] = code
		return
	}

	makeCodes(codes, node.Left, code+"0")
	makeCodes(codes, node.Right, code+"1")
}
