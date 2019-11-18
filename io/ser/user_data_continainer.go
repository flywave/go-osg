package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func checkUDC_UserData(obj interface{}) bool {
	udc := obj.(*model.UserDataContainer)
	return udc.UserData != nil
}

func readUDC_UserData(is *io.OsgIstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	is.Read(is.BEGIN_BRACKET)
	ob := is.ReadObject(nil)
	udc.UserData = ob
	is.Read(is.END_BRACKET)
}

func writeUDC_UserData(os *io.OsgOstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	os.Write(udc.UserData)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

// _descriptions
func checkUDC_Descriptions(obj interface{}) bool {
	udc := obj.(*model.UserDataContainer)
	return len(udc.DescriptionList) > 0
}

func readUDC_Descriptions(is *io.OsgIstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	size := is.ReadSize()
	is.Read(is.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		udc.DescriptionList = append(udc.DescriptionList, str)
	}
	is.Read(is.END_BRACKET)
}

func writeUDC_Descriptions(os *io.OsgOstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	os.Write(len(udc.DescriptionList))
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	for s := range udc.DescriptionList {
		os.Write(s)
		os.Write(os.CRLF)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkUDC_UserObjects(obj interface{}) bool {
	udc := obj.(*model.UserDataContainer)
	return len(udc.ObjectList) > 0
}
func readUDC_UserObjects(is *io.OsgIstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	size := is.ReadSize()
	is.Read(is.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		ob := is.ReadObject(nil)
		udc.ObjectList = append(udc.ObjectList, ob)
	}
	is.Read(is.END_BRACKET)
}

func writeUDC_UserObjects(os *io.OsgOstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	os.Write(len(udc.ObjectList))
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	for o := range udc.ObjectList {
		os.Write(o)
		os.Write(os.CRLF)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func init() {
	fn := func() interface{} {
		udc := model.NewUserDataContainer()
		return &udc
	}
	wrap := io.NewObjectWrapper("UserDataContainer", fn, "osg::UserDataContainer")
	ser1 := io.NewStringSerializer("Name", getObjeName, setObjName)
	ser2 := io.NewEnumSerializer("DataVariance", getDataVariance, setDataVariance)
	ser2.Add("STATIC", model.STATIC)
	ser2.Add("DYNAMIC", model.DYNAMIC)
	ser2.Add("UNSPECIFIED", model.UNSPECIFIED)

	ser3 := io.NewUserSerializer("UDC_UserData", checkUDC_UserData, readUDC_UserData, writeUDC_UserData)
	ser4 := io.NewUserSerializer("UDC_Descriptions", checkUDC_Descriptions, readUDC_Descriptions, writeUDC_Descriptions)
	ser5 := io.NewUserSerializer("UDC_UserObjects", checkUDC_UserObjects, readUDC_UserObjects, writeUDC_UserObjects)

	wrap.AddSerializer(&ser1, io.RW_STRING)
	wrap.AddSerializer(&ser2, io.RW_ENUM)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
