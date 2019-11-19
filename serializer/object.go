package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkUserData(obj interface{}) bool {
	ob := obj.(*model.Object)
	return ob.Udc.UserData != nil && model.IsBaseOfObject(ob.Udc.UserData)
}

func readUserData(is *osg.OsgIstream, obj interface{}) {
	is.Read(is.BEGINBRACKET)
	ob := is.ReadObject(nil)
	if ob != nil {
		obj.(*model.Object).Udc.UserData = ob
	}
	is.Read(is.ENDBRACKET)
}

func writeUserData(os *osg.OsgOstream, obj interface{}) {
	os.Write(os.BEGINBRACKET)
	os.Write(obj.(*model.Object).Udc.UserData)
	os.Write(os.ENDBRACKET)
}

func getObjeName(obj interface{}) interface{} {
	switch ob := obj.(type) {
	case *model.Object:
	case *model.UserDataContainer:
		return &ob.Name
	}
	return nil
}

func setObjName(obj interface{}, val interface{}) {
	switch ob := obj.(type) {
	case *model.Object:
	case *model.UserDataContainer:
		ob.Name = val.(string)
	}
}

func getDataVariance(obj interface{}) interface{} {
	switch ob := obj.(type) {
	case *model.Object:
	case *model.UserDataContainer:
		return &ob.DataVariance
	}
	return nil
}

func setDataVariance(obj interface{}, val interface{}) {
	switch ob := obj.(type) {
	case *model.Object:
	case *model.UserDataContainer:
		ob.DataVariance = val.(int)
	}
}

func getUserDataContainer(obj interface{}) interface{} {
	ob := obj.(*model.Object)
	return ob.Udc
}

func setUserDataContainer(obj interface{}, val interface{}) {
	ob := obj.(*model.Object)
	ob.Udc = val.(*model.UserDataContainer)
}

func init() {
	fn := func() interface{} {
		obj := model.NewObject()
		return &obj
	}
	wrap := osg.NewObjectWrapper("Object", fn, "osg::Object")
	ser1 := osg.NewStringSerializer("Name", getObjeName, setObjName)
	ser2 := osg.NewEnumSerializer("DataVariance", getDataVariance, setDataVariance)
	ser2.Add("STATIC", model.STATIC)
	ser2.Add("DYNAMIC", model.DYNAMIC)
	ser2.Add("UNSPECIFIED", model.UNSPECIFIED)
	ser3 := osg.NewUserSerializer("UserData", checkUserData, readUserData, writeUserData)
	wrap.AddSerializer(&ser1, osg.RWSTRING)
	wrap.AddSerializer(&ser2, osg.RWENUM)
	wrap.AddSerializer(&ser3, osg.RWUSER)

	osg.AddUpdateWrapperVersionProxy(&wrap, 77)
	wrap.MarkSerializerAsRemoved("UserData")
	ser4 := osg.NewObjectSerializer("UserDataContainer", getUserDataContainer, setUserDataContainer)
	wrap.AddSerializer(&ser4, osg.RWOBJECT)

	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
