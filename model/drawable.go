package model

import (
	"reflect"

	"github.com/ungerik/go3d/vec3"
)

const (
	DRAWABLET string = "osg::Drawable"
)

type DrawCallback struct {
	Object
}

type ComputeBoundingBoxCallback struct {
	Object
}

type DrawableInterface interface {
	IsDrawableInterface() bool
}

type Drawable struct {
	Node
	BoundingBox            vec3.Box
	Shape                  *Shape
	SupportsDisplayList    bool
	UseDisplayList         bool
	UseVertexBufferObjects bool

	Callback   *ComputeBoundingBoxCallback
	DwCallback *DrawCallback
}

func (d *Drawable) IsDrawableInterface() bool {
	return true
}

func NewDrawable() Drawable {
	n := NewNode()
	n.Type = DRAWABLET
	return Drawable{Node: n}
}

func IsBaseOfDrawable(obj interface{}) bool {
	if obj == nil {
		return false
	}
	ss := Drawable{}
	baset := reflect.TypeOf(ss)
	t := reflect.TypeOf(obj)
	return t.Implements(baset)
}
