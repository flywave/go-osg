package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
	"github.com/ungerik/go3d/vec3"
)

func SetCenterMode(obj interface{}, pro interface{}) {
	obj.(*model.Lod).Cmode = pro.(uint32)
}

func getCenterMode(obj interface{}) interface{} {
	return &obj.(*model.Lod).Cmode
}

func SetUserCenter(obj interface{}, pro interface{}) {
	obj.(*model.Lod).Center = pro.(vec3.T)
}

func getUserCenter(obj interface{}) interface{} {
	return &obj.(*model.Lod).Center
}

func SetRangeMode(obj interface{}, pro interface{}) {
	obj.(*model.Lod).Rmode = pro.(uint32)
}

func getRangeMode(obj interface{}) interface{} {
	return &obj.(*model.Lod).Rmode
}

func checkUserCenter(obj interface{}) bool {
	lod := obj.(*model.Lod)
	return lod.Cmode == model.USER_DEFINED_CENTER || lod.Cmode == model.UNION_OF_BOUNDING_SPHERE_AND_USER_DEFINED
}
func readUserCenter(is *io.OsgIstream, obj interface{}) {
	lod := obj.(*model.Lod)
	is.Read(&lod.Center)
	is.Read(&lod.Radius)
}

func writeUserCenter(os *io.OsgOstream, obj interface{}) {
	lod := obj.(*model.Lod)
	os.Write(&lod.Center)
	os.Write(&lod.Radius)
	os.Write(os.CRLF)
}

func rangeListChecker(obj interface{}) bool {
	lod := obj.(*model.Lod)
	return len(lod.RangeList) > 0
}

func rangeListReader(is *io.OsgIstream, obj interface{}) {
	lod := obj.(*model.Lod)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	lod.RangeList = make([][2]float32, size, size)
	for i := 0; i < size; i++ {
		var min, max float32
		is.Read(&min)
		is.Read(&max)
		lod.SetRange(i, min, max)
	}
	is.Read(is.END_BRACKET)
}

func rangeListWriter(os *io.OsgOstream, obj interface{}) {
	lod := obj.(*model.Lod)
	size := len(lod.RangeList)
	os.Write(size)
	os.Write(os.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		os.Write(lod.RangeList[i][0])
		os.Write(lod.RangeList[i][1])
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func init() {
	fn := func() interface{} {
		return nil
	}
	wrap := io.NewObjectWrapper("LOD", fn, "osg::Object osg::Node osg::Group osg::LOD")
	ser1 := io.NewEnumSerializer("CenterMode", getCenterMode, SetCenterMode)
	ser1.Add("USE_BOUNDING_SPHERE_CENTER", model.USE_BOUNDING_SPHERE_CENTER)
	ser1.Add("USER_DEFINED_CENTER", model.USER_DEFINED_CENTER)
	ser1.Add("UNION_OF_BOUNDING_SPHERE_AND_USER_DEFINED", model.UNION_OF_BOUNDING_SPHERE_AND_USER_DEFINED)
	wrap.AddSerializer(&ser1, io.RW_ENUM)

	ser2 := io.NewUserSerializer("UserCenter", checkUserCenter, readUserCenter, writeUserCenter)
	wrap.AddSerializer(&ser2, io.RW_ENUM)

	ser3 := io.NewEnumSerializer("RangeMode", getRangeMode, SetRangeMode)
	ser3.Add("DISTANCE_FROM_EYE_POINT", model.DISTANCE_FROM_EYE_POINT)
	ser3.Add("PIXEL_SIZE_ON_SCREEN", model.PIXEL_SIZE_ON_SCREEN)
	wrap.AddSerializer(&ser3, io.RW_ENUM)

	seruser := io.NewUserSerializer("RangeList", rangeListChecker, rangeListReader, rangeListWriter)
	wrap.AddSerializer(&seruser, io.RW_USER)
	io.GetObjectWrapperManager().AddWrap(&wrap)

}
