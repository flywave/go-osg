package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func getTextureDepth(obj interface{}) interface{} {
	td := obj.(*model.Texture)
	return &td.TextureDepth
}

func setTextureDepth(obj interface{}, val interface{}) {
	td := obj.(*model.Texture)
	td.TextureDepth = val.(uint64)
}

func init() {
	fn := func() interface{} {
		td := model.NewTexture3d()
		return &td
	}
	wrap := osg.NewObjectWrapper("Texture3D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture1D")
	ser1 := osg.NewImageSerializer("Image", getImage, setRectImage)
	ser2 := osg.NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	ser3 := osg.NewPropByValSerializer("TextureHeight", false, getTexHeight, setTexHeight)
	ser4 := osg.NewPropByValSerializer("TextureDepth", false, getTextureDepth, setTextureDepth)

	wrap.AddSerializer(&ser1, osg.RWIMAGE)
	wrap.AddSerializer(&ser2, osg.RWUINT)
	wrap.AddSerializer(&ser3, osg.RWUINT)
	wrap.AddSerializer(&ser4, osg.RWUINT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
