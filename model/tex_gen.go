package model

type Coord uint32

const (
	OBJECTLINEAR  = GLOBJECTLINEAR
	EYELINEAR     = GLEYELINEAR
	SPHEREMAP     = GLSPHEREMAP
	NORMALMAP     = GLNORMALMAP
	REFLECTIONMAP = GLREFLECTIONMAP

	S Coord = 0
	T Coord = 1
	R Coord = 2
	Q Coord = 3
)

type TexGen struct {
	StateAttribute
	Mode int
}

func NewTexGen() *TexGen {
	sa := NewStateAttribute()
	return &TexGen{StateAttribute: *sa}
}
