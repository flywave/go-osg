package model

import "errors"

const (
	PAGEDLOD_T string = "osg::PagedLOD"
)

type PerRangeData struct {
	FileName        string
	PriorityOffset  float32
	PriorityScale   float32
	MinExpiryTime   float64
	MinExpiryFrames uint
	TimeStamp       float64
	FrameNumber     uint

	FrameNumberOfLastReleaseGLObjects uint
}

type PagedLod struct {
	Lod
	DataBasePath                   string
	FrameNumberOfLastTraversal     uint
	NumChildrenThatCannotBeExpired uint
	DisableExternalChildrenPaging  bool
	PerRangeDataListType           []PerRangeData
}

func NewPagedLod() PagedLod {
	lod := NewLod()
	lod.Type = PAGEDLOD_T
	return PagedLod{Lod: lod}
}

func (p *PagedLod) AddChild(n *Node) {
	p.Lod.AddChild(n)
	p.PerRangeDataListType = append(p.PerRangeDataListType, PerRangeData{})
}

func (p *PagedLod) AddChild3(n *Node, min float32, max float32) {
	p.Lod.AddChild3(n, min, max)
	p.PerRangeDataListType = append(p.PerRangeDataListType, PerRangeData{})
}

func (p *PagedLod) AddChild5(n *Node, min float32, max float32, filename string, priorityOffset float32, priorityScale float32) {
	p.Lod.AddChild3(n, min, max)
	p.PerRangeDataListType = append(p.PerRangeDataListType, PerRangeData{FileName: filename, PriorityOffset: priorityOffset, PriorityScale: priorityScale})
}

func (p *PagedLod) RemoveChild(pos int, count int) error {
	if p.Lod.RemoveChild2(pos, count) == nil {
		a := p.PerRangeDataListType[:pos]
		b := p.PerRangeDataListType[pos+1+count:]
		p.PerRangeDataListType = append(a, b...)
		return nil
	}
	return errors.New("remove child error")
}
