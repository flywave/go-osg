package ser

import (
	"github.com/flywave/go-osg/io"

	"github.com/flywave/go-osg/model"
)

func InitialBoundChecker(obj interface{}) bool {
	return true
}
func InitialBoundReader(is *io.OsgIstream, obj interface{}) {}

func InitialBoundWriter(os *io.OsgOstream, obj interface{}) {}

func GetStateSet(obj interface{}) interface{} {
	return obj.(*model.Drawable).States
}

func SetStateSet(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).States = pro.(*model.StateSet)
}

func GetCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).Callback
}

func SetCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).Callback = pro.(*model.ComputeBoundingBoxCallback)
}

func GetShape(obj interface{}) interface{} {
	return obj.(*model.Drawable).Shape
}

func SetShape(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).Shape = pro.(*model.Shape)
}

func GetSupportsDisplayList(obj interface{}) interface{} {
	return &obj.(*model.Drawable).SupportsDisplayList
}

func SetSupportsDisplayList(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).SupportsDisplayList = pro.(bool)
}

func GetUseDisplayList(obj interface{}) interface{} {
	return &obj.(*model.Drawable).UseDisplayList
}

func SetUseDisplayList(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).UseDisplayList = pro.(bool)
}

func GetUseVertexBufferObjects(obj interface{}) interface{} {
	return &obj.(*model.Drawable).UseVertexBufferObjects
}

func SetUseVertexBufferObjects(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).UseVertexBufferObjects = pro.(bool)
}

func GetUpdateCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).UpdateCallback
}

func SetUpdateCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).UpdateCallback = pro.(*model.Callback)
}

func GetEventCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).EventCallback
}

func SetEventCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).EventCallback = pro.(*model.Callback)
}

func GetCullCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).CullCallback
}

func SetCullCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).CullCallback = pro.(*model.Callback)
}

func GetDrawCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).DwCallback
}

func SetDrawCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).DwCallback = pro.(*model.DrawCallback)
}

func GetNodeMask(obj interface{}) interface{} {
	return &obj.(*model.Drawable).NodeMask
}

func SetNodeMask(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).NodeMask = pro.(uint32)
}

func GetCullingActive(obj interface{}) interface{} {
	return &obj.(*model.Drawable).CullingActive
}

func SetCullingActive(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).CullingActive = pro.(bool)
}

func init() {
	wrap := io.NewObjectWrapper("Drawable", nil, "osg::Object osg::Node osg::Drawable")
	io.AddUpdateWrapperVersionProxy(&wrap, 154)
	wrap.MarkSerializerAsAdded("osg::Node")

	ser := io.NewObjectSerializer("StateSet", GetStateSet, SetStateSet)
	ser_user := io.NewUserSerializer("InitialBound", InitialBoundChecker, InitialBoundReader, InitialBoundWriter)
	ser_cpm := io.NewObjectSerializer("ComputeBoundingBoxCallback", GetCallback, SetCallback)
	ser_shape := io.NewObjectSerializer("Shape", GetShape, SetShape)
	ser_bool1 := io.NewEnumSerializer("SupportsDisplayList", GetSupportsDisplayList, SetSupportsDisplayList)
	ser_bool2 := io.NewEnumSerializer("UseDisplayList", GetUseDisplayList, SetUseDisplayList)
	ser_bool3 := io.NewEnumSerializer("UseVertexBufferObjects", GetUseVertexBufferObjects, SetUseVertexBufferObjects)
	ser_uc := io.NewObjectSerializer("UpdateCallback", GetUpdateCallback, SetUpdateCallback)
	ser_ec := io.NewObjectSerializer("EventCallback", GetEventCallback, SetEventCallback)
	ser_cc := io.NewObjectSerializer("CullCallback", GetCullCallback, SetCullCallback)
	ser_dc := io.NewObjectSerializer("DrawCallback", GetDrawCallback, SetDrawCallback)

	wrap.AddSerializer(&ser, io.RW_OBJECT)
	wrap.AddSerializer(&ser_user, io.RW_USER)
	wrap.AddSerializer(&ser_cpm, io.RW_OBJECT)
	wrap.AddSerializer(&ser_shape, io.RW_OBJECT)
	wrap.AddSerializer(&ser_bool1, io.RW_BOOL)
	wrap.AddSerializer(&ser_bool2, io.RW_BOOL)
	wrap.AddSerializer(&ser_bool3, io.RW_BOOL)
	wrap.AddSerializer(&ser_uc, io.RW_OBJECT)
	wrap.AddSerializer(&ser_ec, io.RW_OBJECT)
	wrap.AddSerializer(&ser_cc, io.RW_OBJECT)
	wrap.AddSerializer(&ser_dc, io.RW_OBJECT)

	io.AddUpdateWrapperVersionProxy(&wrap, 156)
	wrap.MarkSerializerAsRemoved("UpdateCallback")
	wrap.MarkSerializerAsRemoved("EventCallback")
	wrap.MarkSerializerAsRemoved("CullCallback")
	wrap.MarkSerializerAsRemoved("DrawCallback")

	io.AddUpdateWrapperVersionProxy(&wrap, 142)
	serb1 := io.NewPropByValSerializer("NodeMask", false, GetNodeMask, SetNodeMask)
	wrap.AddSerializer(&serb1, io.RW_UINT)

	io.AddUpdateWrapperVersionProxy(&wrap, 145)
	serb2 := io.NewPropByValSerializer("CullingActive", false, GetCullingActive, SetCullingActive)
	wrap.AddSerializer(&serb2, io.RW_BOOL)

	io.GetObjectWrapperManager().AddWrap(&wrap)
}
