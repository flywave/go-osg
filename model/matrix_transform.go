package model

const (
	MATRIXTRANSFORMT = "osg::MatrixTransform"
)

type MatrixTransform struct {
	Transform
	Matrix [4][4]float32
}

func NewMatrixTransform() *MatrixTransform {
	mt := NewTransform()
	return &MatrixTransform{Transform: *mt}
}
