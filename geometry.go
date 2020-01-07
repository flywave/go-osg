package osg

import (
	"github.com/flywave/go-osg/model"
)

var lookup *IntLookup

func readAttributeBinding(is *OsgIstream) int32 {
	if is.IsBinary() {
		var val int32
		is.Read(&val)
		return val
	} else {
		str := is.ReadString()
		return lookup.GetValue(str)
	}
}

func writeAttributeBinding(os *OsgOstream, val int32) {
	if os.IsBinary() {
		os.Write(val)
	} else {
		os.Write(lookup.GetString(val))
	}
}

func readArray(is *OsgIstream) *model.Array {
	is.PROPERTY.Name = "Array"
	hasArray := false
	var ary *model.Array
	is.Read(is.PROPERTY)
	is.Read(&hasArray)
	if hasArray {
		ary = is.ReadArray()
	}

	is.PROPERTY.Name = "Indices"
	hasIndices := false
	is.Read(is.PROPERTY)
	is.Read(&hasIndices)
	if hasIndices {
		ary2 := is.ReadArray()
		if ary != nil {
			ary.Udc.UserData = ary2
		}
	}

	is.PROPERTY.Name = "Binding"
	is.Read(is.PROPERTY)
	bd := readAttributeBinding(is)
	if ary != nil {
		ary.Binding = bd
	}

	is.PROPERTY.Name = "Normalize"
	var normalizeValue int32 = -1
	is.Read(is.PROPERTY)
	is.Read(&normalizeValue)
	if ary != nil {
		ary.Normalize = normalizeValue != 0
	}
	return ary
}

func writeArray(os *OsgOstream, ary *model.Array) {
	os.PROPERTY.Name = "Array"
	os.Write(os.PROPERTY)
	os.Write(ary != nil)
	var ary2 *model.Array
	if ary != nil {
		os.Write(ary)
		ary2 = ary.Udc.UserData.(*model.Array)
	} else {
		os.Write(os.CRLF)
	}

	os.PROPERTY.Name = "Indices"
	os.Write(os.PROPERTY)
	os.Write(ary2 != nil)
	if ary2 != nil {
		os.Write(ary2)
	} else {
		os.Write(os.CRLF)
	}
	os.PROPERTY.Name = "Binding"
	os.Write(os.PROPERTY)
	writeAttributeBinding(os, ary.Binding)
}

func checkVertexData(geom interface{}) bool {
	return geom.(*model.Geometry).VertexArray != nil
}
func readVertexData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.VertexArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeVertexData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.VertexArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkNormalData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.NormalArray != nil
}
func readNormalData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.NormalArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeNormalData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.NormalArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkColorData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.ColorArray != nil
}
func readColorData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.ColorArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeColorData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.ColorArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkSecondaryColorData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.SecondaryColorArray != nil
}
func readSecondaryColorData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.SecondaryColorArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeSecondaryColorData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.SecondaryColorArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkFogCoordData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.FogCoordArray != nil
}
func readFogCoordData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.FogCoordArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeFogCoordData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.FogCoordArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkTexCoordData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.TexCoordArrayList != nil
}

func readTexCoordData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	is.PROPERTY.Name = "Data"
	geom.TexCoordArrayList = make([]*model.Array, size, size)
	for i := 0; i < size; i++ {
		is.Read(is.PROPERTY)
		is.Read(is.BEGINBRACKET)

		ay := readArray(is)
		geom.SetTexCoordArray(i, ay)
		is.Read(is.ENDBRACKET)
	}
	is.Read(is.ENDBRACKET)
}

func writeTexCoordData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	length := len(geom.TexCoordArrayList)
	os.Write(length)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ary := range geom.TexCoordArrayList {
		os.Write(os.PROPERTY)
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		os.Write(ary)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkVertexAttribData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return len(geom.VertexAttribList) > 0
}

func readVertexAttribData(is *OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	is.PROPERTY.Name = "Data"
	geom.VertexAttribList = make([]*model.Array, size, size)
	for i := 0; i < size; i++ {
		is.Read(is.PROPERTY)
		is.Read(is.BEGINBRACKET)

		ay := readArray(is)
		geom.SetVertexAttribArray(i, ay, model.BINDUNDEFINED)
		is.Read(is.ENDBRACKET)
	}
	is.Read(is.ENDBRACKET)
}

func writeVertexAttribData(os *OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	length := len(geom.VertexAttribList)
	os.Write(length)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ary := range geom.VertexAttribList {
		os.Write(os.PROPERTY)
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		os.Write(ary)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkFastPathHint(g interface{}) bool {
	return false
}

func readFastPathHint(is *OsgIstream, g interface{}) {
	value := false
	if !is.IsBinary() {
		is.Read(&value)
	}
}

func writeFastPathHint(os *OsgOstream, g interface{}) {
}

func getPrimitiveSetList(obj interface{}) interface{} {
	return obj.(*model.Geometry).Primitives
}

func setPrimitiveSetList(obj interface{}, pro interface{}) {
	g := obj.(*model.Geometry)
	g.AddPrimitiveSet(pro)
}

func getTexCoordArrayList(obj interface{}) interface{} {
	return obj.(*model.Geometry).TexCoordArrayList
}

func setTexCoordArrayList(obj interface{}, pro interface{}) {
	pl := obj.(*model.Geometry).TexCoordArrayList
	pl = append(pl, pro.(*model.Array))
}

func getVertexAttribArrayList(obj interface{}) interface{} {
	return obj.(*model.Geometry).VertexAttribList
}

func setVertexAttribArrayList(obj interface{}, pro interface{}) {
	pl := obj.(*model.Geometry).VertexAttribList
	pl = append(pl, pro.(*model.Array))
}

func getVertexData(obj interface{}) interface{} {
	return obj.(*model.Geometry).VertexArray
}

func setVertexData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).VertexArray = pro.(*model.Array)
}

