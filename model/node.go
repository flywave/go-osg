package model

type NodePath []interface{}

const (
	NODET string = "osg::Node"
)

type ComputeBoundingSphereCallback struct {
	Object
}

type NodeInterface interface {
	IsNode() bool
	GetCullingActive() *bool
	SetCullingActive(bool)
	GetNodeMask() *uint32
	SetNodeMask(uint32)
	GetDscriptions() []string
	SetDscriptions([]string)
	GetInitialBound() *Sphere3f
	SetInitialBound(Sphere3f)
	GetStates() *StateSet
	SetStates(*StateSet)
	GetParents() []*Group
	SetParents([]*Group)
	GetSphereCallback() *ComputeBoundingSphereCallback
	SetSphereCallback(*ComputeBoundingSphereCallback)

	Ascend(nv *NodeVisitor)
	Traverse(nv *NodeVisitor)

	GetUpdateCallback() *Callback
	SetUpdateCallback(*Callback)
	GetEventCallback() *Callback
	SetEventCallback(*Callback)
	GetCullCallback() *Callback
	SetCullCallback(*Callback)
	Accept(*NodeVisitor)
}

type Node struct {
	Object
	CullingActive bool
	NodeMask      uint32
	Dscriptions   []string
	InitialBound  Sphere3f
	States        *StateSet
	Parents       []*Group

	SphereCallback *ComputeBoundingSphereCallback
	UpdateCallback *Callback
	EventCallback  *Callback
	CullCallback   *Callback
}

func (n *Node) GetCullingActive() *bool {
	return &n.CullingActive
}
func (n *Node) SetCullingActive(ca bool) {
	n.CullingActive = ca
}
func (n *Node) GetNodeMask() *uint32 {
	return &n.NodeMask
}
func (n *Node) SetNodeMask(v uint32) {
	n.NodeMask = v
}
func (n *Node) GetDscriptions() []string {
	return n.Dscriptions
}
func (n *Node) SetDscriptions(v []string) {
	n.Dscriptions = v
}
func (n *Node) GetInitialBound() *Sphere3f {
	return &n.InitialBound
}
func (n *Node) SetInitialBound(sp Sphere3f) {
	n.InitialBound = sp
}
func (n *Node) GetStates() *StateSet {
	return n.States
}
func (n *Node) SetStates(ss *StateSet) {
	n.States = ss
}
func (n *Node) GetParents() []*Group {
	return n.Parents
}
func (n *Node) SetParents(g []*Group) {
	n.Parents = g
}
func (n *Node) GetSphereCallback() *ComputeBoundingSphereCallback {
	return n.SphereCallback
}
func (n *Node) SetSphereCallback(cb *ComputeBoundingSphereCallback) {
	n.SphereCallback = cb
}
func (n *Node) SetUpdateCallback(cb *Callback) {
	n.UpdateCallback = cb
}
func (n *Node) GetUpdateCallback() *Callback {
	return n.UpdateCallback
}
func (n *Node) SetEventCallback(cb *Callback) {
	n.EventCallback = cb
}
func (n *Node) GetEventCallback() *Callback {
	return n.EventCallback
}

func (n *Node) GetCullCallback() *Callback {
	return n.CullCallback
}
func (n *Node) SetCullCallback(cb *Callback) {
	n.CullCallback = cb
}

func NewNode() *Node {
	obj := NewObject()
	obj.Type = NODET
	return &Node{Object: *obj, NodeMask: 0xffffffff}
}

func (n *Node) Accept(nv *NodeVisitor) {
	if nv.ValidNodeMask(n) {
		nv.PushOntoNodePath(n)
		nv.Apply(n)
		nv.PopFromNodePath()
	}
}

func (n *Node) Ascend(nv *NodeVisitor) {
	for _, g := range n.Parents {
		g.Accept(nv)
	}
}

func (n *Node) Traverse(nv *NodeVisitor) {

}

func (n *Node) IsNode() bool {
	return true
}
