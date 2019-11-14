package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func readValue(is *io.OsgIstream, val interface{}) {}

func readModes(is *io.OsgIstream, val interface{}) {}

func readAttributes(is *io.OsgIstream, val interface{}) {}

func writeValue(os *io.OsgOstream, val interface{}) {}

func writeModes(os *io.OsgOstream, val interface{}) {}

func writeAttributes(os *io.OsgOstream, val interface{}) {}

func checkModeList(obj interface{}) bool {
	return false
}

func readModeList(os *io.OsgOstream, val interface{}) {}

func writeModeList(os *io.OsgOstream, val interface{}) {}

func checkAttributeList(obj interface{}) bool {
	return false
}

func readAttributeList(os *io.OsgOstream, val interface{}) {}

func writeAttributeList(os *io.OsgOstream, val interface{}) {}

func checkTextureModeList(obj interface{}) bool {
	return false
}

func readTextureModeList(is *io.OsgIstream, val interface{}) {}

func writeTextureModeList(os *io.OsgOstream, val interface{}) {}

func checkTextureAttributeList(obj interface{}) bool {
	return false
}

func readTextureAttributeList(is *io.OsgIstream, val interface{}) {}

func writeTextureAttributeList(os *io.OsgOstream, val interface{}) {}

func checkUniformList(obj interface{}) bool {
	return false
}

func readUniformList(is *io.OsgIstream, val interface{}) {}

func writeUniformList(os *io.OsgOstream, val interface{}) {}

func checkDefineList(obj interface{}) bool {
	return false
}

func readDefineList(is *io.OsgIstream, val interface{}) {}

func writeDefineList(obj interface{}) bool {
	return false
}

func init() {
	fn := func() interface{} {
		ss := model.NewStateSet()
		return &ss
	}
	wrap := io.NewObjectWrapper2("StateSet", "flywave::osg::stateset", fn, "osg::Object osg::StateSet")
	ser1 := io.NewUserSerializer("ModeList", checkModeList, readModeList, writeModeList)

}
