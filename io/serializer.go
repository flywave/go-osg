package io

import "github.com/flywave/go-osg/model"

type StringListType []string

func Split(src string, list StringListType, separator uint8) {
	var start int = -1
	for i, c := range src {
		if uint8(c) == separator {
			start = i
		}
	}
	l := len(src)
	for {
		if start == l {
			break
		}
		var end int = -1
		for i, c := range src {
			if uint8(c) == separator {
				end = i
			}
		}

		list = append(list, src[start:end-start])
		if end != l {
			for i, c := range src {
				if uint8(c) == separator {
					start = i
				}
			}
		} else {
			start = end
		}
	}
}

type Value int
type StringToValue map[string]Value
type ValueToString map[Value]string

type IntLookup struct {
	StringToValue StringToValue
	ValueToString ValueToString
}

func (lk *IntLookup) Size() int {
	return len(lk.StringToValue)
}

func (lk *IntLookup) Add(str string, val Value) {
	_, ok := lk.ValueToString[val]
	if ok {
		return
	}
	lk.ValueToString[val] = str
	lk.StringToValue[str] = val
}

func (lk *IntLookup) Add2(str string, newStr string, val Value) {
	_, ok := lk.ValueToString[val]
	if ok {
		return
	}
	lk.ValueToString[val] = str
	lk.StringToValue[newStr] = val
	lk.StringToValue[str] = val
}

func (lk *IntLookup) GetValue(str string) Value {
	v, ok := lk.StringToValue[str]
	if ok {
		var val Value = Value(str[0])
		lk.StringToValue[str] = val
		return val
	}
	return v
}

func (lk *IntLookup) GetString(val Value) string {
	s, ok := lk.ValueToString[val]
	if ok {
		s = string(val)
		lk.ValueToString[val] = s
		return s
	}
	return s
}

type SerType uint32
type Usage uint32

const (
	RW_UNDEFINED       SerType = 0
	RW_USER            SerType = 1
	RW_OBJECT          SerType = 2
	RW_IMAGE           SerType = 3
	RW_LIST            SerType = 4
	RW_BOOL            SerType = 5
	RW_CHAR            SerType = 6
	RW_UCHAR           SerType = 7
	RW_SHORT           SerType = 8
	RW_USHORT          SerType = 9
	RW_INT             SerType = 10
	RW_UINT            SerType = 11
	RW_FLOAT           SerType = 12
	RW_DOUBLE          SerType = 13
	RW_VEC2F           SerType = 14
	RW_VEC2D           SerType = 15
	RW_VEC3F           SerType = 16
	RW_VEC3D           SerType = 17
	RW_VEC4F           SerType = 18
	RW_VEC4D           SerType = 19
	RW_QUAT            SerType = 20
	RW_PLANE           SerType = 21
	RW_MATRIXF         SerType = 22
	RW_MATRIXD         SerType = 23
	RW_MATRIX          SerType = 24
	RW_GLENUM          SerType = 25
	RW_STRING          SerType = 26
	RW_ENUM            SerType = 27
	RW_VEC2B           SerType = 28
	RW_VEC2UB          SerType = 29
	RW_VEC2S           SerType = 30
	RW_VEC2US          SerType = 31
	RW_VEC2I           SerType = 32
	RW_VEC2UI          SerType = 33
	RW_VEC3B           SerType = 34
	RW_VEC3UB          SerType = 35
	RW_VEC3S           SerType = 36
	RW_VEC3US          SerType = 37
	RW_VEC3I           SerType = 38
	RW_VEC3UI          SerType = 39
	RW_VEC4B           SerType = 40
	RW_VEC4UB          SerType = 41
	RW_VEC4S           SerType = 42
	RW_VEC4US          SerType = 43
	RW_VEC4I           SerType = 44
	RW_VEC4UI          SerType = 45
	RW_BOUNDINGBOXF    SerType = 46
	RW_BOUNDINGBOXD    SerType = 47
	RW_BOUNDINGSPHEREF SerType = 48
	RW_BOUNDINGSPHERED SerType = 49
	RW_VECTOR          SerType = 50
	RW_MAP             SerType = 51

	READ_WRITE_PROPERTY Usage = 1
	GET_PROPERTY        Usage = 2
	SET_PROPERTY        Usage = 4
	GET_SET_PROPERTY    Usage = GET_PROPERTY | SET_PROPERTY
)

type BaseSerializer struct {
	First_version int
	Last_version  int
	Usage         Usage
}

func (bs *BaseSerializer) SupportsGetSet() bool {
	return (bs.Usage & GET_SET_PROPERTY) != 0
}

func (bs *BaseSerializer) SupportsGet() bool {
	return (bs.Usage & GET_PROPERTY) != 0
}

func (bs *BaseSerializer) SupportsSet() bool {
	return (bs.Usage & GET_SET_PROPERTY) != 0
}

func NewBaseSerializer(usg Usage) BaseSerializer {
	return BaseSerializer{Usage: usg}
}

type Checker func(interface{}) bool
type Reader func(*OsgIstream, interface{}) bool
type Writer func(*OsgOstream, interface{}) bool

type UserSerializer struct {
	BaseSerializer
	Checker Checker
	Reader  Reader
	Writer  Writer
	Name    string
}
 
func (ser *UserSerializer) Read(is *OsgIstream, obj *model.Object) {}

func (ser *UserSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewUserSerializer(name string, ck Checker, rd Reader, wt Writer) UserSerializer {
	ser := NewBaseSerializer(READ_WRITE_PROPERTY)
	return UserSerializer{BaseSerializer: ser, Name: name, Checker: ck, Reader: rd, Writer: wt}
}



type Getter func()interface{}
type Setter func(interface)


type PropByRefSerializer struct {
	BaseSerializer
	Name         string
	DefaultValue interface{}
	Getter       Getter
	Setter       Setter
}

func (ser *PropByRefSerializer) Read(is *OsgIstream, obj *model.Object) {}

func (ser *PropByRefSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewPropByRefSerializer(name string,def interface{}, gt Getter, st Setter) PropByRefSerializer {
	ser := NewBaseSerializer(READ_WRITE_PROPERTY)
	return PropByRefSerializer{BaseSerializer: ser, Name: name, Getter: gt, Setter: st, DefaultValue: def}
}

type PropByValSerializer struct {
	PropByRefSerializer 
	UseHex       bool
}

func (ser *PropByValSerializer) Read(is *OsgIstream, obj *model.Object) {}

func (ser *PropByValSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewPropByValSerializer(name string,def interface{}, gt Getter, st Setter,hex bool) PropByValSerializer {
	v:=NewPropByRefSerializer(name,def,gt,st)
	return PropByValSerializer{PropByRefSerializer:v,UseHex:hex}
}