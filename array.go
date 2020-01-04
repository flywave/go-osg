package osg

import (
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
		return ay
	}
	wrap := NewObjectWrapper("Array", fn, "osg::Object osg::BufferData osg::Array")
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 147)
		wrap.MarkSerializerAsAdded("osg::BufferData")
		uv.SetLastVersion()
	}

	ser := NewEnumSerializer("Binding", getBinding, setBinding)
	ser.Add("BINDUNDEFINED", model.BINDUNDEFINED)
	ser.Add("BINDOFF", model.BINDOFF)
	ser.Add("BINDOVERALL", model.BINDOVERALL)
	ser.Add("BINDPERPRIMITIVESET", model.BINDPERPRIMITIVESET)
	ser.Add("BINDPERVERTEX", model.BINDPERVERTEX)

	serb1 := NewPropByValSerializer("Normalize", false, getNormalize, setNormalize)
	serb2 := NewPropByValSerializer("PreserveDataType", false, getPreserveDataType, setPreserveDataType)

	wrap.AddSerializer(ser, RWENUM)
	wrap.AddSerializer(serb1, RWBOOL)
	wrap.AddSerializer(serb2, RWBOOL)
	GetObjectWrapperManager().AddWrap(wrap)
}
