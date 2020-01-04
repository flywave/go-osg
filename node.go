package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkInitialBound1(obj interface{}) bool {
	node := obj.(*model.Node)
	return node.InitialBound.Radius > 0
}

func readInitialBound1(is *OsgIstream, obj interface{}) {
	node := obj.(*model.Node)
	is.Read(is.BEGINBRACKET)
	is.PROPERTY.Name = "Center"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Center)
	is.PROPERTY.Name = "Radius"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Radius)
	is.Read(is.ENDBRACKET)
}

func writeInitialBound1(os *OsgOstream, obj interface{}) {
	node := obj.(*model.Node)
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

func checkDescriptions(obj interface{}) bool {
	node := obj.(*model.Node)
	return len(node.Dscriptions) > 0
}

func readDescriptions(is *OsgIstream, obj interface{}) {
	node := obj.(*model.Node)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		node.Dscriptions = append(node.Dscriptions, str)
	}
	is.Read(is.ENDBRACKET)
}

func writeDescriptions(os *OsgOstream, obj interface{}) {
	node := obj.(*model.Node)
	l := len(node.Dscriptions)
	os.Write(l)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for i := 0; i < l; i++ {
		os.Write(&node.Dscriptions[i])
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}
func getSphereCallback(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetSphereCallback()
}

func setSphereCallback(obj interface{}, pro interface{}) {
	obj.(model.NodeInterface).SetSphereCallback(pro.(*model.ComputeBoundingSphereCallback))
}

func getUpdateCallback(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetUpdateCallback()
}

func setUpdateCallback(obj interface{}, pro interface{}) {
	obj.(model.NodeInterface).SetUpdateCallback(pro.(*model.Callback))
}

func getEventCallback(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetEventCallback()
}

func setEventCallback(obj interface{}, pro interface{}) {
	obj.(model.NodeInterface).SetEventCallback(pro.(*model.Callback))
}

func getCullCallback(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetCullCallback()
}

func setCullCallback(obj interface{}, pro interface{}) {
	obj.(model.NodeInterface).SetCullCallback(pro.(*model.Callback))
}

func getCullingActive(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetCullingActive()
}

func setCullingActive(obj interface{}, pro interface{}) {
	obj.(model.NodeInterface).SetCullingActive(pro.(bool))
}
func getNodeMask(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetNodeMask()
}

func setNodeMask(obj interface{}, pro interface{}) {
	obj.(model.NodeInterface).SetNodeMask(pro.(uint32))
}

func getStateSet(obj interface{}) interface{} {
	return obj.(model.NodeInterface).GetStates()
}

func setStateSet(obj interface{}, pro interface{}) {
	t := obj.(model.NodeInterface)
	t.SetStates(pro.(*model.StateSet))
}
func init() {
	fn := func() interface{} {
		nd := model.NewNode()
		return nd
	}
	wrap := NewObjectWrapper("Node", fn, "osg::Object osg::Node")
	seruser := NewUserSerializer("InitialBound", checkInitialBound1, readInitialBound1, writeInitialBound1)
	wrap.AddSerializer(seruser, RWUSER)

	sercpm := NewObjectSerializer("ComputeBoundingSphereCallback", getSphereCallback, setSphereCallback)
	seruc := NewObjectSerializer("UpdateCallback", getUpdateCallback, setUpdateCallback)
	serec := NewObjectSerializer("EventCallback", getEventCallback, setEventCallback)
	sercc := NewObjectSerializer("CullCallback", getCullCallback, setCullCallback)
	serb2 := NewPropByValSerializer("CullingActive", false, getCullingActive, setCullingActive)
	serb1 := NewPropByValSerializer("NodeMask", false, getNodeMask, setNodeMask)

	wrap.AddSerializer(sercpm, RWOBJECT)
	wrap.AddSerializer(seruc, RWOBJECT)
	wrap.AddSerializer(serec, RWOBJECT)
	wrap.AddSerializer(sercc, RWOBJECT)
	wrap.AddSerializer(serb2, RWBOOL)
	wrap.AddSerializer(serb1, RWBOOL)

	seruser2 := NewUserSerializer("Descriptions", checkDescriptions, readDescriptions, writeDescriptions)
	wrap.AddSerializer(seruser2, RWUSER)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 77)
		wrap.MarkSerializerAsRemoved("Descriptions")
		uv.SetLastVersion()
	}

	ser := NewObjectSerializer("StateSet", getStateSet, setStateSet)
	wrap.AddSerializer(ser, RWOBJECT)
	GetObjectWrapperManager().AddWrap(wrap)

}
