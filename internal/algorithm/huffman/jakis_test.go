package huffman

import (
	"testing"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/bits/bitreader"
	huffmantree "github.com/goferwplynie/kompresja/internal/ds/huffmanTree"
)

func TestTreeIntegrity(t *testing.T) {
	bytes := []byte("laldjwidajwjxkdkwa")
	bb := bitbuffer.New()

	nodes := MakeNodes(bytes)
	tree1 := buildTree(nodes)

	codes := make(map[byte][]bool)

	makeCodes(codes, tree1)
	printTree(tree1)

	for _, b := range bytes {
		bb.AddBits(codes[b])
	}

	addMetaData(bb, tree1)

	metadata, _ := extractData(bb.Bytes)
	tree2 := rebuildTree(bitreader.New(metadata.Tree.Bytes))
	if !TreesEqual(tree1, tree2) {
		t.Fatal("nie to samo drzewko :c")
	}

}

func TreesEqual(a, b *huffmantree.Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	aIsLeaf := a.IsLast()
	bIsLeaf := b.IsLast()

	if aIsLeaf != bIsLeaf {
		return false
	}

	if aIsLeaf && bIsLeaf {
		if len(a.Value.Bytes) != len(b.Value.Bytes) {
			return false
		}
		for i := range a.Value.Bytes {
			if a.Value.Bytes[i] != b.Value.Bytes[i] {
				return false
			}
		}
		return true
	}

	return TreesEqual(a.Left, b.Left) && TreesEqual(a.Right, b.Right)
}
