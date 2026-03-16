package osg

import (
	"testing"
)

func TestOsgOptions(t *testing.T) {
	opts := NewOsgOptions()

	if opts.FileType != FileType {
		t.Errorf("NewOsgOptions().FileType = %q, want %q", opts.FileType, FileType)
	}
}

func TestOsgIstreamOptions(t *testing.T) {
	opts := NewOsgIstreamOptions()

	if opts.FileType != FileType {
		t.Errorf("NewOsgIstreamOptions().FileType = %q, want %q", opts.FileType, FileType)
	}
}

func TestOsgIstream_GetFileVersion(t *testing.T) {
	opts := NewOsgIstreamOptions()
	is := NewOsgIstream(opts)

	is.FileVersion = 100

	tests := []struct {
		domain   string
		expected int32
	}{
		{"", 100},
		{"unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			result := is.GetFileVersion(tt.domain)
			if result != tt.expected {
				t.Errorf("GetFileVersion(%q) = %d, want %d", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestOsgIstream_DomainVersionMap(t *testing.T) {
	opts := NewOsgIstreamOptions()
	opts.Domain = "TestDomain:50;AnotherDomain:100"
	is := NewOsgIstream(opts)

	if is.DomainVersionMap["TestDomain"] != 50 {
		t.Errorf("DomainVersionMap[\"TestDomain\"] = %d, want 50", is.DomainVersionMap["TestDomain"])
	}

	if is.DomainVersionMap["AnotherDomain"] != 100 {
		t.Errorf("DomainVersionMap[\"AnotherDomain\"] = %d, want 100", is.DomainVersionMap["AnotherDomain"])
	}
}

func TestOsgIstream_ReadSize(t *testing.T) {
	opts := NewOsgIstreamOptions()
	is := NewOsgIstream(opts)

	if is.ArrayMap == nil {
		t.Error("ArrayMap should not be nil")
	}

	if is.IdentifierMap == nil {
		t.Error("IdentifierMap should not be nil")
	}

	if is.DomainVersionMap == nil {
		t.Error("DomainVersionMap should not be nil")
	}
}

func TestTrimEnclosingSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\t\nworld\t\n", "world"},
		{"", ""},
		{"   ", ""},
		{"no-space", "no-space"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := trimEnclosingSpaces(tt.input)
			if result != tt.expected {
				t.Errorf("trimEnclosingSpaces(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
