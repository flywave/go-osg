package osg

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"os"
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

func NewOsgOstream(opts *OsgOstreamOptions) *OsgOstream {
	p := model.NewObjectProperty()
	bb := model.NewObjectMark()
	bb.Name = "{"
	bb.IndentDelta = INDENT_VALUE
	eb := model.NewObjectMark()
	bb.Name = "}"
	bb.IndentDelta = -INDENT_VALUE
	osg := &OsgOstream{PROPERTY: p, BEGINBRACKET: bb, ENDBRACKET: eb, CRLF: &CrlfType{}, TargetFileVersion: OPENSCENEGRAPHSOVERSION, UseRobustBinaryFormat: true, UseSchemaData: false}
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

func (ostream *OsgOstream) Write(inter interface{}) {
	switch val := inter.(type) {
	case bool:
		ostream.Out.WriteBool(val)
		break
	case int8:
		ostream.Out.WriteChar(val)
		break
	case uint8:
		ostream.Out.WriteUChar(val)
		break
	case int16:
		ostream.Out.WriteShort(val)
		break
	case uint16:
		ostream.Out.WriteUShort(val)
		break
	case int32:
		ostream.Out.WriteInt(val)
		break
	case uint32:
		ostream.Out.WriteUInt(val)
		break
	case int64:
		ostream.Out.WriteLong(val)
		break
	case uint64:
		ostream.Out.WriteULong(val)
		break
	case float32:
		ostream.Out.WriteFloat(val)
		break
	case float64:
		ostream.Out.WriteDouble(val)
		break
	case *[2]float32:
		ostream.Out.WriteFloat(val[0])
		ostream.Out.WriteFloat(val[1])
		break
	case *[2]float64:
		ostream.Out.WriteDouble(val[0])
		ostream.Out.WriteDouble(val[1])
		break
	case *[3]float32:
		ostream.Out.WriteFloat(val[0])
		ostream.Out.WriteFloat(val[1])
		ostream.Out.WriteFloat(val[2])
		break
	case *[3]float64:
		ostream.Out.WriteDouble(val[0])
		ostream.Out.WriteDouble(val[1])
		ostream.Out.WriteDouble(val[2])
		break
	case *[4]float32:
		ostream.Out.WriteFloat(val[0])
		ostream.Out.WriteFloat(val[1])
		ostream.Out.WriteFloat(val[2])
		ostream.Out.WriteFloat(val[3])
		break
	case *[4]float64:
		ostream.Out.WriteDouble(val[0])
		ostream.Out.WriteDouble(val[1])
		ostream.Out.WriteDouble(val[2])
		ostream.Out.WriteDouble(val[3])
		break
	case *[2]int32:
		ostream.Out.WriteInt(val[0])
		ostream.Out.WriteInt(val[1])
		break
	case *[2]int64:
		ostream.Out.WriteLong(val[0])
		ostream.Out.WriteLong(val[1])
		break
	case *[3]int32:
		ostream.Out.WriteInt(val[0])
		ostream.Out.WriteInt(val[1])
		ostream.Out.WriteInt(val[2])
		break
	case *[3]int64:
		ostream.Out.WriteLong(val[0])
		ostream.Out.WriteLong(val[1])
		ostream.Out.WriteLong(val[2])
		break
	case *[4]int32:
		ostream.Out.WriteInt(val[0])
		ostream.Out.WriteInt(val[1])
		ostream.Out.WriteInt(val[2])
		ostream.Out.WriteInt(val[3])
		break
	case *[4]int64:
		ostream.Out.WriteLong(val[0])
		ostream.Out.WriteLong(val[1])
		ostream.Out.WriteLong(val[2])
		ostream.Out.WriteLong(val[3])
		break
	case *[2]uint32:
		ostream.Out.WriteUInt(val[0])
		ostream.Out.WriteUInt(val[1])
		break
	case *[2]uint64:
		ostream.Out.WriteULong(val[0])
		ostream.Out.WriteULong(val[1])
		break
	case *[3]uint32:
		ostream.Out.WriteUInt(val[0])
		ostream.Out.WriteUInt(val[1])
		ostream.Out.WriteUInt(val[2])
		break
	case *[3]uint64:
		ostream.Out.WriteULong(val[0])
		ostream.Out.WriteULong(val[1])
		ostream.Out.WriteULong(val[2])
		break
	case *[4]uint32:
		ostream.Out.WriteUInt(val[0])
		ostream.Out.WriteUInt(val[1])
		ostream.Out.WriteUInt(val[2])
		ostream.Out.WriteUInt(val[3])
		break
	case *[4]uint64:
		ostream.Out.WriteULong(val[0])
		ostream.Out.WriteULong(val[1])
		ostream.Out.WriteULong(val[2])
		ostream.Out.WriteULong(val[3])
		break
	case *[2]uint8:
		ostream.Out.WriteUChar(val[0])
		ostream.Out.WriteUChar(val[1])
		break
	case *[2]uint16:
		ostream.Out.WriteUShort(val[0])
		ostream.Out.WriteUShort(val[1])
		break
	case *[3]uint8:
		ostream.Out.WriteUChar(val[0])
		ostream.Out.WriteUChar(val[1])
		ostream.Out.WriteUChar(val[2])
		break
	case *[3]uint16:
		ostream.Out.WriteUShort(val[0])
		ostream.Out.WriteUShort(val[1])
		ostream.Out.WriteUShort(val[2])
		break
	case *[4]uint8:
		ostream.Out.WriteUChar(val[0])
		ostream.Out.WriteUChar(val[1])
		ostream.Out.WriteUChar(val[2])
		ostream.Out.WriteUChar(val[3])
		break
	case *[4]uint16:
		ostream.Out.WriteUShort(val[0])
		ostream.Out.WriteUShort(val[1])
		ostream.Out.WriteUShort(val[2])
		ostream.Out.WriteUShort(val[3])
		break
	case *[2]int8:
		ostream.Out.WriteChar(val[0])
		ostream.Out.WriteChar(val[1])
		break
	case *[2]int16:
		ostream.Out.WriteShort(val[0])
		ostream.Out.WriteShort(val[1])
		break
	case *[3]int8:
		ostream.Out.WriteChar(val[0])
		ostream.Out.WriteChar(val[1])
		ostream.Out.WriteChar(val[2])
		break
	case *[3]int16:
		ostream.Out.WriteShort(val[0])
		ostream.Out.WriteShort(val[1])
		ostream.Out.WriteShort(val[2])
		break
	case *[4]int8:
		ostream.Out.WriteChar(val[0])
		ostream.Out.WriteChar(val[1])
		ostream.Out.WriteChar(val[2])
		ostream.Out.WriteChar(val[3])
		break
	case *[4]int16:
		ostream.Out.WriteShort(val[0])
		ostream.Out.WriteShort(val[1])
		ostream.Out.WriteShort(val[2])
		ostream.Out.WriteShort(val[3])
		break
	case *[4][4]float32:
		ostream.WriteMatrix4f(val)
		break
	case *[4][4]float64:
		ostream.WriteMatrix4d(val)
		break
	case *[4][4]int32:
		ostream.Write(&val[0])
		ostream.Write(&val[1])
		ostream.Write(&val[2])
		ostream.Write(&val[3])
		break
	case *[4][4]int64:
		ostream.Write(&val[0])
		ostream.Write(&val[1])
		ostream.Write(&val[2])
		ostream.Write(&val[3])
		break
	case *string:
		ostream.Out.WriteString(val)
		break
	case *model.ObjectGlenum:
		ostream.Out.WriteGlenum(val)
		break
	case *model.ObjectProperty:
		ostream.Out.WriteProperty(val)
		break
	case *model.ObjectMark:
		ostream.Out.WriteMark(val)
		break
	case *model.PrimitiveSet:
		ostream.WritePrimitiveSet(val)
		break
	case *CrlfType:
		if ostream.Out.IsBinary() {
			str := "\r\n"
			ostream.Out.WriteString(&str)
		}
		break
	}
}

func (ostream *OsgOstream) GetFileVersion(domain string) int32 {
	if domain == "" {
		return ostream.TargetFileVersion
	}
	v, ok := ostream.DomainVersionMap[domain]
	if ok {
		return v
	}
	return 0
}

func (ostream *OsgOstream) IsBinary() bool {
	return ostream.Out.IsBinary()
}

func (ostream *OsgOstream) findOrCreateArrayId(ay *model.Array, newId *bool) int32 {
	it, ok := ostream.ArrayMap[ay]
	if !ok {
		id := len(ostream.ArrayMap) + 1
		ostream.ArrayMap[ay] = int32(id)
		*newId = true
		return int32(id)
	}
	*newId = false
	return it
}

func (ostream *OsgOstream) findOrCreateObjectId(ob model.ObjectInterface, newId *bool) int32 {
	it, ok := ostream.ObjectMap[ob]
	if !ok {
		id := len(ostream.ObjectMap) + 1
		ostream.ObjectMap[ob] = int32(id)
		*newId = true
		return int32(id)
	}
	*newId = false
	return it
}

func (ostream *OsgOstream) WriteArray(ay *model.Array) {
	if ay == nil {
		return
	}
	isNew := false
	id := ostream.findOrCreateArrayId(ay, &isNew)
	ostream.PROPERTY.Name = "ArrayID"
	ostream.Write(ostream.PROPERTY)
	ostream.Write(id)
	if !isNew {
		ostream.Write(ostream.CRLF)
		return
	}
	ostream.Write(ostream.BEGINBRACKET)
	ostream.Write(ostream.CRLF)
	switch ay.Type {
	case model.ByteArrayType:
		dt := ay.Data.([]int8)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.UByteArrayType:
		dt := ay.Data.([]uint8)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.ShortArrayType:
		dt := ay.Data.([]int16)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.UShortArrayType:
		dt := ay.Data.([]uint16)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.IntArrayType:
		dt := ay.Data.([]int32)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.UIntArrayType:
		dt := ay.Data.([]uint32)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.FloatArrayType:
		dt := ay.Data.([]float32)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.DoubleArrayType:
		dt := ay.Data.([]float64)
		for index, d := range dt {
			ostream.Write(d)
			if index%4 == 0 {
				ostream.Write(ostream.CRLF)
			}
		}
		break
	case model.Vec2bArrayType:
		dt := ay.Data.([][2]int8)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3bArrayType:
		dt := ay.Data.([][3]int8)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4bArrayType:
		dt := ay.Data.([][4]int8)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec2ubArrayType:
		dt := ay.Data.([][2]uint8)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3ubArrayType:
		dt := ay.Data.([][3]uint8)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4ubArrayType:
		dt := ay.Data.([][4]uint8)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break

	case model.Vec2sArrayType:
		dt := ay.Data.([][2]int16)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3sArrayType:
		dt := ay.Data.([][3]int16)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4sArrayType:
		dt := ay.Data.([][4]int16)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break

	case model.Vec2usArrayType:
		dt := ay.Data.([][2]uint16)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3usArrayType:
		dt := ay.Data.([][3]uint16)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4usArrayType:
		dt := ay.Data.([][4]uint16)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break

	case model.Vec2ArrayType:
		dt := ay.Data.([][2]float32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3ArrayType:
		dt := ay.Data.([][3]float32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4ArrayType:
		dt := ay.Data.([][4]float32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break

	case model.Vec2dArrayType:
		dt := ay.Data.([][2]float64)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3dArrayType:
		dt := ay.Data.([][3]float64)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4dArrayType:
		dt := ay.Data.([][4]float64)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break

	case model.Vec2iArrayType:
		dt := ay.Data.([][2]int32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec3iArrayType:
		dt := ay.Data.([][3]int32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4iArrayType:
		dt := ay.Data.([][4]int32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
			ostream.Write(d[2])
			ostream.Write(d[3])
		}
		break

	case model.Vec2uiArrayType:
		dt := ay.Data.([][2]uint32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(d[1])
		}
		break
	case model.Vec3uiArrayType:
		dt := ay.Data.([][3]uint32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
		}
		break
	case model.Vec4uiArrayType:
		dt := ay.Data.([][4]uint32)
		for _, d := range dt {
			ostream.Write(d[0])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[1])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[2])
			ostream.Write(ostream.CRLF)
			ostream.Write(d[3])
			ostream.Write(ostream.CRLF)
		}
		break
	}
	ostream.Write(ostream.ENDBRACKET)
}

func (ostream *OsgOstream) WritePrimitiveSet(ps interface{}) {
	if ps != nil {
		return
	}
	pro := &model.ObjectProperty{Name: "PrimitiveType", MapProperty: true}
	md := &model.ObjectProperty{Name: "PrimitiveType", MapProperty: true}
	switch p := ps.(type) {
	case *model.DrawArrays:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		ostream.Write(pro)
		ostream.Write(md)
		if ostream.TargetFileVersion > 96 {
			ostream.Write(p.NumInstances)
		}
		ostream.Write(p.Count)
		ostream.Write(ostream.CRLF)
		break
	case model.DrawArrayLengths:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		ostream.Write(pro)
		ostream.Write(md)
		if ostream.TargetFileVersion > 96 {
			ostream.Write(p.NumInstances)
		}
		ostream.Write(p.First)
		ostream.Write((int32)(len(p.Data)))
		ostream.Write(ostream.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			ostream.Write(ostream.CRLF)
			ostream.Write(p.Data[i])
			ostream.Write(p.Data[i+1])
			ostream.Write(p.Data[i+2])
			ostream.Write(p.Data[i+3])
			ostream.Write(ostream.CRLF)
		}
		ostream.Write(ostream.ENDBRACKET)
		ostream.Write(ostream.CRLF)
		break
	case model.DrawElementsUByte:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		ostream.Write(pro)
		ostream.Write(md)
		if ostream.TargetFileVersion > 96 {
			ostream.Write(p.NumInstances)
		}
		ostream.Write((int32)(len(p.Data)))
		ostream.Write(ostream.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			ostream.Write(ostream.CRLF)
			ostream.Write(p.Data[i])
			ostream.Write(p.Data[i+1])
			ostream.Write(p.Data[i+2])
			ostream.Write(p.Data[i+3])
			ostream.Write(ostream.CRLF)
		}
		ostream.Write(ostream.ENDBRACKET)
		ostream.Write(ostream.CRLF)
		break
	case model.DrawElementsUShort:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		ostream.Write(pro)
		ostream.Write(md)
		if ostream.TargetFileVersion > 96 {
			ostream.Write(p.NumInstances)
		}
		ostream.Write((int32)(len(p.Data)))
		ostream.Write(ostream.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			ostream.Write(ostream.CRLF)
			ostream.Write(p.Data[i])
			ostream.Write(p.Data[i+1])
			ostream.Write(p.Data[i+2])
			ostream.Write(p.Data[i+3])
			ostream.Write(ostream.CRLF)
		}
		ostream.Write(ostream.ENDBRACKET)
		ostream.Write(ostream.CRLF)
		break
	case model.DrawElementsUInt:
		pro.Value = p.PrimitiveType
		md.Value = p.Mode
		ostream.Write(pro)
		ostream.Write(md)
		if ostream.TargetFileVersion > 96 {
			ostream.Write(p.NumInstances)
		}
		ostream.Write((int32)(len(p.Data)))
		ostream.Write(ostream.BEGINBRACKET)

		for i := 0; i < len(p.Data); i += 4 {
			ostream.Write(ostream.CRLF)
			ostream.Write(p.Data[i])
			ostream.Write(p.Data[i+1])
			ostream.Write(p.Data[i+2])
			ostream.Write(p.Data[i+3])
			ostream.Write(ostream.CRLF)
		}
		ostream.Write(ostream.ENDBRACKET)
		ostream.Write(ostream.CRLF)
		break
	}
}

func (ostream *OsgOstream) WriteImage(img *model.Image) {
	if img == nil {
		ostream.Write("NULL")
		ostream.Write(ostream.CRLF)
	}
	name := img.Type
	nw := false
	id := ostream.findOrCreateObjectId(img, &nw)

	if ostream.TargetFileVersion > 94 {
		ostream.PROPERTY.Name = "ClassName"
		ostream.Write(ostream.PROPERTY)
		ostream.Write(name)
		ostream.Write(ostream.CRLF)
		return
	}
	ostream.PROPERTY.Name = "UniqueID"
	ostream.Write(ostream.PROPERTY)
	ostream.Write(id)
	ostream.Write(ostream.CRLF)
	if nw {
		decision := model.IMAGEEXTERNAL
		switch ostream.WriteImageHint {
		case WRITEINLINEDATA:
			decision = model.IMAGEINLINEDATA
			break
		case WRITEINLINEFILE:
			decision = model.IMAGEINLINEFILE
			break
		case WRITEEXTERNALFILE:
			decision = model.IMAGEWRITEOUT
			break
		case WRITEUSEEXTERNAL:
			decision = model.IMAGEEXTERNAL
			break
		default:
			if img.WriteHint == model.EXTERNALFILE {
				decision = model.IMAGEEXTERNAL
			} else {
				decision = model.IMAGEINLINEDATA
			}
			break
		}
		filename := img.FileName

		if decision == model.IMAGEWRITEOUT ||
			ostream.WriteImageHint == WRITEEXTERNALFILE {
			if filename == "" {
				filename = "image.dds"
			}
		}

		ostream.PROPERTY.Name = "FileName"
		ostream.Write(ostream.PROPERTY)
		ostream.Out.WriteWrappedString(&filename)
		ostream.Write(ostream.CRLF)
		ostream.PROPERTY.Name = "WriteHint"
		ostream.Write(ostream.PROPERTY)
		ostream.Write(img.WriteHint)
		ostream.Write(decision)
		ostream.Write(ostream.CRLF)
		switch decision {
		case model.IMAGEINLINEDATA:
			if ostream.IsBinary() {
				ostream.Write(img.Origin)
				ostream.Write(img.S)
				ostream.Write(img.T)
				ostream.Write(img.R)
				ostream.Write(img.InternalTextureFormat)
				ostream.Write(img.PixelFormat)
				ostream.Write(img.DataType)
				ostream.Write(img.Packing)
				ostream.Write(img.AllocationMode)
				ostream.Write((int32)(len(img.Data)))
				ostream.Out.WriteCharArray(img.Data)
				var num int32 = 0
				ostream.Write(num)
			} else {
				ostream.PROPERTY.Name = "Origin"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.Origin)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "Size"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.S)
				ostream.Write(img.T)
				ostream.Write(img.R)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "InternalTextureFormat"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.InternalTextureFormat)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "PixelFormat"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.PixelFormat)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "DataType"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.DataType)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "Packing"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.Packing)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "AllocationMode"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(img.AllocationMode)
				ostream.Write(ostream.CRLF)

				ostream.PROPERTY.Name = "Data"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(int32(0))
				ostream.Write(ostream.BEGINBRACKET)
				ostream.Write(ostream.CRLF)

				str := base64.StdEncoding.EncodeToString(img.Data)
				ostream.Write(&str)
				ostream.Write(ostream.ENDBRACKET)
				ostream.Write(ostream.CRLF)
			}
			break
		case model.IMAGEINLINEFILE:
			if ostream.IsBinary() {
				fp := img.FileName
				if fp != "" {
					fl, _ := os.Open(fp)
					data := []byte{}
					io.ReadFull(fl, data)
					sz := (int32)(len(data))
					if sz > 0 {
						ostream.Write(sz)
						ostream.Out.WriteCharArray(data)
					} else {
						ostream.Write(sz)
					}
				} else {
					rw := getReaderWriter()
					if rw != nil {
						data := []byte{}
						buf := bytes.NewBuffer(data)
						wt := bufio.NewWriter(buf)
						rw.WriteImageWithWrite(img, wt, ostream.Options)
						ostream.Write(len(data))
						ostream.Out.WriteCharArray(data)
					} else {
						ostream.Write(int(0))
					}
				}
			}
			break
		default:
			break
		}
		ostream.WriteObjectFileds(img, "osg::Object")
	}
}

