package model

const (
	ShapeType string = "osg::Shape"
)

type Shape struct {
	Object
}

func NewShape() Shape {
	obj := NewObject()
	obj.Type = ShapeType
	return Shape{Object: obj}
}
