package io

import (
	"github.com/flywave/go-osg/model"
	"github.com/ungerik/go3d/mat4"
)

type StringToValue map[string]int
type ValueToString map[int]string

type IntLookup struct {
	StringToValue StringToValue
	ValueToString ValueToString
}

func (lk *IntLookup) Size() int {
	return len(lk.StringToValue)
}

func (lk *IntLookup) Add(str string, val int) {
	_, ok := lk.ValueToString[val]
	if ok {
		return
	}
	lk.ValueToString[val] = str
	lk.StringToValue[str] = val
}

func (lk *IntLookup) Add2(str string, newStr string, val int) {
	_, ok := lk.ValueToString[val]
	if ok {
		return
	}
	lk.ValueToString[val] = str
	lk.StringToValue[newStr] = val
	lk.StringToValue[str] = val
}

func (lk *IntLookup) GetValue(str string) int {
	v, ok := lk.StringToValue[str]
	if ok {
		var val int = int(str[0])
		lk.StringToValue[str] = val
		return val
	}
	return v
}

func (lk *IntLookup) GetString(val int) string {
	s, ok := lk.ValueToString[val]
	if ok {
		s = string(val)
		lk.ValueToString[val] = s
		return s
	}
	return s
}

func NewIntLookup() IntLookup {
	return IntLookup{StringToValue: make(StringToValue), ValueToString: make(ValueToString)}
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

type Serializer interface {
	GetSerializerName() string
	Read(is *OsgIstream, obj *model.Object)
	Write(is *OsgOstream, obj *model.Object)
	GetFirstVersion() int
	SetFirstVersion(int)
	GetLastVersion() int
	SetLastVersion(int)
	SupportsSet() bool
	SupportsGet() bool
	SupportsGetSet() bool
}

type BaseSerializer struct {
	FirstVersion int
	LastVersion  int
	Usage        Usage
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

func (bs *BaseSerializer) GetSerializerName() string {
	return ""
}
func (ser *BaseSerializer) Read(is *OsgIstream, obj *model.Object) {
}

func (ser *BaseSerializer) Write(is *OsgOstream, obj *model.Object) {
}

func (ser *BaseSerializer) GetFirstVersion() int {
	return ser.FirstVersion
}
func (ser *BaseSerializer) SetFirstVersion(v int) {
	ser.FirstVersion = v
}
func (ser *BaseSerializer) GetLastVersion() int {
	return ser.LastVersion
}
func (ser *BaseSerializer) SetLastVersion(v int) {
	ser.LastVersion = v
}

func NewBaseSerializer(usg Usage) BaseSerializer {
	return BaseSerializer{Usage: usg}
}

type Checker func(interface{}) bool
type Reader func(*OsgIstream, interface{})
type Writer func(*OsgOstream, interface{})

type UserSerializer struct {
	BaseSerializer
	Checker Checker
	Rd      Reader
	Wt      Writer
	Name    string
}

func (ser *UserSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {
		ok := false
		is.Read(&ok)
		if ok {
			ser.Rd(is, obj)
		}
	} else {
		if is.MatchString(ser.Name) {
			ser.Rd(is, obj)
		}
	}
}

func (ser *UserSerializer) Writer(is *OsgOstream, obj *model.Object) {

}

func (ser *UserSerializer) GetSerializerName() string {
	return ser.Name
}

func NewUserSerializer(name string, ck Checker, rd Reader, wt Writer) UserSerializer {
	ser := NewBaseSerializer(READ_WRITE_PROPERTY)
	return UserSerializer{BaseSerializer: ser, Name: name, Checker: ck, Rd: rd, Wt: wt}
}

type Getter func(interface{}) interface{}
type Setter func(interface{}, interface{})

type TemplateSerializer struct {
	BaseSerializer
	Name   string
	Getter Getter
	Setter Setter
}

func NewTemplateSerializer(name string, gt Getter, st Setter) TemplateSerializer {
	ser := NewBaseSerializer(READ_WRITE_PROPERTY)
	return TemplateSerializer{BaseSerializer: ser, Name: name, Getter: gt, Setter: st}
}

type PropByValSerializer struct {
	TemplateSerializer
	Prop   interface{}
	UseHex bool
}

func (ser *PropByValSerializer) GetSerializerName() string {
	return ser.Name
}

func (ser *PropByValSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
		if is.MatchString(ser.Name) {
			if ser.UseHex {
				is.Read(ser.Prop)
			}
		}
	}
}

func (ser *PropByValSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewPropByValSerializer(name string, hex bool, gt Getter, st Setter) PropByValSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return PropByValSerializer{TemplateSerializer: ser}
}

type PropByRefSerializer struct {
	PropByValSerializer
}

func (ser *PropByRefSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
		if is.MatchString(ser.Name) {
			is.Read(ser.Prop)
		}
	}
}

