package io

import (
	"bufio"
	"encoding/binary"
	"strings"

	"github.com/flywave/go-osg/model"
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
	ReadString() string
	ReadGlenum(value *model.ObjectGlenum)
	ReadProperty(prop *model.ObjectProperty)
	ReadMark(mark *model.ObjectMark)
	ReadCharArray(s *string, size int)
	ReadWrappedString(str *string)
	ReadComponentArray(s *string, numElements int, numComponentsPerElements int, componentSizeInBytes int)
	MatchString(str string) bool
	AdvanceToCurrentEndBracket()
	SetInputSteam(is *OsgIstream)
}

type InputIterator struct {
	In                    bufio.Reader
	InputStream           *OsgIstream
	ByteSwap              int
	SupportBinaryBrackets bool
	Failed                bool
}

func NewInputIterator(in bufio.Reader, bysp int) InputIterator {
	return InputIterator{In: in, ByteSwap: bysp, SupportBinaryBrackets: false, Failed: false}
}

func (it *InputIterator) SetInputSteam(is *OsgIstream) {
	it.InputStream = is
}
func (it *InputIterator) ReadComponentArray(s string, numElements int, numComponentsPerElements int, componentSizeInBytes int) {
}

type BinaryInputIterator struct {
	InputIterator
	Offset         int64
	BeginPositions []int64
	BlockSizes     []int64
}

func NewBinaryInputIterator(in bufio.Reader, bysp int) BinaryInputIterator {
	it := NewInputIterator(in, bysp)
	return BinaryInputIterator{InputIterator: it}
}

func (iter *BinaryInputIterator) readData(val interface{}, size int) {
	binary.Read(&iter.In, binary.LittleEndian, val)
	iter.Offset += int64(size)
}

func (iter *BinaryInputIterator) ReadCharArray(str *string, s int) {
	buf := make([]byte, s)
	iter.readData(buf, s)
	*str = string(buf)
}

func (it *BinaryInputIterator) ReadComponentArray(s *string, numElements int, numComponentsPerElements int, componentSizeInBytes int) {
	size := numElements * numComponentsPerElements * componentSizeInBytes
	if size > 0 {
		var str string
		it.ReadCharArray(&str, size)
		build := strings.Builder{}
		if it.ByteSwap > 0 && componentSizeInBytes > 1 {
			for i := numElements - 1; i >= 0; i-- {
				for j := numComponentsPerElements - 1; j >= 0; j-- {
					build.WriteByte(str[i*j])
				}
			}
		}
		*s = build.String()
	}
}

func (iter *BinaryInputIterator) ReadBool(b *bool) {
	iter.readData(b, model.BOOL_SIZE)
}

func (iter *BinaryInputIterator) ReadChar(b *byte) {
	iter.readData(b, model.CHAR_SIZE)
}

func (iter *BinaryInputIterator) ReadUChar(b *uint8) {
	iter.readData(b, model.CHAR_SIZE)
}

func (iter *BinaryInputIterator) ReadShort(val *int16) {
	iter.readData(val, model.SHORT_SIZE)
}

func (iter *BinaryInputIterator) ReadUShort(val *uint16) {
	iter.readData(val, model.SHORT_SIZE)
}

func (iter *BinaryInputIterator) ReadInt(val *int32) {
	iter.readData(val, model.INT_SIZE)
}

func (iter *BinaryInputIterator) ReadUInt(val *uint32) {
	iter.readData(val, model.INT_SIZE)
}

func (iter *BinaryInputIterator) ReadLong(val *int64) {
	iter.readData(val, model.LONG_SIZE)
}

func (iter *BinaryInputIterator) ReadULong(val *uint64) {
	iter.readData(val, model.LONG_SIZE)
}

func (iter *BinaryInputIterator) ReadFloat(val *float32) {
	iter.readData(val, model.FLOAT_SIZE)
}

func (iter *BinaryInputIterator) ReadDouble(val *float64) {
	iter.readData(val, model.DOUBLE_SIZE)
}

func (iter *BinaryInputIterator) ReadString(val *string) {
	var size int32
	iter.ReadInt(&size)
	iter.ReadCharArray(val, int(size))
}

func (iter *BinaryInputIterator) ReadGlenum(val *model.ObjectGlenum) {
	var c int32
	iter.ReadInt(&c)
	val.Value = int(c)
}

func (iter *BinaryInputIterator) ReadObjectProperty(val *model.ObjectProperty) {
	if val.MapProperty {
		var c int32
		iter.ReadInt(&c)
		val.Value = int(c)
	} else {
		val.Value = 0
	}
}

func (iter *BinaryInputIterator) ReadMark(mark *model.ObjectMark) {
	if iter.SupportBinaryBrackets {
		if mark.Name == "{" {
			iter.BeginPositions = append(iter.BeginPositions, iter.Offset)
			if iter.InputStream.FileVersion > 148 {
				var size int64
				iter.ReadLong(&size)
				iter.BlockSizes = append(iter.BlockSizes, size)
			} else {
				var size int32
				iter.ReadInt(&size)
				iter.BlockSizes = append(iter.BlockSizes, int64(size))
			}
		} else if mark.Name == "}" && len(iter.BlockSizes) > 0 {
			iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
			iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
		}
	}
}

func (iter *BinaryInputIterator) ReadWrappedString(str *string) {
	iter.ReadString(str)
}

func (iter *BinaryInputIterator) AdvanceToCurrentEndBracket() {
	l := len(iter.BeginPositions)
	if iter.SupportBinaryBrackets && l > 0 {
		pos := iter.BeginPositions[l-1]
		bs := len(iter.BlockSizes)
		pos += iter.BlockSizes[bs-1]
		skip := pos - iter.Offset
		iter.Offset = pos
	}
}
