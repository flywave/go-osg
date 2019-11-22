package osg

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/flywave/go-osg/model"
)

const (
	OPENSCENEGRAPHSOVERSION = 160
)

type CrlfType struct{}

type OsgOstreamOptions struct {
	OsgOptions
	UseRobustBinaryFormat bool
	CompressorName        string
	WriteImageHint        string
	Domains               string
	TargetFileVersion     string
}

const (
	WRITEUNKNOWN int32 = 0
	WRITESCENE   int32 = 1
	WRITEIMAGE   int32 = 2
	WRITEOBJECT  int32 = 3

	WRITEUSEIMAGEHINT int32 = 0
	WRITEUSEEXTERNAL  int32 = 1
	WRITEINLINEDATA   int32 = 2
	WRITEINLINEFILE   int32 = 3
	WRITEEXTERNALFILE int32 = 4
)

type OsgOstream struct {
	ArrayMap  map[*model.Array]int32
	ObjectMap map[model.ObjectInterface]int32

	DomainVersionMap      map[string]int32
	WriteImageHint        int32
	UseSchemaData         bool
	UseRobustBinaryFormat bool

	InbuiltSchemaMap map[string]string
	Fields           []string
	SchemaName       string
	CompressorName   string

	Data              []byte
	CompressSource    *bufio.Writer
	UseCompressSource bool
	Options           *OsgOstreamOptions
	Out               OsgOutputIterator
	FileWriter        io.Writer

	FileVersion       int32
	ForceReadingImage bool

	TargetFileVersion int32

	PROPERTY     *model.ObjectProperty
	BEGINBRACKET *model.ObjectMark
	ENDBRACKET   *model.ObjectMark
	CRLF         *CrlfType
}

func NewOsgOstream(opts *OsgOstreamOptions) OsgOstream {
	p := model.NewObjectProperty()
	bb := model.NewObjectMark()
	bb.Name = "{"
	bb.IndentDelta = INDENT_VALUE
	eb := model.NewObjectMark()
	bb.Name = "}"
	bb.IndentDelta = -INDENT_VALUE
	osg := OsgOstream{PROPERTY: &p, BEGINBRACKET: &bb, ENDBRACKET: &eb, CRLF: &CrlfType{}, TargetFileVersion: OPENSCENEGRAPHSOVERSION, UseRobustBinaryFormat: true, UseSchemaData: false}
	osg.Options = opts
	if !opts.UseRobustBinaryFormat {
		osg.UseRobustBinaryFormat = false
	}
	if opts.CompressorName != "" {
		osg.CompressorName = opts.CompressorName
	}

	if opts.WriteImageHint != "" {
		if opts.WriteImageHint == "IncludeData" {
			osg.WriteImageHint = WRITEINLINEDATA
		} else if opts.WriteImageHint == "IncludeFile" {
			osg.WriteImageHint = WRITEINLINEFILE
		} else if opts.WriteImageHint == "UseExternal" {
			osg.WriteImageHint = WRITEUSEEXTERNAL
		} else if opts.WriteImageHint == "WriteOut" {
			osg.WriteImageHint = WRITEEXTERNALFILE
		}
	}
	if opts.Domains != "" {
		ds := strings.Split(opts.Domains, ";")
		for _, str := range ds {
			kv := strings.Split(str, ":")
			if len(kv) > 1 {
				v, _ := strconv.ParseInt(kv[1], 10, 32)
				osg.DomainVersionMap[kv[0]] = int32(v)
			}
		}
	}

	if opts.TargetFileVersion != "" {
		v, _ := strconv.ParseInt(opts.TargetFileVersion, 10, 32)
		if v > 0 && v <= OPENSCENEGRAPHSOVERSION {
			osg.TargetFileVersion = int32(v)
		}
	}

	if osg.TargetFileVersion < 99 {
		osg.UseRobustBinaryFormat = false
	}
	return osg
}