func (ostream *OsgOstream) WriteMatrix4f(mat *[4][4]float32) {
	ostream.Write(ostream.BEGINBRACKET)
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(ostream.ENDBRACKET)
	ostream.Write(ostream.CRLF)
}

func (ostream *OsgOstream) WriteMatrix4d(mat *[4][4]float64) {
	ostream.Write(ostream.BEGINBRACKET)
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(&mat[0])
	ostream.Write(ostream.CRLF)
	ostream.Write(ostream.ENDBRACKET)
	ostream.Write(ostream.CRLF)
}

func (ostream *OsgOstream) WriteObject(obj interface{}) {
	if obj == nil {
		ostream.Write("NULL")
		ostream.Write(ostream.CRLF)
	}
	inter := obj.(model.ObjectInterface)
	name := inter.GetName()
	nw := false
	id := ostream.findOrCreateObjectId(inter, &nw)

	ostream.Write(name)
	ostream.Write(ostream.BEGINBRACKET)
	ostream.Write(ostream.CRLF)
	ostream.PROPERTY.Name = "UniqueID"
	ostream.Write(ostream.PROPERTY)
	ostream.Write(id)
	ostream.Write(ostream.CRLF)
	if nw {
		ostream.WriteObjectFileds(inter, *name)
	}
	ostream.Write(ostream.ENDBRACKET)
	ostream.Write(ostream.CRLF)
}

