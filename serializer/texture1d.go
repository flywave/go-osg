package serializer

import (
	"github.com/flywave/go-osg"
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
	td.TextureWidth = val.(uint64)
}

func getTexHeight(obj interface{}) interface{} {
	td := obj.(*model.Texture)
	return &td.TextureWidth
}

func setTexHeight(obj interface{}, val interface{}) {
	td := obj.(*model.Texture)
	td.TextureWidth = val.(uint64)
}

func init() {
	fn := func() interface{} {
		td := model.NewTexture1d()
		return &td
	}
	wrap := osg.NewObjectWrapper("Texture1D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture1D")
	ser1 := osg.NewImageSerializer("Image", getImage, setRectImage)
	ser2 := osg.NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	wrap.AddSerializer(&ser1, osg.RWIMAGE)
	wrap.AddSerializer(&ser2, osg.RWUINT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
