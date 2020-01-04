package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkUDCUserData(obj interface{}) bool {
	udc := obj.(*model.UserDataContainer)
	return udc.UserData != nil
}

func readUDCUserData(is *OsgIstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	is.Read(is.BEGINBRACKET)
	ob := is.ReadObject(nil)
	udc.UserData = ob
	is.Read(is.ENDBRACKET)
}

func writeUDCUserData(os *OsgOstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.Write(udc.UserData)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

// descriptions
func checkUDCDescriptions(obj interface{}) bool {
	udc := obj.(*model.UserDataContainer)
	return len(udc.DescriptionList) > 0
}

func readUDCDescriptions(is *OsgIstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		udc.DescriptionList = append(udc.DescriptionList, str)
	}
	is.Read(is.ENDBRACKET)
}

func writeUDCDescriptions(os *OsgOstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	os.Write(len(udc.DescriptionList))
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for s := range udc.DescriptionList {
		os.Write(s)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkUDCUserObjects(obj interface{}) bool {
	udc := obj.(*model.UserDataContainer)
	return len(udc.ObjectList) > 0
}
func readUDCUserObjects(is *OsgIstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	for i := 0; i < size; i++ {
		ob := is.ReadObject(nil)
		udc.ObjectList = append(udc.ObjectList, ob)
	}
	is.Read(is.ENDBRACKET)
}

func writeUDCUserObjects(os *OsgOstream, obj interface{}) {
	udc := obj.(*model.UserDataContainer)
	os.Write(len(udc.ObjectList))
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for o := range udc.ObjectList {
		os.Write(o)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func init() {
	fn := func() interface{} {
		udc := model.NewUserDataContainer()
		return udc
	}
	wrap := NewObjectWrapper("UserDataContainer", fn, "osg::UserDataContainer")
	ser1 := NewStringSerializer("Name", getObjeName, setObjName)
	ser2 := NewEnumSerializer("DataVariance", getDataVariance, setDataVariance)
	ser2.Add("STATIC", model.STATIC)
	ser2.Add("DYNAMIC", model.DYNAMIC)
	ser2.Add("UNSPECIFIED", model.UNSPECIFIED)

	ser3 := NewUserSerializer("UDCUserData", checkUDCUserData, readUDCUserData, writeUDCUserData)
	ser4 := NewUserSerializer("UDCDescriptions", checkUDCDescriptions, readUDCDescriptions, writeUDCDescriptions)
	ser5 := NewUserSerializer("UDCUserObjects", checkUDCUserObjects, readUDCUserObjects, writeUDCUserObjects)

	wrap.AddSerializer(ser1, RWSTRING)
	wrap.AddSerializer(ser2, RWENUM)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)
	GetObjectWrapperManager().AddWrap(wrap)
}