func getNormalData(obj interface{}) interface{} {
	return obj.(*model.Geometry).NormalArray
}

func setNormalData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).NormalArray = pro.(*model.Array)

}

func getColorData(obj interface{}) interface{} {
	return obj.(*model.Geometry).ColorArray
}

func setColorData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).ColorArray = pro.(*model.Array)

}

func getSecondaryColorData(obj interface{}) interface{} {
	return obj.(*model.Geometry).SecondaryColorArray

}

func setSecondaryColorData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).SecondaryColorArray = pro.(*model.Array)

}

func getFogCoordData(obj interface{}) interface{} {
	return obj.(*model.Geometry).FogCoordArray

}
func setFogCoordData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).FogCoordArray = pro.(*model.Array)
}

func init() {
	lk := NewIntLookup()
	lookup = lk
	lookup.Add("BINDOFF", 0)
	lookup.Add("BINDOVERALL", 1)
	lookup.Add("BINDPERPRIMITIVESET", 2)
	lookup.Add("BINDPERPRIMITIVE", 3)
	lookup.Add("BINDPERVERTEX", 4)

	fn := func() interface{} {
		g := model.NewGeometry()
		return g
	}

	wrap := NewObjectWrapper("Geometry", fn, "osg::Object osg::Node osg::Drawable osg::Geometry")
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 154)
		wrap.MarkSerializerAsAdded("osg::Node")
		uv.SetLastVersion()
	}
	ps := model.NewPrimitiveSet()
	vser := NewVectorSerializer("PrimitiveSetList", RWOBJECT, ps, getPrimitiveSetList, setPrimitiveSetList)
	wrap.AddSerializer(vser, RWVECTOR)

	ser1 := NewUserSerializer("VertexData", checkVertexData, readVertexData, writeVertexData)
	ser2 := NewUserSerializer("NormalData", checkNormalData, readNormalData, writeNormalData)
	ser3 := NewUserSerializer("ColorData", checkColorData, readColorData, writeColorData)
	ser4 := NewUserSerializer("SecondaryColorData", checkSecondaryColorData, readSecondaryColorData, writeSecondaryColorData)
	ser5 := NewUserSerializer("FogCoordData", checkFogCoordData, readFogCoordData, writeFogCoordData)
	ser6 := NewUserSerializer("TexCoordData", checkTexCoordData, readTexCoordData, writeTexCoordData)
	ser7 := NewUserSerializer("VertexAttribData", checkVertexAttribData, readVertexAttribData, writeVertexAttribData)
	ser8 := NewUserSerializer("FastPathHint", checkFastPathHint, readFastPathHint, writeFastPathHint)

	wrap.AddSerializer(ser1, RWUSER)
	wrap.AddSerializer(ser2, RWUSER)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)
	wrap.AddSerializer(ser6, RWUSER)
	wrap.AddSerializer(ser7, RWUSER)
	wrap.AddSerializer(ser8, RWUSER)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 112)
		wrap.MarkSerializerAsRemoved("VertexData")
		wrap.MarkSerializerAsRemoved("NormalData")
		wrap.MarkSerializerAsRemoved("ColorData")
		wrap.MarkSerializerAsRemoved("SecondaryColorData")
		wrap.MarkSerializerAsRemoved("FogCoordData")
		wrap.MarkSerializerAsRemoved("TexCoordData")
		wrap.MarkSerializerAsRemoved("VertexAttribData")
		wrap.MarkSerializerAsRemoved("FastPathHint")

		ser11 := NewObjectSerializer("VertexData", getVertexData, setVertexData)
		ser21 := NewObjectSerializer("NormalData", getNormalData, setNormalData)
		ser31 := NewObjectSerializer("ColorData", getColorData, setColorData)
		ser41 := NewObjectSerializer("SecondaryColorData", getSecondaryColorData, setSecondaryColorData)
		ser51 := NewObjectSerializer("FogCoordData", getFogCoordData, setFogCoordData)

		wrap.AddSerializer(ser11, RWOBJECT)
		wrap.AddSerializer(ser21, RWOBJECT)
		wrap.AddSerializer(ser31, RWOBJECT)
		wrap.AddSerializer(ser41, RWOBJECT)
		wrap.AddSerializer(ser51, RWOBJECT)

		ay := model.NewArray2()
		vser2 := NewVectorSerializer("TexCoordArrayList", RWOBJECT, ay, getTexCoordArrayList, setTexCoordArrayList)
		vser3 := NewVectorSerializer("VertexAttribArrayList", RWOBJECT, ay, getVertexAttribArrayList, setVertexAttribArrayList)
		wrap.AddSerializer(vser2, RWVECTOR)
		wrap.AddSerializer(vser3, RWVECTOR)
		uv.SetLastVersion()
	}
	GetObjectWrapperManager().AddWrap(wrap)
}
