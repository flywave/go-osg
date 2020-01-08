package osg

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/flywave/go-osg/model"
)

const (
	FileType     string = "Ascii"
	INDENT_VALUE        = 2
)

type OsgOptions struct {
	FileType   string
	Precision  int
	Compressed bool
}

func NewOsgOptions() *OsgOptions {
	return &OsgOptions{FileType: FileType}
}

const (
	READUNKNOWN uint32 = 0
	READSCENE   uint32 = 1
	READIMAGE   uint32 = 2
	READOBJECT  uint32 = 3
)

type OsgIstreamOptions struct {
	OsgOptions
	DbPath            string
	Domain            string
	ForceReadingImage bool
}

func NewOsgIstreamOptions() *OsgIstreamOptions {
	op := NewOsgOptions()
	return &OsgIstreamOptions{OsgOptions: *op}
}

type StreamHeader struct {
	Version       int32
	Type          uint32
	Attributes    int32
	NumDomains    int32
	DomainName    string
	DomainVersion int32
	TypeString    string
	OsgName       string
	OsgVersion    string
}

type OsgIstream struct {
	ArrayMap          map[int32]*model.Array
	IdentifierMap     map[int32]interface{}
	DomainVersionMap  map[string]int32
	FileVersion       int32
	UseSchemaData     bool
	ForceReadingImage bool
	Fields            []string
	In                OsgInputIterator
	Options           *OsgIstreamOptions
	DummyReadObject   *model.Object
	CRLF              CrlfType

	PROPERTY     *model.ObjectProperty
	BEGINBRACKET *model.ObjectMark
	ENDBRACKET   *model.ObjectMark
}

func NewOsgIstream(opt *OsgIstreamOptions) *OsgIstream {
	p := model.NewObjectProperty()
	bb := model.NewObjectMark()
	bb.Name = "{"
	bb.IndentDelta = INDENT_VALUE
	eb := model.NewObjectMark()
	eb.Name = "}"
	eb.IndentDelta = -INDENT_VALUE
	is := &OsgIstream{ArrayMap: make(map[int32]*model.Array), Options: opt, IdentifierMap: make(map[int32]interface{}), DomainVersionMap: make(map[string]int32), PROPERTY: p, BEGINBRACKET: bb, ENDBRACKET: eb}
	if opt.ForceReadingImage {
		is.ForceReadingImage = true
	}
	obj := model.NewObject()
	is.DummyReadObject = obj
	if len(opt.Domain) > 0 {
		domains := strings.Split(opt.Domain, ";")
		for _, str := range domains {
			vals := strings.Split(str, ":")
			if len(vals) > 1 {
				v, _ := strconv.ParseInt(vals[1], 10, 32)
				is.DomainVersionMap[vals[0]] = int32(v)
			}
		}
	}
	return is
}

func (is *OsgIstream) IsBinary() bool {
	return is.In.IsBinary()
}

func (is *OsgIstream) MatchString(str string) bool {
	return is.In.MatchString(str)
}

