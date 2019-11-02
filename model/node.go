package model

// import (_"github.com/ungerik/go3d/vec3")

type Sphere3f struct {
	// Center vec3.T
	Radius float32
}

const (
	NodeType string = "osg::Node"
)

type ComputeBoundingSphereCallback struct {
	Object
}

type Node struct {
	Object
	CullingActive bool
	NodeMask      uint
	Dscriptions   []string
	InitialBound  Sphere3f
	States        *StateSet
	Parents       []*group

	Callback       *ComputeBoundingSphereCallback
	UpdateCallback *Callback
	EventCallback  *Callback
	CullCallback   *Callback
}

func NewNode() Node {
	obj := NewObject()
	obj.ObjectType = NodeType
	return Node{Object: obj}
}

// func Accept(nv *NodeVisitor,){
//   if (nv.ValidNodeMask(n)) {
//     nv.push_onto_node_path(n);
//     nv.apply(n);
//     nv.pop_from_node_path();
//   }
// }
// func (n *node) Ascend(nv *NodeVisitor){
//  }

// func (n *node) Traverse(nv *NodeVisitor){

// }
