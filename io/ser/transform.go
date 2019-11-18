package ser

import (
	"github.com/flywave/go-osg/io"
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
	wrap := io.NewObjectWrapper("Transform", fn, "osg::Object osg::Node osg::Group osg::Transform")
	ser1 := io.NewEnumSerializer("ReferenceFrame", getReferenceFrame, setReferenceFrame)
	ser1.Add("RELATIVE_RF", model.RELATIVE_RF)
	ser1.Add("ABSOLUTE_RF", model.ABSOLUTE_RF)
	ser1.Add("ABSOLUTE_RF_INHERIT_VIEWPOINT", model.ABSOLUTE_RF_INHERIT_VIEWPOINT)
	wrap.AddSerializer(&ser1, io.RW_ENUM)
}
