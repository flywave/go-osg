package io

import (
	"errors"
	"io"

	"github.com/flywave/go-osg/model"
)

const (
	FileType string = "Ascii"
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
	READ_UNKNOWN ReadType = 0
	READ_SCENE   ReadType = 1
	READ_IMAGE   ReadType = 2
	READ_OBJECT  ReadType = 3
)

type OsgIstreamOptions struct {
	OsgOptions
	DbPath            string
	Domain            string
	ForceReadingImage bool
}

type StreamHeader struct {
	Version       int
	Type          ReadType
	Attributes    int
	NumDomains    int
	DomainName    string
	DomainVersion int
	TypeString    string
	OsgName       string
	OsgVersion    string
}

type OsgIstream struct {
	ArrayMap          map[uint]*model.Array
	IdentifierMap     map[uint]interface{}
	DomainVersionMap  map[string]int
	FileVersion       int
	UseSchemaData     bool
	ForceReadingImage bool
	Fields            []string
	In                OsgInputIterator
	Options           OsgIstreamOptions
	DummyReadObject   interface{}
	DataDecompress    io.Reader
	Data              []byte
	CRLF              CrlfType

	PROPERTY      *model.ObjectProperty
	BEGIN_BRACKET *model.ObjectMark
	END_BRACKET   *model.ObjectMark
}

func NewOsgIstream() OsgIstream {
	p := model.NewObjectProperty()
	bb := model.NewObjectMark()
	eb := model.NewObjectMark()
	return OsgIstream{ArrayMap: make(map[uint]*model.Array), IdentifierMap: make(map[uint]*model.Object), DomainVersionMap: make(map[string]int), PROPERTY: &p, BEGIN_BRACKET: &bb, END_BRACKET: &eb}
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
		is.In.ReadUchar(val)
		break
	case *int16:
		is.In.ReadShort(val)
		break
	case *uint16:
		is.In.ReadUshort(val)
		break
	case *int:
		is.In.ReadInt(val)
		break
	case *uint:
		is.In.ReadUint(val)
		break
	case *int64:
		is.In.ReadLong(val)
		break
	case *uint64:
		is.In.ReadUlong(val)
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

func (is *OsgIstream) ReadCharArray(str *string, size int) {
	is.In.ReadCharArray(str, size)
}

func (is *OsgIstream) ReadWrappedString(str *string) {
	is.In.ReadWrappedString(str)
}

func (is *OsgIstream) ReadString() string {
	str := is.In.ReadString()
	return str
}

func (is *OsgIstream) ReadObject() interface{} {
	cls := is.ReadString()
	if cls == "NULL" {
		return nil
	}
	is.Read(is.BEGIN_BRACKET)
	is.PROPERTY.Name = "UniqueID"
	is.Read(is.PROPERTY)
	var id uint32
	is.Read(&id)
	v, ok := is.IdentifierMap[uint(id)]
	if ok {
		is.AdvanceToCurrentEndBracket()
		return v
	}
	obj := is.ReadObjectFields(cls, id)
	is.AdvanceToCurrentEndBracket()
	return obj
}

func (is *OsgIstream) ReadObjectFields(str string, id uint32) interface{} {

}

func (is *OsgIstream) ReadSize() int {
	var size int
	is.Read(&size)
	return size
}

func (is *OsgIstream) ReadImage() *model.Image {
	img := model.NewImage()
	is.Read(&img)
	return &img
}

func (is *OsgIstream) ReadComponentArray(str *string, numElements int, numComponentsPerElements int, componentSizeInBytes int) {
	is.In.ReadComponentArray(str, numElements, numComponentsPerElements, componentSizeInBytes)
}

func (is *OsgIstream) GetFileVersion(domain string) int {
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
	tp := READ_UNKNOWN
	switch it := iter.(type) {
	default:
		if it != nil {
			return tp, errors.New("OsgInputIterator is nil")
		}
	}
	iter.SetInputSteam(is)

	header := StreamHeader{}

	if iter.IsBinary() {
		return tp, nil
	} else {
		is.Read(header.TypeString)
		if header.TypeString == "Scene" {
			header.Type = READ_SCENE
		} else if header.TypeString == "Image" {
			header.Type = READ_IMAGE
		} else if header.TypeString == "Object" {
			header.Type = READ_OBJECT
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
		is.FileVersion = header.Version
		l := len(is.Fields)
		is.Fields = is.Fields[0 : l-1]
		return header.Type, nil
	}
}
