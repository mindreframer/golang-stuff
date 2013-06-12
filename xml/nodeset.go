package xml

import "unsafe"

type NodeSet struct {
	Document
	Nodes []Node
	valid bool
}

func NewNodeSet(document Document, nodes interface{}) (set *NodeSet) {
	set = &NodeSet{valid: true}
	set.Document = document

	switch t := nodes.(type) {
	case []Node:
		set.Nodes = t
	case []unsafe.Pointer:
		if num := len(t); num > 0 {
			set.Nodes = make([]Node, num)
			for i, p := range t {
				set.Nodes[i] = NewNode(p, document)
			}
		}
	default:
		//unexpected param type
		//ignore the data
	}
	return
}

func (set *NodeSet) Length() int {
	return len(set.Nodes)
}

func (set *NodeSet) Remove() {
	if set.valid {
		for _, node := range set.Nodes {
			node.Remove()
		}
		set.valid = false
	}
}
