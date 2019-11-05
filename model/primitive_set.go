package model

type PrimitiveFunctor interface {
	DrawArrays(mode Glenum, first int, count uint64)

	DrawElements8(mode Glenum, count uint64, indices []uint8)

	DrawElements16(mode Glenum, count uint64, indices []uint16)

	DrawElements32(mode Glenum, count uint64, indices []uint32)
}

type PrimitiveIndexFunctor interface {
	DrawArrays(mode Glenum, first int, count uint64)

	DrawElements8(mode Glenum, count uint64, indices []uint8)

	DrawElements16(mode Glenum, count uint64, indices []uint16)

	DrawElements32(mode Glenum, count uint64, indices []uint32)
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
	PrimitiveSetType         string = "osg::PrimitiveSet"
)

type PrimitiveSet struct {
	BufferData
	PrimitiveType PrimitiveTableEnum
	NumInstances  int
	Mode          uint
}

func NewPrimitiveSet() PrimitiveSet {
	bf := NewBufferData()
	bf.Type = PrimitiveSetType
	return PrimitiveSet{BufferData: bf, NumInstances: 0, Mode: 0}
}

const (
	DrawArraysType        string = "osg::DrawArrays"
	DrawArrayLengthsType  string = "osg::DrawArrayLengths"
	DrawElementsUbyteType string = "osg::DrawElementsUByte"
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
		f.DrawArrays(Glenum(d.Mode), d.First, d.Count)
	}
}

func NewDrawArrays() DrawArrays {
	p := NewPrimitiveSet()
	p.Type = DrawArraysType
	return DrawArrays{PrimitiveSet: p, First: 0, Count: 0}
}

type DrawArrayLengths struct {
	PrimitiveSet
	Data  []uint64
	First int
}

func NewDrawArrayLengths() DrawArrayLengths {
	p := NewPrimitiveSet()
	p.Type = DrawArrayLengthsType
	return DrawArrayLengths{PrimitiveSet: p, First: 0}
}

func (dal *DrawArrayLengths) Accept(functor interface{}) {
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		for _, v := range dal.Data {
			f.DrawArrays(Glenum(dal.Mode), dal.First, v)
		}
	}
}

func (dal *DrawArrayLengths) GetNumPrimitives() int {
	l := len(dal.Data)
	switch PrimitiveTableEnum(dal.Mode) {
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

func (dw *DrawElementsUbyte) ResizeElements(size uint64) {
	dw.Data = make([]uint8, size, size)
}

func (dw *DrawElementsUbyte) AddElement(e uint8) {
	dw.Data = append(dw.Data, e)
}

func (dal *DrawElementsUbyte) Accept(functor interface{}) {
	l := len(dal.Data)
	if l == 0 {
		return
	}
	switch f := functor.(type) {
	case PrimitiveFunctor:
	case PrimitiveIndexFunctor:
		f.DrawElements8(Glenum(dal.Mode), uint64(l), dal.Data)
	}
}
