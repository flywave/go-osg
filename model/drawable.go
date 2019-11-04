package model

import "github.com/ungerik/go3d/vec3"

const (
	DrawableType string = "osg::Drawable"
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
	n.Type = DrawableType
	return Drawable{Node: n}
}
