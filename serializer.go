package osg

import "github.com/flywave/go-osg/model"

type StringToValue map[string]int32
type ValueToString map[int32]string

type IntLookup struct {
	StringToValue StringToValue
	ValueToString ValueToString
}

func (lk *IntLookup) Size() int {
	return len(lk.StringToValue)
}

func (lk *IntLookup) Add(str string, val int32) {
	_, ok := lk.ValueToString[val]
	if ok {
		return
	}
	lk.ValueToString[val] = str
	lk.StringToValue[str] = val
}

func (lk *IntLookup) Add2(str string, newStr string, val int32) {
	_, ok := lk.ValueToString[val]
	if ok {
		return
	}
	lk.ValueToString[val] = str
	lk.StringToValue[newStr] = val
	lk.StringToValue[str] = val
}

func (lk *IntLookup) GetValue(str string) int32 {
	v, ok := lk.StringToValue[str]
	if !ok {
		var val int32 = int32(str[0])
		lk.StringToValue[str] = val
		return val
	}
	return v
}

func (lk *IntLookup) GetString(val int32) string {
	s, ok := lk.ValueToString[val]
	if ok {
		s = string(val)
		lk.ValueToString[val] = s
		return s
	}
	return s
}

func NewIntLookup() *IntLookup {
	return &IntLookup{StringToValue: make(StringToValue), ValueToString: make(ValueToString)}
}

type SerType uint32
type Usage uint32

const (
	RWUNDEFINED       SerType = 0
	RWUSER            SerType = 1
	RWOBJECT          SerType = 2
	RWIMAGE           SerType = 3
	RWLIST            SerType = 4
	RWBOOL            SerType = 5
	RWCHAR            SerType = 6
	RWUCHAR           SerType = 7
	RWSHORT           SerType = 8
	RWUSHORT          SerType = 9
	RWINT             SerType = 10
	RWUINT            SerType = 11
	RWFLOAT           SerType = 12
	RWDOUBLE          SerType = 13
	RWVEC2F           SerType = 14
	RWVEC2D           SerType = 15
	RWVEC3F           SerType = 16
	RWVEC3D           SerType = 17
	RWVEC4F           SerType = 18
	RWVEC4D           SerType = 19
	RWQUAT            SerType = 20
	RWPLANE           SerType = 21
	RWMATRIXF         SerType = 22
	RWMATRIXD         SerType = 23
	RWMATRIX          SerType = 24
	RWGLENUM          SerType = 25
	RWSTRING          SerType = 26
	RWENUM            SerType = 27
	RWVEC2B           SerType = 28
	RWVEC2UB          SerType = 29
	RWVEC2S           SerType = 30
	RWVEC2US          SerType = 31
	RWVEC2I           SerType = 32
	RWVEC2UI          SerType = 33
	RWVEC3B           SerType = 34
	RWVEC3UB          SerType = 35
	RWVEC3S           SerType = 36
	RWVEC3US          SerType = 37
	RWVEC3I           SerType = 38
	RWVEC3UI          SerType = 39
	RWVEC4B           SerType = 40
	RWVEC4UB          SerType = 41
	RWVEC4S           SerType = 42
	RWVEC4US          SerType = 43
	RWVEC4I           SerType = 44
	RWVEC4UI          SerType = 45
	RWBOUNDINGBOXF    SerType = 46
	RWBOUNDINGBOXD    SerType = 47
	RWBOUNDINGSPHEREF SerType = 48
	RWBOUNDINGSPHERED SerType = 49
	RWVECTOR          SerType = 50
	RWMAP             SerType = 51

	READWRITEPROPERTY Usage = 1
	GETPROPERTY       Usage = 2
	SETPROPERTY       Usage = 4
	GETSETPROPERTY    Usage = GETPROPERTY | SETPROPERTY
)

type Serializer interface {
	GetSerializerName() string
	Read(is *OsgIstream, obj interface{})
	Write(is *OsgOstream, obj interface{})
	GetFirstVersion() int32
	SetFirstVersion(int32)
	GetLastVersion() int32
	SetLastVersion(int32)
	SupportsSet() bool
	SupportsGet() bool
	SupportsGetSet() bool
	SupportsReadWrite() bool
}

type BaseSerializer struct {
	FirstVersion int32
	LastVersion  int32
	Usage        Usage
}

func (bs *BaseSerializer) SupportsReadWrite() bool {
	return (bs.Usage & READWRITEPROPERTY) != 0
}

func (bs *BaseSerializer) SupportsGetSet() bool {
	return (bs.Usage & GETSETPROPERTY) != 0
}

func (bs *BaseSerializer) SupportsGet() bool {
	return (bs.Usage & GETPROPERTY) != 0
}

func (bs *BaseSerializer) SupportsSet() bool {
	return (bs.Usage & SETPROPERTY) != 0
}

