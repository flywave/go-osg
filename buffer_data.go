package osg

func init() {
	wrap := NewObjectWrapper("BufferData", nil, "osg::Object osg::StateAttribute osg::CullFace")
	GetObjectWrapperManager().AddWrap(wrap)
}
