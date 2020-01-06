package osg

import (
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
	t.TextureWidth = val.(uint32)
}

func getTextureHeightRec(obj interface{}) interface{} {
	t := obj.(*model.TextureRectangle)
	return &t.TextureHeight
}
func setTextureHeightRec(obj interface{}, val interface{}) {
	t := obj.(*model.TextureRectangle)
	t.TextureHeight = val.(uint32)
}

func init() {
	fn := func() interface{} {
		tg := model.NewTextureRectangle()
		return tg
	}
	wrap := NewObjectWrapper("TextureRectangle", fn, "osg::Object osg::StateAttribute osg::Texture osg::TextureRectangle")
	ser1 := NewImageSerializer("Image", getRectImage, setRectImage)
	ser2 := NewPropByValSerializer("TextureWidth", false, getTextureWidthRec, setTextureWidthRec)
	ser3 := NewPropByValSerializer("TextureHeight", false, getTextureHeightRec, setTextureHeightRec)
	wrap.AddSerializer(ser1, RWIMAGE)
	wrap.AddSerializer(ser2, RWINT)
	wrap.AddSerializer(ser3, RWINT)
	GetObjectWrapperManager().AddWrap(wrap)

}
