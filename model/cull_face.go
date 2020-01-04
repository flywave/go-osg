package model

const (
	FRONT               = GLFRONT
	BACK                = GLBACK
	FRONTANDBACK        = GLFRONTANDBACK
	CULLFACET    string = "osg::CullFace"
)

type CullFace struct {
	StateAttribute
	Mode int
}

func NewCullFace() *CullFace {
	att := NewStateAttribute()
	att.Type = CULLFACET
	return &CullFace{StateAttribute: *att, Mode: BACK}
}
