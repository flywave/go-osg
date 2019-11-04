package model

type TraversalMode uint32
type VisitorType uint32

const (
	TRAVERSE_NONE            TraversalMode = 0
	TRAVERSE_PARENTS         TraversalMode = 1
	TRAVERSE_ALL_CHILDREN    TraversalMode = 2
	TRAVERSE_ACTIVE_CHILDREN TraversalMode = 3

	NODE_VISITOR             VisitorType = 0
	UPDATE_VISITOR           VisitorType = 1
	EVENT_VISITOR            VisitorType = 2
	COLLECT_OCCLUDER_VISITOR VisitorType = 3
	CULL_VISITOR             VisitorType = 4
	INTERSECTION_VISITOR     VisitorType = 5

	UNINITIALIZED_FRAME_NUMBER uint32 = 0xffffffff
)

type NodeVisitor struct {
	VisitorType      VisitorType
	TraversalNumber  uint32
	TraversalMode    TraversalMode
	TraversalMask    NodeMask
	NodeMaskOverride NodeMask
	Npath            NodePath
}

func NewNodeVisitor() NodeVisitor {
	return NodeVisitor{VisitorType: NODE_VISITOR, TraversalMode: TRAVERSE_NONE, NodeMaskOverride: 0x0, TraversalMask: 0xffffffff, TraversalNumber: UNINITIALIZED_FRAME_NUMBER}
}

func (v *NodeVisitor) PushOntoNodePath(n *Node) {
	if v.TraversalMode != TRAVERSE_PARENTS {
		v.Npath = append(v.Npath, n)
	} else {
		t := []*Node{n}
		v.Npath = append(t, v.Npath...)
	}
}

func (v *NodeVisitor) PopFromNodePath(n *Node) {
	if v.TraversalMode != TRAVERSE_PARENTS {
		v.Npath = v.Npath[:len(v.Npath)-1]
	} else {
		v.Npath = v.Npath[1 : len(v.Npath)-1]
	}
}

func (v *NodeVisitor) Apply(val interface{}) {
	switch node := val.(type) {
	case *Node:
	case *Lod:
	case *PagedLod:
	case *Group:
	}
}

func (v *NodeVisitor) ValidNodeMask(node *Node) bool {
	return (v.TraversalMask &
		(v.NodeMaskOverride | node.NodeMask)) != 0
}
