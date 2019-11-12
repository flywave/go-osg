package model

type Coord uint32

const (
	OBJECT_LINEAR  = GL_OBJECT_LINEAR
	EYE_LINEAR     = GL_EYE_LINEAR
	SPHERE_MAP     = GL_SPHERE_MAP
	NORMAL_MAP     = GL_NORMAL_MAP
	REFLECTION_MAP = GL_REFLECTION_MAP

	S Coord = 0
	T Coord = 1
	R Coord = 2
	Q Coord = 3
)

type TexGen struct {
	StateAttribute
	Mode int
}

func NewTexGen() TexGen {
	sa := NewStateAttribute()
	return TexGen{StateAttribute: sa}
}
