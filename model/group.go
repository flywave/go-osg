package model

import "errors"

const (
	GROUP_T string = "osg::Group"
)

type Group struct {
	Node
	Children []interface{}
}

func NewGroup() Group {
	n := NewNode()
	n.Type = GROUP_T
	return Group{Node: n}
}

func (g *Group) AddChild(n interface{}) {
	g.Children = append(g.Children, n)
}

func (g *Group) InsertChild(index int, n interface{}) {
	a := g.Children[:index]
	a = append(a, n)
	b := g.Children[index:]
	g.Children = append(a, b...)
}

func (g *Group) GetIndex(n interface{}) int {
	index := -1
	for i, val := range g.Children {
		if val == n {
			index = i
			break
		}
	}
	return index
}

func (g *Group) RemoveChild(n interface{}) error {
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

func (g *Group) ReplaceChild(origChild interface{}, newChild interface{}) error {
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

func (g *Group) SetChild(index int, newChild interface{}) error {
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

func (g *Group) Containsnode(n interface{}) bool {
	index := g.GetIndex(n)
	if index < 0 {
		return false
	}
	return true
}
