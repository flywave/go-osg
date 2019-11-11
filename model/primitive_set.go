package model

type PrimitiveFunctor interface {
	DrawArrays(mode int, first int, count uint64)

	DrawElements8(mode int, count uint64, indices []uint8)

	DrawElements16(mode int, count uint64, indices []uint16)

	DrawElements32(mode int, count uint64, indices []uint32)
}

type PrimitiveIndexFunctor interface {
	DrawArrays(mode int, first int, count uint64)

	DrawElements8(mode int, count uint64, indices []uint8)

	DrawElements16(mode int, count uint64, indices []uint16)

	DrawElements32(mode int, count uint64, indices []uint32)
}

const (
	POINTS                          = GL_POINTS
	LINES                           = GL_LINES
	LINE_STRIP                      = GL_LINE_STRIP
	LINE_LOOP                       = GL_LINE_LOOP
	TRIANGLES                       = GL_TRIANGLES
	TRIANGLE_STRIP                  = GL_TRIANGLE_STRIP
	TRIANGLE_FAN                    = GL_TRIANGLE_FAN
	QUADS                           = GL_QUADS
	QUAD_STRIP                      = GL_QUAD_STRIP
	POLYGON                         = GL_POLYGON
	LINES_ADJACENCY                 = GL_LINES_ADJACENCY
	LINE_STRIP_ADJACENCY            = GL_LINE_STRIP_ADJACENCY
	TRIANGLES_ADJACENCY             = GL_TRIANGLES_ADJACENCY
	TRIANGLE_STRIP_ADJACENCY        = GL_TRIANGLE_STRIP_ADJACENCY
	PATCHES                         = GL_PATCHES
	PRIMITIVESET_T           string = "osg::PrimitiveSet"
)

type PrimitiveSet struct {
	BufferData
	PrimitiveType int
	NumInstances  int
	Mode          uint
}

func NewPrimitiveSet() PrimitiveSet {
	bf := NewBufferData()
	bf.Type = PRIMITIVESET_T
	return PrimitiveSet{BufferData: bf, NumInstances: 0, Mode: 0}
}

const (
	DRAWARRAY_T          string = "osg::DrawArrays"
	DRAWARRAYLENGHT_T    string = "osg::DrawArrayLengths"
	DRWAELEMENTSUBYTE_T  string = "osg::DrawElementsUByte"
	DRWAELEMENTSUSHORT_T string = "osg::DrawElementsUShort"
	DRWAELEMENTSUINT_T   string = "osg::DrawElementsUInt"
)

type DrawArrays struct {
	PrimitiveSet
	First int
	Count uint64
}

func (d *DrawArrays) Accept(functor interface{}) {
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawArrays(int(d.Mode), d.First, d.Count)
	}
}

func NewDrawArrays() DrawArrays {
	p := NewPrimitiveSet()
	p.Type = DRAWARRAY_T
	return DrawArrays{PrimitiveSet: p, First: 0, Count: 0}
}

type DrawArrayLengths struct {
	PrimitiveSet
	Data  []uint64
	First int
}

func NewDrawArrayLengths() DrawArrayLengths {
	p := NewPrimitiveSet()
	p.Type = DRAWARRAYLENGHT_T
	return DrawArrayLengths{PrimitiveSet: p, First: 0}
}

func (dal *DrawArrayLengths) Accept(functor interface{}) {
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		for _, v := range dal.Data {
			f.DrawArrays(int(dal.Mode), dal.First, v)
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
	case LINE_STRIP:
	case LINE_LOOP:
	case TRIANGLE_STRIP:
	case TRIANGLE_FAN:
	case QUAD_STRIP:
	case PATCHES:
	case POLYGON:
		return l
	}
	return 0
}

type DrawElementsUbyte struct {
	PrimitiveSet
	Data []uint8
}

func (dw *DrawElementsUbyte) Size() uint64 {
	return uint64(len(dw.Data))
}

func (dw *DrawElementsUbyte) ResizeElements(size uint64) {
	dw.Data = make([]uint8, size, size)
}

func (dw *DrawElementsUbyte) AddElement(e uint8) {
	dw.Data = append(dw.Data, e)
}

func (dw *DrawElementsUbyte) Accept(functor interface{}) {
	l := len(dw.Data)
	if l == 0 {
		return
	}
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawElements8(int(dw.Mode), uint64(l), dw.Data)
	}
}

func NewDrawElementsUbyte() DrawElementsUbyte {
	p := NewPrimitiveSet()
	p.Type = DRWAELEMENTSUBYTE_T
	return DrawElementsUbyte{PrimitiveSet: p}
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
		f.DrawElements16(int(dw.Mode), uint64(l), dw.Data)
	}
}

func NewDrawElementsUShort() DrawElementsUShort {
	p := NewPrimitiveSet()
	p.Type = DRWAELEMENTSUSHORT_T
	return DrawElementsUShort{PrimitiveSet: p}
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
		f.DrawElements32(int(dw.Mode), uint64(l), dw.Data)
	}
}

func NewDrawElementsUInt() DrawElementsUInt {
	p := NewPrimitiveSet()
	p.Type = DRWAELEMENTSUINT_T
	return DrawElementsUInt{PrimitiveSet: p}
}