func (ostream *OsgOstream) WriteObjectFileds(obj interface{}, name string) {
	wrap := GetObjectWrapperManager().FindWrap(name)
	if wrap == nil {
		return
	}
	ver := ostream.DomainVersionMap[wrap.Domain]
	for _, as := range wrap.Associates {
		if as == nil {
			continue
		}
		if as.FirstVersion <= ver && as.LastVersion >= ver {
			assocWrapper := GetObjectWrapperManager().FindWrap(as.Name)
			if assocWrapper == nil {
				continue
			} else if ostream.UseSchemaData {
				_, ok := ostream.InbuiltSchemaMap[as.Name]
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
						ostream.InbuiltSchemaMap[as.Name] = propertiesStream
					}
				}
			}
			ostream.Fields = append(ostream.Fields, assocWrapper.Name)
			assocWrapper.Write(ostream, obj)
			ostream.Fields = ostream.Fields[:len(ostream.Fields)-1]
		}
	}
}

func (ostream *OsgOstream) Start(out OsgOutputIterator, ty int32) error {
	ostream.Fields = []string{}
	ostream.Fields = append(ostream.Fields, "Start")
	ostream.Out = out
	if out == nil {
		return errors.New("outiterator is nil")
	}
	ostream.Out.SetOutputSteam(ostream)
	if ostream.IsBinary() {
		ostream.Write(ty)
		ostream.Write(ostream.TargetFileVersion)
		var attributes int32 = 0
		if len(ostream.DomainVersionMap) > 0 {
			attributes |= 0x1
		}
		if ostream.UseSchemaData {
			attributes |= 0x2
		}
		if ostream.UseRobustBinaryFormat {
			ostream.Out.SetSupportBinaryBrackets(true)
			attributes |= 0x4
		}
		ostream.Write(attributes)
		size := len(ostream.DomainVersionMap)
		if size > 0 {
			for k, v := range ostream.DomainVersionMap {
				ostream.Write(k)
				ostream.Write(v)
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
		ostream.Write(&typeString)
		ostream.Write(ostream.CRLF)
		ostream.PROPERTY.Name = "#Version"
		ostream.Write(ostream.PROPERTY)
		ostream.Write(OPENSCENEGRAPHSOVERSION)
		ostream.Write(ostream.CRLF)
		ostream.PROPERTY.Name = "#Generator"
		ostream.Write(ostream.PROPERTY)
		str := "FLYWAVE"
		ostream.Write(&str)
		p := "0.1"
		ostream.Write(&p)
		ostream.Write(ostream.CRLF)
		if len(ostream.DomainVersionMap) > 0 {
			for k, v := range ostream.DomainVersionMap {
				ostream.PROPERTY.Name = "#CustomDomain"
				ostream.Write(ostream.PROPERTY)
				ostream.Write(&k)
				ostream.Write(&v)
			}
		}
		ostream.Write(ostream.CRLF)
		ostream.Fields = ostream.Fields[:len(ostream.Fields)-1]
	}
	return nil
}

func (ostream *OsgOstream) Compress(buff *bytes.Buffer) []byte {
	if !ostream.IsBinary() || ostream.CompressorName == "" {
		return buff.Bytes()
	}
	ostream.Fields = []string{}
	if ostream.UseSchemaData {
		ostream.Fields = append(ostream.Fields, "SchemaData")
	}
	var schemaSource string
	var schemaData string
	for k, v := range ostream.InbuiltSchemaMap {
		schemaData += k + "=" + v
		schemaData += "\n"
	}
	sz := len(schemaData)
	schemaSource = strconv.Itoa(sz)
	schemaSource += schemaData

	ostream.Fields = append(ostream.Fields, "Compression")
	compress_wrap := GetObjectWrapperManager().FindCompressor(ostream.CompressorName)
	if compress_wrap == nil {
		return buff.Bytes()
	}
	ostream.Write(&schemaSource)
	var compresseData []byte
	compress_wrap.Compress(buff, compresseData)
	return compresseData
}
