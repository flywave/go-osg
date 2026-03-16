package osg

import (
	"bufio"
	"bytes"
	"testing"
)

func TestNewOsgOstream(t *testing.T) {
	opts := &OsgOstreamOptions{}
	os := NewOsgOstream(opts)

	if os == nil {
		t.Fatal("NewOsgOstream() returned nil")
	}

	if os.TargetFileVersion != OPENSCENEGRAPHSOVERSION {
		t.Errorf("TargetFileVersion = %d, want %d", os.TargetFileVersion, OPENSCENEGRAPHSOVERSION)
	}
}

func TestOsgOstream_GetFileVersion(t *testing.T) {
	opts := &OsgOstreamOptions{}
	os := NewOsgOstream(opts)

	os.TargetFileVersion = 100

	result := os.GetFileVersion("")
	if result != 100 {
		t.Errorf("GetFileVersion(\"\") = %d, want 100", result)
	}
}

func TestOsgOstream_WriteArray_Nil(t *testing.T) {
	opts := &OsgOstreamOptions{}
	os := NewOsgOstream(opts)

	os.WriteArray(nil)
}

func TestOsgOstream_WritePrimitiveSet_Nil(t *testing.T) {
	opts := &OsgOstreamOptions{}
	os := NewOsgOstream(opts)

	os.WritePrimitiveSet(nil)
}

func TestWriteConstants(t *testing.T) {
	if WRITEUNKNOWN != 0 {
		t.Errorf("WRITEUNKNOWN = %d, want 0", WRITEUNKNOWN)
	}
	if WRITESCENE != 1 {
		t.Errorf("WRITESCENE = %d, want 1", WRITESCENE)
	}
	if WRITEIMAGE != 2 {
		t.Errorf("WRITEIMAGE = %d, want 2", WRITEIMAGE)
	}
	if WRITEOBJECT != 3 {
		t.Errorf("WRITEOBJECT = %d, want 3", WRITEOBJECT)
	}
}

func TestWriteImageHintConstants(t *testing.T) {
	if WRITEUSEIMAGEHINT != 0 {
		t.Errorf("WRITEUSEIMAGEHINT = %d, want 0", WRITEUSEIMAGEHINT)
	}
	if WRITEUSEEXTERNAL != 1 {
		t.Errorf("WRITEUSEEXTERNAL = %d, want 1", WRITEUSEEXTERNAL)
	}
	if WRITEINLINEDATA != 2 {
		t.Errorf("WRITEINLINEDATA = %d, want 2", WRITEINLINEDATA)
	}
	if WRITEINLINEFILE != 3 {
		t.Errorf("WRITEINLINEFILE = %d, want 3", WRITEINLINEFILE)
	}
	if WRITEEXTERNALFILE != 4 {
		t.Errorf("WRITEEXTERNALFILE = %d, want 4", WRITEEXTERNALFILE)
	}
}

func TestOutputIterator(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewOutputIterator(wt)

	if iter.IsBinary() {
		t.Error("OutputIterator.IsBinary() should return false")
	}

	if iter.SupportBinaryBrackets {
		t.Error("SupportBinaryBrackets should be false initially")
	}

	iter.SetSupportBinaryBrackets(true)
	if !iter.SupportBinaryBrackets {
		t.Error("SupportBinaryBrackets should be true after SetSupportBinaryBrackets(true)")
	}
}

func TestBinaryOutputIterator(t *testing.T) {
	buf := &bytes.Buffer{}
	wt := bufio.NewWriter(buf)
	iter := NewBinaryOutputIterator(wt)

	if !iter.IsBinary() {
		t.Error("BinaryOutputIterator.IsBinary() should return true")
	}
}

func TestCrlfType(t *testing.T) {
	crlf := &CrlfType{}
	if crlf == nil {
		t.Error("CrlfType should not be nil")
	}
}

func TestOsgOstreamOptions_Defaults(t *testing.T) {
	opts := &OsgOstreamOptions{}

	if opts.UseRobustBinaryFormat {
		t.Error("UseRobustBinaryFormat should be false by default")
	}
}

func TestOsgOstreamOptions_WithValues(t *testing.T) {
	opts := &OsgOstreamOptions{
		UseRobustBinaryFormat: true,
		CompressorName:        "zlib",
		WriteImageHint:        "IncludeData",
		Domains:               "Test:100",
		TargetFileVersion:     "150",
	}

	os := NewOsgOstream(opts)

	if !os.UseRobustBinaryFormat {
		t.Error("UseRobustBinaryFormat should be true")
	}

	if os.CompressorName != "zlib" {
		t.Errorf("CompressorName = %q, want %q", os.CompressorName, "zlib")
	}
}

func TestOsgOstreamOptions_RobustBinaryFormat(t *testing.T) {
	opts := &OsgOstreamOptions{UseRobustBinaryFormat: false}
	os := NewOsgOstream(opts)

	if os.UseRobustBinaryFormat {
		t.Error("UseRobustBinaryFormat should be false when set to false")
	}
}

func TestOsgOstreamOptions_WriteImageHint(t *testing.T) {
	tests := []struct {
		hint     string
		expected int32
	}{
		{"IncludeData", WRITEINLINEDATA},
		{"IncludeFile", WRITEINLINEFILE},
		{"UseExternal", WRITEUSEEXTERNAL},
		{"WriteOut", WRITEEXTERNALFILE},
		{"", WRITEUSEIMAGEHINT},
	}

	for _, tt := range tests {
		t.Run(tt.hint, func(t *testing.T) {
			opts := &OsgOstreamOptions{WriteImageHint: tt.hint}
			os := NewOsgOstream(opts)

			if os.WriteImageHint != tt.expected {
				t.Errorf("WriteImageHint = %d, want %d", os.WriteImageHint, tt.expected)
			}
		})
	}
}

func TestOsgOstreamOptions_TargetFileVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected int32
	}{
		{"100", 100},
		{"150", 150},
		{"", OPENSCENEGRAPHSOVERSION},
		{"1000", OPENSCENEGRAPHSOVERSION},
		{"-1", OPENSCENEGRAPHSOVERSION},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			opts := &OsgOstreamOptions{TargetFileVersion: tt.version}
			os := NewOsgOstream(opts)

			if os.TargetFileVersion != tt.expected {
				t.Errorf("TargetFileVersion = %d, want %d", os.TargetFileVersion, tt.expected)
			}
		})
	}
}

func TestOsgOstreamOptions_Domains(t *testing.T) {
	opts := &OsgOstreamOptions{Domains: "TestDomain:50;AnotherDomain:100"}
	os := NewOsgOstream(opts)

	if os.DomainVersionMap == nil {
		t.Fatal("DomainVersionMap should not be nil")
	}

	if os.DomainVersionMap["TestDomain"] != 50 {
		t.Errorf("DomainVersionMap[\"TestDomain\"] = %d, want 50", os.DomainVersionMap["TestDomain"])
	}

	if os.DomainVersionMap["AnotherDomain"] != 100 {
		t.Errorf("DomainVersionMap[\"AnotherDomain\"] = %d, want 100", os.DomainVersionMap["AnotherDomain"])
	}
}
