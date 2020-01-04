package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkAmbient(pro interface{}) bool {
	return true
}
func readAmbient(is *OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 [4]float32
	is.Read(&val1)
	var val2 [4]float32
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetAmbient(model.GLFRONTANDBACK, val1)
	} else {
		mt.SetAmbient(model.GLFRONT, val1)
		mt.SetAmbient(model.GLBACK, val2)
	}
}

func writeAmbient(os *OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.AmbientFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetAmbient(model.GLFRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetAmbient(model.GLBACK))
	os.Write(os.CRLF)
}

func checkDiffuse(pro interface{}) bool {
	return true
}
func readDiffuse(is *OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 [4]float32
	is.Read(&val1)
	var val2 [4]float32
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetDiffuse(model.GLFRONTANDBACK, val1)
	} else {
		mt.SetDiffuse(model.GLFRONT, val1)
		mt.SetDiffuse(model.GLBACK, val2)
	}
}

func writeDiffuse(os *OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.DiffuseFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetDiffuse(model.GLFRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetDiffuse(model.GLBACK))
	os.Write(os.CRLF)
}

func checkSpecular(pro interface{}) bool {
	return true
}
func readSpecular(is *OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 [4]float32
	is.Read(&val1)
	var val2 [4]float32
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetSpecular(model.GLFRONTANDBACK, val1)
	} else {
		mt.SetSpecular(model.GLFRONT, val1)
		mt.SetSpecular(model.GLBACK, val2)
	}
}

func writeSpecular(os *OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.SpecularFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetSpecular(model.GLFRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetSpecular(model.GLBACK))
	os.Write(os.CRLF)
}

func checkEmission(pro interface{}) bool {
	return true
}
func readEmission(is *OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 [4]float32
	is.Read(&val1)
	var val2 [4]float32
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetEmission(model.GLFRONTANDBACK, val1)
	} else {
		mt.SetEmission(model.GLFRONT, val1)
		mt.SetEmission(model.GLBACK, val2)
	}
}

func writeEmission(os *OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.EmissionFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetEmission(model.GLFRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetEmission(model.GLBACK))
	os.Write(os.CRLF)
}

func checkShininess(pro interface{}) bool {
	return true
}
func readShininess(is *OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 float32
	is.Read(&val1)
	var val2 float32
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetShininess(model.GLFRONTANDBACK, val1)
	} else {
		mt.SetShininess(model.GLFRONT, val1)
		mt.SetShininess(model.GLBACK, val2)
	}
}

func writeShininess(os *OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.ShininessFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetShininess(model.GLFRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetShininess(model.GLBACK))
	os.Write(os.CRLF)
}

func getColorMode(obj interface{}) interface{} {
	mt := obj.(*model.Material)
	return &mt.Cmod
}

func setColorMode(obj interface{}, pro interface{}) {
	mt := obj.(*model.Material)
	mt.Cmod = pro.(int)
}

func init() {
	fn := func() interface{} {
		mt := model.NewMaterial()
		return mt
	}
	wrap := NewObjectWrapper("Material", fn, "osg::Object osg::StateAttribute osg::Material")
	ser1 := NewEnumSerializer("ColorMode", getColorMode, setColorMode)
	ser1.Add("AMBIENT", model.AMBIENT)
	ser1.Add("DIFFUSE", model.DIFFUSE)
	ser1.Add("SPECULAR", model.SPECULAR)
	ser1.Add("EMISSION", model.EMISSION)
	ser1.Add("AMBIENTANDDIFFUSE", model.AMBIENTANDDIFFUSE)
	ser1.Add("MTLOFF", model.MTLOFF)
	wrap.AddSerializer(ser1, RWENUM)

	ser2 := NewUserSerializer("Ambient", checkAmbient, readAmbient, writeAmbient)
	ser3 := NewUserSerializer("Diffuse", checkDiffuse, readDiffuse, writeDiffuse)
	ser4 := NewUserSerializer("Specular", checkSpecular, readSpecular, writeSpecular)
	ser5 := NewUserSerializer("Emission", checkEmission, readEmission, writeEmission)
	ser6 := NewUserSerializer("Shininess", checkShininess, readShininess, writeShininess)

	wrap.AddSerializer(ser2, RWUSER)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)
	wrap.AddSerializer(ser6, RWUSER)
	GetObjectWrapperManager().AddWrap(wrap)
}
