package io

import (
	"bufio"
	"io"
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

func (iter *AsciiInputIterator) ReadString(str string) {
	l := len(iter.PreReadString)
	if l == 0 {
		str = iter.ReadWordConsumer()
	} else {
		str = iter.PreReadString
		iter.PreReadString = iter.ReadWordConsumer()
	}
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

type BinaryInputIterator struct {
	InputIterator
	Offset         uint64
	BeginPositions []int64
	BlockSizes     []int64
}
