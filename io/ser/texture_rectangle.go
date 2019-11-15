package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func getRectImage(obj interface{}) interface{} {
	img := obj.(*model.TextureRectangle)
	return img.Image
}

func setRectImage(obj interface{}, val interface{}) {
	img := obj.(*model.TextureRectangle)
	img.Image = val.(*model.Image)
}

func getTextureWidthRec(obj interface{}) interface{} {
	t := obj.(*model.TextureRectangle)
	return &t.TextureWidth
}
func setTextureWidthRec(obj interface{}, val interface{}) {
	t := obj.(*model.TextureRectangle)
	t.TextureWidth = val.(uint64)
}

func getTextureHeightRec(obj interface{}) interface{} {
	t := obj.(*model.TextureRectangle)
	return &t.TextureHeight
}
func setTextureHeightRec(obj interface{}, val interface{}) {
	t := obj.(*model.TextureRectangle)
	t.TextureHeight = val.(uint64)
}

func init() {
	fn := func() interface{} {
		tg := model.NewTextureRectangle()
		return &tg
	}
	wrap := io.NewObjectWrapper2("TextureRectangle", "flywave::osg::texturerectangle", fn, "osg::Object osg::StateAttribute osg::Texture osg::TextureRectangle")
	ser1 := io.NewImageSerializer("Image", getRectImage, setRectImage)
	ser2 := io.NewPropByValSerializer("TextureWidth", false, getTextureWidthRec, setTextureWidthRec)
	ser3 := io.NewPropByValSerializer("TextureHeight", false, getTextureHeightRec, setTextureHeightRec)
	wrap.AddSerializer(&ser1, io.RW_IMAGE)
	wrap.AddSerializer(&ser2, io.RW_INT)
	wrap.AddSerializer(&ser3, io.RW_INT)
	io.GetObjectWrapperManager().AddWrap(&wrap)

}
