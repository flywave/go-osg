package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkDatabasePath(obj interface{}) bool {
	return true
}

func readDatabasePath(is *OsgIstream, obj interface{}) {
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

func writeDatabasePath(os *OsgOstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	b := len(lod.DataBasePath) > 0
	os.Write(b)
	if b {
		os.Write(&lod.DataBasePath)
	}
	os.Write(os.CRLF)
}

func setFrameNumberOfLastTraversal(obj interface{}, pro interface{}) {
	obj.(*model.PagedLod).FrameNumberOfLastTraversal = pro.(uint32)
}

func getFrameNumberOfLastTraversal(obj interface{}) interface{} {
	return &obj.(*model.PagedLod).FrameNumberOfLastTraversal
}

func setNumChildrenThatCannotBeExpired(obj interface{}, pro interface{}) {
	obj.(*model.PagedLod).NumChildrenThatCannotBeExpired = pro.(uint32)
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

func readRangeDataList(is *OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	size := is.ReadSize()
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
	size = is.ReadSize()
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

func writeRangeDataList(os *OsgOstream, obj interface{}) {
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

func readChildren(is *OsgIstream, obj interface{}) {
	lod := obj.(*model.PagedLod)
	size := is.ReadSize()
	if size > 0 {
		is.Read(is.BEGINBRACKET)
		for i := 0; i < size; i++ {
			ob := is.ReadObject(nil)
			nd, ok := ob.(model.NodeInterface)
			if ok {
				lod.AddChild(nd)
			}
		}
	}
	is.Read(is.ENDBRACKET)
}

func writeChildren(os *OsgOstream, obj interface{}) {
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
		return pl
	}

	wrap := NewObjectWrapper("PagedLOD", fn, "osg::Object osg::Node osg::LOD osg::PagedLOD")
	ser1 := NewUserSerializer("DatabasePath", checkDatabasePath, readDatabasePath, writeDatabasePath)
	ser2 := NewPropByValSerializer("FrameNumberOfLastTraversal", false, getFrameNumberOfLastTraversal, setFrameNumberOfLastTraversal)
	ser3 := NewPropByValSerializer("NumChildrenThatCannotBeExpired", false, getNumChildrenThatCannotBeExpired, setNumChildrenThatCannotBeExpired)
	ser4 := NewPropByValSerializer("DisableExternalChildrenPaging", false, getDisableExternalChildrenPaging, setDisableExternalChildrenPaging)
	ser5 := NewUserSerializer("RangeDataList", checkRangeDataList, readRangeDataList, writeRangeDataList)
	ser6 := NewUserSerializer("Children", checkChildren, readChildren, writeChildren)

	wrap.AddSerializer(ser1, RWUSER)
	wrap.AddSerializer(ser2, RWUINT)
	wrap.AddSerializer(ser3, RWUINT)
	wrap.AddSerializer(ser4, RWBOOL)
	wrap.AddSerializer(ser5, RWUSER)
	wrap.AddSerializer(ser6, RWUSER)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 70)
		wrap.MarkSerializerAsRemoved("FrameNumberOfLastTraversal")
		uv.SetLastVersion()
	}
	GetObjectWrapperManager().AddWrap(wrap)
}
