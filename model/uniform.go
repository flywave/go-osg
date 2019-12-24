package model

import (
	"reflect"

	"github.com/flywave/go3d/vec2"
	"github.com/flywave/go3d/vec3"
	"github.com/flywave/go3d/vec4"
)

type UniformInterface interface {
	IsUniformInterface() bool
}

type UniformBase struct {
	Object
	NameId uint32
}

func (u *UniformBase) IsUniformInterface() bool {
	return true
}

type UIntUniform struct {
	UniformBase
	Value uint
}
type FloatUniform struct {
	UniformBase
	Value float32
}
type DoubleUniform struct {
	UniformBase
	Value float64
}
type Vec2iUniform struct {
	UniformBase
	Value vec2.T
}
type Vec3iUniform struct {
	UniformBase
	Value vec3.T
}
type Vec4iUniform struct {
	UniformBase
	Value vec4.T
}
type Vec2uiUniform struct {
	UniformBase
}
type Vec3uiUniform struct {
	UniformBase
}
type Vec4uiUniform struct {
	UniformBase
}
type Vec2Uniform struct {
	UniformBase
}
type Vec3Uniform struct {
	UniformBase
}
type Vec4Uniformstruct struct {
	UniformBase
}
type Vec2dUniformstruct struct {
	UniformBase
}
type Vec3dUniformstruct struct {
	UniformBase
}
type Vec4dUniform struct {
	UniformBase
}
type PlaneUniform struct {
	UniformBase
}
type MatrixfUniform struct {
	UniformBase
}
type MatrixdUniform struct {
	UniformBase
}

func IsBaseOfUniform(obj interface{}) bool {
	if obj == nil {
		return false
	}
	ss := UniformBase{}
	baset := reflect.TypeOf(ss)
	t := reflect.TypeOf(obj)
	return t.Implements(baset)
}
