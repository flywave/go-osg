package osg

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
	PROPERTY     *model.ObjectProperty
	BEGINBRACKET *model.ObjectMark
	ENDBRACKET   *model.ObjectMark
	CRLF         CrlfType
}

func NewOsgOstream() OsgOstream {
	p := model.NewObjectProperty()
	bb := model.NewObjectMark()
	bb.Name = "{"
	bb.IndentDelta = INDENT_VALUE
	eb := model.NewObjectMark()
	bb.Name = "}"
	bb.IndentDelta = -INDENT_VALUE

	osg := OsgOstream{PROPERTY: &p, BEGINBRACKET: &bb, ENDBRACKET: &eb}
	return osg
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