func (os *OsgOstream) Write(inter interface{}) {
	switch val := inter.(type) {
	case bool:
		os.Out.WriteBool(val)
		break
	case int8:
		os.Out.WriteChar(val)
		break
	case uint8:
		os.Out.WriteUChar(val)
		break
	case int16:
		os.Out.WriteShort(val)
		break
	case uint16:
		os.Out.WriteUShort(val)
		break
	case int32:
		os.Out.WriteInt(val)
		break
	case uint32:
		os.Out.WriteUInt(val)
		break
	case int64:
		os.Out.WriteLong(val)
		break
	case uint64:
		os.Out.WriteULong(val)
		break
	case float32:
		os.Out.WriteFloat(val)
		break
	case float64:
		os.Out.WriteDouble(val)
		break
	case *[2]float32:
		os.Out.WriteFloat(val[0])
		os.Out.WriteFloat(val[1])
		break
	case *[2]float64:
		os.Out.WriteDouble(val[0])
		os.Out.WriteDouble(val[1])
		break
	case *[3]float32:
		os.Out.WriteFloat(val[0])
		os.Out.WriteFloat(val[1])
		os.Out.WriteFloat(val[2])
		break
	case *[3]float64:
		os.Out.WriteDouble(val[0])
		os.Out.WriteDouble(val[1])
		os.Out.WriteDouble(val[2])
		break
	case *[4]float32:
		os.Out.WriteFloat(val[0])
		os.Out.WriteFloat(val[1])
		os.Out.WriteFloat(val[2])
		os.Out.WriteFloat(val[3])
		break
	case *[4]float64:
		os.Out.WriteDouble(val[0])
		os.Out.WriteDouble(val[1])
		os.Out.WriteDouble(val[2])
		os.Out.WriteDouble(val[3])
		break
	case *[2]int32:
		os.Out.WriteInt(val[0])
		os.Out.WriteInt(val[1])
		break
	case *[2]int64:
		os.Out.WriteLong(val[0])
		os.Out.WriteLong(val[1])
		break
	case *[3]int32:
		os.Out.WriteInt(val[0])
		os.Out.WriteInt(val[1])
		os.Out.WriteInt(val[2])
		break
	case *[3]int64:
		os.Out.WriteLong(val[0])
		os.Out.WriteLong(val[1])
		os.Out.WriteLong(val[2])
		break
	case *[4]int32:
		os.Out.WriteInt(val[0])
		os.Out.WriteInt(val[1])
		os.Out.WriteInt(val[2])
		os.Out.WriteInt(val[3])
		break
	case *[4]int64:
		os.Out.WriteLong(val[0])
		os.Out.WriteLong(val[1])
		os.Out.WriteLong(val[2])
		os.Out.WriteLong(val[3])
		break
	case *[2]uint32:
		os.Out.WriteUInt(val[0])
		os.Out.WriteUInt(val[1])
		break
	case *[2]uint64:
		os.Out.WriteULong(val[0])
		os.Out.WriteULong(val[1])
		break
	case *[3]uint32:
		os.Out.WriteUInt(val[0])
		os.Out.WriteUInt(val[1])
		os.Out.WriteUInt(val[2])
		break
	case *[3]uint64:
		os.Out.WriteULong(val[0])
		os.Out.WriteULong(val[1])
		os.Out.WriteULong(val[2])
		break
	case *[4]uint32:
		os.Out.WriteUInt(val[0])
		os.Out.WriteUInt(val[1])
		os.Out.WriteUInt(val[2])
		os.Out.WriteUInt(val[3])
		break
	case *[4]uint64:
		os.Out.WriteULong(val[0])
		os.Out.WriteULong(val[1])
		os.Out.WriteULong(val[2])
		os.Out.WriteULong(val[3])
		break
	case *[2]uint8:
		os.Out.WriteUChar(val[0])
		os.Out.WriteUChar(val[1])
		break
	case *[2]uint16:
		os.Out.WriteUShort(val[0])
		os.Out.WriteUShort(val[1])
		break
	case *[3]uint8:
		os.Out.WriteUChar(val[0])
		os.Out.WriteUChar(val[1])
		os.Out.WriteUChar(val[2])
		break
	case *[3]uint16:
		os.Out.WriteUShort(val[0])
		os.Out.WriteUShort(val[1])
		os.Out.WriteUShort(val[2])
		break
	case *[4]uint8:
		os.Out.WriteUChar(val[0])
		os.Out.WriteUChar(val[1])
		os.Out.WriteUChar(val[2])
		os.Out.WriteUChar(val[3])
		break
	case *[4]uint16:
		os.Out.WriteUShort(val[0])
		os.Out.WriteUShort(val[1])
		os.Out.WriteUShort(val[2])
		os.Out.WriteUShort(val[3])
		break
	case *[2]int8:
		os.Out.WriteChar(val[0])
		os.Out.WriteChar(val[1])
		break
	case *[2]int16:
		os.Out.WriteShort(val[0])
		os.Out.WriteShort(val[1])
		break
	case *[3]int8:
		os.Out.WriteChar(val[0])
		os.Out.WriteChar(val[1])
		os.Out.WriteChar(val[2])
		break
	case *[3]int16:
		os.Out.WriteShort(val[0])
		os.Out.WriteShort(val[1])
		os.Out.WriteShort(val[2])
		break
	case *[4]int8:
		os.Out.WriteChar(val[0])
		os.Out.WriteChar(val[1])
		os.Out.WriteChar(val[2])
		os.Out.WriteChar(val[3])
		break
	case *[4]int16:
		os.Out.WriteShort(val[0])
		os.Out.WriteShort(val[1])
		os.Out.WriteShort(val[2])
		os.Out.WriteShort(val[3])
		break
	case *[4][4]float32:
		os.WriteMatrix4f(val)
		break
	case *[4][4]float64:
		os.WriteMatrix4d(val)
		break
	case *[4][4]int32:
		os.Write(&val[0])
		os.Write(&val[1])
		os.Write(&val[2])
		os.Write(&val[3])
		break
	case *[4][4]int64:
		os.Write(&val[0])
		os.Write(&val[1])
		os.Write(&val[2])
		os.Write(&val[3])
		break
	case *string:
		os.Out.WriteString(val)
		break
	case *model.ObjectGlenum:
		os.Out.WriteGlenum(val)
		break
	case *model.ObjectProperty:
		os.Out.WriteProperty(val)
		break
	case *model.ObjectMark:
		os.Out.WriteMark(val)
		break
	case *model.PrimitiveSet:
		os.WritePrimitiveSet(val)
		break
	case *CrlfType:
		if os.Out.IsBinary() {
			str := "\r\n"
			os.Out.WriteString(&str)
		}
		break
	}
}

