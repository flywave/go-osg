package io

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
	FileType     string = "Ascii"
	INDENT_VALUE        = 2
)

type OsgOptions struct {
	FileType   string
	Precision  int
	Compressed bool
}

func NewOsgOptions() OsgOptions {
	return OsgOptions{FileType: FileType}
}

type ReadType int

const (
	READUNKNOWN ReadType = 0
	READSCENE   ReadType = 1
	READIMAGE   ReadType = 2
	READOBJECT  ReadType = 3
)

type OsgIstreamOptions struct {
	OsgOptions
	DbPath            string
	Domain            string
	ForceReadingImage bool
}

func NewOsgIstreamOptions() OsgIstreamOptions {
	op := NewOsgOptions()
	return OsgIstreamOptions{OsgOptions: op}
}

type StreamHeader struct {
	Version       int32
	Type          ReadType
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
	DataDecompress    io.Reader
	Data              []byte
	CRLF              CrlfType

	PROPERTY     *model.ObjectProperty
	BEGINBRACKET *model.ObjectMark
	ENDBRACKET   *model.ObjectMark
}

func NewOsgIstream(opt *OsgIstreamOptions) OsgIstream {
	p := model.NewObjectProperty()
	bb := model.NewObjectMark()
	bb.Name = "{"
	bb.IndentDelta = INDENT_VALUE
	eb := model.NewObjectMark()
	bb.Name = "}"
	bb.IndentDelta = -INDENT_VALUE
	is := OsgIstream{ArrayMap: make(map[int32]*model.Array), Options: opt, IdentifierMap: make(map[int32]interface{}), DomainVersionMap: make(map[string]int32), PROPERTY: &p, BEGINBRACKET: &bb, ENDBRACKET: &eb}
	if opt.ForceReadingImage {
		is.ForceReadingImage = true
	}
	obj := model.NewObject()
	is.DummyReadObject = &obj
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
	case *byte:
		is.In.ReadChar(val)
		break
	// case *uint8:
	// 	is.In.ReadUChar(val)
	// 	break
	case *int16:
		is.In.ReadShort(val)
		break
	case *uint16:
		is.In.ReadUShort(val)
		break
	case *int32:
		is.In.ReadInt(val)
		break
	case *uint32:
		is.In.ReadUInt(val)
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
			// levelSize := is.ReadSize() - 1
			is.Read(is.BEGINBRACKET)
			var encodedData string
			is.ReadWrappedString(&encodedData)
			d, e := base64.StdEncoding.DecodeString(encodedData)
			if e == nil {
				imgdata.Data = d
			}
			is.Read(is.ENDBRACKET)
		}
		nm := model.NewImage()
		img = &nm
		img.Origin = imgdata.Origin
		img.S = imgdata.S
		img.T = imgdata.T
		img.R = imgdata.R
		img.InternalTextureFormat = imgdata.InternalFormat
		img.PixelFormat = imgdata.PixelFormat
		img.DataType = imgdata.DataType
		img.Packing = imgdata.Packing
		img.Data = imgdata.Data
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
					rd := bytes.NewBuffer(dt)
					img = rw.ReadImage(rd, &opts)
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
	if loadedFromCache && name != "" {
		rw := getReaderWriter()
		f, _ := os.Open(name)
		img = rw.ReadImage(f, &opts)

		img2 := is.ReadObjectFields("osg::Image", id, nil)
		is.IdentifierMap[id] = img2
		return img2.(*model.Image)
	} else {
		img2 := is.ReadObjectFields("osg::Image", id, nil)
		img = img2.(*model.Image)
		img.Name = name
		img.WriteHint = writeHint
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
	var id int32
	is.Read(&id)
	v, ok := is.IdentifierMap[id]
	if ok {
		is.AdvanceToCurrentEndBracket()
		return v
	}
	obj = is.ReadObjectFields(cls, id, obj)
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
		inst := wrap.CreateInstanceFunc()
		obj = &inst
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
	var size int
	is.Read(&size)
	return size
}

func (is *OsgIstream) ReadComponentArray(str []byte, numElements int, numComponentsPerElements int, componentSizeInBytes int) {
	is.In.ReadComponentArray(str, numElements, numComponentsPerElements, componentSizeInBytes)
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

func (is *OsgIstream) Start(iter OsgInputIterator) (ReadType, error) {
	is.In = iter
	is.Fields = []string{}
	is.Fields = append(is.Fields, "Start")
	tp := READUNKNOWN
	switch it := iter.(type) {
	default:
		if it != nil {
			return tp, errors.New("OsgInputIterator is nil")
		}
	}
	iter.SetInputSteam(is)

	header := StreamHeader{}

	if iter.IsBinary() {
		is.Read((*int)(&header.Type))
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
		is.Read(header.TypeString)
		if header.TypeString == "Scene" {
			header.Type = READSCENE
		} else if header.TypeString == "Image" {
			header.Type = READIMAGE
		} else if header.TypeString == "Object" {
			header.Type = READOBJECT
		}
		v := model.ObjectProperty{Name: "#Version"}
		is.Read(&v)
		g := model.ObjectProperty{Name: "#Generator"}
		is.Read(&g)
		is.Read(header.OsgName)
		is.Read(header.OsgVersion)
		for {
			if is.MatchString("#CustomDomain") {
				header.DomainName = ""
				is.Read(header.DomainName)
				is.Read(header.DomainVersion)
				is.DomainVersionMap[header.DomainName] = header.DomainVersion
				break
			}
		}
	}
	is.FileVersion = header.Version
	l := len(is.Fields)
	is.Fields = is.Fields[:l-1]
	return header.Type, nil
}

func (is *OsgIstream) Decompress() {
	if is.IsBinary() {
		return
	}
	is.Fields = []string{}
	compressorName := is.ReadString()
	if compressorName != "0" {
		is.Fields = append(is.Fields, compressorName)
	}
	compressor := GetObjectWrapperManager().FindCompressor(compressorName)
	if compressor == nil {
		panic("inputstream: Failed to decompress stream, No such compressor.")
	}
	var src []byte
	compressor.DeCompress(is.In.GetIterator(), src)
	bufReader := bytes.NewBuffer(src)
	is.In.SetIterator(bufio.NewReader(bufReader))
	is.Fields = is.Fields[:len(is.Fields)-1]
	if is.UseSchemaData {
		is.Fields = append(is.Fields, "SchemaData")
	}
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
