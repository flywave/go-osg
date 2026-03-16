package osg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"
)

func TestBinaryOutputIterator_IsBinary(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	if !iter.IsBinary() {
		t.Errorf("IsBinary() = false, want true")
	}
}

func TestBinaryOutputIterator_WriteBool(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	iter.WriteBool(true)
	wt.Flush()

	data := buf.Bytes()
	if len(data) != 1 || data[0] != 1 {
		t.Errorf("WriteBool(true) = %v, want [1]", data)
	}
}

func TestBinaryOutputIterator_WriteInt(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	iter.WriteInt(12345)
	wt.Flush()

	data := buf.Bytes()
	var val int32
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &val)

	if val != 12345 {
		t.Errorf("WriteInt(12345) = %d, want 12345", val)
	}
}

func TestBinaryOutputIterator_WriteFloat(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	iter.WriteFloat(3.14)
	wt.Flush()

	data := buf.Bytes()
	var val float32
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &val)

	if val != 3.14 {
		t.Errorf("WriteFloat(3.14) = %f, want 3.14", val)
	}
}

func TestBinaryOutputIterator_WriteString(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	str := "Hello"
	iter.WriteString(&str)
	wt.Flush()

	if buf.Len() == 0 {
		t.Errorf("WriteString() produced empty output")
	}
}

func TestBinaryOutputIterator_WriteCharArray(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	data := []byte{1, 2, 3, 4, 5}
	iter.WriteCharArray(data)
	wt.Flush()

	result := buf.Bytes()
	if !bytes.Equal(result, data) {
		t.Errorf("WriteCharArray() = %v, want %v", result, data)
	}
}

func TestBinaryOutputIterator_WriteMultipleTypes(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	iter.WriteBool(true)
	iter.WriteInt(42)
	iter.WriteFloat(2.5)
	wt.Flush()

	if buf.Len() == 0 {
		t.Errorf("WriteMultipleTypes() produced empty output")
	}
}

func TestOutputIterator_SetSupportBinaryBrackets(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewOutputIterator(wt)

	if iter.SupportBinaryBrackets {
		t.Error("SupportBinaryBrackets should be false initially")
	}

	iter.SetSupportBinaryBrackets(true)
	if !iter.SupportBinaryBrackets {
		t.Error("SupportBinaryBrackets should be true after SetSupportBinaryBrackets(true)")
	}
}

func TestMarkHelper(t *testing.T) {
	mh := MakeMarkHelper()

	if mh.buff == nil {
		t.Error("MakeMarkHelper() buff should not be nil")
	}

	if mh.Stream == nil {
		t.Error("MakeMarkHelper() Stream should not be nil")
	}

	mh.Stream.Write([]byte("test"))
	mh.Stream.Flush()

	data := mh.GetBuff()
	if !bytes.Equal(data, []byte("test")) {
		t.Errorf("GetBuff() = %v, want %v", data, []byte("test"))
	}
}
