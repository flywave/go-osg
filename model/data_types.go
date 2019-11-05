package model

const (
	OSG_HEADER_LOW  = 0x6C910EA1
	OSG_HEADER_HIGH = 0x1AFB4545

	IMAGE_INLINE_DATA int = 0
	IMAGE_INLINE_FILE int = 1
	IMAGE_EXTERNAL    int = 2
	IMAGE_WRITE_OUT   int = 3
)

type ObjectProperty struct {
	Name        string
	Value       int
	MapProperty bool
}

func NewObjectProperty() ObjectProperty {
	return ObjectProperty{}
}

type ObjectMark struct {
	Name        string
	IndentDelta int
}

func NewObjectMark() ObjectMark {
	return ObjectMark{}
}
