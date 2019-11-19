package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func init() {
	fn := func() interface{} {
		td := model.NewTexture2d()
		return &td
	}
	wrap := osg.NewObjectWrapper("Texture2D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture1D")
	ser1 := osg.NewImageSerializer("Image", getImage, setRectImage)
	ser2 := osg.NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	ser3 := osg.NewPropByValSerializer("TextureHeight", false, getTexHeight, setTexHeight)

	wrap.AddSerializer(&ser1, osg.RWIMAGE)
	wrap.AddSerializer(&ser2, osg.RWUINT)
	wrap.AddSerializer(&ser3, osg.RWUINT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
