package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func getNumInstances(obj interface{}) interface{} {
	return &obj.(*model.PrimitiveSet).NumInstances
}

func setNumInstances(obj interface{}, val interface{}) {
	obj.(*model.PrimitiveSet).NumInstances = val.(int32)
}

func getPrimitMode(obj interface{}) interface{} {
	return &obj.(*model.PrimitiveSet).Mode
}

func setPrimitMode(obj interface{}, val interface{}) {
	obj.(*model.PrimitiveSet).Mode = val.(int32)
}

func getDAFirst(obj interface{}) interface{} {
	return &obj.(*model.DrawArrays).First
}

func setDAFirst(obj interface{}, val interface{}) {
	obj.(*model.DrawArrays).First = val.(int32)
}

func getDACount(obj interface{}) interface{} {
	return &obj.(*model.DrawArrays).Count
}

func setDACount(obj interface{}, val interface{}) {
	obj.(*model.DrawArrays).Count = val.(int32)
}

func getDLFirst(obj interface{}) interface{} {
	return &obj.(*model.DrawArrays).First
}

func setDLFirst(obj interface{}, val interface{}) {
	obj.(*model.DrawArrays).First = val.(int32)
}

func getDAData(obj interface{}) interface{} {
	return &obj.(*model.DrawArrayLengths).Data
}

func setDAData(obj interface{}, val interface{}) {
	obj.(*model.DrawArrayLengths).Data = val.([]int32)
}

func getDADataUByte(obj interface{}) interface{} {
	return &obj.(*model.DrawElementsUByte).Data
}

func setDADataUByte(obj interface{}, val interface{}) {
	obj.(*model.DrawElementsUByte).Data = val.([]uint8)
}

func getDADataUShort(obj interface{}) interface{} {
	return &obj.(*model.DrawElementsUShort).Data
}

func setDADataUShort(obj interface{}, val interface{}) {
	obj.(*model.DrawElementsUShort).Data = val.([]uint16)
}

func getDADataUInt(obj interface{}) interface{} {
	return &obj.(*model.DrawElementsUInt).Data
}

func setDADataUInt(obj interface{}, val interface{}) {
	obj.(*model.DrawElementsUInt).Data = val.([]uint32)
}

