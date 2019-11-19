package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkPosX(tex interface{}) bool {
	return true
}

func readPosX(is *osg.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVEX, image)
		is.Read(is.ENDBRACKET)
	}
}

func writePosX(os *osg.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVEX)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGINBRACKET)
		os.Write(image)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegX(tex interface{}) bool {
	return true
}
func readNegX(is *osg.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVEX, image)
		is.Read(is.ENDBRACKET)
	}
}

func writeNegX(os *osg.OsgOstream, tex interface{}) {
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
func readPosY(is *osg.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVEY, image)
		is.Read(is.ENDBRACKET)
	}
}

func writePosY(os *osg.OsgOstream, tex interface{}) {
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
func readNegY(is *osg.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVEY, image)
		is.Read(is.ENDBRACKET)
	}
}

func writeNegY(os *osg.OsgOstream, tex interface{}) {
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
func readPosZ(is *osg.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVEZ, image)
		is.Read(is.ENDBRACKET)
	}
}

func writePosZ(os *osg.OsgOstream, tex interface{}) {
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
func readNegZ(is *osg.OsgIstream, tex interface{}) {

	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGINBRACKET)
		image := is.ReadImage(true)
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVEZ, image)
		is.Read(is.ENDBRACKET)
	}
}

func writeNegZ(os *osg.OsgOstream, tex interface{}) {
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
	t.TextureWidth = val.(uint64)
}

func getTextureHeight(obj interface{}) interface{} {
	t := obj.(*model.TextureCubeMap)
	return &t.TextureHeight
}
func setTextureHeight(obj interface{}, val interface{}) {
	t := obj.(*model.TextureCubeMap)
	t.TextureHeight = val.(uint64)
}

func init() {
	fn := func() interface{} {
		tg := model.NewTextureCubeMap()
		return &tg
	}
	wrap := osg.NewObjectWrapper2("TextureCubeMap", "flywave::osg::texturecubemap", fn, "osg::Object osg::StateAttribute osg::Texture osg::TextureCubeMap")
	ser1 := osg.NewUserSerializer("PosX", checkPosX, readPosX, writePosX)
	ser2 := osg.NewUserSerializer("NegX", checkNegX, readNegX, writeNegX)
	ser3 := osg.NewUserSerializer("PosY", checkPosY, readPosY, writePosY)
	ser4 := osg.NewUserSerializer("NegY", checkNegY, readNegY, writeNegY)
	ser5 := osg.NewUserSerializer("PosZ", checkPosZ, readPosZ, writePosZ)
	ser6 := osg.NewUserSerializer("NegZ", checkNegZ, readNegZ, writeNegZ)
	ser7 := osg.NewPropByValSerializer("TextureWidth", false, getTextureWidth, setTextureWidth)
	ser8 := osg.NewPropByValSerializer("TextureHeight", false, getTextureHeight, setTextureHeight)
	wrap.AddSerializer(&ser1, osg.RWUSER)
	wrap.AddSerializer(&ser2, osg.RWUSER)
	wrap.AddSerializer(&ser3, osg.RWUSER)
	wrap.AddSerializer(&ser4, osg.RWUSER)
	wrap.AddSerializer(&ser5, osg.RWUSER)
	wrap.AddSerializer(&ser6, osg.RWUSER)

	wrap.AddSerializer(&ser7, osg.RWINT)
	wrap.AddSerializer(&ser8, osg.RWINT)

	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
