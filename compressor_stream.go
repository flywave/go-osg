package osg

import (
	"compress/zlib"
	"io"
)

type CompressorStream struct {
	Name string
}

func (stream *CompressorStream) Compress(st io.Writer, src []byte) {
	w := zlib.NewWriter(st)
	w.Write(src)
	w.Close()
}

func (stream *CompressorStream) DeCompress(st io.Reader) ([]byte, error) {
	r, _ := zlib.NewReader(st)
	var src []byte
	for {
		buf := make([]byte, 4096)
		n, e := io.ReadFull(r, buf)
		if e != nil {
			return nil, e
		}
		if n != 0 {
			src = append(src, buf[0:n]...)
		}
	}
	return src, nil
}
