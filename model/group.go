package model

import "errors"

const (
	GROUPT string = "osg::Group"
)

type GroupInterface interface {
	SetChildren(c []NodeInterface)
	GetChildren() []NodeInterface
	AddChild(n NodeInterface)
	InsertChild(index int, n NodeInterface)
	GetIndex(n NodeInterface) int
	RemoveChild(n NodeInterface) error
	RemoveChild2(pos int, count int) error
	ReplaceChild(origChild NodeInterface, newChild NodeInterface) error
	SetChild(index int, newChild NodeInterface) error
	Containsnode(n NodeInterface) bool
}

type Group struct {
	Node
	Children []NodeInterface
}

func (g *Group) GetChildren() []NodeInterface {
	return g.Children
}

func (g *Group) SetChildren(c []NodeInterface) {
	g.Children = c
}

func (g *Group) Traverse(nv *NodeVisitor) {
	for _, c := range g.Children {
		c.Accept(nv)
	}
}

func NewGroup() *Group {
	n := NewNode()
	n.Type = GROUPT
	return &Group{Node: *n}
}

func (g *Group) AddChild(n NodeInterface) {
	g.Children = append(g.Children, n)
}

func (g *Group) InsertChild(index int, n NodeInterface) {
	a := g.Children[:index]
	a = append(a, n)
	b := g.Children[index:]
	g.Children = append(a, b...)
}

func (g *Group) GetIndex(n NodeInterface) int {
	index := -1
	for i, val := range g.Children {
		if val == n {
			index = i
			break
		}
	}
	return index
}

func (g *Group) RemoveChild(n NodeInterface) error {
	index := g.GetIndex(n)
	if index < 0 {
		return errors.New("have no this child")
	}

	a := g.Children[:index]
	a = append(a, n)
	b := g.Children[index+1:]
	g.Children = append(a, b...)
	return nil
}

func (g *Group) RemoveChild2(pos int, count int) error {
	if pos < 0 {
		return errors.New("pos out of range")
	}

	l := len(g.Children)
	if pos > l-1 || pos+count > l {
		return errors.New("pos out of range")
	}

	a := g.Children[:pos]
	b := g.Children[pos+1+count:]
	g.Children = append(a, b...)
	return nil
}

func (g *Group) ReplaceChild(origChild NodeInterface, newChild NodeInterface) error {
	index := g.GetIndex(origChild)
	if index < 0 {
		return errors.New("out of range")
	}

	a := g.Children[:index]
	a = append(a, newChild)
	b := g.Children[index+1:]
	g.Children = append(a, b...)
	return nil
}

func (g *Group) SetChild(index int, newChild NodeInterface) error {
	if index < 0 {
		return errors.New("out of range")
	}
	le := len(g.Children)
	if index >= le {
		return errors.New("out of range")
	}
	g.Children[index] = newChild

	return nil
}

func (g *Group) Containsnode(n NodeInterface) bool {
	index := g.GetIndex(n)
	if index < 0 {
		return false
	}
	return true
}
