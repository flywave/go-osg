package model

import (
	"errors"

	"github.com/ungerik/go3d/vec3"
)

const (
	USE_BOUNDING_SPHERE_CENTER                       = 0
	USER_DEFINED_CENTER                              = 1
	UNION_OF_BOUNDING_SPHERE_AND_USER_DEFINED        = 2
	DISTANCE_FROM_EYE_POINT                          = 0
	PIXEL_SIZE_ON_SCREEN                             = 1
	LOD_T                                     string = "osg::Lod"
)

type RangeListType [][2]float32

type Lod struct {
	Group
	Cmode     uint32
	Center    vec3.T
	Radius    float32
	Rmode     uint32
	RangeList RangeListType
}

func NewLod() Lod {
	g := NewGroup()
	g.Type = LOD_T
	return Lod{Group: g, Cmode: USE_BOUNDING_SPHERE_CENTER, Rmode: DISTANCE_FROM_EYE_POINT}
}

func (lod *Lod) AddChild(n interface{}) {
	lod.Group.AddChild(n)
	rl := len(lod.RangeList)
	if len(lod.Group.Children) > rl {
		f := [2]float32{0, 0}
		if rl > 0 {
			last := lod.RangeList[rl-1][1]
			f[0] = last
			f[1] = last
		}
		lod.RangeList = append(lod.RangeList, f)
	}
}

func (lod *Lod) AddChild3(n interface{}, min float32, max float32) {
	lod.Group.AddChild(n)
	rl := len(lod.RangeList)
	if len(lod.Group.Children) > rl {
		f := [2]float32{min, max}
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

func (lod *Lod) SetRange(childNo int, min float32, max float32) {
	l := len(lod.RangeList)
	if childNo >= l {
		s := childNo + 1 - l
		ls := make([][2]float32, s, s)
		lod.RangeList = append(lod.RangeList, ls...)
	}
	f := [2]float32{min, max}
	lod.RangeList[childNo] = f
}
