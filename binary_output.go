package osg

import (
	"bufio"
	"encoding/binary"

	"github.com/flywave/go-osg/model"
)

type OsgOutputIterator interface {
	IsBinary() bool
	WriteBool(b bool)
	WriteChar(c int8)
	WriteUChar(c uint8)
	WriteShort(s int16)
	WriteUShort(us uint16)
	WriteInt(i int32)
	WriteUInt(i uint32)
	WriteLong(l int64)
	WriteULong(ul uint64)
	WriteFloat(f float32)
	WriteDouble(d float64)
	WriteString(*string)
	WriteGlenum(value *model.ObjectGlenum)
	WriteProperty(prop *model.ObjectProperty)
	WriteMark(mark *model.ObjectMark)
	WriteCharArray([]byte)
	GetIterator() *bufio.Writer
	SetIterator(*bufio.Writer)
	SetSupportBinaryBrackets(sbb bool)
}

type OutputIterator struct {
	RootStream            *bufio.Writer
	Out                   *bufio.Writer
	OutputStream          *OsgOstream
	SupportBinaryBrackets bool
}

func NewOutputIterator(wt *bufio.Writer) OutputIterator {
	return OutputIterator{SupportBinaryBrackets: false, Out: wt, RootStream: wt}
}

func (it *OutputIterator) IsBinary() bool {
	return false
}

func (it *OutputIterator) SetSupportBinaryBrackets(sbb bool) {
	it.SupportBinaryBrackets = sbb
}

func (it *OutputIterator) SetOutputSteam(os *OsgOstream) {
	it.OutputStream = os
}

func (iter *OutputIterator) GetIterator() *bufio.Writer {
	return iter.Out
}

func (iter *OutputIterator) SetIterator(bw *bufio.Writer) {
	iter.Out = bw
}

type BinaryOutputIterator struct {
	OutputIterator
}

func NewBinaryOutputIterator(wt *bufio.Writer) BinaryOutputIterator {
	ot := NewOutputIterator(wt)
	return BinaryOutputIterator{OutputIterator: ot}
}

func (it *BinaryOutputIterator) writerData(iter interface{}) {
	binary.Write(it.Out, binary.LittleEndian, iter)
}

func (it *BinaryOutputIterator) WriteBool(b bool) {
	it.writerData(b)
}

func (it *BinaryOutputIterator) WriteChar(val int8) {
	it.writerData(val)
}

func (it *BinaryOutputIterator) WriteShort(val int16) {
	it.writerData(val)

}

func (it *BinaryOutputIterator) WriteInt(val int32) {
	it.writerData(val)

}

func (it *BinaryOutputIterator) WriteLong(val int64) {
	it.writerData(val)

}

func (it *BinaryOutputIterator) WriteUChar(val uint8) {
	it.writerData(val)

}
func (it *BinaryOutputIterator) WriteUShort(val uint16) {
	it.writerData(val)
}

func (it *BinaryOutputIterator) WriteUInt(val uint32) {
	it.writerData(val)
}

func (it *BinaryOutputIterator) WriteULong(val uint64) {
	it.writerData(val)
}

func (it *BinaryOutputIterator) WriteFloat(val float32) {
	it.writerData(val)
}

func (it *BinaryOutputIterator) WriteDouble(val float64) {
	it.writerData(val)
}

func (it *BinaryOutputIterator) WriteString(val *string) {
	str := *val
	it.writerData([]byte(str))
}

func (it *BinaryOutputIterator) WriteGlenum(value *model.ObjectGlenum) {
	it.writerData(value.Value)
}

func (it *BinaryOutputIterator) WriteProperty(value *model.ObjectProperty) {
	if value.MapProperty {
		it.writerData(value.Value)
	}
}

func (it *BinaryOutputIterator) WriteMark(mark *model.ObjectMark) {
	if it.SupportBinaryBrackets {
		if it.OutputStream != nil && it.OutputStream.FileVersion > 148 {

		}
	}
}

func (it *BinaryOutputIterator) WriteCharArray(s []byte) {
	it.writerData(s)
}
