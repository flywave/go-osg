package model

const (
	POSITIONATTITUDETRANSFORMT string = "osg::PositionAttitudeTransform"
)

type PositionAttitudeTransform struct {
	Transform
	Position [3]float64
	Attitude [4]float64
	Scale    [3]float64
}

func NewPositionAttitudeTransform() *PositionAttitudeTransform {
	t := NewTransform()
	t.Type = POSITIONATTITUDETRANSFORMT
	return &PositionAttitudeTransform{Transform: *t}
}
