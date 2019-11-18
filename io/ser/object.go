package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func checkUserData(obj interface{}) bool {
	ob := obj.(*model.Object)
	return ob.Udc.UserData != nil && model.TypeBaseOfObject(ob.Udc.UserData)
}

func readUserData(is *io.OsgIstream, obj interface{}) {
	is.Read(is.BEGIN_BRACKET)
	ob := is.ReadObject(nil)
	if ob != nil {
		obj.(*model.Object).Udc.UserData = ob
	}
	is.Read(is.END_BRACKET)
}

func writeUserData(os *io.OsgOstream, obj interface{}) {
	os.Write(os.BEGIN_BRACKET)
	os.Write(obj.(*model.Object).Udc.UserData)
	os.Write(os.END_BRACKET)
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
	wrap := io.NewObjectWrapper("Object", fn, "osg::Object")
	ser1 := io.NewStringSerializer("Name", getObjeName, setObjName)
	ser2 := io.NewEnumSerializer("DataVariance", getDataVariance, setDataVariance)
	ser2.Add("STATIC", model.STATIC)
	ser2.Add("DYNAMIC", model.DYNAMIC)
	ser2.Add("UNSPECIFIED", model.UNSPECIFIED)
	ser3 := io.NewUserSerializer("UserData", checkUserData, readUserData, writeUserData)
	wrap.AddSerializer(&ser1, io.RW_STRING)
	wrap.AddSerializer(&ser2, io.RW_ENUM)
	wrap.AddSerializer(&ser3, io.RW_USER)

	io.AddUpdateWrapperVersionProxy(&wrap, 77)
	wrap.MarkSerializerAsRemoved("UserData")
	ser4 := io.NewObjectSerializer("UserDataContainer", getUserDataContainer, setUserDataContainer)
	wrap.AddSerializer(&ser4, io.RW_OBJECT)

	io.GetObjectWrapperManager().AddWrap(&wrap)
}
