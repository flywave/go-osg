package model

const (
	ALPHAFUNC_T string = "osg::AlphaFunc"
)

type AlphaFunc struct {
	StateAttribute
	ReferenceValue float32
}

func NewAlphaFunc() AlphaFunc {
	att := NewStateAttribute()
	att.Type = ALPHAFUNC_T
	return AlphaFunc{StateAttribute: att, ReferenceValue: 1.0}
}
