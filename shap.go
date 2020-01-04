package osg

func init() {
	wrap := NewObjectWrapper("Shape", nil, "osg::Object osg::Shape")
	GetObjectWrapperManager().AddWrap(wrap)
}
