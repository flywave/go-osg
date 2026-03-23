package osg

import (
	"github.com/flywave/go-osg/model"
)

func getSAUpdateCallback(obj interface{}) interface{} {
	return &obj.(*model.StateAttribute).UpdateCallback
}
func setSAUpdateCallback(obj interface{}, val interface{}) {
	if val == nil {
		obj.(*model.StateAttribute).UpdateCallback = nil
	} else {
		obj.(*model.StateAttribute).UpdateCallback = val.(*model.Callback)
	}
}

func getSAEventCallback(obj interface{}) interface{} {
	return &obj.(*model.StateAttribute).EventCallback
}

func setSAEventCallback(obj interface{}, val interface{}) {
	if val == nil {
		obj.(*model.StateAttribute).EventCallback = nil
	} else {
		obj.(*model.StateAttribute).EventCallback = val.(*model.Callback)
	}

}

func init() {
	wrap := NewObjectWrapper("StateAttribute", nil, "osg::Object osg::StateAttribute")
	ser1 := NewObjectSerializer("UpdateCallback", getSAUpdateCallback, setSAUpdateCallback)
	ser2 := NewObjectSerializer("EventCallback", getSAEventCallback, setSAEventCallback)

	wrap.AddSerializer(ser1, RWOBJECT)
	wrap.AddSerializer(ser2, RWOBJECT)
	GetObjectWrapperManager().AddWrap(wrap)
}
