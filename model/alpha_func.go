package model

const (
	ALPHAFUNCT string = "osg::AlphaFunc"
)

type AlphaFunc struct {
	StateAttribute
	ReferenceValue float32
	ComparisonFunc int
}

func NewAlphaFunc() *AlphaFunc {
	att := NewStateAttribute()
	att.Type = ALPHAFUNCT
	return &AlphaFunc{StateAttribute: *att, ReferenceValue: 1.0, ComparisonFunc: GLALWAYS}
}
