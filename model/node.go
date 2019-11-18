package model

import "reflect"

type NodePath []interface{}

const (
	NODET string = "osg::Node"
)

type ComputeBoundingSphereCallback struct {
	Object
}
type NodeInterface interface {
	IsNode() bool
}

type Node struct {
	Object
	CullingActive bool
	NodeMask      uint32
	Dscriptions   []string
	InitialBound  Sphere3f
	States        *StateSet
	Parents       []*Group

	Callback       *ComputeBoundingSphereCallback
	UpdateCallback *Callback
	EventCallback  *Callback
	CullCallback   *Callback
}

func NewNode() Node {
	obj := NewObject()
	obj.Type = NODET
	return Node{Object: obj, NodeMask: 0xffffffff}
}

func (n *Node) Accept(nv *NodeVisitor) {
	if nv.ValidNodeMask(n) {
		nv.PushOntoNodePath(n)
		nv.Apply(n)
		nv.PopFromNodePath(n)
	}
}

func (n *Node) Ascend(nv *NodeVisitor) {

}

func (n *Node) Traverse(nv *NodeVisitor) {

}

func (n *Node) IsNode() bool {
	return true
}

func IsBaseOfNode(obj interface{}) bool {
	if obj == nil {
		return false
	}
	no := NewNode()
	baset := reflect.TypeOf(no)
	t := reflect.TypeOf(obj)
	return t.Implements(baset)
}