func (bs *BaseSerializer) GetSerializerName() string {
	return ""
}
func (ser *BaseSerializer) Read(is *OsgIstream, obj interface{}) {
}

func (ser *BaseSerializer) Write(is *OsgOstream, obj interface{}) {
}

func (ser *BaseSerializer) GetFirstVersion() int32 {
	return ser.FirstVersion
}
func (ser *BaseSerializer) SetFirstVersion(v int32) {
	ser.FirstVersion = v
}
func (ser *BaseSerializer) GetLastVersion() int32 {
	return ser.LastVersion
}
func (ser *BaseSerializer) SetLastVersion(v int32) {
	ser.LastVersion = v
}

func NewBaseSerializer(usg Usage) *BaseSerializer {
	return &BaseSerializer{Usage: usg, FirstVersion: 0, LastVersion: INT32_MAX}
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

func (ser *UserSerializer) Read(is *OsgIstream, obj interface{}) {
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

func (ser *UserSerializer) Writer(is *OsgOstream, obj interface{}) {

}

func (ser *UserSerializer) GetSerializerName() string {
	return ser.Name
}

func NewUserSerializer(name string, ck Checker, rd Reader, wt Writer) *UserSerializer {
	ser := NewBaseSerializer(READWRITEPROPERTY)
	return &UserSerializer{BaseSerializer: *ser, Name: name, Checker: ck, Rd: rd, Wt: wt}
}

type Getter func(interface{}) interface{}
type Setter func(interface{}, interface{})

type TemplateSerializer struct {
	BaseSerializer
	Name   string
	Getter Getter
	Setter Setter
}

func (ser *TemplateSerializer) GetSerializerName() string {
	return ser.Name
}

func NewTemplateSerializer(name string, gt Getter, st Setter) *TemplateSerializer {
	ser := NewBaseSerializer(READWRITEPROPERTY)
	return &TemplateSerializer{BaseSerializer: *ser, Name: name, Getter: gt, Setter: st}
}

type PropByValSerializer struct {
	TemplateSerializer
	UseHex bool
}

func (ser *PropByValSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		is.Read(ser.Getter(obj))
	} else {
		if is.MatchString(ser.Name) {
			if ser.UseHex {

			} else {
				is.Read(ser.Getter(obj))
			}
		}
	}
}

func (ser *PropByValSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewPropByValSerializer(name string, hex bool, gt Getter, st Setter) *PropByValSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &PropByValSerializer{TemplateSerializer: *ser}
}

type PropByRefSerializer struct {
	PropByValSerializer
}

func (ser *PropByRefSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		is.Read(ser.Getter(obj))
	} else {
		if is.MatchString(ser.Name) {
			is.Read(ser.Getter(obj))
		}
	}
}

func (ser *PropByRefSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewPropByRefSerializer(name string, gt Getter, st Setter) *PropByRefSerializer {
	ser := NewPropByValSerializer(name, false, gt, st)
	return &PropByRefSerializer{PropByValSerializer: *ser}
}

type MatrixSerializer struct {
	TemplateSerializer
}

func (ser *MatrixSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		is.Read(ser.Getter(obj))
	} else {
		if is.MatchString(ser.Name) {
			is.Read(ser.Getter(obj))
		}
	}
}

func (ser *MatrixSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewMatrixSerializer(name string, gt Getter, st Setter) *MatrixSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &MatrixSerializer{TemplateSerializer: *ser}
}

type GlenumSerializer struct {
	TemplateSerializer
}

func (ser *GlenumSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		is.Read(ser.Getter(obj))
	} else {
		if is.MatchString(ser.Name) {
			is.Read(ser.Getter(obj))
		}
	}
}

func (ser *GlenumSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewGlenumSerializer(name string, gt Getter, st Setter) *GlenumSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &GlenumSerializer{TemplateSerializer: *ser}
}

type StringSerializer struct {
	TemplateSerializer
}

func (ser *StringSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		is.Read(ser.Getter(obj))
	} else {
		if is.MatchString(ser.Name) {
			is.Read(ser.Getter(obj))
		}
	}
}

func (ser *StringSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewStringSerializer(name string, gt Getter, st Setter) *StringSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &StringSerializer{TemplateSerializer: *ser}
}

type ObjectSerializer struct {
	TemplateSerializer
}

func (ser *ObjectSerializer) Read(is *OsgIstream, obj interface{}) {
	hasObj := false
	if is.IsBinary() {
		is.Read(&hasObj)
		if hasObj {
			ser.Setter(obj, is.ReadObject(nil))
		}
	} else {
		if is.MatchString(ser.Name) {
			is.Read(&hasObj)
			if hasObj {
				is.Read(is.BEGINBRACKET)
				ser.Setter(obj, is.ReadObject(nil))
				is.Read(is.ENDBRACKET)
			}
		}
	}
}

