package osg

import (
	"github.com/flywave/go-osg/model"
)

func setCenterMode(obj interface{}, pro interface{}) {
	obj.(*model.Lod).Cmode = pro.(uint32)
}

func getCenterMode(obj interface{}) interface{} {
	return &obj.(*model.Lod).Cmode
}

func setRangeMode(obj interface{}, pro interface{}) {
	obj.(*model.Lod).Rmode = pro.(uint32)
}

func getRangeMode(obj interface{}) interface{} {
	return &obj.(*model.Lod).Rmode
}

func checkUserCenter(obj interface{}) bool {
	lod := obj.(*model.Lod)
	return lod.Cmode == model.USERDEFINEDCENTER || lod.Cmode == model.UNIONOFBOUNDINGSPHEREANDUSERDEFINED
}
func readUserCenter(is *OsgIstream, obj interface{}) {
	lod := obj.(*model.Lod)
	is.Read(&lod.Center)
	is.Read(&lod.Radius)
}

func writeUserCenter(os *OsgOstream, obj interface{}) {
	lod := obj.(*model.Lod)
	os.Write(&lod.Center)
	os.Write(&lod.Radius)
	os.Write(os.CRLF)
}

func rangeListChecker(obj interface{}) bool {
	lod := obj.(*model.Lod)
	return len(lod.RangeList) > 0
}

func rangeListReader(is *OsgIstream, obj interface{}) {
	lod := obj.(*model.Lod)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	lod.RangeList = make([][2]float32, size, size)
	for i := 0; i < size; i++ {
		var min, max float32
		is.Read(&min)
		is.Read(&max)
		lod.SetRange(i, min, max)
	}
	is.Read(is.ENDBRACKET)
}

func rangeListWriter(os *OsgOstream, obj interface{}) {
	lod := obj.(*model.Lod)
	size := len(lod.RangeList)
	os.Write(size)
	os.Write(os.BEGINBRACKET)
	for i := 0; i < size; i++ {
		os.Write(lod.RangeList[i][0])
		os.Write(lod.RangeList[i][1])
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
	wrap.AddSerializer(&ser1, RWENUM)

	ser2 := NewUserSerializer("UserCenter", checkUserCenter, readUserCenter, writeUserCenter)
	wrap.AddSerializer(&ser2, RWENUM)

	ser3 := NewEnumSerializer("RangeMode", getRangeMode, setRangeMode)
	ser3.Add("DISTANCEFROMEYEPOINT", model.DISTANCEFROMEYEPOINT)
	ser3.Add("PIXELSIZEONSCREEN", model.PIXELSIZEONSCREEN)
	wrap.AddSerializer(&ser3, RWENUM)

	seruser := NewUserSerializer("RangeList", rangeListChecker, rangeListReader, rangeListWriter)
	wrap.AddSerializer(&seruser, RWUSER)
	GetObjectWrapperManager().AddWrap(&wrap)

}
