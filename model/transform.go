package model

const (
	RELATIVERF                 = 0
	ABSOLUTERF                 = 1
	ABSOLUTERFINHERITVIEWPOINT = 2

	TRANSFORMT string = "osg::Transform"
)

type Transform struct {
	Group
	ReferenceFrame int
}

func NewTransform() *Transform {
	g := NewGroup()
	g.Type = TRANSFORMT
	return &Transform{Group: *g, ReferenceFrame: RELATIVERF}
}
