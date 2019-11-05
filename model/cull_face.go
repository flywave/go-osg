package model

const (
	FRONT                 = GL_FRONT
	BACK                  = GL_BACK
	FRONT_AND_BACK        = GL_FRONT_AND_BACK
	CullFaceType   string = "osg::CullFace"
)

type CullFace struct {
	StateAttribute
	Mode Glenum
}

func NewCullFace() CullFace {
	att := NewStateAttribute()
	att.Type = CullFaceType
	return CullFace{StateAttribute: att, Mode: BACK}
}
