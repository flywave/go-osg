package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func getBinding(obj interface{}) interface{} {
	return obj.(*model.Array).Binding
}

func setBinding(obj interface{}, fc interface{}) {
	obj.(*model.Array).Binding = *fc.(*int32)
}

func getNormalize(obj interface{}) interface{} {
	return obj.(*model.Array).Normalize
}

func setNormalize(obj interface{}, fc interface{}) {
	obj.(*model.Array).Normalize = *fc.(*bool)
}

func getPreserveDataType(obj interface{}) interface{} {
	return obj.(*model.Array).PreserveDataType
}

func setPreserveDataType(obj interface{}, fc interface{}) {
	obj.(*model.Array).PreserveDataType = *fc.(*bool)
}

func init() {
	fn := func() interface{} {
		ay := model.NewArray2()
		return &ay
	}
	wrap := osg.NewObjectWrapper("Array", fn, "osg::Object osg::BufferData osg::Array")
	osg.AddUpdateWrapperVersionProxy(&wrap, 147)
	wrap.MarkSerializerAsAdded("osg::BufferData")

	ser := osg.NewEnumSerializer("Binding", getBinding, setBinding)
	ser.Add("BINDUNDEFINED", model.BINDUNDEFINED)
	ser.Add("BINDOFF", model.BINDOFF)
	ser.Add("BINDOVERALL", model.BINDOVERALL)
	ser.Add("BINDPERPRIMITIVESET", model.BINDPERPRIMITIVESET)
	ser.Add("BINDPERVERTEX", model.BINDPERVERTEX)

	serb1 := osg.NewPropByValSerializer("Normalize", false, getNormalize, setNormalize)
	serb2 := osg.NewPropByValSerializer("PreserveDataType", false, getPreserveDataType, setPreserveDataType)

	wrap.AddSerializer(&ser, osg.RWENUM)
	wrap.AddSerializer(&serb1, osg.RWBOOL)
	wrap.AddSerializer(&serb2, osg.RWBOOL)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
