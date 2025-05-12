package huffmantree

type NodeValue struct {
	Bytes []byte
	Freq  int
}

type Node struct {
	Value NodeValue
	Left  *Node
	Right *Node
}

func NewValue(bytes []byte, freq ...int) NodeValue {
	var frequency int
	if len(freq) > 0 {
		frequency = freq[0]
	}
	return NodeValue{
		Bytes: bytes,
		Freq:  frequency,
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
