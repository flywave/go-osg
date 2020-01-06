package model

import "errors"

const (
	PAGEDLODT string = "osg::PagedLOD"
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
	FrameNumberOfLastTraversal     uint32
	NumChildrenThatCannotBeExpired uint32
	DisableExternalChildrenPaging  bool
	PerRangeDataList               []PerRangeData
}

func NewPagedLod() *PagedLod {
	lod := NewLod()
	lod.Type = PAGEDLODT
	return &PagedLod{Lod: *lod}
}

func (p *PagedLod) expandPerRangeDataTo(pos int) {
	l := len(p.PerRangeDataList)
	if pos > l-1 {
		s := l - 1 - pos
		pl := make([]PerRangeData, s, s)
		p.PerRangeDataList = append(p.PerRangeDataList, pl...)
	}
}

func (p *PagedLod) Accept(nv *NodeVisitor) {
	if nv.ValidNodeMask(p) {
		nv.PushOntoNodePath(p)
		nv.Apply(p)
		nv.PopFromNodePath()
	}
}

func (p *PagedLod) AddChild(n NodeInterface) {
	p.Lod.AddChild(n)
	p.PerRangeDataList = append(p.PerRangeDataList, PerRangeData{})
}

func (p *PagedLod) AddChild3(n NodeInterface, min float32, max float32) {
	p.Lod.AddChild3(n, min, max)
	p.PerRangeDataList = append(p.PerRangeDataList, PerRangeData{})
}

func (p *PagedLod) AddChild5(n NodeInterface, min float32, max float32, filename string, priorityOffset float32, priorityScale float32) {
	p.Lod.AddChild3(n, min, max)
	p.PerRangeDataList = append(p.PerRangeDataList, PerRangeData{FileName: filename, PriorityOffset: priorityOffset, PriorityScale: priorityScale})
}

func (p *PagedLod) RemoveChild(pos int, count int) error {
	if p.Lod.RemoveChild2(pos, count) == nil {
		a := p.PerRangeDataList[:pos]
		b := p.PerRangeDataList[pos+1+count:]
		p.PerRangeDataList = append(a, b...)
		return nil
	}
	return errors.New("remove child error")
}

func (p *PagedLod) SetFileName(index int, name string) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].FileName = name
}

func (p *PagedLod) GetFileName(index int) string {
	return p.PerRangeDataList[index].FileName
}

func (p *PagedLod) SetPriorityOffset(index int, offset float32) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].PriorityOffset = offset

}

func (p *PagedLod) GetPriorityOffset(index int) float32 {
	return p.PerRangeDataList[index].PriorityOffset
}

func (p *PagedLod) SetPriorityScale(index int, scale float32) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].PriorityScale = scale
}

func (p *PagedLod) GetPriorityScale(index int) float32 {
	return p.PerRangeDataList[index].PriorityScale
}

func (p *PagedLod) SetMinimumExpiryTime(index int, time float64) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].MinExpiryTime = time
}

func (p *PagedLod) GetMinimumExpiryTime(index int) float64 {
	return p.PerRangeDataList[index].MinExpiryTime
}

func (p *PagedLod) SetMinimumExpiryFrames(index int, frm uint) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].MinExpiryFrames = frm
}

func (p *PagedLod) GetMinimumExpiryFrames(index int) uint {
	return p.PerRangeDataList[index].MinExpiryFrames
}

func (p *PagedLod) SetTimeStamp(index int, stamp float64) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].TimeStamp = stamp
}

func (p *PagedLod) GetTimeStamp(index int) float64 {
	return p.PerRangeDataList[index].TimeStamp
}

func (p *PagedLod) SetFrameNumber(index int, frm uint) {
	p.expandPerRangeDataTo(index)
	p.PerRangeDataList[index].FrameNumber = frm
}

func (p *PagedLod) GetFrameNumber(index int) uint {
	return p.PerRangeDataList[index].FrameNumber
}
