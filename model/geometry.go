package model

import "errors"

const (
	GEOMETRYT string = "osg::Geometry"
)

type Geometry struct {
	Drawable
	Primitives          []interface{}
	VertexArray         *Array
	NormalArray         *Array
	FaceVrray           *Array
	ColorArray          *Array
	SecondaryColorArray *Array
	FogCoordArray       *Array
	TexCoordArrayList   []*Array
	VertexAttribList    []*Array
}

func NewGeometry() *Geometry {
	dw := NewDrawable()
	dw.Type = GEOMETRYT
	return &Geometry{Drawable: *dw}
}

func (g *Geometry) AddPrimitiveSet(p interface{}) {
	g.Primitives = append(g.Primitives, p)
}

func (g *Geometry) SetPrimitiveSet(i int, p interface{}) error {
	l := len(g.Primitives)
	if i > l-1 {
		return errors.New("out of range")
	}
	g.Primitives[i] = p
	return nil
}

func (g *Geometry) InsertPrimitiveSet(i int, p interface{}) error {
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

func (g *Geometry) SetTexCoordArrayBinding(i int, array *Array, binding int32) error {
	l := len(g.TexCoordArrayList)
	if int(i) > l-1 {
		return errors.New("out of range")
	}
	if binding != BINDUNDEFINED {
		array.Binding = binding
	} else {
		array.Binding = BINDPERVERTEX
	}
	g.TexCoordArrayList[i] = array
	return nil
}

func (g *Geometry) SetTexCoordArray(i int, array *Array) error {
	return g.SetTexCoordArrayBinding(i, array, BINDUNDEFINED)
}

func (g *Geometry) SetVertexAttribArray(i int, array *Array, binding int32) error {
	l := len(g.VertexAttribList)
	if i > l-1 {
		return errors.New("out of range")
	}
	if binding != BINDUNDEFINED {
		array.Binding = binding
	} else {
		array.Binding = BINDPERVERTEX

	}
	g.VertexAttribList[i] = array
	return nil
}

func (g *Geometry) Accept(nv *NodeVisitor) {
	nv.Geos = append(nv.Geos, g)
	if nv.ValidNodeMask(g) {
		nv.PushOntoNodePath(g)
		nv.Apply(g)
		nv.PopFromNodePath()
	}
}
