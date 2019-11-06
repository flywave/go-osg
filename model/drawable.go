package model

import "github.com/ungerik/go3d/vec3"

const (
	DRAWABLE_T string = "osg::Drawable"
)

type DrawCallback struct {
	Object
}

type ComputeBoundingBoxCallback struct {
	Object
}

type Drawable struct {
	Node
	BoundingBox            vec3.Box
	Shape                  *Shape
	SupportsDisplayList    bool
	UseLisplayList         bool
	UseVertexBufferObjects bool

	Callback   *ComputeBoundingBoxCallback
	DwCallback *DrawCallback
}

func NewDrawable() Drawable {
	n := NewNode()
	n.Type = DRAWABLE_T
	return Drawable{Node: n}
}
