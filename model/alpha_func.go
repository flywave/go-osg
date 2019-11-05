package model

const (
	AlphaFuncType string = "osg::AlphaFunc"
)

type AlphaFunc struct {
	StateAttribute
	ReferenceValue float32
}

func NewAlphaFunc() AlphaFunc {
	att := NewStateAttribute()
	att.Type = AlphaFuncType
	return AlphaFunc{StateAttribute: att, ReferenceValue: 1.0}
}
