package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func checkDatabasePath(obj interface{}) bool {
	return true
}

func readDatabasePath(is *io.OsgIstream, obj interface{}) {
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

func writeDatabasePath(os *io.OsgOstream, obj interface{}) {
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

func readRangeDataList(is *io.OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	size := 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	lod.PerRangeDataList = make([]model.PerRangeData, size, size)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		lod.SetFileName(i, str)
	}
	is.Read(is.END_BRACKET)
	is.PROPERTY.Name = "PriorityList"
	is.Read(is.PROPERTY)
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)

	for i := 0; i < size; i++ {
		var off, scale float32
		is.Read(&off)
		is.Read(&scale)
		lod.SetPriorityOffset(i, off)
		lod.SetPriorityOffset(i, scale)
	}
	is.Read(is.END_BRACKET)
}

func writeRangeDataList(os *io.OsgOstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	l := len(lod.PerRangeDataList)
	os.Write(l)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	for _, per := range lod.PerRangeDataList {
		os.Write(per.FileName)
	}
	os.Write(os.CRLF)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)

	os.PROPERTY.Name = "PriorityList"
	os.Write(os.PROPERTY)
	os.Write(l)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	for _, per := range lod.PerRangeDataList {
		os.Write(per.PriorityOffset)
		os.Write(per.PriorityScale)
		os.Write(os.CRLF)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkChildren(obj interface{}) bool {
	lod := obj.(*model.PagedLod)
	return len(lod.Children) > 0
}

func readChildren(is *io.OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	var size int
	is.Read(&size)
	if size > 0 {
		is.Read(is.BEGIN_BRACKET)
		for i := 0; i < size; i++ {
			obj = is.ReadObject()
			if model.IsBaseOfNode(lod) {
				lod.AddChild(obj)
			}
		}
	}
	is.Read(is.END_BRACKET)
}

func writeChildren(os *io.OsgOstream, obj interface{}) {
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
		os.Write(os.BEGIN_BRACKET)
		os.Write(os.CRLF)
		for i, per := range lod.PerRangeDataList {
			if len(per.FileName) > 0 {
				continue
			}
			if i < len(lod.Children) {
				os.Write(lod.Children[i])
			}
		}
		os.Write(os.END_BRACKET)
	}
}

func init() {
	fn := func() interface{} {
		pl := model.NewPagedLod()
		return &pl
	}

	wrap := io.NewObjectWrapper("PagedLOD", fn, "osg::Object osg::Node osg::LOD osg::PagedLOD")
	ser1 := io.NewUserSerializer("DatabasePath", checkDatabasePath, readDatabasePath, writeDatabasePath)
	ser2 := io.NewPropByValSerializer("FrameNumberOfLastTraversal", false, getFrameNumberOfLastTraversal, setFrameNumberOfLastTraversal)
	ser3 := io.NewPropByValSerializer("NumChildrenThatCannotBeExpired", false, getNumChildrenThatCannotBeExpired, setNumChildrenThatCannotBeExpired)
	ser4 := io.NewPropByValSerializer("DisableExternalChildrenPaging", false, getDisableExternalChildrenPaging, setDisableExternalChildrenPaging)
	ser5 := io.NewUserSerializer("RangeDataList", checkRangeDataList, readRangeDataList, writeRangeDataList)
	ser6 := io.NewUserSerializer("Children", checkChildren, readChildren, writeChildren)

	wrap.AddSerializer(&ser1, io.RW_USER)
	wrap.AddSerializer(&ser2, io.RW_UINT)
	wrap.AddSerializer(&ser3, io.RW_UINT)
	wrap.AddSerializer(&ser4, io.RW_BOOL)
	wrap.AddSerializer(&ser5, io.RW_USER)
	wrap.AddSerializer(&ser6, io.RW_USER)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
