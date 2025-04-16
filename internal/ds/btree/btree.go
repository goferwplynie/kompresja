package btree

type Node struct {
	Value string
	Freq  int
	Left  *Node
	Right *Node
}

func New(value string) *Node {
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