func (ser *ObjectSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewObjectSerializer(name string, gt Getter, st Setter) *ObjectSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &ObjectSerializer{TemplateSerializer: *ser}
}

type ImageSerializer struct {
	TemplateSerializer
}

func (ser *ImageSerializer) Read(is *OsgIstream, obj interface{}) {
	hasObj := false
	if is.IsBinary() {
		is.Read(&hasObj)
		if hasObj {
			ser.Setter(obj, is.ReadImage(false))
		}
	} else {
		if is.MatchString(ser.Name) {
			is.Read(&hasObj)
			if hasObj {
				is.Read(is.BEGINBRACKET)
				ser.Setter(obj, is.ReadImage(false))
				is.Read(is.ENDBRACKET)
			}
		}
	}
}

func (ser *ImageSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewImageSerializer(name string, gt Getter, st Setter) *ImageSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &ImageSerializer{TemplateSerializer: *ser}
}

type EnumSerializer struct {
	TemplateSerializer
	LookUp    *IntLookup
	EnumValue *int32
}

func (ser *EnumSerializer) Add(str string, val int32) {
	ser.LookUp.Add(str, val)
}

func (ser *EnumSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		is.Read(ser.Getter(obj))
	} else {
		if is.MatchString(ser.Name) {
			var str string
			is.Read(&str)
			ser.Setter(obj, ser.LookUp.GetValue(str))
		}
	}
}

func (ser *EnumSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewEnumSerializer(name string, gt Getter, st Setter) *EnumSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &EnumSerializer{TemplateSerializer: *ser, LookUp: NewIntLookup()}
}

type VectorSerializer struct {
	TemplateSerializer
	ElementType SerType
	Element     interface{}
}

func (ser *VectorSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		size := is.ReadSize()
		switch ser.Element.(type) {
		case *model.PrimitiveSet:
			for i := 0; i < size; i++ {
				ser.Setter(obj, is.ReadPrimitiveSet())
			}
		case *model.Array:
			for i := 0; i < size; i++ {
				ser.Setter(obj, is.ReadArray())
			}
		}
	} else {
		if is.MatchString(ser.Name) {
			size := is.ReadSize()
			is.Read(is.BEGINBRACKET)
			if size > 0 {
				switch ser.Element.(type) {
				case *model.PrimitiveSet:
					for i := 0; i < size; i++ {
						ser.Setter(obj, is.ReadPrimitiveSet())
					}
				case *model.Array:
					for i := 0; i < size; i++ {
						ser.Setter(obj, is.ReadArray())
					}
				}
			}
			is.Read(is.ENDBRACKET)
		}
	}
}

func (ser *VectorSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewVectorSerializer(name string, ty SerType, element interface{}, gt Getter, st Setter) *VectorSerializer {
	ser := NewTemplateSerializer(name, gt, st)
	return &VectorSerializer{TemplateSerializer: *ser, ElementType: ty, Element: element}
}

type IsAVectorSerializer struct {
	TemplateSerializer
	ElementType     SerType
	NumElementOnRow uint
}

func (ser *IsAVectorSerializer) Read(is *OsgIstream, obj interface{}) {
	if is.IsBinary() {
		var size int32
		is.Read(&size)
		vec := ser.genVect(is, size)
		ser.Setter(obj, vec)
	} else {
		if is.MatchString(ser.Name) {
			var size int32
			is.Read(&size)
			vec := ser.genVect(is, size)
			ser.Setter(obj, vec)
		}
	}
}

func (ser *IsAVectorSerializer) genVect(is *OsgIstream, size int32) interface{} {
	switch ser.ElementType {
	case RWUCHAR:
		vec := make([]uint8, int(size), int(size))
		for i := range vec {
			is.Read(&vec[i])
		}
		return vec
	case RWUSHORT:
		vec := make([]uint16, int(size), int(size))
		for i := range vec {
			is.Read(&vec[i])
		}
		return vec
	case RWINT:
		vec := make([]int32, int(size), int(size))
		for i := range vec {
			is.Read(&vec[i])
		}
		return vec
	case RWUINT:
		vec := make([]uint32, int(size), int(size))
		for i := range vec {
			is.Read(&vec[i])
		}
		return vec
	default:
		return nil
	}
}

func (ser *IsAVectorSerializer) Writer(is *OsgOstream, obj interface{}) {}

func NewIsAVectorSerializer(name string, ty SerType, nrow uint, gt Getter, st Setter) *IsAVectorSerializer {
	ser := NewTemplateSerializer(name, nil, nil)
	return &IsAVectorSerializer{TemplateSerializer: *ser, ElementType: ty, NumElementOnRow: nrow}
}
