package io

type OsgOstreamOptions struct {
	OsgOptions
	UseRobustBinaryFormat bool
	CompressorName        string
	WriteImageHint        string
	Domains               string
	TargetFileVersion     string
}

type OsgOstream struct{}


func (os *OsgOstream) Write(inter interface{}) {
}

func (os *OsgOstream) GetFileVersion(domain string) int {
	return 0
}