func (is *OsgIstream) Read(inter interface{}) {
	switch val := inter.(type) {
	case *bool:
		is.In.ReadBool(val)
		break
	case *int8:
		is.In.ReadChar(val)
		break
	case *uint8:
		is.In.ReadUChar(val)
		break
	case *int16:
		is.In.ReadShort(val)
		break
	case *uint16:
		is.In.ReadUShort(val)
		break
	case *int:
		var t int32
		is.In.ReadInt(&t)
		*val = int(t)
		break
	case *int32:
		is.In.ReadInt((*int32)(val))
		break
	case *uint:
		var t uint32
		is.In.ReadUInt(&t)
		*val = uint(t)
		break
	case *uint32:
		is.In.ReadUInt((*uint32)(val))
		break
	case *int64:
		is.In.ReadLong(val)
		break
	case *uint64:
		is.In.ReadULong(val)
		break
	case *float32:
		is.In.ReadFloat(val)
		break
	case *float64:
		is.In.ReadDouble(val)
		break
	case *[2]float32:
		is.In.ReadFloat(&val[0])
		is.In.ReadFloat(&val[1])
		break
	case *[2]float64:
		is.In.ReadDouble(&val[0])
		is.In.ReadDouble(&val[1])
		break
	case *[3]float32:
		is.In.ReadFloat(&val[0])
		is.In.ReadFloat(&val[1])
		is.In.ReadFloat(&val[2])
		break
	case *[3]float64:
		is.In.ReadDouble(&val[0])
		is.In.ReadDouble(&val[1])
		is.In.ReadDouble(&val[2])
		break
	case *[4]float32:
		is.In.ReadFloat(&val[0])
		is.In.ReadFloat(&val[1])
		is.In.ReadFloat(&val[2])
		is.In.ReadFloat(&val[3])
		break
	case *[4]float64:
		is.In.ReadDouble(&val[0])
		is.In.ReadDouble(&val[1])
		is.In.ReadDouble(&val[2])
		is.In.ReadDouble(&val[3])
		break
	case *[2]int32:
		is.In.ReadInt(&val[0])
		is.In.ReadInt(&val[1])
		break
	case *[2]int64:
		is.In.ReadLong(&val[0])
		is.In.ReadLong(&val[1])
		break
	case *[3]int32:
		is.In.ReadInt(&val[0])
		is.In.ReadInt(&val[1])
		is.In.ReadInt(&val[2])
		break
	case *[3]int64:
		is.In.ReadLong(&val[0])
		is.In.ReadLong(&val[1])
		is.In.ReadLong(&val[2])
		break
	case *[4]int32:
		is.In.ReadInt(&val[0])
		is.In.ReadInt(&val[1])
		is.In.ReadInt(&val[2])
		is.In.ReadInt(&val[3])
		break
	case *[4]int64:
		is.In.ReadLong(&val[0])
		is.In.ReadLong(&val[1])
		is.In.ReadLong(&val[2])
		is.In.ReadLong(&val[3])
		break
	case *[2]uint32:
		is.In.ReadUInt(&val[0])
		is.In.ReadUInt(&val[1])
		break
	case *[2]uint64:
		is.In.ReadULong(&val[0])
		is.In.ReadULong(&val[1])
		break
	case *[3]uint32:
		is.In.ReadUInt(&val[0])
		is.In.ReadUInt(&val[1])
		is.In.ReadUInt(&val[2])
		break
	case *[3]uint64:
		is.In.ReadULong(&val[0])
		is.In.ReadULong(&val[1])
		is.In.ReadULong(&val[2])
		break
	case *[4]uint32:
		is.In.ReadUInt(&val[0])
		is.In.ReadUInt(&val[1])
		is.In.ReadUInt(&val[2])
		is.In.ReadUInt(&val[3])
		break
	case *[4]uint64:
		is.In.ReadULong(&val[0])
		is.In.ReadULong(&val[1])
		is.In.ReadULong(&val[2])
		is.In.ReadULong(&val[3])
		break
	case *[2]uint8:
		is.In.ReadUChar(&val[0])
		is.In.ReadUChar(&val[1])
		break
	case *[2]uint16:
		is.In.ReadUShort(&val[0])
		is.In.ReadUShort(&val[1])
		break
	case *[3]uint8:
		is.In.ReadUChar(&val[0])
		is.In.ReadUChar(&val[1])
		is.In.ReadUChar(&val[2])
		break
	case *[3]uint16:
		is.In.ReadUShort(&val[0])
		is.In.ReadUShort(&val[1])
		is.In.ReadUShort(&val[2])
		break
	case *[4]uint8:
		is.In.ReadUChar(&val[0])
		is.In.ReadUChar(&val[1])
		is.In.ReadUChar(&val[2])
		is.In.ReadUChar(&val[3])
		break
	case *[4]uint16:
		is.In.ReadUShort(&val[0])
		is.In.ReadUShort(&val[1])
		is.In.ReadUShort(&val[2])
		is.In.ReadUShort(&val[3])
		break
	case *[2]int8:
		is.In.ReadChar(&val[0])
		is.In.ReadChar(&val[1])
		break
	case *[2]int16:
		is.In.ReadShort(&val[0])
		is.In.ReadShort(&val[1])
		break
	case *[3]int8:
		is.In.ReadChar(&val[0])
		is.In.ReadChar(&val[1])
		is.In.ReadChar(&val[2])
		break
	case *[3]int16:
		is.In.ReadShort(&val[0])
		is.In.ReadShort(&val[1])
		is.In.ReadShort(&val[2])
		break
	case *[4]int8:
		is.In.ReadChar(&val[0])
		is.In.ReadChar(&val[1])
		is.In.ReadChar(&val[2])
		is.In.ReadChar(&val[3])
		break
	case *[4]int16:
		is.In.ReadShort(&val[0])
		is.In.ReadShort(&val[1])
		is.In.ReadShort(&val[2])
		is.In.ReadShort(&val[3])
		break
	case *[4][4]float32:
		is.ReadMatrix4f(val)
		break
	case *[4][4]float64:
		is.ReadMatrix4d(val)
		break
	case *[4][4]int32:
		is.Read(&val[0])
		is.Read(&val[1])
		is.Read(&val[2])
		is.Read(&val[3])
		break
	case *[4][4]int64:
		is.Read(&val[0])
		is.Read(&val[1])
		is.Read(&val[2])
		is.Read(&val[3])
		break
	case *string:
		st := is.In.ReadString()
		*val = st
		break
	case *model.ObjectGlenum:
		is.In.ReadGlenum(val)
		break
	case *model.ObjectProperty:
		is.In.ReadProperty(val)
		break
	case *model.ObjectMark:
		is.In.ReadMark(val)
		break
	}
}

