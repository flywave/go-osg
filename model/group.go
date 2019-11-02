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
	n.ObjectType = GroupType
	return Group{Node: n}
}

func AppendChild(g interface{}, n interface{}) error {
	switch t := g.(type) {
	case Group:
		var v *Node
		switch nd := n.(type) {
		case *Node:
			v = nd
		case Node:
			v = &nd
		default:
			return errors.New("type error")
		}
		t.Children = append(g.Children, v)
		break
	case Lod:
		switch nd := n.(type) {
		case MinMaxPair:
			t.RangeList = append(g.RangeList, nd)
		default:
			return errors.New("type error")
		}
		break
	}
}

// func InsertChild(g interface{},index int,n *Node)error{
// 		switch t:=g.(type){
// 		case Group:{
// 			var v *Node;
// 			switch nd:=n.(type){
// 			case *Node:
// 				v = nd;
//  			break
// 		case Node:
// 			v=&nd;
// 			break
// 		default:
// 			return errors.New("type error")
// 			}
// 			a:=t.Children[:index-1]
// 			a = append(a,v)
// 			b:=t.Children[index:]
// 			t.Children = append(a,b...)
// 			break
// 		}
// 		case Lod:{
// 			switch nd:=n.(type){
// 			case MinMaxPair:
// 			t.RangeList = append(g.RangeList,nd)
// 			break
// 		default:
// 			return errors.New("type error")
// 			}
// 			break
// 		}
// 		}
// 		return nil
// }

// func  (g *Group)GetIndex(g interface{},n *Node)int{
// 	index := -1;
// 	for i,val:range Children{
// 			if(val==n){
// 				index = i;
// 				break;
// 			}
// 	}
// 	return index;
// }

// func RemoveChild(g interface{},n *Node)bool{
// 	index: = g.GetIndex(n)
// 	if(index<0)return false
// 	a:=Children[:index-1]
// 	a = append(a,n)
// 	b:=Children[index:]
// 	Children = append(a,b...)
// 	return true;
// }

// func RemoveChild(g interface{},pos int,count int)bool{
// 	l:=len(Children)
// 	if(pos>l-1)return false;
// 	if(pos+count>l)return false
// 	a:=Children[:index-1]
// 	b:=Children[index+count:]
// 	Children = append(a,b...)
// 	return true;
// }

// func ReplaceChild(g interface{},origChild *Node,newChild *Node)bool{
// 	index: = g.GetIndex(origChild)
// 	if(index<0)return false
// 	a:=Children[:index-1]
// 	a = append(a,newChild)
// 	b:=Children[index-1:]
// 	Children = append(a,b...)
// 	return true;
// }

// func SetChild(g interface{},pos int,newChild *Node)bool{
// 	GetIndex[pos] = newChild
// 	return true;
// }

// func Containsnode(g interface{},n *Node)bool{
// 	index: = g.GetIndex(n)
// 	if(index<0)return false
// 	return true
// }
