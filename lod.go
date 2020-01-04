package osg

import (
	"github.com/flywave/go-osg/model"
)

func setCenterMode(obj interface{}, pro interface{}) {
	obj.(model.LodInterface).SetCmode(pro.(uint32))
}

func getCenterMode(obj interface{}) interface{} {
	return obj.(model.LodInterface).GetCmode()
}

func setRangeMode(obj interface{}, pro interface{}) {
	obj.(model.LodInterface).SetRmode(pro.(uint32))
}

func getRangeMode(obj interface{}) interface{} {
	return obj.(model.LodInterface).GetRmode()
}

func checkUserCenter(obj interface{}) bool {
	lod := obj.(model.LodInterface)
	md := *lod.GetCmode()
	return md == model.USERDEFINEDCENTER || md == model.UNIONOFBOUNDINGSPHEREANDUSERDEFINED
}
func readUserCenter(is *OsgIstream, obj interface{}) {
	lod := obj.(model.LodInterface)
	ct := [3]float64{}
	var r float64
	is.Read(&ct)
	is.Read(&r)
	lod.SetCenter([3]float32{float32(ct[0]), float32(ct[1]), float32(ct[2])})
	lod.SetRadius(float32(r))
}

func writeUserCenter(os *OsgOstream, obj interface{}) {
	lod := obj.(model.LodInterface)
	os.Write(*lod.GetCenter())
	os.Write(*lod.GetRadius())
	os.Write(os.CRLF)
}

func rangeListChecker(obj interface{}) bool {
	lod := obj.(model.LodInterface)
	return len(lod.GetRangeList()) > 0
}

func rangeListReader(is *OsgIstream, obj interface{}) {
	lod := obj.(model.LodInterface)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	lod.SetRangeList(make([][2]float32, size, size))
	for i := 0; i < size; i++ {
		var min, max float32
		is.Read(&min)
		is.Read(&max)
		lod.SetRange(i, min, max)
	}
	is.Read(is.ENDBRACKET)
}

func rangeListWriter(os *OsgOstream, obj interface{}) {
	lod := obj.(model.LodInterface)
	size := len(lod.GetRangeList())
	os.Write(size)
	os.Write(os.BEGINBRACKET)
	for i := 0; i < size; i++ {
		os.Write(lod.GetRangeList()[i][0])
		os.Write(lod.GetRangeList()[i][1])
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func init() {
	fn := func() interface{} {
		return nil
	}
	wrap := NewObjectWrapper("LOD", fn, "osg::Object osg::Node osg::Group osg::LOD")
	ser1 := NewEnumSerializer("CenterMode", getCenterMode, setCenterMode)
	ser1.Add("USEBOUNDINGSPHERECENTER", model.USEBOUNDINGSPHERECENTER)
	ser1.Add("USERDEFINEDCENTER", model.USERDEFINEDCENTER)
	ser1.Add("UNIONOFBOUNDINGSPHEREANDUSERDEFINED", model.UNIONOFBOUNDINGSPHEREANDUSERDEFINED)
	wrap.AddSerializer(ser1, RWENUM)

	ser2 := NewUserSerializer("UserCenter", checkUserCenter, readUserCenter, writeUserCenter)
	wrap.AddSerializer(ser2, RWENUM)

	ser3 := NewEnumSerializer("RangeMode", getRangeMode, setRangeMode)
	ser3.Add("DISTANCEFROMEYEPOINT", model.DISTANCEFROMEYEPOINT)
	ser3.Add("PIXELSIZEONSCREEN", model.PIXELSIZEONSCREEN)
	wrap.AddSerializer(ser3, RWENUM)

	seruser := NewUserSerializer("RangeList", rangeListChecker, rangeListReader, rangeListWriter)
	wrap.AddSerializer(seruser, RWUSER)
	GetObjectWrapperManager().AddWrap(wrap)

}
