package io

import (
	"io"

	"github.com/flywave/go-osg/model"
)

type OsgOptions struct {
	FileType   string
	Precision  int
	Compressed bool
}

type OsgIstreamOptions struct {
	OsgOptions
	DbPath            string
	Domain            string
	ForceReadingImage bool
}

type OsgIstream struct {
	ArrayMap          map[uint]*model.Array
	IdentifierMap     map[uint]*model.Object
	DomainVersionMap  map[uint]string
	FileVersion       int
	UseSchemaData     bool
	ForceReadingImage bool
	Fields            []string
	In                OsgInputIterator
	Options           OsgIstreamOptions
	DummyReadObject   *model.Object
	DataDecompress    io.Reader
	Data              []byte
}

func (is *OsgIstream) IsBinary() bool {
	return false
}

func (is *OsgIstream) MatchString(str string) bool {
	return false
}

func (is *OsgIstream) Read(inter interface{}) {
	switch val := inter.(type) {
	case *bool:
		is.In.ReadBool(val)
	}
}

type OsgOstream struct{}
