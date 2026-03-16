package osg

import (
	"compress/zlib"
	"io"
)

type CompressorStream struct {
	Name string
}

func (stream *CompressorStream) Compress(st io.Writer, src []byte) error {
	w := zlib.NewWriter(st)
	if _, err := w.Write(src); err != nil {
		return err
	}
	return w.Close()
}

func (stream *CompressorStream) DeCompress(st io.Reader) ([]byte, error) {
	r, err := zlib.NewReader(st)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
