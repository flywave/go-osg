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
	BINDUNDEFINED       = -1
	BINDOFF             = 0
	BINDOVERALL         = 1
	BINDPERPRIMITIVESET = 2
	BINDPERVERTEX       = 4
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
	return Array{BufferData: buf, DataSize: 0, DataType: 0, Normalize: false, PreserveDataType: false, Binding: BINDUNDEFINED}
}
