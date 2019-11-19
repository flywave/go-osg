package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkInitialBound(obj interface{}) bool {
	node := obj.(*model.Drawable)
	return node.InitialBound.Radius > 0
}

func readInitialBound(is *osg.OsgIstream, obj interface{}) {
	node := obj.(*model.Drawable)
	is.Read(is.BEGINBRACKET)
	is.PROPERTY.Name = "Center"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Center)
	is.PROPERTY.Name = "Radius"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Radius)
	is.Read(is.ENDBRACKET)
}

func writeInitialBound(os *osg.OsgOstream, obj interface{}) {
	node := obj.(*model.Drawable)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Center"
	os.Write(os.PROPERTY)
	os.Write(node.InitialBound.Center)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Radius"
	os.Write(os.PROPERTY)
	os.Write(node.InitialBound.Radius)
	os.Write(os.CRLF)
	os.Write(os.ENDBRACKET)
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
	wrap := osg.NewObjectWrapper("Drawable", nil, "osg::Object osg::Node osg::Drawable")
	osg.AddUpdateWrapperVersionProxy(&wrap, 154)
	wrap.MarkSerializerAsAdded("osg::Node")

	ser := osg.NewObjectSerializer("StateSet", getStateSet, setStateSet)
	seruser := osg.NewUserSerializer("InitialBound", checkInitialBound, readInitialBound, writeInitialBound)
	sercpm := osg.NewObjectSerializer("ComputeBoundingBoxCallback", getCallback, setCallback)
	sershape := osg.NewObjectSerializer("Shape", getShape, setShape)
	serbool1 := osg.NewEnumSerializer("SupportsDisplayList", getSupportsDisplayList, setSupportsDisplayList)
	serbool2 := osg.NewEnumSerializer("UseDisplayList", getUseDisplayList, setUseDisplayList)
	serbool3 := osg.NewEnumSerializer("UseVertexBufferObjects", getUseVertexBufferObjects, setUseVertexBufferObjects)
	seruc := osg.NewObjectSerializer("UpdateCallback", getUpdateCallback, setUpdateCallback)
	serec := osg.NewObjectSerializer("EventCallback", getEventCallback, setEventCallback)
	sercc := osg.NewObjectSerializer("CullCallback", getCullCallback, setCullCallback)
	serdc := osg.NewObjectSerializer("DrawCallback", getDrawCallback, setDrawCallback)

	wrap.AddSerializer(&ser, osg.RWOBJECT)
	wrap.AddSerializer(&seruser, osg.RWUSER)
	wrap.AddSerializer(&sercpm, osg.RWOBJECT)
	wrap.AddSerializer(&sershape, osg.RWOBJECT)
	wrap.AddSerializer(&serbool1, osg.RWBOOL)
	wrap.AddSerializer(&serbool2, osg.RWBOOL)
	wrap.AddSerializer(&serbool3, osg.RWBOOL)
	wrap.AddSerializer(&seruc, osg.RWOBJECT)
	wrap.AddSerializer(&serec, osg.RWOBJECT)
	wrap.AddSerializer(&sercc, osg.RWOBJECT)
	wrap.AddSerializer(&serdc, osg.RWOBJECT)

	osg.AddUpdateWrapperVersionProxy(&wrap, 156)
	wrap.MarkSerializerAsRemoved("UpdateCallback")
	wrap.MarkSerializerAsRemoved("EventCallback")
	wrap.MarkSerializerAsRemoved("CullCallback")
	wrap.MarkSerializerAsRemoved("DrawCallback")

	osg.AddUpdateWrapperVersionProxy(&wrap, 142)
	serb1 := osg.NewPropByValSerializer("NodeMask", false, getNodeMask, setNodeMask)
	wrap.AddSerializer(&serb1, osg.RWUINT)

	osg.AddUpdateWrapperVersionProxy(&wrap, 145)
	serb2 := osg.NewPropByValSerializer("CullingActive", false, getCullingActive, setCullingActive)
	wrap.AddSerializer(&serb2, osg.RWBOOL)

	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