func init() {
	fn1 := func() interface{} {
		prv := model.NewPrimitiveSet()
		return &prv
	}

	fn2 := func() interface{} {
		dl := model.NewDrawArrays()
		return &dl
	}

	fn3 := func() interface{} {
		dl := model.NewDrawArrayLengths()
		return &dl
	}

	fn4 := func() interface{} {
		dl := model.NewDrawElementsUByte()
		return &dl
	}

	fn5 := func() interface{} {
		dl := model.NewDrawElementsUShort()
		return &dl
	}

	fn6 := func() interface{} {
		dl := model.NewDrawElementsUInt()
		return &dl
	}

	wrap1 := osg.NewObjectWrapper("PrimitiveSet", fn1, "osg::Object osg::BufferData osg::PrimitiveSet")
	osg.AddUpdateWrapperVersionProxy(&wrap1, 147)
	wrap1.MarkSerializerAsAdded("osg::BufferData")
	ser1 := osg.NewPropByValSerializer("NumInstances", false, getNumInstances, setNumInstances)
	ser2 := osg.NewEnumSerializer("Mode", getPrimitMode, setPrimitMode)
	ser2.Add("POINTS", model.GLPOINTS)
	ser2.Add("LINES", model.GLLINES)
	ser2.Add("LINESTRIP", model.GLLINESTRIP)
	ser2.Add("LINELOOP", model.GLLINELOOP)
	ser2.Add("TRIANGLES", model.GLTRIANGLES)
	ser2.Add("TRIANGLESTRIP", model.GLTRIANGLESTRIP)
	ser2.Add("TRIANGLEFAN", model.GLTRIANGLEFAN)
	ser2.Add("QUADS", model.GLQUADS)
	ser2.Add("QUADSTRIP", model.GLQUADSTRIP)
	ser2.Add("POLYGON", model.GLPOLYGON)
	ser2.Add("LINESADJACENCY", model.GLLINESADJACENCY)
	ser2.Add("LINESTRIPADJACENCY", model.GLLINESTRIPADJACENCY)
	ser2.Add("TRIANGLESADJACENCY", model.GLTRIANGLESADJACENCY)
	ser2.Add("TRIANGLESTRIPADJACENCY", model.GLTRIANGLESTRIPADJACENCY)
	ser2.Add("PATCHES", model.GLPATCHES)
	wrap1.AddSerializer(&ser1, osg.RWINT)
	wrap1.AddSerializer(&ser2, osg.RWENUM)
	osg.GetObjectWrapperManager().AddWrap(&wrap1)

	wrap2 := osg.NewObjectWrapper("DrawArrays", fn2, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawArrays")
	osg.AddUpdateWrapperVersionProxy(&wrap2, 147)
	wrap2.MarkSerializerAsAdded("osg::BufferData")
	ser3 := osg.NewPropByValSerializer("First", false, getDAFirst, setDAFirst)
	ser4 := osg.NewPropByValSerializer("Count", false, getDACount, setDACount)
	wrap2.AddSerializer(&ser3, osg.RWINT)
	wrap2.AddSerializer(&ser4, osg.RWUINT)
	osg.GetObjectWrapperManager().AddWrap(&wrap2)

	wrap3 := osg.NewObjectWrapper2("DrawArrayLengths", "flywave::osg::drawarraylengths", fn3, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawArrays")
	osg.AddUpdateWrapperVersionProxy(&wrap2, 147)
	wrap3.MarkSerializerAsAdded("osg::BufferData")
	ser5 := osg.NewPropByValSerializer("First", false, getDLFirst, setDLFirst)

	var tydata int32
	ser6 := osg.NewVectorSerializer("Data", osg.RWOBJECT, &tydata, getDAData, setDAData)
	wrap3.AddSerializer(&ser5, osg.RWINT)
	wrap3.AddSerializer(&ser6, osg.RWVECTOR)
	osg.GetObjectWrapperManager().AddWrap(&wrap3)

	wrap4 := osg.NewObjectWrapper2("DrawElementsUByte", "flywave::osg::drawelementsubyte", fn4, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawElementsUByte")
	osg.AddUpdateWrapperVersionProxy(&wrap4, 147)
	wrap4.MarkSerializerAsAdded("osg::BufferData")
	var tydata2 uint8
	ser7 := osg.NewVectorSerializer("Data", osg.RWOBJECT, &tydata2, getDADataUByte, setDADataUByte)
	wrap4.AddSerializer(&ser7, osg.RWUCHAR)
	osg.GetObjectWrapperManager().AddWrap(&wrap4)

	wrap5 := osg.NewObjectWrapper2("DrawElementsUShort", "flywave::osg::drawelementsushort", fn5, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawElementsUShort")
	osg.AddUpdateWrapperVersionProxy(&wrap5, 147)
	wrap5.MarkSerializerAsAdded("osg::BufferData")
	var tydata3 uint16
	ser8 := osg.NewVectorSerializer("Data", osg.RWOBJECT, &tydata3, getDADataUShort, setDADataUShort)
	wrap5.AddSerializer(&ser8, osg.RWUSHORT)
	osg.GetObjectWrapperManager().AddWrap(&wrap5)

	wrap6 := osg.NewObjectWrapper2("DrawElementsUInt", "flywave::osg::drawelementsuint", fn6, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawElementsUInt")
	osg.AddUpdateWrapperVersionProxy(&wrap6, 147)
	wrap6.MarkSerializerAsAdded("osg::BufferData")
	var tydata4 uint32
	ser9 := osg.NewVectorSerializer("Data", osg.RWOBJECT, &tydata4, getDADataUInt, setDADataUInt)
	wrap5.AddSerializer(&ser9, osg.RWUINT)
	osg.GetObjectWrapperManager().AddWrap(&wrap6)
}
