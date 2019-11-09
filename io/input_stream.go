package io

import (
	"bufio"
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
	ReadString(s *string)
	ReadGlenum(value *model.ObjectGlenum)
	ReadProperty(prop *model.ObjectProperty)
	ReadMark(mark *model.ObjectMark)
	ReadCharArray(s string, size *uint)
	ReadWrappedString(str string)
	MatchString(str string)
	AdvanceToCurrentEndBracket()
}

type InputIterator struct {
	In                    bufio.Reader
	InputStream           *OsgIstream
	ByteSwap              int
	SupportBinaryBrackets bool
	Failed                bool
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

func (iter *AsciiInputIterator) ReadCharArray(str string, s *uint) {
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

type BinaryInputIterator struct {
	InputIterator
	Offset         uint64
	BeginPositions []int64
	BlockSizes     []int64
}

type OsgOptions struct {
	FileType   string
	Precision  int
	Compressed bool
}

type OsgIstreamOptions struct {
	OsgOptions
	DbPath            string
	Domain            string
	ForceReadingImage bool
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
}

func NewOsgIstream() OsgIstream {
	return OsgIstream{ArrayMap: make(map[uint]*model.Array), IdentifierMap: make(map[uint]*model.Object), DomainVersionMap: make(map[string]int)}
}

func (is *OsgIstream) IsBinary() bool {
	return false
}

func (is *OsgIstream) MatchString(str string) bool {
	return false
}

func (is *OsgIstream) Read(inter interface{}) {
	switch val := inter.(type) {
	case *bool:
		is.In.ReadBool(val)
	}
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
