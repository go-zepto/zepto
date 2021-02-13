package where

type Node struct {
	Parent *Node
	Nodes  []*Node
	Type   WhereType
	Key    string
	Value  interface{}
	IsLast bool
}

func (n *Node) LastChild() *Node {
	l := len(n.Nodes)
	if l == 0 {
		return nil
	}
	return n.Nodes[l-1]
}

func (n *Node) AddChild(child *Node) {
	n.Nodes = append(n.Nodes, child)
}
