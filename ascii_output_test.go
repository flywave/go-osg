package osg

import (
	"bufio"
	"bytes"
	"testing"
)

func TestAsciiOutputIterator_IsBinary(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewAsciiOutputIterator(wt)

	if iter.IsBinary() {
		t.Errorf("IsBinary() = true, want false")
	}
}

func TestAsciiOutputIterator_WriteBool(t *testing.T) {
	tests := []struct {
		value    bool
		expected string
	}{
		{true, "TRUE"},
		{false, "FALSE"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			buf := &bytes.Buffer{}
			wt := bufio.NewWriter(buf)
			iter := NewAsciiOutputIterator(wt)

			iter.WriteBool(tt.value)
			wt.Flush()

			result := buf.String()
			if result != tt.expected+" " {
				t.Errorf("WriteBool(%v) = %q, want %q", tt.value, result, tt.expected+" ")
			}
		})
	}
}

func TestAsciiOutputIterator_WriteInt(t *testing.T) {
	tests := []struct {
		value    int32
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{-456, "-456"},
		{2147483647, "2147483647"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			buf := &bytes.Buffer{}
			wt := bufio.NewWriter(buf)
			iter := NewAsciiOutputIterator(wt)

			iter.WriteInt(tt.value)
			wt.Flush()

			result := buf.String()
			if result != tt.expected+" " {
				t.Errorf("WriteInt(%d) = %q, want %q", tt.value, result, tt.expected+" ")
			}
		})
	}
}

func TestAsciiOutputIterator_WriteFloat(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewAsciiOutputIterator(wt)

	iter.WriteFloat(3.14159)
	wt.Flush()

	result := buf.String()
	if result == "" {
		t.Errorf("WriteFloat() produced empty output")
	}
}

func TestAsciiOutputIterator_WriteString(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"Hello", "Hello"},
		{"World", "World"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			buf := &bytes.Buffer{}
			wt := bufio.NewWriter(buf)
			iter := NewAsciiOutputIterator(wt)

			iter.WriteString(&tt.value)
			wt.Flush()

			result := buf.String()
			if result != tt.expected+" " {
				t.Errorf("WriteString(%q) = %q, want %q", tt.value, result, tt.expected+" ")
			}
		})
	}
}

func TestAsciiOutputIterator_WriteWrappedString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "simple string",
			value:    "Hello",
			expected: `"Hello" `,
		},
		{
			name:     "string with quote",
			value:    `Hello "World"`,
			expected: `"Hello \"World\"" `,
		},
		{
			name:     "string with backslash",
			value:    `Hello\World`,
			expected: `"Hello\\World" `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			wt := bufio.NewWriter(buf)
			iter := NewAsciiOutputIterator(wt)

			iter.WriteWrappedString(&tt.value)
			wt.Flush()

			result := buf.String()
			if result != tt.expected {
				t.Errorf("WriteWrappedString(%q) = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestAsciiOutputIterator_WriteCharArray(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewAsciiOutputIterator(wt)

	data := []byte{1, 2, 3, 4, 5}
	iter.WriteCharArray(data)
	wt.Flush()

	result := buf.String()
	if result != "" {
		t.Errorf("WriteCharArray() should produce empty output in ASCII mode, got %q", result)
	}
}

func TestAsciiOutputIterator_Indent(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewAsciiOutputIterator(wt)

	iter.Indent = 4
	iter.ReadyForIndent = true
	iter.WriteBool(true)
	wt.Flush()

	result := buf.String()
	if result != "    TRUE " {
		t.Errorf("Indent handling failed, got %q, want %q", result, "    TRUE ")
	}
}

func TestAsciiOutputIterator_MultipleWrites(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewAsciiOutputIterator(wt)

	iter.WriteInt(1)
	iter.WriteInt(2)
	iter.WriteInt(3)
	wt.Flush()

	result := buf.String()
	expected := "1 2 3 "
	if result != expected {
		t.Errorf("Multiple writes = %q, want %q", result, expected)
	}
}
