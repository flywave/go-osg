package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkInitialBound1(obj interface{}) bool {
	node := obj.(*model.Node)
	return node.InitialBound.Radius > 0
}

func readInitialBound1(is *osg.OsgIstream, obj interface{}) {
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

func writeInitialBound1(os *osg.OsgOstream, obj interface{}) {
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

func readDescriptions(is *osg.OsgIstream, obj interface{}) {
	node := obj.(*model.Node)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		node.Dscriptions = append(node.Dscriptions, str)
	}
	is.Read(is.ENDBRACKET)
}

func writeDescriptions(os *osg.OsgOstream, obj interface{}) {
	node := obj.(*model.Node)
	l := len(node.Dscriptions)
	os.Write(l)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for i := 0; i < l; i++ {
		os.WriteWrappedString(node.Dscriptions[i])
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}
func getCallback1(obj interface{}) interface{} {
	return obj.(*model.Node).Callback
}

func setCallback1(obj interface{}, pro interface{}) {
	obj.(*model.Node).Callback = pro.(*model.ComputeBoundingSphereCallback)
}

func getUpdateCallback1(obj interface{}) interface{} {
	return obj.(*model.Node).UpdateCallback
}

func setUpdateCallback1(obj interface{}, pro interface{}) {
	obj.(*model.Node).UpdateCallback = pro.(*model.Callback)
}

func getEventCallback1(obj interface{}) interface{} {
	return obj.(*model.Node).EventCallback
}

func setEventCallback1(obj interface{}, pro interface{}) {
	obj.(*model.Node).EventCallback = pro.(*model.Callback)
}

func getCullCallback1(obj interface{}) interface{} {
	return obj.(*model.Node).CullCallback
}

func setCullCallback1(obj interface{}, pro interface{}) {
	obj.(*model.Node).CullCallback = pro.(*model.Callback)
}

func getCullingActive1(obj interface{}) interface{} {
	return &obj.(*model.Node).CullingActive
}

func setCullingActive1(obj interface{}, pro interface{}) {
	obj.(*model.Node).CullingActive = pro.(bool)
}
func getNodeMask1(obj interface{}) interface{} {
	return &obj.(*model.Node).NodeMask
}

func setNodeMask1(obj interface{}, pro interface{}) {
	obj.(*model.Node).NodeMask = pro.(uint32)
}

func getStateSet1(obj interface{}) interface{} {
	return obj.(*model.Node).States
}

func setStateSet1(obj interface{}, pro interface{}) {
	obj.(*model.Node).States = pro.(*model.StateSet)
}
func init() {
	fn := func() interface{} {
		nd := model.NewNode()
		return &nd
	}
	wrap := osg.NewObjectWrapper("Node", fn, "osg::Object osg::Node")
	seruser := osg.NewUserSerializer("InitialBound", checkInitialBound1, readInitialBound1, writeInitialBound1)
	wrap.AddSerializer(&seruser, osg.RWUSER)

	sercpm := osg.NewObjectSerializer("ComputeBoundingSphereCallback", getCallback1, setCallback1)
	seruc := osg.NewObjectSerializer("UpdateCallback", getUpdateCallback1, setUpdateCallback1)
	serec := osg.NewObjectSerializer("EventCallback", getEventCallback1, setEventCallback1)
	sercc := osg.NewObjectSerializer("CullCallback", getCullCallback1, setCullCallback1)
	serb2 := osg.NewPropByValSerializer("CullingActive", false, getCullingActive1, setCullingActive1)
	serb1 := osg.NewPropByValSerializer("NodeMask", false, getNodeMask1, setNodeMask1)

	wrap.AddSerializer(&sercpm, osg.RWOBJECT)
	wrap.AddSerializer(&seruc, osg.RWOBJECT)
	wrap.AddSerializer(&serec, osg.RWOBJECT)
	wrap.AddSerializer(&sercc, osg.RWOBJECT)
	wrap.AddSerializer(&serb2, osg.RWBOOL)
	wrap.AddSerializer(&serb1, osg.RWBOOL)

	seruser2 := osg.NewUserSerializer("Descriptions", checkDescriptions, readDescriptions, writeDescriptions)
	wrap.AddSerializer(&seruser2, osg.RWUSER)
	osg.AddUpdateWrapperVersionProxy(&wrap, 77)
	wrap.MarkSerializerAsRemoved("Descriptions")

	ser := osg.NewObjectSerializer("StateSet", getStateSet1, setStateSet1)
	wrap.AddSerializer(&ser, osg.RWOBJECT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)

}
