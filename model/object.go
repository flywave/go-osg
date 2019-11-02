package model

type DataVariance int32

const (
	DYNAMIC     DataVariance = 0
	STATIC      DataVariance = 1
	UNSPECIFIED DataVariance = 2
	ObjectType  string       = "osg::Object"
)

type Object struct {
	name         string
	Type         string
	Propertys    map[string]string
	DataVariance DataVariance
	Udc          *UserDataContainer
}

func NewObject() Object {
	return &object{Type: ObjectType, DataVariance: UNSPECIFIED}
}

type Callback struct {
	Object
	Callback *Callback
}
