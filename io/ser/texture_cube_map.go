package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func checkPosX(tex interface{}) bool {
	return true
}

func readPosX(is *io.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGIN_BRACKET)
		image := is.ReadImage()
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVE_X, image)
		is.Read(is.END_BRACKET)
	}
}

func writePosX(os *io.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVE_X)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGIN_BRACKET)
		os.Write(image)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegX(tex interface{}) bool {
	return true
}
func readNegX(is *io.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGIN_BRACKET)
		image := is.ReadImage()
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVE_X, image)
		is.Read(is.END_BRACKET)
	}
}

func writeNegX(os *io.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.NEGATIVE_X)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGIN_BRACKET)
		os.Write(image)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkPosY(tex interface{}) bool {
	return true
}
func readPosY(is *io.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGIN_BRACKET)
		image := is.ReadImage()
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVE_Y, image)
		is.Read(is.END_BRACKET)
	}
}

func writePosY(os *io.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVE_Y)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGIN_BRACKET)
		os.Write(image)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegY(tex interface{}) bool {
	return true
}
func readNegY(is *io.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGIN_BRACKET)
		image := is.ReadImage()
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVE_Y, image)
		is.Read(is.END_BRACKET)
	}
}

func writeNegY(os *io.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.NEGATIVE_Y)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGIN_BRACKET)
		os.Write(image)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkPosZ(tex interface{}) bool {
	return true
}
func readPosZ(is *io.OsgIstream, tex interface{}) {
	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGIN_BRACKET)
		image := is.ReadImage()
		tex.(*model.TextureCubeMap).SetImage(model.POSITIVE_Z, image)
		is.Read(is.END_BRACKET)
	}
}

func writePosZ(os *io.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.POSITIVE_Z)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGIN_BRACKET)
		os.Write(image)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.CRLF)
}

func checkNegZ(tex interface{}) bool {
	return true
}
func readNegZ(is *io.OsgIstream, tex interface{}) {

	var hasImage bool
	is.Read(&hasImage)
	if hasImage {
		is.Read(is.BEGIN_BRACKET)
		image := is.ReadImage()
		tex.(*model.TextureCubeMap).SetImage(model.NEGATIVE_Z, image)
		is.Read(is.END_BRACKET)
	}
}

func writeNegZ(os *io.OsgOstream, tex interface{}) {
	image := tex.(*model.TextureCubeMap).GetImage(model.NEGATIVE_Z)
	os.Write(image != nil)
	if image != nil {
		os.Write(os.BEGIN_BRACKET)
		os.Write(image)
		os.Write(os.END_BRACKET)
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
	wrap := io.NewObjectWrapper2("TextureCubeMap", "flywave::osg::texturecubemap", fn, "osg::Object osg::StateAttribute osg::Texture osg::TextureCubeMap")
	ser1 := io.NewUserSerializer("PosX", checkPosX, readPosX, writePosX)
	ser2 := io.NewUserSerializer("NegX", checkNegX, readNegX, writeNegX)
	ser3 := io.NewUserSerializer("PosY", checkPosY, readPosY, writePosY)
	ser4 := io.NewUserSerializer("NegY", checkNegY, readNegY, writeNegY)
	ser5 := io.NewUserSerializer("PosZ", checkPosZ, readPosZ, writePosZ)
	ser6 := io.NewUserSerializer("NegZ", checkNegZ, readNegZ, writeNegZ)
	ser7 := io.NewPropByValSerializer("TextureWidth", false, getTextureWidth, setTextureWidth)
	ser8 := io.NewPropByValSerializer("TextureHeight", false, getTextureHeight, setTextureHeight)
	wrap.AddSerializer(&ser1, io.RW_USER)
	wrap.AddSerializer(&ser2, io.RW_USER)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)
	wrap.AddSerializer(&ser6, io.RW_USER)

	wrap.AddSerializer(&ser7, io.RW_INT)
	wrap.AddSerializer(&ser8, io.RW_INT)

	io.GetObjectWrapperManager().AddWrap(&wrap)
}
