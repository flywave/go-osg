package osg

import (
	"github.com/flywave/go-osg/model"
)

func getTextureDepth(obj interface{}) interface{} {
	td := obj.(*model.Texture)
	return &td.TextureDepth
}

func setTextureDepth(obj interface{}, val interface{}) {
	td := obj.(*model.Texture)
	td.TextureDepth = val.(uint32)
}

func init() {
	fn := func() interface{} {
		td := model.NewTexture3d()
		return td
	}
	wrap := NewObjectWrapper("Texture3D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture3D")
	ser1 := NewImageSerializer("Image", getImage, setImage)
	ser2 := NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	ser3 := NewPropByValSerializer("TextureHeight", false, getTexHeight, setTexHeight)
	ser4 := NewPropByValSerializer("TextureDepth", false, getTextureDepth, setTextureDepth)

	wrap.AddSerializer(ser1, RWIMAGE)
	wrap.AddSerializer(ser2, RWUINT)
	wrap.AddSerializer(ser3, RWUINT)
	wrap.AddSerializer(ser4, RWUINT)
	GetObjectWrapperManager().AddWrap(wrap)
}
