package io

import "strings"

const (
	FEATURE_NONE               = 0
	FEATURE_READ_OBJECT        = 1 << 0
	FEATURE_READ_IMAGE         = 1 << 1
	FEATURE_READ_HEIGHT_FIELD  = 1 << 2
	FEATURE_READ_NODE          = 1 << 3
	FEATURE_READ_SHADER        = 1 << 4
	FEATURE_WRITE_OBJECT       = 1 << 5
	FEATURE_WRITE_IMAGE        = 1 << 6
	FEATURE_WRITE_HEIGHT_FIELD = 1 << 7
	FEATURE_WRITE_NODE         = 1 << 8
	FEATURE_WRITE_SHADER       = 1 << 9
	FEATURE_READ_SCRIPT        = 1 << 10
	FEATURE_WRITE_SCRIPT       = 1 << 11
	FEATURE_ALL                = FEATURE_READ_OBJECT | FEATURE_READ_IMAGE |
		FEATURE_READ_HEIGHT_FIELD | FEATURE_READ_NODE |
		FEATURE_READ_SHADER | FEATURE_READ_SCRIPT |
		FEATURE_WRITE_OBJECT | FEATURE_WRITE_IMAGE |
		FEATURE_WRITE_HEIGHT_FIELD | FEATURE_WRITE_NODE |
		FEATURE_WRITE_SHADER | FEATURE_WRITE_SCRIPT

	NOT_IMPLEMENTED             = 0
	FILE_NOT_HANDLED            = 1
	FILE_NOT_FOUND              = 2
	ERROR_IN_READING_FILE       = 3
	FILE_LOADED                 = 4
	FILE_LOADED_FROM_CACHE      = 5
	FILE_REQUESTED              = 6
	INSUFFICIENT_MEMORY_TO_LOAD = 7

	ERROR_IN_WRITING_FILE = 2
	FILE_SAVED            = 3

	READ   = 0
	WRITE  = 1
	CREATE = 2
)

type ReadWrite struct {
	SupportedProtocal   map[string]string
	SupportedExtensions map[string]string
	SupportedOptions    map[string]string
}

func NewReadWrite() ReadWrite {
	rw := ReadWrite{SupportedProtocal: make(map[string]string), SupportedExtensions: make(map[string]string), SupportedOptions: make(map[string]string)}
	rw.SupportedExtensions["osg2"] = "OpenSceneGraph extendable format"
	rw.SupportedExtensions["osgt"] = "OpenSceneGraph extendable ascii format"
	rw.SupportedExtensions["osgb"] = "OpenSceneGraph extendable binary format"

	rw.SupportedExtensions["jpg"] = "jpg image format"
	rw.SupportedExtensions["jpeg"] = "jpeg image format"
	rw.SupportedExtensions["png"] = "png image format"
	rw.SupportedExtensions["bmp"] = "bitmap image format"

	rw.SupportedOptions["Ascii"] =
		"Import/Export option: Force reading/writing ascii file"
	rw.SupportedOptions["ForceReadingImage"] =
		"Import option: Load an empty image instead if required file missed"
	rw.SupportedOptions["SchemaData"] =
		"Export option: Record inbuilt schema data into a binary file"
	rw.SupportedOptions["SchemaFile=<file>"] =
		"Import/Export option: Use/Record an ascii schema file"
	rw.SupportedOptions["Compressor=<name>"] =
		"Export option: Use an inbuilt or user-defined compressor"
	rw.SupportedOptions["WriteImageHint=<hint>"] =
		"Export option: Hint of writing image to stream: <IncludeData> writes Image::data[) directly; <IncludeFile> writes the image file itself to stream; <UseExternal> writes only the filename; <WriteOut> writes Image::data[) to disk as external file."
	return rw
}

func (rw *ReadWrite) AcceptsExtension(ext string) bool {
	e := strings.ToLower(ext)
	_, ok := rw.SupportedExtensions[e]
	return ok
}
func (rw *ReadWrite) AcceptsProtocol(pro string) bool {
	p := strings.ToLower(pro)
	_, ok := rw.SupportedProtocal[p]
	return ok
}

func (rw *ReadWrite) SupportExtension(fmt string, desc string) {
	e := strings.ToLower(fmt)
	rw.SupportedExtensions[e] = desc
}

func (rw *ReadWrite) SupportProtocol(fmt string, desc string) {
	e := strings.ToLower(fmt)
	rw.SupportedProtocal[e] = desc
}

func (rw *ReadWrite) SupportOption(fmt string, desc string) {
	rw.SupportedOptions[fmt] = desc
}
