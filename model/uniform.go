package model

import (
	"github.com/ungerik/go3d/vec2"
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
)

type UniformBase struct {
	Object
	NameId uint32
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
