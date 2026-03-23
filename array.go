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

func getArrayData(obj interface{}) interface{} {
	return obj.(*model.Array).Data
}

func setArrayData(obj interface{}, fc interface{}) {
	obj.(*model.Array).Data = fc
}

// FIX: Add getters and setters for Type, DataType, DataSize
func getType(obj interface{}) interface{} {
	return obj.(*model.Array).Type
}

func setType(obj interface{}, fc interface{}) {
	obj.(*model.Array).Type = *fc.(*model.ArrayTable)
}

func getArrayDataType(obj interface{}) interface{} {
	return obj.(*model.Array).DataType
}

func setArrayDataType(obj interface{}, fc interface{}) {
	obj.(*model.Array).DataType = *fc.(*int32)
}

func getDataSize(obj interface{}) interface{} {
	return obj.(*model.Array).DataSize
}

func setDataSize(obj interface{}, fc interface{}) {
	obj.(*model.Array).DataSize = *fc.(*int32)
}

func registerArrayWrapper(name string, arrayType int, elemType int, elemSize int) {
	fn := func() interface{} {
		return model.NewArray2()
	}
	wrap := NewObjectWrapper(name, fn, "osg::Object osg::BufferData osg::Array "+name)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 117)
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

	// FIX: Add Type, DataType, DataSize serializers for FileVersion >= 112
	serType := NewPropByRefSerializer("Type", getType, setType)
	serDataType := NewPropByRefSerializer("DataType", getArrayDataType, setArrayDataType)
	serDataSize := NewPropByRefSerializer("DataSize", getDataSize, setDataSize)

	wrap.AddSerializer(ser, RWENUM)
	wrap.AddSerializer(serb1, RWBOOL)
	wrap.AddSerializer(serb2, RWBOOL)
	wrap.AddSerializer(serType, RWINT)
	wrap.AddSerializer(serDataType, RWINT)
	wrap.AddSerializer(serDataSize, RWINT)

	switch elemSize {
	case 1:
		serData := NewIsAVectorSerializer("Array", RWFLOAT, 1, getArrayData, setArrayData)
		wrap.AddSerializer(serData, RWVECTOR)
	case 2:
		serData := NewIsAVectorSerializer("Array", RWFLOAT, 2, getArrayData, setArrayData)
		wrap.AddSerializer(serData, RWVECTOR)
	case 3:
		serData := NewIsAVectorSerializer("Array", RWFLOAT, 3, getArrayData, setArrayData)
		wrap.AddSerializer(serData, RWVECTOR)
	case 4:
		serData := NewIsAVectorSerializer("Array", RWFLOAT, 4, getArrayData, setArrayData)
		wrap.AddSerializer(serData, RWVECTOR)
	}

	GetObjectWrapperManager().AddWrap(wrap)
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

	// FIX: Add Type, DataType, DataSize serializers for FileVersion >= 112
	serType := NewPropByRefSerializer("Type", getType, setType)
	serDataType := NewPropByRefSerializer("DataType", getArrayDataType, setArrayDataType)
	serDataSize := NewPropByRefSerializer("DataSize", getDataSize, setDataSize)

	wrap.AddSerializer(ser, RWENUM)
	wrap.AddSerializer(serb1, RWBOOL)
	wrap.AddSerializer(serb2, RWBOOL)
	wrap.AddSerializer(serType, RWINT)
	wrap.AddSerializer(serDataType, RWINT)
	wrap.AddSerializer(serDataSize, RWINT)
	GetObjectWrapperManager().AddWrap(wrap)

	registerArrayWrapper("FloatArray", model.IDFLOATARRAY, model.GLFLOAT, 1)
	registerArrayWrapper("Vec2Array", model.IDVEC2ARRAY, model.GLFLOAT, 2)
	registerArrayWrapper("Vec3Array", model.IDVEC3ARRAY, model.GLFLOAT, 3)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 117)
		wrap.MarkSerializerAsAdded("osg::Vec3Array")
		uv.SetLastVersion()
	}
	registerArrayWrapper("Vec4Array", model.IDVEC4ARRAY, model.GLFLOAT, 4)
	registerArrayWrapper("DoubleArray", model.IDDOUBLEARRAY, model.GLDOUBLE, 1)
	registerArrayWrapper("Vec2dArray", model.IDVEC2DARRAY, model.GLDOUBLE, 2)
	registerArrayWrapper("Vec3dArray", model.IDVEC3DARRAY, model.GLDOUBLE, 3)
	registerArrayWrapper("Vec4dArray", model.IDVEC4DARRAY, model.GLDOUBLE, 4)
	registerArrayWrapper("ByteArray", model.IDBYTEARRAY, model.GLBYTE, 1)
	registerArrayWrapper("Vec2bArray", model.IDVEC2BARRAY, model.GLBYTE, 2)
	registerArrayWrapper("Vec3bArray", model.IDVEC3BARRAY, model.GLBYTE, 3)
	registerArrayWrapper("Vec4bArray", model.IDVEC4BARRAY, model.GLBYTE, 4)
	registerArrayWrapper("UByteArray", model.IDUBYTEARRAY, model.GLUNSIGNEDBYTE, 1)
	registerArrayWrapper("Vec2ubArray", model.IDVEC2UBARRAY, model.GLUNSIGNEDBYTE, 2)
	registerArrayWrapper("Vec3ubArray", model.IDVEC3UBARRAY, model.GLUNSIGNEDBYTE, 3)
	registerArrayWrapper("Vec4ubArray", model.IDVEC4UBARRAY, model.GLUNSIGNEDBYTE, 4)
	registerArrayWrapper("ShortArray", model.IDSHORTARRAY, model.GLSHORT, 1)
	registerArrayWrapper("Vec2sArray", model.IDVEC2SARRAY, model.GLSHORT, 2)
	registerArrayWrapper("Vec3sArray", model.IDVEC3SARRAY, model.GLSHORT, 3)
	registerArrayWrapper("Vec4sArray", model.IDVEC4SARRAY, model.GLSHORT, 4)
	registerArrayWrapper("UShortArray", model.IDUSHORTARRAY, model.GLUNSIGNEDSHORT, 1)
	registerArrayWrapper("Vec2usArray", model.IDVEC2USARRAY, model.GLUNSIGNEDSHORT, 2)
	registerArrayWrapper("Vec3usArray", model.IDVEC3USARRAY, model.GLUNSIGNEDSHORT, 3)
	registerArrayWrapper("Vec4usArray", model.IDVEC4USARRAY, model.GLUNSIGNEDSHORT, 4)
	registerArrayWrapper("IntArray", model.IDINTARRAY, model.GLINT, 1)
	registerArrayWrapper("Vec2iArray", model.IDVEC2IARRAY, model.GLINT, 2)
	registerArrayWrapper("Vec3iArray", model.IDVEC3IARRAY, model.GLINT, 3)
	registerArrayWrapper("Vec4iArray", model.IDVEC4IARRAY, model.GLINT, 4)
	registerArrayWrapper("UIntArray", model.IDUINTARRAY, model.GLUNSIGNEDINT, 1)
	registerArrayWrapper("Vec2uiArray", model.IDVEC2UIARRAY, model.GLUNSIGNEDINT, 2)
	registerArrayWrapper("Vec3uiArray", model.IDVEC3UIARRAY, model.GLUNSIGNEDINT, 3)
	registerArrayWrapper("Vec4uiArray", model.IDVEC4UIARRAY, model.GLUNSIGNEDINT, 4)
}
