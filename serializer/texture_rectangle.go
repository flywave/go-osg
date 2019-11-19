package serializer

import (
	"github.com/flywave/go-osg"
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
	wrap := osg.NewObjectWrapper2("TextureRectangle", "flywave::osg::texturerectangle", fn, "osg::Object osg::StateAttribute osg::Texture osg::TextureRectangle")
	ser1 := osg.NewImageSerializer("Image", getRectImage, setRectImage)
	ser2 := osg.NewPropByValSerializer("TextureWidth", false, getTextureWidthRec, setTextureWidthRec)
	ser3 := osg.NewPropByValSerializer("TextureHeight", false, getTextureHeightRec, setTextureHeightRec)
	wrap.AddSerializer(&ser1, osg.RWIMAGE)
	wrap.AddSerializer(&ser2, osg.RWINT)
	wrap.AddSerializer(&ser3, osg.RWINT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)

}
