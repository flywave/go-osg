package io

type OsgOstreamOptions struct {
	OsgOptions
	UseRobustBinaryFormat bool
	CompressorName        string
	WriteImageHint        string
	Domains               string
	TargetFileVersion     string
}
