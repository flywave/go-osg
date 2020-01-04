package model

const (
	SHAPET      string = "osg::Shape"
	SHAPEMODELT string = "osg::ShadeModel"
)

type Shape struct {
	Object
}

func NewShape() *Shape {
	obj := NewObject()
	obj.Type = SHAPET
	return &Shape{Object: *obj}
}
