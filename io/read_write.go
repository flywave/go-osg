package io

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"os"
	"strings"

	"github.com/flywave/go-osg/model"
)

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

type ReadResult struct {
	Status  int
	Message string
	Object  interface{}
}

func (res *ReadResult) GetObject() *model.Object {
	switch o := res.Object.(type) {
	case *model.Object:
		return o
	}
	return nil
}

func (res *ReadResult) GetImage() *model.Image {
	switch o := res.Object.(type) {
	case *model.Image:
		return o
	}
	return nil
}

func (res *ReadResult) GetNode() *model.Node {
	switch o := res.Object.(type) {
	case *model.Node:
		return o
	}
	return nil
}

func (res *ReadResult) StatusMessage() string {
	var description string
	switch res.Status {
	case NOT_IMPLEMENTED:
		description += "not implemented"
		break
	case FILE_NOT_HANDLED:
		description += "file not handled"
		break
	case ERROR_IN_WRITING_FILE:
		description += "write error"
		break
	case FILE_SAVED:
		description += "file saved"
		break
	}

	if len(res.Message) != 0 {
		description += " (" + res.Message + ")"
	}
	return description
}

type ReadWrite struct {
	SupportedProtocal   map[string]string
	SupportedExtensions map[string]string
	SupportedOptions    map[string]string
}

var rw *ReadWrite

func getReaderWriter() *ReadWrite {
	if rw == nil {
		rw = &ReadWrite{}
	}
	return rw
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

func (rw *ReadWrite) ReadInputIterator(reader *bufio.Reader, op *OsgIstreamOptions) OsgInputIterator {
	extensionIsAscii := false
	if op != nil {
		if op.FileType == "Ascii" {
			extensionIsAscii = true
		}
	}
	if extensionIsAscii {
		head := make([]byte, 6, 6)
		reader.Read(head)
		if string(head) == "#Ascii" {
			rd := NewAsciiInputIterator(reader)
			return &rd
		} else {
			return nil
		}
	} else {
		rd := NewBinaryInputIterator(reader, 1)
		return &rd
	}
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

func (rw *ReadWrite) PrepareReading(fname string, op *OsgIstreamOptions) (*OsgIstreamOptions, error) {
	subs := strings.Split(fname, ".")
	ext := subs[len(subs)-1]
	if !rw.AcceptsExtension(ext) {
		return nil, errors.New("not support")
	}
	if op == nil {
		o := NewOsgIstreamOptions()
		op = &o
	}

	if ext == "osgt" {
		op.FileType = "Ascii"
	} else if ext == "osgb" {
		op.FileType = "Binary"
	} else if ext == "jpg" || ext == "jpeg" {
		op.FileType = "JPEG"
	} else if ext == "png" {
		op.FileType = "PNG"
	} else if ext == "bmp" {
		op.FileType = "BMP"
	} else {
		op.FileType = ""
	}
	return op, nil
}

func (rw *ReadWrite) OpenReader(fname string) *bufio.Reader {
	f, e := os.Open(fname)
	if e != nil {
		return nil
	}
	rd := bufio.NewReader(f)
	return rd
}

func (rw *ReadWrite) OpenWriter(fname string) *bufio.Writer {
	return nil
}

func (rw *ReadWrite) ReadObject(fname string, opt *OsgIstreamOptions) *ReadResult {
	lopt, e := rw.PrepareReading(fname, opt)
	if e != nil {
		return nil
	}
	in := rw.OpenReader(fname)
	return rw.ReadObjectWithReader(bufio.NewReader(in), lopt)
}

func (rw *ReadWrite) ReadObjectWithReader(rd *bufio.Reader, opt *OsgIstreamOptions) *ReadResult {
	iter := rw.ReadInputIterator(rd, opt)
	is := NewOsgIstream(opt)
	t, e := is.Start(iter)

	if e != nil || t != READ_UNKNOWN {
		return &ReadResult{Status: FILE_NOT_HANDLED}
	}
	ty, e := is.Start(iter)
	if e != nil {
		if ty == READ_UNKNOWN {
			return &ReadResult{Status: FILE_NOT_HANDLED}
		}
		is.Decompress()
		obj := is.ReadObject(nil)
		if obj == nil {
			return &ReadResult{Status: FILE_NOT_HANDLED}
		}
		return &ReadResult{Object: obj}
	}
	return nil
}

func (rw *ReadWrite) ReadImage(data []byte, opts *OsgIstreamOptions) *model.Image {
	img := model.NewImage()
	rd := bytes.NewBuffer(data)
	mg, ty, e := image.Decode(rd)
	if e == nil {
		img.S = int32(mg.Bounds().Max.X - mg.Bounds().Min.X)
		img.T = int32(mg.Bounds().Max.Y - mg.Bounds().Min.Y)
	}
}
