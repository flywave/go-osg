package model

type Binding int

const (
	BIND_UNDEFINED         Binding = -1
	BIND_OFF               Binding = 0
	BIND_OVERALL           Binding = 1
	BIND_PER_PRIMITIVE_SET Binding = 2
	BIND_PER_VERTEX        Binding = 4
)

type Array struct {
	Type             ArrayTable
	data_size        int
	data_type        int
	Binding          Binding
	Normalize        bool
	PreserveDataType bool
}
