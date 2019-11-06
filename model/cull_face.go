package model

const (
	FRONT                 = GL_FRONT
	BACK                  = GL_BACK
	FRONT_AND_BACK        = GL_FRONT_AND_BACK
	CULLFACE_T     string = "osg::CullFace"
)

type CullFace struct {
	StateAttribute
	Mode int
}

func NewCullFace() CullFace {
	att := NewStateAttribute()
	att.Type = CULLFACE_T
	return CullFace{StateAttribute: att, Mode: BACK}
}
