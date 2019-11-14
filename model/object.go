package model

import "reflect"

type DataVariance int32

var Type_Mapping map[string]interface{}

const (
	DYNAMIC            = 0
	STATIC             = 1
	UNSPECIFIED        = 2
	OBJECT_T    string = "osg::Object"
)

type ObjectInterface interface {
	IsObject() bool
}

type Object struct {
	Name         string
	Type         string
	Propertys    map[string]string
	DataVariance int
	Udc          UserDataContainer
}

func NewObject() Object {
	return Object{Type: OBJECT_T, DataVariance: UNSPECIFIED, Propertys: make(map[string]string), Udc: NewUserDataContainer()}
}

func (obj *Object) IsObject() bool {
	return true
}

type Callback struct {
	Object
	Callback *Callback
}

func TypeBaseOfObject(obj interface{}) bool {
	no := NewObject()
	baset := reflect.TypeOf(no)
	t := reflect.TypeOf(obj)
	return t.Implements(baset)
}
