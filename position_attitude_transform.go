package osg

import "github.com/flywave/go-osg/model"

func getPositionAttitudeTransformPosition(obj interface{}) interface{} {
	return &obj.(*model.PositionAttitudeTransform).Position
}

func setPositionAttitudeTransformPosition(obj interface{}, val interface{}) {
	obj.(*model.PositionAttitudeTransform).Position = val.([3]float64)
}

func getPositionAttitudeTransformAttitude(obj interface{}) interface{} {
	return &obj.(*model.PositionAttitudeTransform).Attitude
}

func setPositionAttitudeTransformAttitude(obj interface{}, val interface{}) {
	obj.(*model.PositionAttitudeTransform).Attitude = val.([4]float64)
}

func getPositionAttitudeTransformScale(obj interface{}) interface{} {
	return &obj.(*model.PositionAttitudeTransform).Scale
}

func setPositionAttitudeTransformScale(obj interface{}, val interface{}) {
	obj.(*model.PositionAttitudeTransform).Scale = val.([3]float64)
}

func init() {
	fn := func() interface{} {
		return model.NewPositionAttitudeTransform()
	}
	wrap := NewObjectWrapper("PositionAttitudeTransform", fn, "osg::Object osg::Node osg::Group osg::Transform osg::PositionAttitudeTransform")
	ser1 := NewPropByValSerializer("Position", false, getPositionAttitudeTransformPosition, setPositionAttitudeTransformPosition)
	ser2 := NewPropByValSerializer("Attitude", false, getPositionAttitudeTransformAttitude, setPositionAttitudeTransformAttitude)
	ser3 := NewPropByValSerializer("Scale", false, getPositionAttitudeTransformScale, setPositionAttitudeTransformScale)
	wrap.AddSerializer(ser1, RWDOUBLE|0xF0000000)
	wrap.AddSerializer(ser2, RWDOUBLE|0xF0000000)
	wrap.AddSerializer(ser3, RWDOUBLE|0xF0000000)
	GetObjectWrapperManager().AddWrap(wrap)
}
