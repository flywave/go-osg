package model

import (
	"errors"
)

const (
	USEBOUNDINGSPHERECENTER                    = 0
	USERDEFINEDCENTER                          = 1
	UNIONOFBOUNDINGSPHEREANDUSERDEFINED        = 2
	DISTANCEFROMEYEPOINT                       = 0
	PIXELSIZEONSCREEN                          = 1
	LODT                                string = "osg::Lod"
)

type LodInterface interface {
	GetCmode() *uint32
	SetCmode(uint32)
	GetCenter() *[3]float32
	SetCenter([3]float32)
	GetRadius() *float32
	SetRadius(float32)
	GetRangeList() RangeListType
	SetRangeList(RangeListType)
	AddChild(NodeInterface)
	AddChild3(NodeInterface, float32, float32)
	RemoveChild2(int, int) error
	SetRange(int, float32, float32)
	GetRmode() *uint32
	SetRmode(uint32)
}

type RangeListType [][2]float32

type Lod struct {
	Group
	Cmode     uint32
	Center    [3]float32
	Radius    float32
	Rmode     uint32
	RangeList RangeListType
}

func NewLod() *Lod {
	g := NewGroup()
	g.Type = LODT
	return &Lod{Group: *g, Cmode: USEBOUNDINGSPHERECENTER, Rmode: DISTANCEFROMEYEPOINT}
}

func (lod *Lod) GetRmode() *uint32 {
	return &lod.Rmode
}
func (lod *Lod) SetRmode(c uint32) {
	lod.Rmode = c
}

func (lod *Lod) GetCmode() *uint32 {
	return &lod.Cmode
}
func (lod *Lod) SetCmode(c uint32) {
	lod.Cmode = c
}

func (lod *Lod) GetCenter() *[3]float32 {
	return &lod.Center
}
func (lod *Lod) SetCenter(ct [3]float32) {
	lod.Center = ct
}
func (lod *Lod) GetRadius() *float32 {
	return &lod.Radius
}
func (lod *Lod) SetRadius(r float32) {
	lod.Radius = r
}
func (lod *Lod) GetRangeList() RangeListType {
	return lod.RangeList
}
func (lod *Lod) SetRangeList(rl RangeListType) {
	lod.RangeList = rl
}

func (lod *Lod) AddChild(n NodeInterface) {
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

func (lod *Lod) AddChild3(n NodeInterface, min float32, max float32) {
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
