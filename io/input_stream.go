package io

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/flywave/go-osg/model"
)

const (
	FileType string = "Ascii"
)

type OsgInputIterator interface {
	IsBinary() bool
	ReadBool(b *bool)
	ReadChar(c *int8)
	ReadUchar(c *uint8)
	ReadShort(s *int16)
	ReadUshort(us *uint16)
	ReadInt(i *int)
	ReadUint(i *uint)
	ReadLong(l *int64)
	ReadUlong(ul *uint64)
	ReadFloat(f *float32)
	ReadDouble(d *float64)
	ReadString(s string)
	ReadGlenum(value *model.ObjectGlenum)
	ReadProperty(prop *model.ObjectProperty)
	ReadMark(mark *model.ObjectMark)
	ReadCharArray(s string, size int)
	ReadWrappedString(str string)
	ReadComponentArray(s string, numElements int, numComponentsPerElements int, componentSizeInBytes int)
	MatchString(str string) bool
	AdvanceToCurrentEndBracket()
	SetInputSteam(is *OsgIstream)
	ReadObject() interface{}
}

type InputIterator struct {
	In                    bufio.Reader
	InputStream           *OsgIstream
	ByteSwap              int
	SupportBinaryBrackets bool
	Failed                bool
}

func (it *InputIterator) SetInputSteam(is *OsgIstream) {
	it.InputStream = is
}
func (it *InputIterator) ReadComponentArray(s string, numElements int, numComponentsPerElements int, componentSizeInBytes int) {

}

type AsciiInputIterator struct {
	InputIterator
	PreReadString string
	MarkString    string
}

func NewAsciiInputIterator(rd io.Reader) AsciiInputIterator {
	it := InputIterator{In: *bufio.NewReader(rd), ByteSwap: 0, SupportBinaryBrackets: false, Failed: false}
	return AsciiInputIterator{InputIterator: it}
}

func (iter *AsciiInputIterator) ReadWordConsumer() string {
	var b strings.Builder
	for {
		c, e := iter.In.ReadByte()
		if e == nil {
			break
		}
		switch c {
		case ' ':
			break
		case '\n':
			break
		case '\r':
			_, e = iter.In.ReadByte()
			break
		default:
			b.WriteRune(rune(c))
		}
	}
	return b.String()
}

func (iter *AsciiInputIterator) IsBinary() bool {
	return false
}

func (iter *AsciiInputIterator) ReadBool(b *bool) {
	str := ""
	iter.ReadString(str)
	if str == "TRUE" {
		*b = true
	} else {
		*b = false
	}
}

func (iter *AsciiInputIterator) ReadChar(c *int8) {
	var s int16
	iter.ReadShort(&s)
	*c = int8(s)
}

func (iter *AsciiInputIterator) ReadUchar(c *uint8) {
	var s int16
	iter.ReadShort(&s)
	*c = uint8(s)
}

func (iter *AsciiInputIterator) ReadShort(s *int16) {
	iter.ReadShort(s)
}

func (iter *AsciiInputIterator) ReadUshort(us *uint16) {
	iter.ReadUshort(us)
}

func (iter *AsciiInputIterator) ReadInt(i *int) {
	var str string
	iter.ReadString(str)
	res, e := strconv.ParseInt(str, 10, 32)
	if e == nil {
		*i = int(res)
	}
}

func (iter *AsciiInputIterator) ReadUint(i *uint) {
	var str string
	iter.ReadString(str)
	res, e := strconv.ParseUint(str, 10, 32)
	if e == nil {
		*i = uint(res)
	}
}

func (iter *AsciiInputIterator) ReadLong(l *int64) {
	var str string
	iter.ReadString(str)
	res, e := strconv.ParseInt(str, 10, 32)
	if e == nil {
		*l = res
	}
}

func (iter *AsciiInputIterator) ReadUlong(ul *uint64) {
	var str string
	iter.ReadString(str)
	res, e := strconv.ParseUint(str, 10, 32)
	if e == nil {
		*ul = res
	}
}

func (iter *AsciiInputIterator) ReadFloat(f *float32) {
	var str string
	iter.ReadString(str)
	res, e := strconv.ParseFloat(str, 64)
	if e == nil {
		*f = float32(res)
	}
}

