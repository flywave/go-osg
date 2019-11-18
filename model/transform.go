package model

const (
	RELATIVE_RF                   = 0
	ABSOLUTE_RF                   = 1
	ABSOLUTE_RF_INHERIT_VIEWPOINT = 2

	TRANSFORM_T string = "osg::Transform"
)

type Transform struct {
	Group
	ReferenceFrame int
}

func NewTransform() Transform {
	g := NewGroup()
	g.Type = TRANSFORM_T
	return Transform{Group: g, ReferenceFrame: RELATIVE_RF}
}
