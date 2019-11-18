package model

type BufferData struct {
	Object
	BufferIndex uint32
}

func NewBufferData() BufferData {
	obj := NewObject()
	return BufferData{Object: obj}
}

type Binding int32

const (
	BIND_UNDEFINED         = -1
	BIND_OFF               = 0
	BIND_OVERALL           = 1
	BIND_PER_PRIMITIVE_SET = 2
	BIND_PER_VERTEX        = 4
)

type Array struct {
	BufferData
	Type             ArrayTable
	DataSize         int32
	DataType         int32
	Binding          int32
	Normalize        bool
	PreserveDataType bool
}

func NewArray() Array {
	buf := NewBufferData()
	return Array{BufferData: buf, DataSize: 0, DataType: 0, Normalize: false, PreserveDataType: false, Binding: BIND_UNDEFINED}
}
