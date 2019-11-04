package model

const (
	NEVER                = GL_NEVER
	LESS                 = GL_LESS
	EQUAL                = GL_EQUAL
	LEQUAL               = GL_LEQUAL
	GREATER              = GL_GREATER
	NOTEQUAL             = GL_NOTEQUAL
	GEQUAL               = GL_GEQUAL
	ALWAYS               = GL_ALWAYS
	AlphaFuncType string = "osg::AlphaFunc"
)

type AlphaFunc struct {
	StateAttribute
	ReferenceValue float32
}
