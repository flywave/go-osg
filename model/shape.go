package model

const (
	SHAPE_T      string = "osg::Shape"
	SHAPEMODEL_T string = "osg::ShadeModel"
)

type Shape struct {
	Object
}

func NewShape() Shape {
	obj := NewObject()
	obj.Type = SHAPE_T
	return Shape{Object: obj}
}
