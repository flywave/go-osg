package model

type ReferenceFrame uint32

const (
	RELATIVE_RF                   ReferenceFrame = 0
	ABSOLUTE_RF                   ReferenceFrame = 1
	ABSOLUTE_RF_INHERIT_VIEWPOINT ReferenceFrame = 2

	TRANSFORM_T string = "osg::Transform"
)

type Transform struct {
	Group
	ReferenceFrame ReferenceFrame
}

func NewTransform() Transform {
	g := NewGroup()
	g.Type = TRANSFORM_T
	return Transform{Group: g, ReferenceFrame: RELATIVE_RF}
}
