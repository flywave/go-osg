package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func getReferenceFrame(obj interface{}) interface{} {
	tran := obj.(*model.Transform)
	return &tran.ReferenceFrame
}

func setReferenceFrame(obj interface{}, val interface{}) {
	tran := obj.(*model.Transform)
	tran.ReferenceFrame = val.(int)
}

func init() {
	fn := func() interface{} {
		td := model.NewTransform()
		return &td
	}
	wrap := osg.NewObjectWrapper("Transform", fn, "osg::Object osg::Node osg::Group osg::Transform")
	ser1 := osg.NewEnumSerializer("ReferenceFrame", getReferenceFrame, setReferenceFrame)
	ser1.Add("RELATIVERF", model.RELATIVERF)
	ser1.Add("ABSOLUTERF", model.ABSOLUTERF)
	ser1.Add("ABSOLUTERFINHERITVIEWPOINT", model.ABSOLUTERFINHERITVIEWPOINT)
	wrap.AddSerializer(&ser1, osg.RWENUM)
}
