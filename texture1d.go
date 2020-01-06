package osg

import (
	"github.com/flywave/go-osg/model"
)

func getImage(obj interface{}) interface{} {
	td := obj.(*model.Texture)
	return td.Image
}

func setImage(obj interface{}, val interface{}) {
	td := obj.(*model.Texture)
	td.Image = val.(*model.Image)
}

func getTexWidth(obj interface{}) interface{} {
	td := obj.(*model.Texture)
	return &td.TextureWidth
}

func setTexWidth(obj interface{}, val interface{}) {
	td := obj.(*model.Texture)
	td.TextureWidth = val.(uint32)
}

func getTexHeight(obj interface{}) interface{} {
	td := obj.(*model.Texture)
	return &td.TextureWidth
}

func setTexHeight(obj interface{}, val interface{}) {
	td := obj.(*model.Texture)
	td.TextureWidth = val.(uint32)
}

func init() {
	fn := func() interface{} {
		td := model.NewTexture1d()
		return td
	}
	wrap := NewObjectWrapper("Texture1D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture1D")
	ser1 := NewImageSerializer("Image", getImage, setImage)
	ser2 := NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	wrap.AddSerializer(ser1, RWIMAGE)
	wrap.AddSerializer(ser2, RWUINT)
	GetObjectWrapperManager().AddWrap(wrap)
}
