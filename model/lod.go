package model

import (
	"errors"

	"github.com/ungerik/go3d/vec3"
)

type CenterMode uint32
type RangeMode uint32

const (
	USE_BOUNDING_SPHERE_CENTER                CenterMode = 0
	USER_DEFINED_CENTER                       CenterMode = 1
	UNION_OF_BOUNDING_SPHERE_AND_USER_DEFINED CenterMode = 2
	DISTANCE_FROM_EYE_POINT                   RangeMode  = 0
	PIXEL_SIZE_ON_SCREEN                      RangeMode  = 1
	LOD_T                                     string     = "osg::Lod"
)

type MinMaxPair []float32

type RangeListType []MinMaxPair

type Lod struct {
	Group
	Cmode     CenterMode
	Center    vec3.T
	Radius    float32
	Rmode     RangeMode
	RangeList RangeListType
}

func NewLod() Lod {
	g := NewGroup()
	g.Type = LOD_T
	return Lod{Group: g}
}

func (lod *Lod) AddChild(n *Node) {
	lod.Group.AddChild(n)
	rl := len(lod.RangeList)
	if len(lod.Group.Children) > rl {
		f := []float32{0, 0}
		if rl > 0 {
			last := lod.RangeList[rl-1][1]
			f[0] = last
			f[1] = last
		}
		lod.RangeList = append(lod.RangeList, f)
	}
}

func (lod *Lod) AddChild3(n *Node, min float32, max float32) {
	lod.Group.AddChild(n)
	rl := len(lod.RangeList)
	if len(lod.Group.Children) > rl {
		f := []float32{min, max}
		lod.RangeList = append(lod.RangeList, f)
	}
}

func (lod *Lod) RemoveChild2(pos int, count int) error {
	if lod.Group.RemoveChild2(pos, count) == nil {
		l := len(lod.RangeList)
		if pos > l-1 || pos+count > l {
			return errors.New("pos out of range")
		}

		a := lod.RangeList[:pos]
		b := lod.RangeList[pos+1+count:]
		lod.RangeList = append(a, b...)
		return nil
	}
	return errors.New("remove child error")
}

func (lod *Lod) SetRange(childNo uint, min float32, max float32) {
	f := []float32{min, max}
	lod.RangeList[childNo] = f
}
