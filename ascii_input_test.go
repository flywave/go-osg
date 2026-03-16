package osg

import (
	"bufio"
	"strings"
	"testing"
)

func TestAsciiInputIterator_IsBinary(t *testing.T) {
	input := ""
	rd := strings.NewReader(input)
	bufRd := bufio.NewReader(rd)
	iter := NewAsciiInputIterator(bufRd)

	if iter.IsBinary() {
		t.Errorf("IsBinary() = true, want false")
	}
}

func TestAsciiInputIterator_ReadBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"TRUE", true},
		{"FALSE", false},
		{"1", false},
		{"0", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rd := strings.NewReader(tt.input)
			bufRd := bufio.NewReader(rd)
			iter := NewAsciiInputIterator(bufRd)

			var val bool
			iter.ReadBool(&val)

			if val != tt.expected {
				t.Errorf("ReadBool() = %v, want %v", val, tt.expected)
			}
		})
	}
}

func TestAsciiInputIterator_ReadInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int32
	}{
		{"123", 123},
		{"-456", -456},
		{"0", 0},
		{"2147483647", 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rd := strings.NewReader(tt.input)
			bufRd := bufio.NewReader(rd)
			iter := NewAsciiInputIterator(bufRd)

			var val int32
			iter.ReadInt(&val)

			if val != tt.expected {
				t.Errorf("ReadInt() = %d, want %d", val, tt.expected)
			}
		})
	}
}

func TestAsciiInputIterator_ReadFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float32
	}{
		{"1.5", 1.5},
		{"-2.5", -2.5},
		{"0.0", 0.0},
		{"3.14159", 3.14159},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rd := strings.NewReader(tt.input)
			bufRd := bufio.NewReader(rd)
			iter := NewAsciiInputIterator(bufRd)

			var val float32
			iter.ReadFloat(&val)

			if val != tt.expected {
				t.Errorf("ReadFloat() = %f, want %f", val, tt.expected)
			}
		})
	}
}

func TestAsciiInputIterator_ReadString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"World", "World"},
		{"Test", "Test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rd := strings.NewReader(tt.input)
			bufRd := bufio.NewReader(rd)
			iter := NewAsciiInputIterator(bufRd)

			val := iter.ReadString()

			if val != tt.expected {
				t.Errorf("ReadString() = %q, want %q", val, tt.expected)
			}
		})
	}
}

func TestAsciiInputIterator_ReadWrappedString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "simple string",
			input:    "Hello",
			contains: "Hello",
		},
		{
			name:     "single word quoted",
			input:    `"Hello"`,
			contains: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := strings.NewReader(tt.input)
			bufRd := bufio.NewReader(rd)
			iter := NewAsciiInputIterator(bufRd)

			var val string
			iter.ReadWrappedString(&val)

			if !strings.Contains(val, tt.contains) {
				t.Errorf("ReadWrappedString() = %q, should contain %q", val, tt.contains)
			}
		})
	}
}

func TestAsciiInputIterator_MatchString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected bool
	}{
		{
			name:     "match",
			input:    "Hello World",
			pattern:  "Hello",
			expected: true,
		},
		{
			name:     "no match",
			input:    "Hello World",
			pattern:  "World",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := strings.NewReader(tt.input)
			bufRd := bufio.NewReader(rd)
			iter := NewAsciiInputIterator(bufRd)

			result := iter.MatchString(tt.pattern)

			if result != tt.expected {
				t.Errorf("MatchString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAsciiInputIterator_AdvanceToCurrentEndBracket(t *testing.T) {
	input := "{ nested { } } } rest"
	rd := strings.NewReader(input)
	bufRd := bufio.NewReader(rd)
	iter := NewAsciiInputIterator(bufRd)

	iter.AdvanceToCurrentEndBracket()

	remaining, _ := bufRd.ReadString(0)
	if !strings.HasPrefix(strings.TrimSpace(remaining), "rest") {
		t.Errorf("AdvanceToCurrentEndBracket() did not advance correctly, remaining: %q", remaining)
	}
}

func TestAsciiInputIterator_SkipWhitespace(t *testing.T) {
	input := "   \n\r  Hello"
	rd := strings.NewReader(input)
	bufRd := bufio.NewReader(rd)
	iter := NewAsciiInputIterator(bufRd)

	val := iter.ReadString()

	if val != "Hello" {
		t.Errorf("skip() then ReadString() = %q, want %q", val, "Hello")
	}
}

func TestAsciiInputIterator_ReadMultipleValues(t *testing.T) {
	input := "1 2 3 4 5"
	rd := strings.NewReader(input)
	bufRd := bufio.NewReader(rd)
	iter := NewAsciiInputIterator(bufRd)

	var vals [5]int32
	for i := 0; i < 5; i++ {
		iter.ReadInt(&vals[i])
	}

	expected := [5]int32{1, 2, 3, 4, 5}
	if vals != expected {
		t.Errorf("ReadMultipleValues() = %v, want %v", vals, expected)
	}
}

func TestAsciiInputIterator_ReadCharArray(t *testing.T) {
	rd := strings.NewReader("")
	bufRd := bufio.NewReader(rd)
	iter := NewAsciiInputIterator(bufRd)

	result := iter.ReadCharArray(5)
	if len(result) != 0 {
		t.Errorf("ReadCharArray() should return empty slice for ASCII mode, got %v", result)
	}
}
