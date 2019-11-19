package serializer

import (
	"github.com/flywave/go-osg"
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
	wrap := osg.NewObjectWrapper2("StateAttribute", "flywave::osg::stateattribute", nil, "osg::Object osg::StateAttribute")
	ser1 := osg.NewObjectSerializer("UpdateCallback", getSAUpdateCallback, setSAUpdateCallback)
	ser2 := osg.NewObjectSerializer("EventCallback", getSAEventCallback, setSAEventCallback)

	wrap.AddSerializer(&ser1, osg.RWOBJECT)
	wrap.AddSerializer(&ser2, osg.RWOBJECT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
