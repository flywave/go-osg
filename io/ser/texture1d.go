package ser

import (
	"github.com/flywave/go-osg/io"

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
	wrap := io.NewObjectWrapper("Texture1D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture1D")
	ser1 := io.NewImageSerializer("Image", getImage, setRectImage)
	ser2 := io.NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	wrap.AddSerializer(&ser1, io.RW_IMAGE)
	wrap.AddSerializer(&ser2, io.RW_UINT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
