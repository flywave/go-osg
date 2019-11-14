package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func getSAUpdateCallback(obj interface{}) interface{} {
	return &obj.(*model.StateAttribute).UpdateCallback
}
func setSAUpdateCallback(obj interface{}, val interface{}) {
	obj.(*model.StateAttribute).UpdateCallback = val.(*model.Callback)
}

func getSAEventCallback(obj interface{}) interface{} {
	return &obj.(*model.StateAttribute).EventCallback
}

func setSAEventCallback(obj interface{}, val interface{}) {
	obj.(*model.StateAttribute).EventCallback = val.(*model.Callback)

}

func init() {
	wrap := io.NewObjectWrapper2("StateAttribute", "flywave::osg::stateattribute", nil, "osg::Object osg::StateAttribute")
	ser1 := io.NewObjectSerializer("UpdateCallback", getSAUpdateCallback, setSAUpdateCallback)
	ser2 := io.NewObjectSerializer("EventCallback", getSAEventCallback, setSAEventCallback)

	wrap.AddSerializer(&ser1, io.RW_OBJECT)
	wrap.AddSerializer(&ser2, io.RW_OBJECT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
