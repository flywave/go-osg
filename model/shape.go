package model

const (
	SHAPE_T string = "osg::Shape"
)

type Shape struct {
	Object
}

func NewShape() Shape {
	obj := NewObject()
	obj.Type = SHAPE_T
	return Shape{Object: obj}
}
