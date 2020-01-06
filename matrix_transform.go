package osg

import "github.com/flywave/go-osg/model"

func getMatrixf(obj interface{}) interface{} {
	return &obj.(*model.MatrixTransform).Matrix
}

func setMatrixf(obj interface{}, mat interface{}) {
	obj.(*model.MatrixTransform).Matrix = *mat.(*[4][4]float32)
}

func init() {
	fn := func() interface{} {
		mt := model.NewMatrixTransform()
		return mt
	}
	wrap := NewObjectWrapper("MatrixTransform", fn, "osg::Object osg::Node osg::Group osg::Transform osg::MatrixTransform")
	ser := NewMatrixSerializer("Matrix", getMatrixf, setMatrixf)
	wrap.AddSerializer(ser, RWMATRIXF)
	GetObjectWrapperManager().AddWrap(wrap)
}
