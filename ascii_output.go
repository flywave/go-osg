package osg

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/flywave/go-osg/model"
)

type AsciiOutputIterator struct {
	OutputIterator
	ReadyForIndent bool
	Indent         int
}

func NewAsciiOutputIterator(wt *bufio.Writer) *AsciiOutputIterator {
	ot := NewOutputIterator(wt)
	return &AsciiOutputIterator{OutputIterator: *ot, ReadyForIndent: false, Indent: 0}
}

func (it *AsciiOutputIterator) writeBlank() {
	it.Out.WriteString(" ")
}

func (it *AsciiOutputIterator) indentIfRequired() {
	if it.ReadyForIndent {
		str := ""
		for i := 0; i < it.Indent; i++ {
			str += " "
		}
		it.Out.WriteString(str)
		it.ReadyForIndent = false
	}
}

func (it *AsciiOutputIterator) WriteBool(b bool) {
	it.indentIfRequired()
	if b {
		it.Out.WriteString("TRUE ")
	} else {
		it.Out.WriteString("FALSE ")
	}
}

func (it *AsciiOutputIterator) WriteChar(val int8) {
	it.Out.WriteString(strconv.Itoa(int(val)))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteShort(val int16) {
	it.Out.WriteString(strconv.Itoa(int(val)))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteInt(val int32) {
	it.Out.WriteString(strconv.Itoa(int(val)))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteLong(val int64) {
	it.Out.WriteString(strconv.Itoa(int(val)))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteUChar(val uint8) {
	it.Out.WriteString(strconv.FormatUint(uint64(val), 10))
	it.writeBlank()
}
func (it *AsciiOutputIterator) WriteUShort(val uint16) {
	it.Out.WriteString(strconv.FormatUint(uint64(val), 10))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteUInt(val uint32) {
	it.Out.WriteString(strconv.FormatUint(uint64(val), 10))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteULong(val uint64) {
	it.Out.WriteString(strconv.FormatUint(uint64(val), 10))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteFloat(val float32) {
	it.Out.WriteString(strconv.FormatFloat(float64(val), 'f', -1, 32))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteDouble(val float64) {
	it.Out.WriteString(strconv.FormatFloat(val, 'f', -1, 64))
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteString(val *string) {
	it.Out.WriteString(*val)
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteGlenum(value *model.ObjectGlenum) {
	str := GetObjectWrapperManager().FindLookup("GL").GetString(value.Value)
	it.WriteString(&str)
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteProperty(value *model.ObjectProperty) {
	str := GetObjectWrapperManager().FindLookup(value.Name).GetString(value.Value)
	it.WriteString(&str)
	it.writeBlank()
}

func (it *AsciiOutputIterator) WriteCharArray(s []byte) {}

func (it *AsciiOutputIterator) WriteMark(mark *model.ObjectMark) {
	it.Indent += int(mark.IndentDelta)
	it.indentIfRequired()
	it.WriteString(&mark.Name)
	it.writeBlank()
}
func (it *AsciiOutputIterator) WriteWrappedString(str *string) {
	b := strings.Builder{}
	for _, c := range *str {
		if byte(c) == '"' {
			b.WriteByte('\\')
		} else if byte(c) == '\\' {
			b.WriteByte('\\')
		}
		b.WriteByte(byte(c))
	}
	s := "\"" + b.String()
	s += "\""
	it.WriteString(&s)
}
