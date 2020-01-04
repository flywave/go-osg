package model

const (
	OSGHEADERLOW  = 0x6C910EA1
	OSGHEADERHIGH = 0x1AFB4545

	IMAGEINLINEDATA int32 = 0
	IMAGEINLINEFILE int32 = 1
	IMAGEEXTERNAL   int32 = 2
	IMAGEWRITEOUT   int32 = 3
)

type ObjectGlenum struct {
	Value int32
}

func NewObjectGlenum() *ObjectGlenum {
	return &ObjectGlenum{}
}

type ObjectProperty struct {
	Name        string
	Value       int32
	MapProperty bool
}

func NewObjectProperty() *ObjectProperty {
	return &ObjectProperty{}
}

type ObjectMark struct {
	Name        string
	IndentDelta int32
}

func NewObjectMark() *ObjectMark {
	return &ObjectMark{}
}
