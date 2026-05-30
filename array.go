package osg

import (
	"github.com/flywave/go-osg/model"
)

func getBinding(obj interface{}) interface{} {
	return &obj.(*model.Array).Binding
}

func setBinding(obj interface{}, fc interface{}) {
	obj.(*model.Array).Binding = *fc.(*int32)
}

func getNormalize(obj interface{}) interface{} {
	return &obj.(*model.Array).Normalize
}

func setNormalize(obj interface{}, fc interface{}) {
	obj.(*model.Array).Normalize = *fc.(*bool)
}

func getPreserveDataType(obj interface{}) interface{} {
	return &obj.(*model.Array).PreserveDataType
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

func getType(obj interface{}) interface{} {
	return &obj.(*model.Array).Type
}

func setType(obj interface{}, fc interface{}) {
	obj.(*model.Array).Type = *fc.(*model.ArrayTable)
}

func getArrayDataType(obj interface{}) interface{} {
	return &obj.(*model.Array).DataType
}

func setArrayDataType(obj interface{}, fc interface{}) {
	obj.(*model.Array).DataType = *fc.(*int32)
}

func getDataSize(obj interface{}) interface{} {
	return &obj.(*model.Array).DataSize
}

func setDataSize(obj interface{}, fc interface{}) {
	obj.(*model.Array).DataSize = *fc.(*int32)
}

func getElementSize(obj interface{}) interface{} {
	return &obj.(*model.Array).ElementSize
}

func setElementSize(obj interface{}, fc interface{}) {
	obj.(*model.Array).ElementSize = *fc.(*int32)
}

func getTotalDataSize(obj interface{}) interface{} {
	return &obj.(*model.Array).TotalDataSize
}

func setTotalDataSize(obj interface{}, fc interface{}) {
	obj.(*model.Array).TotalDataSize = *fc.(*int32)
}

func getNumElements(obj interface{}) interface{} {
	return &obj.(*model.Array).NumElements
}

func setNumElements(obj interface{}, fc interface{}) {
	obj.(*model.Array).NumElements = *fc.(*int32)
}

func registerArrayWrapper(name string, arrayType int, elemType SerType, dataSize int) {
	fn := func() interface{} {
		ty := model.ArrayTable(arrayType)
		var dt int32
		switch elemType {
		case RWFLOAT, RWVEC2F, RWVEC3F, RWVEC4F:
			dt = model.GLFLOAT
		case RWDOUBLE, RWVEC2D, RWVEC3D, RWVEC4D:
			dt = model.GLDOUBLE
		case RWCHAR, RWVEC2B, RWVEC3B, RWVEC4B:
			dt = model.GLBYTE
		case RWUCHAR, RWVEC2UB, RWVEC3UB, RWVEC4UB:
			dt = model.GLUNSIGNEDBYTE
		case RWSHORT, RWVEC2S, RWVEC3S, RWVEC4S:
			dt = model.GLSHORT
		case RWUSHORT, RWVEC2US, RWVEC3US, RWVEC4US:
			dt = model.GLUNSIGNEDSHORT
		case RWINT, RWVEC2I, RWVEC3I, RWVEC4I:
			dt = model.GLINT
		case RWUINT, RWVEC2UI, RWVEC3UI, RWVEC4UI:
			dt = model.GLUNSIGNEDINT
		default:
			dt = model.GLFLOAT
		}
		dsize := int32(dataSize)
		return model.NewArray(ty, dt, dsize)
	}
	wrap := NewObjectWrapper(name, fn, "osg::Object osg::BufferData osg::Array osg::"+name)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 147)
		wrap.MarkSerializerAsAdded("osg::BufferData")
		uv.SetLastVersion()
	}

	serData := NewIsAVectorSerializer("vector", elemType, 1, getArrayData, setArrayData)
	serData.GetDataSize = func(obj interface{}) int32 {
		if arr, ok := obj.(*model.Array); ok {
			return arr.DataSize
		}
		return 0
	}
	wrap.AddSerializer(serData, RWVECTOR)

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

	wrap.AddSerializer(ser, RWENUM)
	wrap.AddSerializer(serb1, RWBOOL)
	wrap.AddSerializer(serb2, RWBOOL)
	GetObjectWrapperManager().AddWrap(wrap)

	registerArrayWrapper("FloatArray", model.IDFLOATARRAY, RWFLOAT, 1)
	registerArrayWrapper("Vec2Array", model.IDVEC2ARRAY, RWVEC2F, 2)
	registerArrayWrapper("Vec3Array", model.IDVEC3ARRAY, RWVEC3F, 3)
	registerArrayWrapper("Vec4Array", model.IDVEC4ARRAY, RWVEC4F, 4)
	registerArrayWrapper("DoubleArray", model.IDDOUBLEARRAY, RWDOUBLE, 1)
	registerArrayWrapper("Vec2dArray", model.IDVEC2DARRAY, RWVEC2D, 2)
	registerArrayWrapper("Vec3dArray", model.IDVEC3DARRAY, RWVEC3D, 3)
	registerArrayWrapper("Vec4dArray", model.IDVEC4DARRAY, RWVEC4D, 4)
	registerArrayWrapper("ByteArray", model.IDBYTEARRAY, RWCHAR, 1)
	registerArrayWrapper("Vec2bArray", model.IDVEC2BARRAY, RWVEC2B, 2)
	registerArrayWrapper("Vec3bArray", model.IDVEC3BARRAY, RWVEC3B, 3)
	registerArrayWrapper("Vec4bArray", model.IDVEC4BARRAY, RWVEC4B, 4)
	registerArrayWrapper("UByteArray", model.IDUBYTEARRAY, RWUCHAR, 1)
	registerArrayWrapper("Vec2ubArray", model.IDVEC2UBARRAY, RWVEC2UB, 2)
	registerArrayWrapper("Vec3ubArray", model.IDVEC3UBARRAY, RWVEC3UB, 3)
	registerArrayWrapper("Vec4ubArray", model.IDVEC4UBARRAY, RWVEC4UB, 4)
	registerArrayWrapper("ShortArray", model.IDSHORTARRAY, RWSHORT, 1)
	registerArrayWrapper("Vec2sArray", model.IDVEC2SARRAY, RWVEC2S, 2)
	registerArrayWrapper("Vec3sArray", model.IDVEC3SARRAY, RWVEC3S, 3)
	registerArrayWrapper("Vec4sArray", model.IDVEC4SARRAY, RWVEC4S, 4)
	registerArrayWrapper("UShortArray", model.IDUSHORTARRAY, RWUSHORT, 1)
	registerArrayWrapper("Vec2usArray", model.IDVEC2USARRAY, RWVEC2US, 2)
	registerArrayWrapper("Vec3usArray", model.IDVEC3USARRAY, RWVEC3US, 3)
	registerArrayWrapper("Vec4usArray", model.IDVEC4USARRAY, RWVEC4US, 4)
	registerArrayWrapper("IntArray", model.IDINTARRAY, RWINT, 1)
	registerArrayWrapper("Vec2iArray", model.IDVEC2IARRAY, RWVEC2I, 2)
	registerArrayWrapper("Vec3iArray", model.IDVEC3IARRAY, RWVEC3I, 3)
	registerArrayWrapper("Vec4iArray", model.IDVEC4IARRAY, RWVEC4I, 4)
	registerArrayWrapper("UIntArray", model.IDUINTARRAY, RWUINT, 1)
	registerArrayWrapper("Vec2uiArray", model.IDVEC2UIARRAY, RWVEC2UI, 2)
	registerArrayWrapper("Vec3uiArray", model.IDVEC3UIARRAY, RWVEC3UI, 3)
	registerArrayWrapper("Vec4uiArray", model.IDVEC4UIARRAY, RWVEC4UI, 4)
}
