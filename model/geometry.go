package model

import "errors"

const (
	GeometryType string = "osg::Geometry"
)

type Geometry struct {
	Drawable
	Primitives          []interface{}
	ArrayList           *Array
	VertexArray         *Array
	FaceVrray           *Array
	ColorArray          *Array
	SecondaryColorArray *Array
	FogCoordArray       *Array
	TexCoordArray       []*Array
	VertexAttribList    []*Array
}

func NewGeometry() Geometry {
	dw := NewDrawable()
	dw.Type = GeometryType
	return Geometry{Drawable: dw}
}

func (g *Geometry) AddPrimitiveSet(p *PrimitiveSet) {
	g.Primitives = append(g.Primitives, p)
}

func (g *Geometry) SetPrimitiveSet(i int, p *PrimitiveSet) error {
	l := len(g.Primitives)
	if i > l-1 {
		return errors.New("out of range")
	}
	g.Primitives[i] = p
	return nil
}

func (g *Geometry) InsertPrimitiveSet(i int, p *PrimitiveSet) error {
	l := len(g.Primitives)
	if i > l-1 {
		return errors.New("out of range")
	}
	a := g.Primitives[:i]
	b := g.Primitives[i:]
	a = append(a, p)
	g.Primitives = append(a, b...)
	return nil
}

func (g *Geometry) RemovePrimitiveSet(i int, count int) error {
	l := len(g.Primitives)
	if i > l-1 || i < 0 || count < 0 || i+count >= l-1 {
		return errors.New("out of range")
	}
	a := g.Primitives[:i]
	b := g.Primitives[i+count:]
	g.Primitives = append(a, b...)
	return nil
}

func (g *Geometry) SetTexCoordArrayBinding(i int, array *Array, binding Binding) error {
	l := len(g.TexCoordArray)
	if i > l-1 {
		return errors.New("out of range")
	}
	if binding != BIND_UNDEFINED {
		array.Binding = binding
	} else {
		array.Binding = BIND_PER_VERTEX
	}
	g.TexCoordArray[i] = array
	return nil
}

func (g *Geometry) SetTexCoordArray(i int, array *Array) error {
	return g.SetTexCoordArrayBinding(i, array, BIND_UNDEFINED)
}

func (g *Geometry) SetVertexAttribArray(i int, array *Array, binding Binding) error {
	l := len(g.VertexAttribList)
	if i > l-1 {
		return errors.New("out of range")
	}
	if binding != BIND_UNDEFINED {
		array.Binding = binding
	} else {
		array.Binding = BIND_PER_VERTEX

	}
	g.VertexAttribList[i] = array
	return nil
}

func (g *Geometry) Accept(inter interface{}) {
	for _, pri := range g.Primitives {
		switch p := pri.(type) {
		case *DrawArrays:
		case *DrawArrayLengths:
			p.Accept(inter)
		}
	}
}
