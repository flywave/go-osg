package model

type DataVariance int32

var TypeMapping map[string]interface{}

const (
	DYNAMIC            = 0
	STATIC             = 1
	UNSPECIFIED        = 2
	OBJECTT     string = "osg::Object"
)

type ObjectInterface interface {
	IsObject() bool
	GetName() *string
	SetName(name string)
	GetProperty() map[string]string
	GetDataVariance() *int32
	SetDataVariance(int32)
	GetUserDataContainer() *UserDataContainer
	SetUserDataContainer(udc *UserDataContainer)
}

type Object struct {
	Name         string
	Type         string
	Propertys    map[string]string
	DataVariance int32
	Udc          *UserDataContainer
}

func NewObject() *Object {
	udc := NewUserDataContainer()
	return &Object{Type: OBJECTT, DataVariance: UNSPECIFIED, Propertys: make(map[string]string), Udc: udc}
}

func (obj *Object) GetUserDataContainer() *UserDataContainer {
	return obj.Udc
}
func (obj *Object) SetUserDataContainer(udc *UserDataContainer) {
	obj.Udc = udc
}

func (obj *Object) GetDataVariance() *int32 {
	return &obj.DataVariance
}

func (obj *Object) SetDataVariance(d int32) {
	obj.DataVariance = d
}

func (obj *Object) IsObject() bool {
	return true
}

func (obj *Object) GetName() *string {
	return &obj.Name
}

func (obj *Object) SetName(name string) {
	obj.Name = name
}

func (obj *Object) GetProperty() map[string]string {
	return obj.Propertys
}

type Callback struct {
	Object
	Callback *Callback
}

type ValueObject struct {
	Object
	Value interface{}
}

func NewValueObject() *ValueObject {
	ob := NewObject()
	ob.Type = "osg:ValueObject"
	return &ValueObject{Object: *ob}
}
