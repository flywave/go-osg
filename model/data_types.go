package model

const (
	OSG_HEADER_LOW  = 0x6C910EA1
	OSG_HEADER_HIGH = 0x1AFB4545

	IMAGE_INLINE_DATA int32 = 0
	IMAGE_INLINE_FILE int32 = 1
	IMAGE_EXTERNAL    int32 = 2
	IMAGE_WRITE_OUT   int32 = 3
)

type ObjectGlenum struct {
	Value int32
}

func NewObjectGlenum() ObjectGlenum {
	return ObjectGlenum{}
}

type ObjectProperty struct {
	Name        string
	Value       int32
	MapProperty bool
}

func NewObjectProperty() ObjectProperty {
	return ObjectProperty{}
}

type ObjectMark struct {
	Name        string
	IndentDelta int32
}

func NewObjectMark() ObjectMark {
	return ObjectMark{}
}
