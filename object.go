package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkUserData(obj interface{}) bool {
	ob := obj.(*model.Object)
	_, ok := ob.Udc.UserData.(*model.Object)
	return ob.Udc.UserData != nil && ok
}

func readUserData(is *OsgIstream, obj interface{}) {
	is.Read(is.BEGINBRACKET)
	ob := is.ReadObject(nil)
	if ob != nil {
		obj.(*model.Object).Udc.UserData = ob
	}
	is.Read(is.ENDBRACKET)
}

func writeUserData(os *OsgOstream, obj interface{}) {
	os.Write(os.BEGINBRACKET)
	os.Write(obj.(*model.Object).Udc.UserData)
	os.Write(os.ENDBRACKET)
}

func getObjeName(obj interface{}) interface{} {
	ob := obj.(model.ObjectInterface)
	return ob.GetName()
}

func setObjName(obj interface{}, val interface{}) {
	ob := obj.(model.ObjectInterface)
	ob.SetName(val.(string))
}

func getDataVariance(obj interface{}) interface{} {
	ob := obj.(model.ObjectInterface)
	return ob.GetDataVariance()
}

func setDataVariance(obj interface{}, val interface{}) {
	obj.(model.ObjectInterface).SetDataVariance(val.(int32))
}

func getUserDataContainer(obj interface{}) interface{} {
	return obj.(model.ObjectInterface).GetUserDataContainer()
}

func setUserDataContainer(obj interface{}, val interface{}) {
	obj.(model.ObjectInterface).SetUserDataContainer(val.(*model.UserDataContainer))
}

func init() {
	fn := func() interface{} {
		obj := model.NewObject()
		return obj
	}
	wrap := NewObjectWrapper("Object", fn, "osg::Object")
	ser1 := NewStringSerializer("Name", getObjeName, setObjName)
	ser2 := NewEnumSerializer("DataVariance", getDataVariance, setDataVariance)
	ser2.Add("STATIC", model.STATIC)
	ser2.Add("DYNAMIC", model.DYNAMIC)
	ser2.Add("UNSPECIFIED", model.UNSPECIFIED)
	ser3 := NewUserSerializer("UserData", checkUserData, readUserData, writeUserData)
	wrap.AddSerializer(ser1, RWSTRING)
	wrap.AddSerializer(ser2, RWENUM)
	wrap.AddSerializer(ser3, RWUSER)

	{
		uv := AddUpdateWrapperVersionProxy(wrap, 77)
		wrap.MarkSerializerAsRemoved("UserData")
		ser4 := NewObjectSerializer("UserDataContainer", getUserDataContainer, setUserDataContainer)
		wrap.AddSerializer(ser4, RWOBJECT)
		uv.SetLastVersion()
	}
	GetObjectWrapperManager().AddWrap(wrap)
}
