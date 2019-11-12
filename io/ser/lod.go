package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func SetCenterMode(obj interface{}, pro interface{}) {
	obj.(*model.Lod).Cmode = pro.(uint32)
}

func GetCenterMode(obj interface{}) interface{} {
	return &obj.(*model.Lod).Cmode
}

func init() {
	fn := func() interface{} {
		return nil
	}
	wrap := io.NewObjectWrapper("LOD", fn, "osg::Object osg::Node osg::Group osg::LOD")
	serbool1 := io.NewEnumSerializer("CenterMode", GetCenterMode, SetCenterMode)
	wrap.AddSerializer(&serbool1, io.RW_OBJECT)
}
