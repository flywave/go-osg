package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkDatabasePath(obj interface{}) bool {
	return true
}

func readDatabasePath(is *osg.OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	hasp := false
	is.Read(&hasp)
	if !hasp {
		if len(is.Options.DbPath) > 0 {
			lod.DataBasePath = is.Options.DbPath
		}
	} else {
		is.ReadWrappedString(&lod.DataBasePath)
	}
}

func writeDatabasePath(os *osg.OsgOstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	b := len(lod.DataBasePath) > 0
	os.Write(b)
	if b {
		os.WriteWrappedString(lod.DataBasePath)
	}
	os.Write(os.CRLF)
}

func setFrameNumberOfLastTraversal(obj interface{}, pro interface{}) {
	obj.(*model.PagedLod).FrameNumberOfLastTraversal = pro.(uint)
}

func getFrameNumberOfLastTraversal(obj interface{}) interface{} {
	return &obj.(*model.PagedLod).FrameNumberOfLastTraversal
}

func setNumChildrenThatCannotBeExpired(obj interface{}, pro interface{}) {
	obj.(*model.PagedLod).NumChildrenThatCannotBeExpired = pro.(uint)
}

func getNumChildrenThatCannotBeExpired(obj interface{}) interface{} {
	return &obj.(*model.PagedLod).NumChildrenThatCannotBeExpired
}

func setDisableExternalChildrenPaging(obj interface{}, pro interface{}) {
	obj.(*model.PagedLod).DisableExternalChildrenPaging = pro.(bool)
}

func getDisableExternalChildrenPaging(obj interface{}) interface{} {
	return &obj.(*model.PagedLod).DisableExternalChildrenPaging
}

func checkRangeDataList(obj interface{}) bool {
	lod := obj.(*model.PagedLod)
	return len(lod.PerRangeDataList) > 0
}

func readRangeDataList(is *osg.OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	size := 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	lod.PerRangeDataList = make([]model.PerRangeData, size, size)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		lod.SetFileName(i, str)
	}
	is.Read(is.ENDBRACKET)
	is.PROPERTY.Name = "PriorityList"
	is.Read(is.PROPERTY)
	is.Read(&size)
	is.Read(is.BEGINBRACKET)

	for i := 0; i < size; i++ {
		var off, scale float32
		is.Read(&off)
		is.Read(&scale)
		lod.SetPriorityOffset(i, off)
		lod.SetPriorityOffset(i, scale)
	}
	is.Read(is.ENDBRACKET)
}

func writeRangeDataList(os *osg.OsgOstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	l := len(lod.PerRangeDataList)
	os.Write(l)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for _, per := range lod.PerRangeDataList {
		os.Write(per.FileName)
	}
	os.Write(os.CRLF)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)

	os.PROPERTY.Name = "PriorityList"
	os.Write(os.PROPERTY)
	os.Write(l)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for _, per := range lod.PerRangeDataList {
		os.Write(per.PriorityOffset)
		os.Write(per.PriorityScale)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkChildren(obj interface{}) bool {
	lod := obj.(*model.PagedLod)
	return len(lod.Children) > 0
}

func readChildren(is *osg.OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	var size int
	is.Read(&size)
	if size > 0 {
		is.Read(is.BEGINBRACKET)
		for i := 0; i < size; i++ {
			ob := is.ReadObject(nil)
			if model.IsBaseOfNode(lod) {
				lod.AddChild(ob)
			}
		}
	}
	is.Read(is.ENDBRACKET)
}

func writeChildren(os *osg.OsgOstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	size := len(lod.PerRangeDataList)
	dynamicLoadedSize := 0
	for _, per := range lod.PerRangeDataList {
		if len(per.FileName) > 0 {
			dynamicLoadedSize++
		}
	}
	realSize := size - dynamicLoadedSize
	if realSize > 0 {
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		for i, per := range lod.PerRangeDataList {
			if len(per.FileName) > 0 {
				continue
			}
			if i < len(lod.Children) {
				os.Write(lod.Children[i])
			}
		}
		os.Write(os.ENDBRACKET)
	}
}

func init() {
	fn := func() interface{} {
		pl := model.NewPagedLod()
		return &pl
	}

	wrap := osg.NewObjectWrapper("PagedLOD", fn, "osg::Object osg::Node osg::LOD osg::PagedLOD")
	ser1 := osg.NewUserSerializer("DatabasePath", checkDatabasePath, readDatabasePath, writeDatabasePath)
	ser2 := osg.NewPropByValSerializer("FrameNumberOfLastTraversal", false, getFrameNumberOfLastTraversal, setFrameNumberOfLastTraversal)
	ser3 := osg.NewPropByValSerializer("NumChildrenThatCannotBeExpired", false, getNumChildrenThatCannotBeExpired, setNumChildrenThatCannotBeExpired)
	ser4 := osg.NewPropByValSerializer("DisableExternalChildrenPaging", false, getDisableExternalChildrenPaging, setDisableExternalChildrenPaging)
	ser5 := osg.NewUserSerializer("RangeDataList", checkRangeDataList, readRangeDataList, writeRangeDataList)
	ser6 := osg.NewUserSerializer("Children", checkChildren, readChildren, writeChildren)

	wrap.AddSerializer(&ser1, osg.RWUSER)
	wrap.AddSerializer(&ser2, osg.RWUINT)
	wrap.AddSerializer(&ser3, osg.RWUINT)
	wrap.AddSerializer(&ser4, osg.RWBOOL)
	wrap.AddSerializer(&ser5, osg.RWUSER)
	wrap.AddSerializer(&ser6, osg.RWUSER)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
