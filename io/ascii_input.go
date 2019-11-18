package io

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/flywave/go-osg/model"
)

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
	str := iter.ReadString()
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
	str := iter.ReadString()
	res, e := strconv.ParseInt(str, 10, 32)
	if e == nil {
		*s = int16(res)
	}
}

func (iter *AsciiInputIterator) ReadUshort(us *uint16) {
	var s int16
	iter.ReadShort(&s)
	*us = uint16(s)
}

func (iter *AsciiInputIterator) ReadInt(i *int) {
	str := iter.ReadString()
	res, e := strconv.ParseInt(str, 10, 32)
	if e == nil {
		*i = int(res)
	}
}

func (iter *AsciiInputIterator) ReadUint(i *uint) {
	str := iter.ReadString()
	res, e := strconv.ParseUint(str, 10, 32)
	if e == nil {
		*i = uint(res)
	}
}

func (iter *AsciiInputIterator) ReadLong(l *int64) {
	str := iter.ReadString()
	res, e := strconv.ParseInt(str, 10, 64)
	if e == nil {
		*l = res
	}
}

func (iter *AsciiInputIterator) ReadUlong(ul *uint64) {
	str := iter.ReadString()
	res, e := strconv.ParseUint(str, 10, 64)
	if e == nil {
		*ul = res
	}
}

func (iter *AsciiInputIterator) ReadFloat(f *float32) {
	str := iter.ReadString()
	res, e := strconv.ParseFloat(str, 64)
	if e == nil {
		*f = float32(res)
	}
}

func (iter *AsciiInputIterator) ReadDouble(d *float64) {
	str := iter.ReadString()
	res, e := strconv.ParseFloat(str, 64)
	if e == nil {
		*d = res
	}
}

func (iter *AsciiInputIterator) ReadString() string {
	var str string
	l := len(iter.PreReadString)
	if l == 0 {
		str = iter.ReadWordConsumer()
	} else {
		str = iter.PreReadString
		iter.PreReadString = iter.ReadWordConsumer()
	}
	return str
}

func (iter *AsciiInputIterator) ReadGlenum(value *model.ObjectGlenum) {
	str := iter.ReadString()
	value.Value = GetObjectWrapperManager().FindLookup("GL").GetValue(str)
}

func (iter *AsciiInputIterator) ReadProperty(prop *model.ObjectProperty) {
	str := iter.ReadString()
	value := 0
	if prop.MapProperty {
		value = GetObjectWrapperManager().FindLookup(prop.Name).GetValue(str)
	} else {
		prop.Name = str
	}
	prop.Value = value
}

func (iter *AsciiInputIterator) ReadCharArray(str *string, s int) {
	return
}

func (iter *AsciiInputIterator) ReadMark(mark *model.ObjectMark) {
	iter.MarkString = iter.ReadString()
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
		str := iter.ReadString()
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

func (iter *AsciiInputIterator) ReadWrappedString(str *string) {
	var bd strings.Builder
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
				bd.WriteByte(ch)
			}
			break
		}
	} else {
		for {
			if (ch != ' ') && (ch != 0) && (ch != '\n') {
				bd.WriteByte(ch)
				iter.getCharacter(&ch)
			} else {
				break
			}
		}
	}
	*str = bd.String()
}

func (iter *AsciiInputIterator) getCharacter(c *byte) {
	l := len(iter.PreReadString)
	if l == 0 {
		iter.PreReadString = iter.ReadWordConsumer()
	}
	*c = iter.PreReadString[0]
	iter.PreReadString = iter.PreReadString[1:]
}
