package osg

import (
	"bufio"
	"encoding/binary"

	"github.com/flywave/go-osg/model"
)

type OsgInputIterator interface {
	IsBinary() bool
	ReadBool(b *bool)
	ReadChar(c *int8)
	ReadUChar(c *uint8)
	ReadShort(s *int16)
	ReadUShort(us *uint16)
	ReadInt(i *int32)
	ReadUInt(i *uint32)
	ReadLong(l *int64)
	ReadULong(ul *uint64)
	ReadFloat(f *float32)
	ReadDouble(d *float64)
	ReadString() string
	ReadGlenum(value *model.ObjectGlenum)
	ReadProperty(prop *model.ObjectProperty)
	ReadMark(mark *model.ObjectMark)
	ReadCharArray(size int) []byte
	ReadWrappedString(str *string)
	MatchString(str string) bool
	AdvanceToCurrentEndBracket()
	SetInputSteam(is *OsgIstream)
	GetIterator() *bufio.Reader
	SetIterator(*bufio.Reader)
	SetSupportBinaryBrackets(sbb bool)
}

type InputIterator struct {
	In                    *bufio.Reader
	InputStream           *OsgIstream
	ByteSwap              int
	SupportBinaryBrackets bool
	Failed                bool
}

func NewInputIterator(in *bufio.Reader, bysp int) *InputIterator {
	return &InputIterator{In: in, ByteSwap: bysp, SupportBinaryBrackets: false, Failed: false}
}

func (it *InputIterator) SetSupportBinaryBrackets(sbb bool) {
	it.SupportBinaryBrackets = sbb
}
func (it *InputIterator) SetInputSteam(is *OsgIstream) {
	it.InputStream = is
}

func (iter *InputIterator) MatchString(str string) bool {
	return false
}
func (iter *InputIterator) GetIterator() *bufio.Reader {
	return iter.In
}

func (iter *InputIterator) SetIterator(bf *bufio.Reader) {
	iter.In = bf
}

type BinaryInputIterator struct {
	InputIterator
	Offset         int64
	BeginPositions []int64
	BlockSizes     []int64
}

func NewBinaryInputIterator(in *bufio.Reader, bysp int) *BinaryInputIterator {
	it := NewInputIterator(in, bysp)
	bf := make([]byte, 8)
	in.Read(bf)
	return &BinaryInputIterator{InputIterator: *it}
}

func (iter *BinaryInputIterator) IsBinary() bool {
	return true
}

func (iter *BinaryInputIterator) readData(val interface{}, size int) {
	binary.Read(iter.In, binary.LittleEndian, val)
	iter.Offset += int64(size)
}

func (iter *BinaryInputIterator) ReadCharArray(s int) []byte {
	arry := make([]byte, s)
	iter.readData(arry, s)
	return arry
}

func (iter *BinaryInputIterator) ReadBool(b *bool) {
	iter.readData(b, model.BOOLSIZE)
}

func (iter *BinaryInputIterator) ReadChar(b *int8) {
	iter.readData(b, model.CHARSIZE)
}

func (iter *BinaryInputIterator) ReadUChar(b *uint8) {
	iter.readData(b, model.CHARSIZE)
}

func (iter *BinaryInputIterator) ReadShort(val *int16) {
	iter.readData(val, model.SHORTSIZE)
}

func (iter *BinaryInputIterator) ReadUShort(val *uint16) {
	iter.readData(val, model.SHORTSIZE)
}

func (iter *BinaryInputIterator) ReadInt(val *int32) {
	iter.readData(val, model.INTSIZE)
}

func (iter *BinaryInputIterator) ReadUInt(val *uint32) {
	iter.readData(val, model.INTSIZE)
}

func (iter *BinaryInputIterator) ReadLong(val *int64) {
	iter.readData(val, model.LONGSIZE)
}

func (iter *BinaryInputIterator) ReadULong(val *uint64) {
	iter.readData(val, model.LONGSIZE)
}

func (iter *BinaryInputIterator) ReadFloat(val *float32) {
	iter.readData(val, model.FLOATSIZE)
}

func (iter *BinaryInputIterator) ReadDouble(val *float64) {
	iter.readData(val, model.DOUBLESIZE)
}

func (iter *BinaryInputIterator) ReadString() string {
	var size int32
	iter.ReadInt(&size)
	return string(iter.ReadCharArray(int(size)))
}

func (iter *BinaryInputIterator) ReadGlenum(val *model.ObjectGlenum) {
	var c int32
	iter.ReadInt(&c)
	val.Value = c
}

func (iter *BinaryInputIterator) ReadProperty(val *model.ObjectProperty) {
	if val.MapProperty {
		var c int32
		iter.ReadInt(&c)
		val.Value = c
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
		} else if mark.Name == "}" && len(iter.BeginPositions) > 0 {
			iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
			iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
		}
	}
}

func (iter *BinaryInputIterator) ReadWrappedString(str *string) {
	*str = iter.ReadString()
}

func (iter *BinaryInputIterator) AdvanceToCurrentEndBracket() {
	l := len(iter.BeginPositions)
	if iter.SupportBinaryBrackets && l > 0 {
		pos := iter.BeginPositions[l-1]
		bs := len(iter.BlockSizes)
		pos += iter.BlockSizes[bs-1]
		skip := pos - iter.Offset
		iter.Offset = pos
		iter.In.Discard(int(skip))
		iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
		iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
	}
}
