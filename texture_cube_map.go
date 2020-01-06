package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkPosX(tex interface{}) bool {
	return true
}

func readPosX(is *OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVEX, image)
		is.Read(is.ENDBRACKET)
	}
}

func writePosX(os *OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVEX)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.WriteImage(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegX(tex interface{}) bool {
	return true
}
func readNegX(is *OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVEX, image)
		is.Read(is.ENDBRACKET)
	}
}

func writeNegX(os *OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.NEGATIVEX)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.Write(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkPosY(tex interface{}) bool {
	return true
}
func readPosY(is *OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVEY, image)
		is.Read(is.ENDBRACKET)
	}
}

func writePosY(os *OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVEY)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.Write(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegY(tex interface{}) bool {
	return true
}
func readNegY(is *OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVEY, image)
		is.Read(is.ENDBRACKET)
	}
}

func writeNegY(os *OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.NEGATIVEY)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.Write(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkPosZ(tex interface{}) bool {
	return true
}
func readPosZ(is *OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVEZ, image)
		is.Read(is.ENDBRACKET)
	}
}

func writePosZ(os *OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVEZ)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.Write(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegZ(tex interface{}) bool {
	return true
}
func readNegZ(is *OsgIstream, tex interface{}) {

	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVEZ, image)
		is.Read(is.ENDBRACKET)
	}
}

func writeNegZ(os *OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.NEGATIVEZ)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.Write(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func getTextureWidth(obj interface{}) interface{} {
	t := obj.(*model.TextureCubeMap)
	return &t.TextureWidth
}
func setTextureWidth(obj interface{}, val interface{}) {
	t := obj.(*model.TextureCubeMap)
	t.TextureWidth = val.(uint32)
}

func getTextureHeight(obj interface{}) interface{} {
	t := obj.(*model.TextureCubeMap)
	return &t.TextureHeight
}
func setTextureHeight(obj interface{}, val interface{}) {
	t := obj.(*model.TextureCubeMap)
	t.TextureHeight = val.(uint32)
}

func init() {
	fn := func() interface{} {
		tg := model.NewTextureCubeMap()
		return tg
	}
	wrap := NewObjectWrapper("TextureCubeMap", fn, "osg::Object osg::StateAttribute osg::Texture osg::TextureCubeMap")
	ser1 := NewUserSerializer("PosX", checkPosX, readPosX, writePosX)
	ser2 := NewUserSerializer("NegX", checkNegX, readNegX, writeNegX)
	ser3 := NewUserSerializer("PosY", checkPosY, readPosY, writePosY)
	ser4 := NewUserSerializer("NegY", checkNegY, readNegY, writeNegY)
	ser5 := NewUserSerializer("PosZ", checkPosZ, readPosZ, writePosZ)
	ser6 := NewUserSerializer("NegZ", checkNegZ, readNegZ, writeNegZ)
	ser7 := NewPropByValSerializer("TextureWidth", false, getTextureWidth, setTextureWidth)
	ser8 := NewPropByValSerializer("TextureHeight", false, getTextureHeight, setTextureHeight)
	wrap.AddSerializer(ser1, RWUSER)
	wrap.AddSerializer(ser2, RWUSER)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)
	wrap.AddSerializer(ser6, RWUSER)

	wrap.AddSerializer(ser7, RWINT)
	wrap.AddSerializer(ser8, RWINT)

	GetObjectWrapperManager().AddWrap(wrap)
}
