package osg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"
)

func TestBinaryInputIterator_IsBinary(t *testing.T) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint32(data[0:4], OSG_HEADER_LOW)
	binary.LittleEndian.PutUint32(data[4:8], OSG_HEADER_HIGH)

	rd := bytes.NewReader(data)
	bufRd := bufio.NewReader(rd)
	iter := NewBinaryInputIterator(bufRd)

	if !iter.IsBinary() {
		t.Errorf("IsBinary() = false, want true")
	}
}

func TestBinaryInputIterator_ReadBool(t *testing.T) {
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:4], OSG_HEADER_LOW)
	binary.LittleEndian.PutUint32(data[4:8], OSG_HEADER_HIGH)
	binary.LittleEndian.PutUint32(data[8:12], 1)

	rd := bytes.NewReader(data)
	bufRd := bufio.NewReader(rd)
	iter := NewBinaryInputIterator(bufRd)

	var val bool
	iter.ReadBool(&val)

	if !val {
		t.Errorf("ReadBool() = %v, want true", val)
	}
}

func TestBinaryInputIterator_ReadInt(t *testing.T) {
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:4], OSG_HEADER_LOW)
	binary.LittleEndian.PutUint32(data[4:8], OSG_HEADER_HIGH)
	binary.LittleEndian.PutUint32(data[8:12], 12345)

	rd := bytes.NewReader(data)
	bufRd := bufio.NewReader(rd)
	iter := NewBinaryInputIterator(bufRd)

	var val int32
	iter.ReadInt(&val)

	if val != 12345 {
		t.Errorf("ReadInt() = %d, want 12345", val)
	}
}

func TestBinaryInputIterator_ReadFloat(t *testing.T) {
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:4], OSG_HEADER_LOW)
	binary.LittleEndian.PutUint32(data[4:8], OSG_HEADER_HIGH)
	binary.LittleEndian.PutUint32(data[8:12], 0x41480000)

	rd := bytes.NewReader(data)
	bufRd := bufio.NewReader(rd)
	iter := NewBinaryInputIterator(bufRd)

	var val float32
	iter.ReadFloat(&val)

	expected := float32(12.5)
	if val != expected {
		t.Errorf("ReadFloat() = %f, want %f", val, expected)
	}
}

func TestBinaryInputIterator_ReadString(t *testing.T) {
	testStr := "Hello"
	data := make([]byte, 17)
	binary.LittleEndian.PutUint32(data[0:4], OSG_HEADER_LOW)
	binary.LittleEndian.PutUint32(data[4:8], OSG_HEADER_HIGH)
	binary.LittleEndian.PutUint32(data[8:12], uint32(len(testStr)))
	copy(data[12:17], testStr)

	rd := bytes.NewReader(data)
	bufRd := bufio.NewReader(rd)
	iter := NewBinaryInputIterator(bufRd)

	val := iter.ReadString()

	if val != testStr {
		t.Errorf("ReadString() = %q, want %q", val, testStr)
	}
}

func TestBinaryInputIterator_ReadCharArray(t *testing.T) {
	testData := []byte{1, 2, 3, 4, 5}
	data := make([]byte, 13)
	binary.LittleEndian.PutUint32(data[0:4], OSG_HEADER_LOW)
	binary.LittleEndian.PutUint32(data[4:8], OSG_HEADER_HIGH)
	copy(data[8:13], testData)

	rd := bytes.NewReader(data)
	bufRd := bufio.NewReader(rd)
	iter := NewBinaryInputIterator(bufRd)

	val := iter.ReadCharArray(5)

	if !bytes.Equal(val, testData) {
		t.Errorf("ReadCharArray() = %v, want %v", val, testData)
	}
}
