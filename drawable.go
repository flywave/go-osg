package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkInitialBound(obj interface{}) bool {
	node := obj.(*model.Drawable)
	return node.InitialBound.Radius > 0
}

func readInitialBound(is *OsgIstream, obj interface{}) {
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

func writeInitialBound(os *OsgOstream, obj interface{}) {
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

func getBoxCallback(obj interface{}) interface{} {
	return obj.(model.DrawableInterface).GetBoxCallback()
}

func setBoxCallback(obj interface{}, pro interface{}) {
	obj.(model.DrawableInterface).SetBoxCallback(pro.(*model.ComputeBoundingBoxCallback))
}

func getShape(obj interface{}) interface{} {
	return obj.(model.DrawableInterface).GetShape()
}

func setShape(obj interface{}, pro interface{}) {
	obj.(model.DrawableInterface).SetShape(pro.(*model.Shape))
}

func getSupportsDisplayList(obj interface{}) interface{} {
	return obj.(model.DrawableInterface).GetSupportsDisplayList()
}

func setSupportsDisplayList(obj interface{}, pro interface{}) {
	obj.(model.DrawableInterface).SetSupportsDisplayList(pro.(bool))
}

func getUseDisplayList(obj interface{}) interface{} {
	return obj.(model.DrawableInterface).GetUseDisplayList()
}

func setUseDisplayList(obj interface{}, pro interface{}) {
	obj.(model.DrawableInterface).SetUseDisplayList(pro.(bool))
}

func getUseVertexBufferObjects(obj interface{}) interface{} {
	return obj.(model.DrawableInterface).GetUseVertexBufferObjects()
}

func setUseVertexBufferObjects(obj interface{}, pro interface{}) {
	obj.(model.DrawableInterface).SetUseVertexBufferObjects(pro.(bool))
}

func getDrawCallback(obj interface{}) interface{} {
	return obj.(model.DrawableInterface).GetDwCallback()
}

func setDrawCallback(obj interface{}, pro interface{}) {
	obj.(model.DrawableInterface).SetDwCallback(pro.(*model.DrawCallback))
}

func init() {
	wrap := NewObjectWrapper("Drawable", nil, "osg::Object osg::Node osg::Drawable")
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 154)
		wrap.MarkSerializerAsAdded("osg::Node")
		uv.SetLastVersion()
	}

	ser := NewObjectSerializer("StateSet", getStateSet, setStateSet)
	seruser := NewUserSerializer("InitialBound", checkInitialBound, readInitialBound, writeInitialBound)
	sercpm := NewObjectSerializer("ComputeBoundingBoxCallback", getBoxCallback, setBoxCallback)
	sershape := NewObjectSerializer("Shape", getShape, setShape)
	serbool1 := NewEnumSerializer("SupportsDisplayList", getSupportsDisplayList, setSupportsDisplayList)
	serbool2 := NewEnumSerializer("UseDisplayList", getUseDisplayList, setUseDisplayList)
	serbool3 := NewEnumSerializer("UseVertexBufferObjects", getUseVertexBufferObjects, setUseVertexBufferObjects)
	seruc := NewObjectSerializer("UpdateCallback", getUpdateCallback, setUpdateCallback)
	serec := NewObjectSerializer("EventCallback", getEventCallback, setEventCallback)
	sercc := NewObjectSerializer("CullCallback", getCullCallback, setCullCallback)
	serdc := NewObjectSerializer("DrawCallback", getDrawCallback, setDrawCallback)

	wrap.AddSerializer(ser, RWOBJECT)
	wrap.AddSerializer(seruser, RWUSER)
	wrap.AddSerializer(sercpm, RWOBJECT)
	wrap.AddSerializer(sershape, RWOBJECT)
	wrap.AddSerializer(serbool1, RWBOOL)
	wrap.AddSerializer(serbool2, RWBOOL)
	wrap.AddSerializer(serbool3, RWBOOL)
	wrap.AddSerializer(seruc, RWOBJECT)
	wrap.AddSerializer(serec, RWOBJECT)
	wrap.AddSerializer(sercc, RWOBJECT)
	wrap.AddSerializer(serdc, RWOBJECT)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 156)
		wrap.MarkSerializerAsRemoved("StateSet")
		wrap.MarkSerializerAsRemoved("UpdateCallback")
		wrap.MarkSerializerAsRemoved("EventCallback")
		wrap.MarkSerializerAsRemoved("CullCallback")
		wrap.MarkSerializerAsRemoved("DrawCallback")
		uv.SetLastVersion()
	}
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 142)
		serb1 := NewPropByValSerializer("NodeMask", false, getNodeMask, setNodeMask)
		wrap.AddSerializer(serb1, RWUINT)
		uv.SetLastVersion()
	}
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 145)
		serb2 := NewPropByValSerializer("CullingActive", false, getCullingActive, setCullingActive)
		wrap.AddSerializer(serb2, RWBOOL)
		uv.SetLastVersion()
	}
	GetObjectWrapperManager().AddWrap(wrap)
}
