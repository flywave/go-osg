package ser

import (
	"github.com/flywave/go-osg/io"

	"github.com/flywave/go-osg/model"
)

func checkInitialBound(obj interface{}) bool {
	node := obj.(*model.Drawable)
	return node.InitialBound.Radius > 0
}

func readInitialBound(is *io.OsgIstream, obj interface{}) {
	node := obj.(*model.Drawable)
	is.Read(is.BEGIN_BRACKET)
	is.PROPERTY.Name = "Center"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Center)
	is.PROPERTY.Name = "Radius"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Radius)
	is.Read(is.END_BRACKET)
}

func writeInitialBound(os *io.OsgOstream, obj interface{}) {
	node := obj.(*model.Drawable)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Center"
	os.Write(os.PROPERTY)
	os.Write(node.InitialBound.Center)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Radius"
	os.Write(os.PROPERTY)
	os.Write(node.InitialBound.Radius)
	os.Write(os.CRLF)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

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
	seruser := io.NewUserSerializer("InitialBound", checkInitialBound, readInitialBound, writeInitialBound)
	sercpm := io.NewObjectSerializer("ComputeBoundingBoxCallback", getCallback, setCallback)
	sershape := io.NewObjectSerializer("Shape", getShape, setShape)
	serbool1 := io.NewEnumSerializer("SupportsDisplayList", getSupportsDisplayList, setSupportsDisplayList)
	serbool2 := io.NewEnumSerializer("UseDisplayList", getUseDisplayList, setUseDisplayList)
	serbool3 := io.NewEnumSerializer("UseVertexBufferObjects", getUseVertexBufferObjects, setUseVertexBufferObjects)
	seruc := io.NewObjectSerializer("UpdateCallback", getUpdateCallback, setUpdateCallback)
	serec := io.NewObjectSerializer("EventCallback", getEventCallback, setEventCallback)
	sercc := io.NewObjectSerializer("CullCallback", getCullCallback, setCullCallback)
	serdc := io.NewObjectSerializer("DrawCallback", getDrawCallback, setDrawCallback)

	wrap.AddSerializer(&ser, io.RW_OBJECT)
	wrap.AddSerializer(&seruser, io.RW_USER)
	wrap.AddSerializer(&sercpm, io.RW_OBJECT)
	wrap.AddSerializer(&sershape, io.RW_OBJECT)
	wrap.AddSerializer(&serbool1, io.RW_BOOL)
	wrap.AddSerializer(&serbool2, io.RW_BOOL)
	wrap.AddSerializer(&serbool3, io.RW_BOOL)
	wrap.AddSerializer(&seruc, io.RW_OBJECT)
	wrap.AddSerializer(&serec, io.RW_OBJECT)
	wrap.AddSerializer(&sercc, io.RW_OBJECT)
	wrap.AddSerializer(&serdc, io.RW_OBJECT)

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
