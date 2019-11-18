package ser

import (
	"github.com/flywave/go-osg/io"

	"github.com/flywave/go-osg/model"
)

func checkInitialBound1(obj interface{}) bool {
	node := obj.(*model.Node)
	return node.InitialBound.Radius > 0
}

func readInitialBound1(is *io.OsgIstream, obj interface{}) {
	node := obj.(*model.Node)
	is.Read(is.BEGIN_BRACKET)
	is.PROPERTY.Name = "Center"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Center)
	is.PROPERTY.Name = "Radius"
	is.Read(is.PROPERTY)
	is.Read(&node.InitialBound.Radius)
	is.Read(is.END_BRACKET)
}

func writeInitialBound1(os *io.OsgOstream, obj interface{}) {
	node := obj.(*model.Node)
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

func checkDescriptions(obj interface{}) bool {
	node := obj.(*model.Node)
	return len(node.Dscriptions) > 0
}

func readDescriptions(is *io.OsgIstream, obj interface{}) {
	node := obj.(*model.Node)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		var str string
		is.ReadWrappedString(&str)
		node.Dscriptions = append(node.Dscriptions, str)
	}
	is.Read(is.END_BRACKET)
}

func writeDescriptions(os *io.OsgOstream, obj interface{}) {
	node := obj.(*model.Node)
	l := len(node.Dscriptions)
	os.Write(l)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	for i := 0; i < l; i++ {
		os.WriteWrappedString(node.Dscriptions[i])
	}
	os.Write(os.END_BRACKET)
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
	wrap := io.NewObjectWrapper("Node", fn, "osg::Object osg::Node")
	seruser := io.NewUserSerializer("InitialBound", checkInitialBound1, readInitialBound1, writeInitialBound1)
	wrap.AddSerializer(&seruser, io.RW_USER)

	sercpm := io.NewObjectSerializer("ComputeBoundingSphereCallback", getCallback1, setCallback1)
	seruc := io.NewObjectSerializer("UpdateCallback", getUpdateCallback1, setUpdateCallback1)
	serec := io.NewObjectSerializer("EventCallback", getEventCallback1, setEventCallback1)
	sercc := io.NewObjectSerializer("CullCallback", getCullCallback1, setCullCallback1)
	serb2 := io.NewPropByValSerializer("CullingActive", false, getCullingActive1, setCullingActive1)
	serb1 := io.NewPropByValSerializer("NodeMask", false, getNodeMask1, setNodeMask1)

	wrap.AddSerializer(&sercpm, io.RW_OBJECT)
	wrap.AddSerializer(&seruc, io.RW_OBJECT)
	wrap.AddSerializer(&serec, io.RW_OBJECT)
	wrap.AddSerializer(&sercc, io.RW_OBJECT)
	wrap.AddSerializer(&serb2, io.RW_BOOL)
	wrap.AddSerializer(&serb1, io.RW_BOOL)

	seruser2 := io.NewUserSerializer("Descriptions", checkDescriptions, readDescriptions, writeDescriptions)
	wrap.AddSerializer(&seruser2, io.RW_USER)
	io.AddUpdateWrapperVersionProxy(&wrap, 77)
	wrap.MarkSerializerAsRemoved("Descriptions")

	ser := io.NewObjectSerializer("StateSet", getStateSet1, setStateSet1)
	wrap.AddSerializer(&ser, io.RW_OBJECT)
	io.GetObjectWrapperManager().AddWrap(&wrap)

}