func (is *OsgIstream) ReadCharArray(size int) []byte {
	return is.In.ReadCharArray(size)
}

func (is *OsgIstream) ReadWrappedString(str *string) {
	is.In.ReadWrappedString(str)
}

func (is *OsgIstream) ReadString() string {
	str := is.In.ReadString()
	return str
}

type imagedata struct {
	Origin         int32
	S              int32
	T              int32
	R              int32
	InternalFormat int32

	PixelFormat int32
	DataType    int32
	Packing     int32
	Mode        int32
	Size        int32
	Data        []byte
}

func (is *OsgIstream) ReadMatrix4f(mat *[4][4]float32) {
	is.Read(is.BEGINBRACKET)
	is.Read(&mat[0])
	is.Read(&mat[1])
	is.Read(&mat[2])
	is.Read(&mat[3])
	is.Read(is.ENDBRACKET)
}

func (is *OsgIstream) ReadMatrix4d(mat *[4][4]float64) {
	is.Read(is.BEGINBRACKET)
	is.Read(&mat[0])
	is.Read(&mat[1])
	is.Read(&mat[2])
	is.Read(&mat[3])
	is.Read(is.ENDBRACKET)
}

func (is *OsgIstream) ReadArray() *model.Array {
	is.PROPERTY.Name = "ArrayID"
	is.Read(is.PROPERTY)
	var id int32
	is.Read(&id)
	ay, ok := is.ArrayMap[id]
	if ok {
		return ay
	}
	ty := model.ObjectProperty{Name: "ArrayType", MapProperty: true}
	is.Read(&ty)
	var size int32
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	var arry *model.Array
	switch ty.Value {
	case model.IDBYTEARRAY:
		{
			arry = model.NewArray(model.ByteArrayType, model.GLBYTE, 1)
			data := make([]int8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadChar(&data[i])
			}
			arry.Data = data
		}
	case model.IDUBYTEARRAY:
		{
			arry = model.NewArray(model.UByteArrayType, model.GLUNSIGNEDBYTE, 1)
			data := make([]byte, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUChar(&data[i])
			}
			arry.Data = data
		}
	case model.IDSHORTARRAY:
		{
			arry = model.NewArray(model.ShortArrayType, model.GLSHORT, 1)
			data := make([]int16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadShort(&data[i])
			}
			arry.Data = data
		}
	case model.IDUSHORTARRAY:
		{
			arry = model.NewArray(model.UShortArrayType, model.GLUNSIGNEDSHORT, 1)
			data := make([]uint16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUShort(&data[i])
			}
			arry.Data = data
		}
	case model.IDINTARRAY:
		{
			arry = model.NewArray(model.IntArrayType, model.GLINT, 1)
			data := make([]int32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadInt(&data[i])
			}
			arry.Data = data
		}
	case model.IDUINTARRAY:
		{
			arry = model.NewArray(model.UIntArrayType, model.GLUNSIGNEDINT, 1)
			data := make([]uint32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUInt(&data[i])
			}
			arry.Data = data
		}
	case model.IDFLOATARRAY:
		{
			arry = model.NewArray(model.FloatArrayType, model.GLFLOAT, 1)
			data := make([]float32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadFloat(&data[i])
			}
			arry.Data = data
		}
	case model.IDDOUBLEARRAY:
		{
			arry = model.NewArray(model.DoubleArrayType, model.GLDOUBLE, 1)
			data := make([]float64, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadDouble(&data[i])
			}
			arry.Data = data
		}
	case model.IDVEC2BARRAY:
		{
			arry = model.NewArray(model.Vec2bArrayType, model.GLBYTE, 2)
			data := make([][2]int8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadChar(&data[i][0])
				is.In.ReadChar(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3BARRAY:
		{
			arry = model.NewArray(model.Vec3bArrayType, model.GLBYTE, 3)
			data := make([][3]int8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadChar(&data[i][0])
				is.In.ReadChar(&data[i][1])
				is.In.ReadChar(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4BARRAY:
		{
			arry = model.NewArray(model.Vec4bArrayType, model.GLBYTE, 4)
			data := make([][4]int8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadChar(&data[i][0])
				is.In.ReadChar(&data[i][1])
				is.In.ReadChar(&data[i][2])
				is.In.ReadChar(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2UBARRAY:
		{
			arry = model.NewArray(model.Vec2ubArrayType, model.GLUNSIGNEDBYTE, 2)
			data := make([][2]uint8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUChar(&data[i][0])
				is.In.ReadUChar(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3UBARRAY:
		{
			arry = model.NewArray(model.Vec3ubArrayType, model.GLUNSIGNEDBYTE, 3)
			data := make([][3]uint8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUChar(&data[i][0])
				is.In.ReadUChar(&data[i][1])
				is.In.ReadUChar(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4UBARRAY:
		{
			arry = model.NewArray(model.Vec4ubArrayType, model.GLUNSIGNEDBYTE, 4)
			data := make([][4]uint8, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUChar(&data[i][0])
				is.In.ReadUChar(&data[i][1])
				is.In.ReadUChar(&data[i][2])
				is.In.ReadUChar(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2SARRAY:
		{
			arry = model.NewArray(model.Vec2sArrayType, model.GLSHORT, 2)
			data := make([][2]int16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadShort(&data[i][0])
				is.In.ReadShort(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3SARRAY:
		{
			arry = model.NewArray(model.Vec3sArrayType, model.GLSHORT, 3)
			data := make([][3]int16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadShort(&data[i][0])
				is.In.ReadShort(&data[i][1])
				is.In.ReadShort(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4SARRAY:
		{
			arry = model.NewArray(model.Vec4sArrayType, model.GLSHORT, 4)
			data := make([][4]int16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadShort(&data[i][0])
				is.In.ReadShort(&data[i][1])
				is.In.ReadShort(&data[i][2])
				is.In.ReadShort(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2USARRAY:
		{
			arry = model.NewArray(model.Vec2usArrayType, model.GLUNSIGNEDSHORT, 2)
			data := make([][2]uint16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUShort(&data[i][0])
				is.In.ReadUShort(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3USARRAY:
		{
			arry = model.NewArray(model.Vec3usArrayType, model.GLUNSIGNEDSHORT, 3)
			data := make([][3]uint16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUShort(&data[i][0])
				is.In.ReadUShort(&data[i][1])
				is.In.ReadUShort(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4USARRAY:
		{
			arry = model.NewArray(model.Vec4usArrayType, model.GLUNSIGNEDSHORT, 4)
			data := make([][4]uint16, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUShort(&data[i][0])
				is.In.ReadUShort(&data[i][1])
				is.In.ReadUShort(&data[i][2])
				is.In.ReadUShort(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2ARRAY:
		{
			arry = model.NewArray(model.Vec2ArrayType, model.GLFLOAT, 2)
			data := make([][2]float32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadFloat(&data[i][0])
				is.In.ReadFloat(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3ARRAY:
		{
			arry = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
			data := make([][3]float32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadFloat(&data[i][0])
				is.In.ReadFloat(&data[i][1])
				is.In.ReadFloat(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4ARRAY:
		{
			arry = model.NewArray(model.Vec4ArrayType, model.GLFLOAT, 4)
			data := make([][4]float32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadFloat(&data[i][0])
				is.In.ReadFloat(&data[i][1])
				is.In.ReadFloat(&data[i][2])
				is.In.ReadFloat(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2DARRAY:
		{
			arry = model.NewArray(model.Vec2dArrayType, model.GLDOUBLE, 2)
			data := make([][2]float64, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadDouble(&data[i][0])
				is.In.ReadDouble(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3DARRAY:
		{
			arry = model.NewArray(model.Vec3dArrayType, model.GLDOUBLE, 3)
			data := make([][3]float64, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadDouble(&data[i][0])
				is.In.ReadDouble(&data[i][1])
				is.In.ReadDouble(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4DARRAY:
		{
			arry = model.NewArray(model.Vec4dArrayType, model.GLDOUBLE, 4)
			data := make([][4]float64, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadDouble(&data[i][0])
				is.In.ReadDouble(&data[i][1])
				is.In.ReadDouble(&data[i][2])
				is.In.ReadDouble(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2IARRAY:
		{
			arry = model.NewArray(model.Vec2iArrayType, model.GLINT, 2)
			data := make([][2]int32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadInt(&data[i][0])
				is.In.ReadInt(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3IARRAY:
		{
			arry = model.NewArray(model.Vec3iArrayType, model.GLINT, 3)
			data := make([][3]int32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadInt(&data[i][0])
				is.In.ReadInt(&data[i][1])
				is.In.ReadInt(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4IARRAY:
		{
			arry = model.NewArray(model.Vec4iArrayType, model.GLINT, 4)
			data := make([][4]int32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadInt(&data[i][0])
				is.In.ReadInt(&data[i][1])
				is.In.ReadInt(&data[i][2])
				is.In.ReadInt(&data[i][3])
			}
			arry.Data = data
		}
	case model.IDVEC2UIARRAY:
		{
			arry = model.NewArray(model.Vec2uiArrayType, model.GLUNSIGNEDINT, 2)
			data := make([][2]uint32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUInt(&data[i][0])
				is.In.ReadUInt(&data[i][1])
			}
			arry.Data = data
		}
	case model.IDVEC3UIARRAY:
		{
			arry = model.NewArray(model.Vec3uiArrayType, model.GLUNSIGNEDINT, 3)
			data := make([][3]uint32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUInt(&data[i][0])
				is.In.ReadUInt(&data[i][1])
				is.In.ReadUInt(&data[i][2])
			}
			arry.Data = data
		}
	case model.IDVEC4UIARRAY:
		{
			arry = model.NewArray(model.Vec3uiArrayType, model.GLUNSIGNEDINT, 4)
			data := make([][4]uint32, size, size)
			for i := 0; i < int(size); i++ {
				is.In.ReadUInt(&data[i][0])
				is.In.ReadUInt(&data[i][1])
				is.In.ReadUInt(&data[i][2])
				is.In.ReadUInt(&data[i][3])
			}
			arry.Data = data
		}
	}
	is.Read(is.ENDBRACKET)
	return arry
}

func (is *OsgIstream) ReadPrimitiveSet() interface{} {
	if is.FileVersion >= 112 {
		obj := is.ReadObject(nil)
		return obj
	} else {
		ty := model.ObjectProperty{Name: "PrimitiveType", Value: 0, MapProperty: true}
		md := model.ObjectProperty{Name: "PrimitiveType", Value: 0, MapProperty: true}
		var numInstances, first, count, size int32
		is.Read(&ty)
		is.Read(&md)
		if is.FileVersion > 96 {
			is.Read(&numInstances)
		}
		switch ty.Value {
		case model.IDDRAWARRAYS:
			is.Read(&first)
			is.Read(&count)
			da := model.NewDrawArrays()
			da.Mode = md.Value
			da.First = first
			da.Count = count
			da.NumInstances = numInstances
			return da
		case model.IDDRAWARRAYLENGTH:
			is.Read(&first)
			is.Read(&size)
			is.Read(is.BEGINBRACKET)
			dl := model.NewDrawArrayLengths()
			dl.Mode = md.Value
			dl.First = first
			var value int32
			for i := 0; i < int(size); i++ {
				is.Read(&value)
				dl.Data = append(dl.Data, value)
			}
			is.Read(is.ENDBRACKET)
			dl.NumInstances = numInstances
			return dl
		case model.IDDRAWELEMENTSUBYTE:
			is.Read(&size)
			is.Read(is.BEGINBRACKET)
			de := model.NewDrawElementsUByte()
			var value uint8
			for i := 0; i < int(size); i++ {
				is.Read(&value)
				de.Data = append(de.Data, value)
			}
			is.Read(is.ENDBRACKET)
			de.NumInstances = numInstances
			return de
		case model.IDDRAWELEMENTSUSHORT:
			is.Read(&size)
			is.Read(is.BEGINBRACKET)
			ds := model.NewDrawElementsUShort()
			var value uint16
			for i := 0; i < int(size); i++ {
				is.Read(&value)
				ds.Data = append(ds.Data, value)
			}
			is.Read(is.ENDBRACKET)
			ds.NumInstances = numInstances
			return &ds
		case model.IDDRAWELEMENTSUINT:
			is.Read(&size)
			is.Read(is.BEGINBRACKET)
			duint := model.NewDrawElementsUInt()
			var value uint32
			for i := 0; i < int(size); i++ {
				is.Read(&value)
				duint.Data = append(duint.Data, value)
			}
			is.Read(is.ENDBRACKET)
			duint.NumInstances = numInstances
			return duint
		}
	}
	return nil
}

func (is *OsgIstream) ReadImage(readFromExternal bool) *model.Image {
	className := "osg::Image"
	var name string
	var id int32 = 0
	var writeHint int32 = 0
	var decision int32 = model.IMAGEEXTERNAL
	imgdata := imagedata{}
	loadedFromCache := false
	opts := OsgIstreamOptions{}
	var img *model.Image
	if is.FileVersion > 94 {
		is.PROPERTY.Name = "ClassName"
		is.Read(is.PROPERTY)
		is.Read(&className)
	}
	is.PROPERTY.Name = "UniqueID"
	is.Read(is.PROPERTY)
	is.Read(&id)

	idfy, ok := is.IdentifierMap[id]
	if ok {
		return idfy.(*model.Image)
	}
	is.PROPERTY.Name = "FileName"
	is.Read(is.PROPERTY)
	is.ReadWrappedString(&name)

	is.PROPERTY.Name = "WriteHint"
	is.Read(is.PROPERTY)
	is.Read(&writeHint)
	is.Read(&decision)

	switch decision {
	case model.IMAGEINLINEDATA:
		if is.IsBinary() {
			is.Read(&imgdata.Origin)
			is.Read(&imgdata.S)
			is.Read(&imgdata.T)
			is.Read(&imgdata.R)
			is.Read(&imgdata.InternalFormat)
			is.Read(&imgdata.PixelFormat)
			is.Read(&imgdata.DataType)
			is.Read(&imgdata.Packing)
			is.Read(&imgdata.Mode)
			is.Read(&imgdata.Size)
			if imgdata.Size > 0 {
				var numMipmaps uint32 = 0
				is.Read(&numMipmaps)
				imgdata.Data = is.ReadCharArray(int(imgdata.Size))
			}
		} else {
			is.PROPERTY.Name = "Origin"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.Origin)
			is.PROPERTY.Name = "Size"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.S)
			is.Read(&imgdata.T)
			is.Read(&imgdata.R)
			is.PROPERTY.Name = "InternalTextureFormat"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.InternalFormat)
			is.PROPERTY.Name = "PixelFormat"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.PixelFormat)
			is.PROPERTY.Name = "DataType"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.DataType)
			is.PROPERTY.Name = "Packing"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.Packing)
			is.PROPERTY.Name = "AllocationMode"
			is.Read(is.PROPERTY)
			is.Read(&imgdata.Mode)
			is.PROPERTY.Name = "Data"
			is.Read(is.PROPERTY)
			is.ReadSize() //	// levelSize :=is.ReadSize()-1
			is.Read(is.BEGINBRACKET)
			var encodedData string
			is.ReadWrappedString(&encodedData)
			d, e := base64.StdEncoding.DecodeString(encodedData)
			if e == nil {
				imgdata.Data = d
			}
			is.Read(is.ENDBRACKET)
		}
		img = model.NewImage()
		img.Origin = imgdata.Origin
		img.S = imgdata.S
		img.T = imgdata.T
		img.R = imgdata.R
		img.InternalTextureFormat = imgdata.InternalFormat
		img.PixelFormat = imgdata.PixelFormat
		img.DataType = imgdata.DataType
		img.Packing = imgdata.Packing
		img.Data = ([]uint8)(imgdata.Data)
		img.AllocationMode = model.USENEWDELETE
		readFromExternal = false
		break
	case model.IMAGEINLINEFILE:
		if is.IsBinary() {
			size := is.ReadSize()
			if size > 0 {
				dt := is.ReadCharArray(size)
				rw := getReaderWriter()
				if rw != nil {
					sub := strings.Split(name, ".")
					opts.FileType = sub[len(sub)-1]
					buf := bytes.NewBuffer(dt)
					img = rw.ReadImageWithReader(buf, &opts).GetImage()
				}
			}
		}
		break
	case model.IMAGEEXTERNAL:
		break
	case model.IMAGEWRITEOUT:
		break
	default:
		break
	}
	if readFromExternal && name != "" {
		rw := getReaderWriter()
		img = rw.ReadImage(name, &opts).GetImage()
	}
	if loadedFromCache {
		img2 := is.ReadObjectFields("osg::object", id, img)
		return img2.(*model.Image)
	} else {
		img2 := is.ReadObjectFields("osg::object", id, img)
		img = img2.(*model.Image)
		img.Name = name
		img.WriteHint = writeHint
		is.IdentifierMap[id] = img
		return img
	}
}

func (is *OsgIstream) ReadObject(obj interface{}) interface{} {
	cls := is.ReadString()
	if cls == "NULL" {
		return nil
	}
	is.Read(is.BEGINBRACKET)
	is.PROPERTY.Name = "UniqueID"
	is.Read(is.PROPERTY)
	var id uint32
	is.Read(&id)
	v, ok := is.IdentifierMap[int32(id)]
	if ok {
		is.AdvanceToCurrentEndBracket()
		return v
	}
	obj = is.ReadObjectFields(cls, int32(id), obj)
	is.AdvanceToCurrentEndBracket()
	return obj
}

func (is *OsgIstream) ReadObjectFields(className string, id int32, obj interface{}) interface{} {
	wrap := GetObjectWrapperManager().FindWrap(className)
	if wrap == nil {
		return nil
	}
	ver := is.GetFileVersion(wrap.Domain)
	if obj == nil {
		obj = wrap.CreateInstanceFunc()
	}
	is.IdentifierMap[id] = obj
	for _, ass := range wrap.Associates {
		if ass.FirstVersion > ver {
			continue
		}
		if ver <= ass.LastVersion {
			asswrap := GetObjectWrapperManager().FindWrap(ass.Name)
			if asswrap == nil {
				continue
			}
			is.Fields = append(is.Fields, asswrap.Name)
			asswrap.Read(is, obj)
			is.Fields = is.Fields[:len(is.Fields)-1]
		}
	}
	return obj
}

func (is *OsgIstream) ReadSize() int {
	var size uint32
	is.Read(&size)
	return int(size)
}

func (is *OsgIstream) GetFileVersion(domain string) int32 {
	if len(domain) == 0 {
		return is.FileVersion
	}
	v, ok := is.DomainVersionMap[domain]
	if ok {
		return v
	}
	return 0
}

func (is *OsgIstream) AdvanceToCurrentEndBracket() {
	is.In.AdvanceToCurrentEndBracket()
}

func (is *OsgIstream) Start(iter OsgInputIterator) (uint32, error) {
	is.In = iter
	is.Fields = []string{}
	is.Fields = append(is.Fields, "Start")
	tp := READUNKNOWN
	if iter == nil {
		return tp, errors.New("OsgInputIterator is nil")
	}
	iter.SetInputSteam(is)
	header := StreamHeader{}
	if iter.IsBinary() {
		is.Read(&header.Type)
		is.Read(&header.Version)
		is.Read(&header.Attributes)
		if header.Attributes&0x4 > 0 {
			is.In.SetSupportBinaryBrackets(true)
		}
		if header.Attributes&0x2 > 0 {
			is.UseSchemaData = true
		}
		if header.Attributes&0x1 > 0 {
			is.Read(&header.NumDomains)
			var i int32
			for i = 0; i < header.NumDomains; i++ {
				is.Read(&header.DomainName)
				is.Read(&header.DomainVersion)
				is.DomainVersionMap[header.DomainName] = header.DomainVersion
			}
		}
	} else {
		is.Read(&header.TypeString)
		if header.TypeString == "Scene" {
			header.Type = READSCENE
		} else if header.TypeString == "Image" {
			header.Type = READIMAGE
		} else if header.TypeString == "Object" {
			header.Type = READOBJECT
		}
		is.PROPERTY.Name = "#Version"
		is.Read(is.PROPERTY)
		is.Read(&header.Version)
		is.Read(header.Version)
		is.PROPERTY.Name = "#Generator"
		is.Read(is.PROPERTY)
		is.Read(&header.OsgName)
		is.Read(&header.OsgVersion)
		for {
			if is.MatchString("#CustomDomain") {
				header.DomainName = ""
				is.Read(&header.DomainName)
				is.Read(&header.DomainVersion)
				is.DomainVersionMap[header.DomainName] = header.DomainVersion
			} else {
				break
			}
		}
	}
	is.FileVersion = header.Version
	l := len(is.Fields)
	is.Fields = is.Fields[:l-1]
	return header.Type, nil
}

func (is *OsgIstream) Decompress() error {
	if !is.IsBinary() {
		return nil
	}
	is.Fields = []string{}
	compressorName := is.ReadString()
	if compressorName != "0" {
		is.Fields = append(is.Fields, compressorName)
		compressor := GetObjectWrapperManager().FindCompressor(compressorName)
		if compressor == nil {
			return errors.New("inputstream: Failed to decompress stream, No such compressor.")
		}
		src, e := compressor.DeCompress(is.In.GetIterator())
		if e != nil {
			return e
		}
		bufReader := bytes.NewBuffer(src)
		is.In.SetIterator(bufio.NewReader(bufReader))
		is.Fields = is.Fields[:len(is.Fields)-1]
	}
	if is.UseSchemaData {
		is.Fields = append(is.Fields, "SchemaData")
	}
	return nil
}

func trimEnclosingSpaces(str string) string {
	if str == "" {
		return str
	}
	return strings.TrimSpace(str)
}

func (is *OsgIstream) ReadSchema() {
	schem := bytes.NewBufferString(is.ReadString())
	rd := bufio.NewReader(schem)
	for {
		l, prx, err := rd.ReadLine()
		if !prx && err != nil {
			break
		}
		vs := strings.Split(string(l), "=")
		if len(vs) < 2 {
			continue
		}
		is.SetWrapperSchema(trimEnclosingSpaces(vs[0]), trimEnclosingSpaces(vs[1]))
	}
}

func (is *OsgIstream) SetWrapperSchema(name string, prop string) {
	wrap := GetObjectWrapperManager().FindWrap(name)
	if wrap == nil {
		return
	}
	var methods []string
	var types []SerType

	schema := strings.Split(prop, " ")
	for _, str := range schema {
		keyAndValue := strings.Split(str, ":")
		methods = append(methods, keyAndValue[0])
		if len(keyAndValue) > 1 {
			v, _ := strconv.ParseInt(keyAndValue[1], 10, 32)
			types = append(types, SerType(v))
		} else {
			types = append(types, RWUNDEFINED)
		}
	}
	wrap.ReadSchema(methods, types)
}

func (is *OsgIstream) ResetSchema() {
	for _, v := range GetObjectWrapperManager().Wraps {
		v.ResetSchema()
	}
}
