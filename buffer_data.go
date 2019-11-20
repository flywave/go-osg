package osg

func init() {
	wrap := NewObjectWrapper2("BufferData", " model.BufferData", nil, "osg::Object osg::StateAttribute osg::CullFace")
	GetObjectWrapperManager().AddWrap(&wrap)
}
