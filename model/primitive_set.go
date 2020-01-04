package model

type PrimitiveFunctor interface {
	DrawArrays(mode int32, first int32, count int32)

	DrawElements8(mode int32, count int32, indices []uint8)

	DrawElements16(mode int32, count int32, indices []uint16)

	DrawElements32(mode int32, count int32, indices []uint32)
}

type PrimitiveIndexFunctor interface {
	DrawArrays(mode int32, first int32, count int32)

	DrawElements8(mode int32, count int32, indices []uint8)

	DrawElements16(mode int32, count int32, indices []uint16)

	DrawElements32(mode int32, count int32, indices []uint32)
}

const (
	POINTS                        = GLPOINTS
	LINES                         = GLLINES
	LINESTRIP                     = GLLINESTRIP
	LINELOOP                      = GLLINELOOP
	TRIANGLES                     = GLTRIANGLES
	TRIANGLESTRIP                 = GLTRIANGLESTRIP
	TRIANGLEFAN                   = GLTRIANGLEFAN
	QUADS                         = GLQUADS
	QUADSTRIP                     = GLQUADSTRIP
	POLYGON                       = GLPOLYGON
	LINESADJACENCY                = GLLINESADJACENCY
	LINESTRIPADJACENCY            = GLLINESTRIPADJACENCY
	TRIANGLESADJACENCY            = GLTRIANGLESADJACENCY
	TRIANGLESTRIPADJACENCY        = GLTRIANGLESTRIPADJACENCY
	PATCHES                       = GLPATCHES
	PRIMITIVESETT          string = "osg::PrimitiveSet"
)

type PrimitiveSetInterface interface {
	Accept(interface{})
}

type PrimitiveSet struct {
	BufferData
	PrimitiveType int32
	NumInstances  int32
	Mode          int32
}

func NewPrimitiveSet() *PrimitiveSet {
	bf := NewBufferData()
	bf.Type = PRIMITIVESETT
	return &PrimitiveSet{BufferData: *bf, NumInstances: 0, Mode: 0}
}

const (
	DRAWARRAYT          string = "osg::DrawArrays"
	DRAWARRAYLENGHTT    string = "osg::DrawArrayLengths"
	DRWAELEMENTSUBYTET  string = "osg::DrawElementsUByte"
	DRWAELEMENTSUSHORTT string = "osg::DrawElementsUShort"
	DRWAELEMENTSUINTT   string = "osg::DrawElementsUInt"
)

type DrawArrays struct {
	PrimitiveSet
	First int32
	Count int32
}

func (d *DrawArrays) Accept(functor interface{}) {
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawArrays(d.Mode, d.First, d.Count)
	}
}

func NewDrawArrays() *DrawArrays {
	p := NewPrimitiveSet()
	p.Type = DRAWARRAYT
	p.PrimitiveType = IDDRAWARRAYS
	return &DrawArrays{PrimitiveSet: *p, First: 0, Count: 0}
}

type DrawArrayLengths struct {
	PrimitiveSet
	Data  []int32
	First int32
}

func NewDrawArrayLengths() *DrawArrayLengths {
	p := NewPrimitiveSet()
	p.Type = DRAWARRAYLENGHTT
	p.PrimitiveType = IDDRAWARRAYLENGTH
	return &DrawArrayLengths{PrimitiveSet: *p, First: 0}
}

func (dal *DrawArrayLengths) Accept(functor interface{}) {
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		for _, v := range dal.Data {
			f.DrawArrays(dal.Mode, dal.First, v)
		}
	}
}

func (dal *DrawArrayLengths) GetNumPrimitives() int {
	l := len(dal.Data)
	switch dal.Mode {
	case POINTS:
		return l
	case LINES:
		return l / 2
	case TRIANGLES:
		return l / 3
	case QUADS:
		return l / 4
	case LINESTRIP:
	case LINELOOP:
	case TRIANGLESTRIP:
	case TRIANGLEFAN:
	case QUADSTRIP:
	case PATCHES:
	case POLYGON:
		return l
	}
	return 0
}

type DrawElementsUByte struct {
	PrimitiveSet
	Data []uint8
}

func (dw *DrawElementsUByte) Size() uint64 {
	return uint64(len(dw.Data))
}

func (dw *DrawElementsUByte) ResizeElements(size uint64) {
	dw.Data = make([]uint8, size, size)
}

func (dw *DrawElementsUByte) AddElement(e uint8) {
	dw.Data = append(dw.Data, e)
}

func (dw *DrawElementsUByte) Accept(functor interface{}) {
	l := len(dw.Data)
	if l == 0 {
		return
	}
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawElements8(dw.Mode, int32(l), dw.Data)
	}
}

func NewDrawElementsUByte() *DrawElementsUByte {
	p := NewPrimitiveSet()
	p.Type = DRWAELEMENTSUBYTET
	p.PrimitiveType = IDDRAWELEMENTSUBYTE
	return &DrawElementsUByte{PrimitiveSet: *p}
}

type DrawElementsUShort struct {
	PrimitiveSet
	Data []uint16
}

func (dw *DrawElementsUShort) Size() uint64 {
	return uint64(len(dw.Data))
}

func (dw *DrawElementsUShort) ResizeElements(size uint64) {
	dw.Data = make([]uint16, size, size)
}

func (dw *DrawElementsUShort) AddElement(e uint16) {
	dw.Data = append(dw.Data, e)
}

func (dw *DrawElementsUShort) Accept(functor interface{}) {
	l := len(dw.Data)
	if l == 0 {
		return
	}
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawElements16(dw.Mode, int32(l), dw.Data)
	}
}

func NewDrawElementsUShort() *DrawElementsUShort {
	p := NewPrimitiveSet()
	p.Type = DRWAELEMENTSUSHORTT
	p.PrimitiveType = IDDRAWELEMENTSUSHORT
	return &DrawElementsUShort{PrimitiveSet: *p}
}

type DrawElementsUInt struct {
	PrimitiveSet
	Data []uint32
}

func (dw *DrawElementsUInt) Size() uint64 {
	return uint64(len(dw.Data))
}

func (dw *DrawElementsUInt) ResizeElements(size uint64) {
	dw.Data = make([]uint32, size, size)
}

func (dw *DrawElementsUInt) AddElement(e uint32) {
	dw.Data = append(dw.Data, e)
}

func (dw *DrawElementsUInt) Accept(functor interface{}) {
	l := len(dw.Data)
	if l == 0 {
		return
	}
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawElements32(dw.Mode, int32(l), dw.Data)
	}
}

func NewDrawElementsUInt() *DrawElementsUInt {
	p := NewPrimitiveSet()
	p.Type = DRWAELEMENTSUINTT
	p.PrimitiveType = IDDRAWELEMENTSUINT
	return &DrawElementsUInt{PrimitiveSet: *p}
}
