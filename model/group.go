package model

import "errors"

const (
	GroupType string = "osg::Group"
)

type Group struct {
	Node
	Children []*Node
}

func NewGroup() Group {
	n := NewNode()
	n.Type = GroupType
	return Group{Node: n}
}

func (g *Group) AddChild(n *Node) {
	g.Children = append(g.Children, n)
}

func (g *Group) InsertChild(index int, n *Node) {
	a := g.Children[:index]
	a = append(a, n)
	b := g.Children[index:]
	g.Children = append(a, b...)
}

func (g *Group) GetIndex(n *Node) int {
	index := -1
	for i, val := range g.Children {
		if val == n {
			index = i
			break
		}
	}
	return index
}

func (g *Group) RemoveChild(n *Node) error {
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

func (g *Group) ReplaceChild(origChild *Node, newChild *Node) error {
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

func (g *Group) SetChild(index int, newChild *Node) error {
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

func (g *Group) Containsnode(n *Node) bool {
	index := g.GetIndex(n)
	if index < 0 {
		return false
	}
	return true
}
