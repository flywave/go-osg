package ser

import (
	"github.com/flywave/go-osg/io"

	"github.com/flywave/go-osg/model"
)

func init() {
	fn := func() interface{} {
		td := model.NewTexture1d()
		return &td
	}
	wrap := io.NewObjectWrapper("Texture2D", fn, "osg::Object osg::StateAttribute osg::Texture osg::Texture1D")
	ser1 := io.NewImageSerializer("Image", getImage, setRectImage)
	ser2 := io.NewPropByValSerializer("TextureWidth", false, getTexWidth, setTexWidth)
	ser3 := io.NewPropByValSerializer("TextureHeight", false, getTexHeight, setTexHeight)

	wrap.AddSerializer(&ser1, io.RW_IMAGE)
	wrap.AddSerializer(&ser2, io.RW_UINT)
	wrap.AddSerializer(&ser3, io.RW_UINT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