func (ser *PropByRefSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewPropByRefSerializer(name string, gt Getter, st Setter) PropByRefSerializer {
	ser := NewPropByValSerializer(name, false, gt, st)
	return PropByRefSerializer{PropByValSerializer: ser}
}

type MatrixSerializer struct {
	TemplateSerializer
	int mat4.T
}

func (ser *MatrixSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *MatrixSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func (ser *MatrixSerializer) GetSerializerName() string {
	return ""
}

func NewMatrixSerializer(name string, gt Getter, st Setter) MatrixSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return MatrixSerializer{TemplateSerializer: ser}
}

type GlenumSerializer struct {
	TemplateSerializer
	Int int
}

func (ser *GlenumSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *GlenumSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func (ser *GlenumSerializer) GetSerializerName() string {
	return ""
}

func NewGlenumSerializer(name string, gt Getter, st Setter) GlenumSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return GlenumSerializer{TemplateSerializer: ser}
}

type StringSerializer struct {
	TemplateSerializer
}

func (ser *StringSerializer) Read(is *OsgIstream, obj *model.Object) {}

func (ser *StringSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func (ser *StringSerializer) GetSerializerName() string {
	return ""
}
func NewStringSerializer(name string, gt Getter, st Setter) StringSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return StringSerializer{TemplateSerializer: ser}
}

type ObjectSerializer struct {
	TemplateSerializer
}

func (ser *ObjectSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *ObjectSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewObjectSerializer(name string, gt Getter, st Setter) ObjectSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return ObjectSerializer{TemplateSerializer: ser}
}

type ImageSerializer struct {
	TemplateSerializer
}

func (ser *ImageSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *ImageSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewImageSerializer(name string, gt Getter, st Setter) ImageSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return ImageSerializer{TemplateSerializer: ser}
}

type EnumSerializer struct {
	TemplateSerializer
	LookUp    IntLookup
	EnumValue *int
}

func (ser *EnumSerializer) Add(str string, val int) {
	ser.LookUp.Add(str, val)
}

func (ser *EnumSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
		if is.MatchString(ser.Name) {
			var str string
			is.Read(str)
			*ser.EnumValue = ser.LookUp.GetValue(str)
		}
	}
}

func (ser *EnumSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewEnumSerializer(name string, gt Getter, st Setter) EnumSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return EnumSerializer{TemplateSerializer: ser, LookUp: NewIntLookup()}
}

type ListSerializer struct {
	TemplateSerializer
}

func (ser *ListSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *ListSerializer) Writer(is *OsgOstream, obj *model.Object) {}

type ConstGetter func() []interface{}

type VectorSerializer struct {
	TemplateSerializer
	ElementType     SerType
	NumElementOnRow uint
}

func (ser *VectorSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *VectorSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewVectorSerializer(name string, ty SerType, nrow uint, gt Getter, st Setter) VectorSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return VectorSerializer{TemplateSerializer: ser, ElementType: ty, NumElementOnRow: nrow}
}

type IsAVectorSerializer struct {
	TemplateSerializer
	ElementType     SerType
	NumElementOnRow uint
}

func (ser *IsAVectorSerializer) Read(is *OsgIstream, obj *model.Object) {
	if is.IsBinary() {

	} else {
	}
}

func (ser *IsAVectorSerializer) Writer(is *OsgOstream, obj *model.Object) {}

func NewIsAVectorSerializer(name string, ty SerType, nrow uint) IsAVectorSerializer {
	ser := NewTemplateSerializer(name, nil, nil)
	return IsAVectorSerializer{TemplateSerializer: ser, ElementType: ty, NumElementOnRow: nrow}
}
