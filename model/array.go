package model

type BufferData struct {
	Object
	BufferIndex uint32
}

func NewBufferData() BufferData {
	obj := NewObject()
	return BufferData{Object: obj}
}

type Binding int

const (
	BIND_UNDEFINED         Binding = -1
	BIND_OFF               Binding = 0
	BIND_OVERALL           Binding = 1
	BIND_PER_PRIMITIVE_SET Binding = 2
	BIND_PER_VERTEX        Binding = 4
)

type Array struct {
	BufferData
	Type             ArrayTable
	DataSize         int
	DataType         int
	Binding          Binding
	Normalize        bool
	PreserveDataType bool
}

func NewArray() Array {
	buf := NewBufferData()
	return Array{BufferData: buf, DataSize: 0, DataType: 0, Normalize: false, PreserveDataType: false, Binding: BIND_UNDEFINED}
}
