package huffmantree

type NodeValue struct {
	Chars string
	Freq  int
}

type Node struct {
	Value NodeValue
	Left  *Node
	Right *Node
}

func NewValue(chars string, freq int) NodeValue {
	return NodeValue{
		Chars: chars,
		Freq:  freq,
	}
}

func NewNode(value NodeValue) *Node {
	return &Node{
		Value: value,
		Left:  nil,
		Right: nil,
	}
}

func (n *Node) AddLeft(node *Node) {
	n.Left = node
}

func (n *Node) AddRight(node *Node) {
	n.Right = node
}

func (n *Node) IsLast() bool {
	if n.Left == nil && n.Right == nil {
		return true
	}

	return false
}
