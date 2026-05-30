package osg

func init() {
	wrap := NewObjectWrapper("BufferData", nil, "osg::Object osg::BufferData")
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 147)
		ser := NewObjectSerializer("BufferObject", func(obj interface{}) interface{} { return nil }, func(obj interface{}, val interface{}) {})
		wrap.AddSerializer(ser, RWOBJECT)
		uv.SetLastVersion()
	}
	GetObjectWrapperManager().AddWrap(wrap)
}
