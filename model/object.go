package model

import "reflect"

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
}

type Object struct {
	Name         string
	Type         string
	Propertys    map[string]string
	DataVariance int
	Udc          *UserDataContainer
}

func NewObject() Object {
	udc := NewUserDataContainer()
	return Object{Type: OBJECTT, DataVariance: UNSPECIFIED, Propertys: make(map[string]string), Udc: &udc}
}

func (obj *Object) IsObject() bool {
	return true
}

type Callback struct {
	Object
	Callback *Callback
}

func IsBaseOfObject(obj interface{}) bool {
	no := NewObject()
	baset := reflect.TypeOf(no)
	t := reflect.TypeOf(obj)
	return t.Implements(baset)
}

type ValueObject struct {
	Object
	Value interface{}
}

func NewValueObject() ValueObject {
	ob := NewObject()
	ob.Type = "osg:ValueObject"
	return ValueObject{Object: ob}
}