func (iter *AsciiInputIterator) ReadDouble(d *float64) {
	var str string
	iter.ReadString(str)
	res, e := strconv.ParseFloat(str, 64)
	if e == nil {
		*d = res
	}
}

func (iter *AsciiInputIterator) ReadString(str string) {
	l := len(iter.PreReadString)
	if l == 0 {
		str = iter.ReadWordConsumer()
	} else {
		str = iter.PreReadString
		iter.PreReadString = iter.ReadWordConsumer()
	}
}

func (iter *AsciiInputIterator) ReadGlenum(value *model.ObjectGlenum) {
}

func (iter *AsciiInputIterator) ReadProperty(prop *model.ObjectProperty) {}

func (iter *AsciiInputIterator) ReadCharArray(str string, s int) {
}

func (iter *AsciiInputIterator) ReadMark() {
	iter.ReadString(iter.MarkString)
}

func (iter *AsciiInputIterator) MatchString(str string) bool {
	l := len(iter.PreReadString)
	if l == 0 {
		iter.PreReadString = iter.ReadWordConsumer()
	}
	if iter.PreReadString == str {
		iter.PreReadString = ""
		return true
	}
	return false
}

func (iter *AsciiInputIterator) AdvanceToCurrentEndBracket() {
	blocks := 0
	for {
		var str string
		iter.ReadString(str)
		if len(str) == 0 {
			break
		}
		if str == "}" {
			if blocks <= 0 {
				return
			} else {
				blocks--
			}
		} else if str == "{" {
			blocks++
		}
	}
}

func (iter *AsciiInputIterator) ReadWrappedString(str string) {
	var ch byte
	iter.getCharacter(&ch)
	for {
		if ch == ' ' || (ch == '\n') || (ch == '\r') {
			iter.getCharacter(&ch)
			break
		}
	}
	if ch == '"' {
		iter.getCharacter(&ch)
		for {
			if ch != '"' {
				if ch == '\\' {
					iter.getCharacter(&ch)
				}
				str = string(append([]byte(str), ch))
			}
			break
		}
	} else {
		for {
			if (ch != ' ') && (ch != 0) && (ch != '\n') {
				str = string(append([]byte(str), ch))
				iter.getCharacter(&ch)
			} else {
				break
			}
		}
	}
}

func (iter *AsciiInputIterator) getCharacter(c *byte) {
	l := len(iter.PreReadString)
	if l == 0 {
		iter.PreReadString = iter.ReadWordConsumer()
	}
	*c = iter.PreReadString[0]
	iter.PreReadString = iter.PreReadString[1:]
}

func (iter *AsciiInputIterator) ReadObject() interface{} {
	return nil
}

type BinaryInputIterator struct {
	InputIterator
	Offset         uint64
	BeginPositions []int64
	BlockSizes     []int64
}

func (iter *BinaryInputIterator) ReadObject() interface{} {
	return nil
}

func (iter *BinaryInputIterator) ReadCharArray(str string, s int) {
}

func (it *BinaryInputIterator) ReadComponentArray(s string, numElements int, numComponentsPerElements int, componentSizeInBytes int) {
	size := numElements * numComponentsPerElements * componentSizeInBytes
	if size > 0 {
		it.ReadCharArray(s, size)
		build := strings.Builder{}
		if it.ByteSwap > 0 && componentSizeInBytes > 1 {
			for i := numElements - 1; i >= 0; i-- {
				for j := numComponentsPerElements - 1; j >= 0; j-- {
					build.WriteByte(s[i*j])
				}
			}
		}
		s = build.String()
	}
}

type OsgOptions struct {
	FileType   string
	Precision  int
	Compressed bool
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
	IdentifierMap     map[uint]*model.Object
	DomainVersionMap  map[string]int
	FileVersion       int
	UseSchemaData     bool
	ForceReadingImage bool
	Fields            []string
	In                OsgInputIterator
	Options           OsgIstreamOptions
	DummyReadObject   *model.Object
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
	case string:
		is.In.ReadString(val)
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

func (is *OsgIstream) ReadCharArray(str string, size int) {
	is.In.ReadCharArray(str, size)
}

func (is *OsgIstream) ReadWrappedString(str string) {
	is.In.ReadWrappedString(str)
}

func (is *OsgIstream) ReadObject() interface{} {
	return is.In.ReadObject()
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
