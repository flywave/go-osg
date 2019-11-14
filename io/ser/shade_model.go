package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func getShadeMode(obj interface{}) interface{} {
	return &obj.(*model.ShadeModel).Mode
}

func setShadeMode(obj interface{}, val interface{}) {
	obj.(*model.ShadeModel).Mode = val.(int)
}

func init() {
	fn := func() interface{} {
		sm := model.NewShadeModel()
		return &sm
	}
	wrap := io.NewObjectWrapper2("ShadeModel", "flywave::osg::shademodel", fn, "osg::Object osg::StateAttribute osg::ShadeModel")
	ser := io.NewEnumSerializer("Mode", getShadeMode, setShadeMode)
	ser.Add("FLAT", model.FLAT)
	ser.Add("SMOOTH", model.SMOOTH)
	wrap.AddSerializer(&ser, io.RW_ENUM)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
