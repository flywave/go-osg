package model

type DataVariance int32

var Type_Mapping map[string]interface{}

const (
	DYNAMIC     DataVariance = 0
	STATIC      DataVariance = 1
	UNSPECIFIED DataVariance = 2
	OBJECT_T    string       = "osg::Object"
)

type Object struct {
	Name         string
	Type         string
	Propertys    map[string]string
	DataVariance DataVariance
	Udc          UserDataContainer
}

func NewObject() Object {
	return Object{Type: OBJECT_T, DataVariance: UNSPECIFIED, Propertys: make(map[string]string), Udc: NewUserDataContainer()}
}

type Callback struct {
	Object
	Callback *Callback
}
