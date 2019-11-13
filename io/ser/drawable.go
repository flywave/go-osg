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

func getStateSet(obj interface{}) interface{} {
	return obj.(*model.Drawable).States
}

func setStateSet(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).States = pro.(*model.StateSet)
}

func getCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).Callback
}

func setCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).Callback = pro.(*model.ComputeBoundingBoxCallback)
}

func getShape(obj interface{}) interface{} {
	return obj.(*model.Drawable).Shape
}

func setShape(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).Shape = pro.(*model.Shape)
}

func getSupportsDisplayList(obj interface{}) interface{} {
	return &obj.(*model.Drawable).SupportsDisplayList
}

func setSupportsDisplayList(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).SupportsDisplayList = pro.(bool)
}

func getUseDisplayList(obj interface{}) interface{} {
	return &obj.(*model.Drawable).UseDisplayList
}

func setUseDisplayList(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).UseDisplayList = pro.(bool)
}

func getUseVertexBufferObjects(obj interface{}) interface{} {
	return &obj.(*model.Drawable).UseVertexBufferObjects
}

func setUseVertexBufferObjects(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).UseVertexBufferObjects = pro.(bool)
}

func getUpdateCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).UpdateCallback
}

func setUpdateCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).UpdateCallback = pro.(*model.Callback)
}

func getEventCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).EventCallback
}

func setEventCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).EventCallback = pro.(*model.Callback)
}

func getCullCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).CullCallback
}

func setCullCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).CullCallback = pro.(*model.Callback)
}

func getDrawCallback(obj interface{}) interface{} {
	return obj.(*model.Drawable).DwCallback
}

func setDrawCallback(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).DwCallback = pro.(*model.DrawCallback)
}

func getNodeMask(obj interface{}) interface{} {
	return &obj.(*model.Drawable).NodeMask
}

func setNodeMask(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).NodeMask = pro.(uint32)
}

func getCullingActive(obj interface{}) interface{} {
	return &obj.(*model.Drawable).CullingActive
}

func setCullingActive(obj interface{}, pro interface{}) {
	obj.(*model.Drawable).CullingActive = pro.(bool)
}

func init() {
	wrap := io.NewObjectWrapper("Drawable", nil, "osg::Object osg::Node osg::Drawable")
	io.AddUpdateWrapperVersionProxy(&wrap, 154)
	wrap.MarkSerializerAsAdded("osg::Node")

	ser := io.NewObjectSerializer("StateSet", getStateSet, setStateSet)
	ser_user := io.NewUserSerializer("InitialBound", InitialBoundChecker, InitialBoundReader, InitialBoundWriter)
	ser_cpm := io.NewObjectSerializer("ComputeBoundingBoxCallback", getCallback, setCallback)
	ser_shape := io.NewObjectSerializer("Shape", getShape, setShape)
	ser_bool1 := io.NewEnumSerializer("SupportsDisplayList", getSupportsDisplayList, setSupportsDisplayList)
	ser_bool2 := io.NewEnumSerializer("UseDisplayList", getUseDisplayList, setUseDisplayList)
	ser_bool3 := io.NewEnumSerializer("UseVertexBufferObjects", getUseVertexBufferObjects, setUseVertexBufferObjects)
	ser_uc := io.NewObjectSerializer("UpdateCallback", getUpdateCallback, setUpdateCallback)
	ser_ec := io.NewObjectSerializer("EventCallback", getEventCallback, setEventCallback)
	ser_cc := io.NewObjectSerializer("CullCallback", getCullCallback, setCullCallback)
	ser_dc := io.NewObjectSerializer("DrawCallback", getDrawCallback, setDrawCallback)

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
	serb1 := io.NewPropByValSerializer("NodeMask", false, getNodeMask, setNodeMask)
	wrap.AddSerializer(&serb1, io.RW_UINT)

	io.AddUpdateWrapperVersionProxy(&wrap, 145)
	serb2 := io.NewPropByValSerializer("CullingActive", false, getCullingActive, setCullingActive)
	wrap.AddSerializer(&serb2, io.RW_BOOL)

	io.GetObjectWrapperManager().AddWrap(&wrap)
}
