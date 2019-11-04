package model

type Geometry struct {
	Drawable
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
	return Geometry{Drawable: dw}
}
