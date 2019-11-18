package io

import "github.com/flywave/go-osg/model"

type CrlfType struct{}

type OsgOstreamOptions struct {
	OsgOptions
	UseRobustBinaryFormat bool
	CompressorName        string
	WriteImageHint        string
	Domains               string
	TargetFileVersion     string
}

type OsgOstream struct {
	PROPERTY      *model.ObjectProperty
	BEGIN_BRACKET *model.ObjectMark
	END_BRACKET   *model.ObjectMark
	CRLF          CrlfType
}

func (os *OsgOstream) Write(inter interface{}) {
}

func (os *OsgOstream) GetFileVersion(domain string) int32 {
	return 0
}

func (os *OsgOstream) IsBinary() bool {
	return true
}

func (os *OsgOstream) WriteWrappedString(str string) {
}
