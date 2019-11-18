package io

import (
	"bufio"
	"errors"
	"image"
	"io"
	"os"
	"strings"

	"github.com/flywave/go-osg/model"
)

const (
	FEATURENONE             = 0
	FEATUREREADOBJECT       = 1 << 0
	FEATUREREADIMAGE        = 1 << 1
	FEATUREREADHEIGHTFIELD  = 1 << 2
	FEATUREREADNODE         = 1 << 3
	FEATUREREADSHADER       = 1 << 4
	FEATUREWRITEOBJECT      = 1 << 5
	FEATUREWRITEIMAGE       = 1 << 6
	FEATUREWRITEHEIGHTFIELD = 1 << 7
	FEATUREWRITENODE        = 1 << 8
	FEATUREWRITESHADER      = 1 << 9
	FEATUREREADSCRIPT       = 1 << 10
	FEATUREWRITESCRIPT      = 1 << 11
	FEATUREALL              = FEATUREREADOBJECT | FEATUREREADIMAGE |
		FEATUREREADHEIGHTFIELD | FEATUREREADNODE |
		FEATUREREADSHADER | FEATUREREADSCRIPT |
		FEATUREWRITEOBJECT | FEATUREWRITEIMAGE |
		FEATUREWRITEHEIGHTFIELD | FEATUREWRITENODE |
		FEATUREWRITESHADER | FEATUREWRITESCRIPT

	NOTIMPLEMENTED           = 0
	FILENOTHANDLED           = 1
	FILENOTFOUND             = 2
	ERRORINREADINGFILE       = 3
	FILELOADED               = 4
	FILELOADEDFROMCACHE      = 5
	FILEREQUESTED            = 6
	INSUFFICIENTMEMORYTOLOAD = 7

	ERRORINWRITINGFILE = 2
	FILESAVED          = 3

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
	case NOTIMPLEMENTED:
		description += "not implemented"
		break
	case FILENOTHANDLED:
		description += "file not handled"
		break
	case ERRORINWRITINGFILE:
		description += "write error"
		break
	case FILESAVED:
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

	if e != nil || t != READUNKNOWN {
		return &ReadResult{Status: FILENOTHANDLED}
	}
	ty, e := is.Start(iter)
	if e != nil {
		if ty == READUNKNOWN {
			return &ReadResult{Status: FILENOTHANDLED}
		}
		is.Decompress()
		obj := is.ReadObject(nil)
		if obj == nil {
			return &ReadResult{Status: FILENOTHANDLED}
		}
		return &ReadResult{Object: obj}
	}
	return nil
}

func (rw *ReadWrite) ReadImage(rd io.Reader, opts *OsgIstreamOptions) *model.Image { //TODO process image
	img := model.NewImage()
	mg, _, e := image.Decode(rd)
	if e == nil {
		img.S = int32(mg.Bounds().Max.X - mg.Bounds().Min.X)
		img.RowLength = img.S
		img.T = int32(mg.Bounds().Max.Y - mg.Bounds().Min.Y)
		img.R = 1
		img.PixelFormat = 0
	}
	return &img
}
