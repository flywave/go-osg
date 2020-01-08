package model

type TraversalMode uint32
type VisitorType uint32

const (
	TRAVERSENONE           TraversalMode = 0
	TRAVERSEPARENTS        TraversalMode = 1
	TRAVERSEALLCHILDREN    TraversalMode = 2
	TRAVERSEACTIVECHILDREN TraversalMode = 3

	NODEVISITOR            VisitorType = 0
	UPDATEVISITOR          VisitorType = 1
	EVENTVISITOR           VisitorType = 2
	COLLECTOCCLUDERVISITOR VisitorType = 3
	CULLVISITOR            VisitorType = 4
	INTERSECTIONVISITOR    VisitorType = 5

	UNINITIALIZEDFRAMENUMBER uint32 = 0xffffffff
)

type NodeVisitor struct {
	VisitorType      VisitorType
	TraversalNumber  uint32
	TraversalMode    TraversalMode
	TraversalMask    uint32
	NodeMaskOverride uint32
	Npath            NodePath
	Geos             []*Geometry
}

func NewNodeVisitor() *NodeVisitor {
	return &NodeVisitor{VisitorType: NODEVISITOR, TraversalMode: TRAVERSEALLCHILDREN, NodeMaskOverride: 0x0, TraversalMask: 0xffffffff, TraversalNumber: UNINITIALIZEDFRAMENUMBER}
}

func (v *NodeVisitor) PushOntoNodePath(n NodeInterface) {
	if v.TraversalMode != TRAVERSEPARENTS {
		v.Npath = append(v.Npath, n)
	} else {
		t := []interface{}{n}
		v.Npath = append(t, v.Npath...)
	}
}

func (v *NodeVisitor) PopFromNodePath() {
	if v.TraversalMode != TRAVERSEPARENTS {
		v.Npath = v.Npath[:len(v.Npath)-1]
	} else {
		v.Npath = v.Npath[1 : len(v.Npath)-1]
	}
}

func (v *NodeVisitor) Traverse(node NodeInterface) {
	if v.TraversalMode == TRAVERSEPARENTS {
		node.Ascend(v)
	} else {
		node.Traverse(v)
	}
}

func (v *NodeVisitor) Apply(val NodeInterface) {
	switch node := val.(type) {
	case *Node:
		v.Traverse(node)
	case *PagedLod:
		v.Traverse(node)
	case *Geode:
		v.Traverse(node)
	case *Group:
		v.Traverse(node)
	}
}

func (v *NodeVisitor) ValidNodeMask(node NodeInterface) bool {
	msk := node.GetNodeMask()
	return (v.TraversalMask & (v.NodeMaskOverride | *msk)) != 0
}
