package io

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

func (stream *CompressorStream) DeCompress(st io.Reader, src []byte) {
	r, _ := zlib.NewReader(st)
	io.ReadFull(r, src)
}
