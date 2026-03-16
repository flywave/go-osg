package osg

import (
	"bytes"
	"testing"
)

func TestCompressorStream_Compress(t *testing.T) {
	stream := &CompressorStream{Name: "zlib"}

	tests := []struct {
		name    string
		input   []byte
		wantLen bool
	}{
		{
			name:    "empty data",
			input:   []byte{},
			wantLen: true,
		},
		{
			name:    "small data",
			input:   []byte{1, 2, 3, 4, 5},
			wantLen: true,
		},
		{
			name:    "large data",
			input:   make([]byte, 10000),
			wantLen: true,
		},
		{
			name:    "repetitive data",
			input:   bytes.Repeat([]byte{0xAA}, 1000),
			wantLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := stream.Compress(&buf, tt.input)
			if err != nil {
				t.Errorf("Compress() error = %v", err)
				return
			}

			if tt.wantLen && buf.Len() == 0 && len(tt.input) > 0 {
				t.Errorf("Compress() produced empty output for non-empty input")
			}
		})
	}
}

func TestCompressorStream_DeCompress(t *testing.T) {
	stream := &CompressorStream{Name: "zlib"}

	t.Run("valid compressed data", func(t *testing.T) {
		var buf bytes.Buffer
		stream.Compress(&buf, []byte{1, 2, 3, 4, 5})
		reader := bytes.NewReader(buf.Bytes())
		got, err := stream.DeCompress(reader)
		if err != nil {
			t.Errorf("DeCompress() error = %v", err)
			return
		}
		if !bytes.Equal(got, []byte{1, 2, 3, 4, 5}) {
			t.Errorf("DeCompress() = %v, want %v", got, []byte{1, 2, 3, 4, 5})
		}
	})

	t.Run("empty compressed data", func(t *testing.T) {
		var buf bytes.Buffer
		stream.Compress(&buf, []byte{})
		reader := bytes.NewReader(buf.Bytes())
		got, err := stream.DeCompress(reader)
		if err != nil {
			t.Errorf("DeCompress() error = %v", err)
			return
		}
		if len(got) != 0 {
			t.Errorf("DeCompress() = %v, want empty", got)
		}
	})

	t.Run("invalid compressed data", func(t *testing.T) {
		reader := bytes.NewReader([]byte{0, 1, 2, 3, 4, 5})
		_, err := stream.DeCompress(reader)
		if err == nil {
			t.Errorf("DeCompress() expected error for invalid data")
		}
	})
}

func TestCompressorStream_RoundTrip(t *testing.T) {
	stream := &CompressorStream{Name: "zlib"}

	testData := [][]byte{
		{},
		{1},
		{1, 2, 3, 4, 5},
		bytes.Repeat([]byte{0xFF}, 100),
		[]byte("Hello, World! This is a test string for compression."),
	}

	for i, data := range testData {
		t.Run("roundtrip", func(t *testing.T) {
			var compressed bytes.Buffer
			err := stream.Compress(&compressed, data)
			if err != nil {
				t.Errorf("Test %d: Compress() error = %v", i, err)
				return
			}

			decompressed, err := stream.DeCompress(&compressed)
			if err != nil {
				t.Errorf("Test %d: DeCompress() error = %v", i, err)
				return
			}

			if !bytes.Equal(data, decompressed) {
				t.Errorf("Test %d: RoundTrip failed, got %v, want %v", i, decompressed, data)
			}
		})
	}
}
