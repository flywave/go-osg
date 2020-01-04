package osg

import (
	"github.com/flywave/go-osg/model"
)

func init() {
	fn := func() interface{} {
		td := model.NewTexture2d()
		return td
	}
	wrap := NewObjectWrapper("Texture2D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture2D")
	ser1 := NewImageSerializer("Image", getImage, setImage)
	ser2 := NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	ser3 := NewPropByValSerializer("TextureHeight", false, getTexHeight, setTexHeight)

	wrap.AddSerializer(ser1, RWIMAGE)
	wrap.AddSerializer(ser2, RWUINT)
	wrap.AddSerializer(ser3, RWUINT)
	GetObjectWrapperManager().AddWrap(wrap)
}