func (os *OsgOstream) GetFileVersion(domain string) int32 {
	if domain == "" {
		return os.TargetFileVersion
	}
	v, ok := os.DomainVersionMap[domain]
	if ok {
		return v
	}
	return 0
}

func (os *OsgOstream) IsBinary() bool {
	return os.Out.IsBinary()
}

func (os *OsgOstream) findOrCreateArrayId(ay *model.Array, newId *bool) int32 {
	it, ok := os.ArrayMap[ay]
	if !ok {
		id := len(os.ArrayMap) + 1
		os.ArrayMap[ay] = int32(id)
		*newId = true
		return int32(id)
	}
	*newId = false
	return it
}

func (os *OsgOstream) findOrCreateObjectId(ob model.ObjectInterface, newId *bool) int32 {
	it, ok := os.ObjectMap[ob]
	if !ok {
		id := len(os.ObjectMap) + 1
		os.ObjectMap[ob] = int32(id)
		*newId = true
		return int32(id)
	}
	*newId = false
	return it
}

func (os *OsgOstream) WriteArray(ay *model.Array) {
	if ay == nil {
		return
	}
	isNew := false
	id := os.findOrCreateArrayId(ay, &isNew)
	os.PROPERTY.Name = "ArrayID"
	os.Write(os.PROPERTY)
	os.Write(id)
	if !isNew {
		os.Write(os.CRLF)
		return
	}
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	switch ay.Type {
	case model.ByteArrayType:
		dt := ay.Data.([]int8)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.UByteArrayType:
		dt := ay.Data.([]uint8)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.ShortArrayType:
		dt := ay.Data.([]int16)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.UShortArrayType:
		dt := ay.Data.([]uint16)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.IntArrayType:
		dt := ay.Data.([]int32)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.UIntArrayType:
		dt := ay.Data.([]uint32)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.FloatArrayType:
		dt := ay.Data.([]float32)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.DoubleArrayType:
		dt := ay.Data.([]float64)
		for index, d := range dt {
			os.Write(d)
			if index%4 == 0 {
				os.Write(os.CRLF)
			}
		}
		break
	case model.Vec2bArrayType:
		dt := ay.Data.([][2]int8)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3bArrayType:
		dt := ay.Data.([][3]int8)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4bArrayType:
		dt := ay.Data.([][4]int8)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break
	case model.Vec2ubArrayType:
		dt := ay.Data.([][2]uint8)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3ubArrayType:
		dt := ay.Data.([][3]uint8)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4ubArrayType:
		dt := ay.Data.([][4]uint8)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break

	case model.Vec2sArrayType:
		dt := ay.Data.([][2]int16)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3sArrayType:
		dt := ay.Data.([][3]int16)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4sArrayType:
		dt := ay.Data.([][4]int16)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break

	case model.Vec2usArrayType:
		dt := ay.Data.([][2]uint16)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3usArrayType:
		dt := ay.Data.([][3]uint16)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4usArrayType:
		dt := ay.Data.([][4]uint16)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break

	case model.Vec2ArrayType:
		dt := ay.Data.([][2]float32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3ArrayType:
		dt := ay.Data.([][3]float32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4ArrayType:
		dt := ay.Data.([][4]float32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break

	case model.Vec2dArrayType:
		dt := ay.Data.([][2]float64)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3dArrayType:
		dt := ay.Data.([][3]float64)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4dArrayType:
		dt := ay.Data.([][4]float64)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break

	case model.Vec2iArrayType:
		dt := ay.Data.([][2]int32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(os.CRLF)
		}
		break
	case model.Vec3iArrayType:
		dt := ay.Data.([][3]int32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4iArrayType:
		dt := ay.Data.([][4]int32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
			os.Write(d[2])
			os.Write(d[3])
		}
		break

	case model.Vec2uiArrayType:
		dt := ay.Data.([][2]uint32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(d[1])
		}
		break
	case model.Vec3uiArrayType:
		dt := ay.Data.([][3]uint32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
		}
		break
	case model.Vec4uiArrayType:
		dt := ay.Data.([][4]uint32)
		for _, d := range dt {
			os.Write(d[0])
			os.Write(os.CRLF)
			os.Write(d[1])
			os.Write(os.CRLF)
			os.Write(d[2])
			os.Write(os.CRLF)
			os.Write(d[3])
			os.Write(os.CRLF)
		}
		break
	}
	os.Write(os.ENDBRACKET)
}

func (os *OsgOstream) WritePrimitiveSet(ps interface{}) {
	if ps != nil {
		return
	}
	pro := &model.ObjectProperty{Name: "PrimitiveType", MapProperty: true}
	md := &model.ObjectProperty{Name: "PrimitiveType", MapProperty: true}
	switch p := ps.(type) {
	case *model.DrawArrays:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		os.Write(pro)
		os.Write(md)
		if os.TargetFileVersion > 96 {
			os.Write(p.NumInstances)
		}
		os.Write(p.Count)
		os.Write(os.CRLF)
		break
	case model.DrawArrayLengths:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		os.Write(pro)
		os.Write(md)
		if os.TargetFileVersion > 96 {
			os.Write(p.NumInstances)
		}
		os.Write(p.First)
		os.Write((int32)(len(p.Data)))
		os.Write(os.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			os.Write(os.CRLF)
			os.Write(p.Data[i])
			os.Write(p.Data[i+1])
			os.Write(p.Data[i+2])
			os.Write(p.Data[i+3])
			os.Write(os.CRLF)
		}
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
		break
	case model.DrawElementsUByte:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		os.Write(pro)
		os.Write(md)
		if os.TargetFileVersion > 96 {
			os.Write(p.NumInstances)
		}
		os.Write((int32)(len(p.Data)))
		os.Write(os.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			os.Write(os.CRLF)
			os.Write(p.Data[i])
			os.Write(p.Data[i+1])
			os.Write(p.Data[i+2])
			os.Write(p.Data[i+3])
			os.Write(os.CRLF)
		}
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
		break
	case model.DrawElementsUShort:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		os.Write(pro)
		os.Write(md)
		if os.TargetFileVersion > 96 {
			os.Write(p.NumInstances)
		}
		os.Write((int32)(len(p.Data)))
		os.Write(os.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			os.Write(os.CRLF)
			os.Write(p.Data[i])
			os.Write(p.Data[i+1])
			os.Write(p.Data[i+2])
			os.Write(p.Data[i+3])
			os.Write(os.CRLF)
		}
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
		break
	case model.DrawElementsUInt:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		os.Write(pro)
		os.Write(md)
		if os.TargetFileVersion > 96 {
			os.Write(p.NumInstances)
		}
		os.Write((int32)(len(p.Data)))
		os.Write(os.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			os.Write(os.CRLF)
			os.Write(p.Data[i])
			os.Write(p.Data[i+1])
			os.Write(p.Data[i+2])
			os.Write(p.Data[i+3])
			os.Write(os.CRLF)
		}
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
		break
	}
}

func (os *OsgOstream) WriteImage(inter *model.Image) {
}

func (os *OsgOstream) WriteMatrix4f(mat *[4][4]float32) {
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func (os *OsgOstream) WriteMatrix4d(mat *[4][4]float64) {
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(&mat[0])
	os.Write(os.CRLF)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func (os *OsgOstream) WriteObject(obj interface{}) {
	if obj == nil {
		os.Write("NULL")
		os.Write(os.CRLF)
	}
	inter := obj.(model.ObjectInterface)
	name := inter.GetName()
	nw := false
	id := os.findOrCreateObjectId(inter, &nw)

	os.Write(name)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "UniqueID"
	os.Write(os.PROPERTY)
	os.Write(id)
	os.Write(os.CRLF)
	if nw {
		os.WriteObjectFileds(inter, *name)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func (os *OsgOstream) WriteObjectFileds(obj interface{}, name string) {
	wrap := GetObjectWrapperManager().FindWrap(name)
	if wrap == nil {
		return
	}
	ver := os.DomainVersionMap[wrap.Domain]
	for _, as := range wrap.Associates {
		if as == nil {
			continue
		}
		if as.FirstVersion <= ver && as.LastVersion >= ver {
			assocWrapper := GetObjectWrapperManager().FindWrap(as.Name)
			if assocWrapper == nil {
				continue
			} else if os.UseSchemaData {
				_, ok := os.InbuiltSchemaMap[as.Name]
				if ok {
					prop := []string{}
					tys := []SerType{}
					assocWrapper.WriteSchema(prop, tys)
					size := len(prop)
					if size > len(tys) {
						size = len(tys)
					}
					if size > 0 {
						var propertiesStream string
						for i := 0; i < size; i++ {
							propertiesStream += prop[i] + ":"
							propertiesStream += strconv.Itoa(int(tys[i]))
							propertiesStream += " "
						}
						os.InbuiltSchemaMap[as.Name] = propertiesStream
					}
				}
			}
			os.Fields = append(os.Fields, assocWrapper.Name)
			assocWrapper.Write(os, obj)
			os.Fields = os.Fields[:len(os.Fields)-1]
		}
	}
}

func (os *OsgOstream) Start(out OsgOutputIterator, ty int32) {
	os.Fields = []string{}
	os.Fields = append(os.Fields, "Start")
	os.Out = out
	if out == nil {
		panic("outiterator is nil")
	}
	os.Out.SetOutputSteam(os)
	if os.IsBinary() {
		os.Write(ty)
		os.Write(os.TargetFileVersion)
		var attributes int32 = 0
		if len(os.DomainVersionMap) > 0 {
			attributes |= 0x1
		}
		if os.UseSchemaData {
			attributes |= 0x2
		}
		if os.UseRobustBinaryFormat {
			os.Out.SetSupportBinaryBrackets(true)
			attributes |= 0x4
		}
		os.Write(attributes)
		size := len(os.DomainVersionMap)
		if size > 0 {
			for k, v := range os.DomainVersionMap {
				os.Write(k)
				os.Write(v)
			}
		}
	} else {
		typeString := "Unknown"
		switch ty {
		case WRITESCENE:
			typeString = "Scene"
			break
		case WRITEIMAGE:
			typeString = "Image"
			break
		case WRITEOBJECT:
			typeString = "Object"
			break
		}
		os.Write(&typeString)
		os.Write(os.CRLF)
		os.PROPERTY.Name = "#Version"
		os.Write(os.PROPERTY)
		os.Write(OPENSCENEGRAPHSOVERSION)
		os.Write(os.CRLF)
		os.PROPERTY.Name = "#Generator"
		os.Write(os.PROPERTY)
		str := "FLYWAVE"
		os.Write(&str)
		p := "0.1"
		os.Write(&p)
		os.Write(os.CRLF)
		if len(os.DomainVersionMap) > 0 {
			for k, v := range os.DomainVersionMap {
				os.PROPERTY.Name = "#CustomDomain"
				os.Write(os.PROPERTY)
				os.Write(&k)
				os.Write(&v)
			}
		}
		os.Write(os.CRLF)
		os.Fields = os.Fields[:len(os.Fields)-1]
	}
}

func (os *OsgOstream) Compress() []byte {
	if !os.IsBinary() || os.CompressorName == "" {
		return os.Data
	}
	os.Fields = []string{}
	if os.UseSchemaData {
		os.Fields = append(os.Fields, "SchemaData")
	}
	var schemaSource string
	var schemaData string
	for k, v := range os.InbuiltSchemaMap {
		schemaData += k + "=" + v
		schemaData += "\n"
	}
	sz := len(schemaData)
	schemaSource = strconv.Itoa(sz)
	schemaSource += schemaData

	os.Fields = append(os.Fields, "Compression")
	compress_wrap := GetObjectWrapperManager().FindCompressor(os.CompressorName)
	if compress_wrap == nil {
		return os.Data
	}
	os.Write(&schemaSource)
	var compresseData []byte
	compress_wrap.Compress(os.Out.GetIterator(), compresseData)
	return compresseData
}
