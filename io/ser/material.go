package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
	"github.com/ungerik/go3d/vec4"
)

func checkAmbient(pro interface{}) bool {
	return true
}
func readAmbient(is *io.OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 vec4.T
	is.Read(&val1)
	var val2 vec4.T
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetAmbient(model.GL_FRONT_AND_BACK, val1)
	} else {
		mt.SetAmbient(model.GL_FRONT, val1)
		mt.SetAmbient(model.GL_BACK, val2)
	}
}

func writeAmbient(os *io.OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.AmbientFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetAmbient(model.GL_FRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetAmbient(model.GL_BACK))
	os.Write(os.CRLF)
}

func checkDiffuse(pro interface{}) bool {
	return true
}
func readDiffuse(is *io.OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 vec4.T
	is.Read(&val1)
	var val2 vec4.T
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetDiffuse(model.GL_FRONT_AND_BACK, val1)
	} else {
		mt.SetDiffuse(model.GL_FRONT, val1)
		mt.SetDiffuse(model.GL_BACK, val2)
	}
}

func writeDiffuse(os *io.OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.DiffuseFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetDiffuse(model.GL_FRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetDiffuse(model.GL_BACK))
	os.Write(os.CRLF)
}

func checkSpecular(pro interface{}) bool {
	return true
}
func readSpecular(is *io.OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 vec4.T
	is.Read(&val1)
	var val2 vec4.T
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetSpecular(model.GL_FRONT_AND_BACK, val1)
	} else {
		mt.SetSpecular(model.GL_FRONT, val1)
		mt.SetSpecular(model.GL_BACK, val2)
	}
}

func writeSpecular(os *io.OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.SpecularFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetSpecular(model.GL_FRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetSpecular(model.GL_BACK))
	os.Write(os.CRLF)
}

func checkEmission(pro interface{}) bool {
	return true
}
func readEmission(is *io.OsgIstream, pro interface{}) {
	mt := pro.(*model.Material)
	frontAndBack := false
	is.Read(&frontAndBack)
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	var val1 vec4.T
	is.Read(&val1)
	var val2 vec4.T
	is.PROPERTY.Name = "Front"
	is.Read(is.PROPERTY)
	is.Read(&val2)
	if frontAndBack {
		mt.SetEmission(model.GL_FRONT_AND_BACK, val1)
	} else {
		mt.SetEmission(model.GL_FRONT, val1)
		mt.SetEmission(model.GL_BACK, val2)
	}
}

func writeEmission(os *io.OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.EmissionFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetEmission(model.GL_FRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetEmission(model.GL_BACK))
	os.Write(os.CRLF)
}

func checkShininess(pro interface{}) bool {
	return true
}
func readShininess(is *io.OsgIstream, pro interface{}) {
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
		mt.SetShininess(model.GL_FRONT_AND_BACK, val1)
	} else {
		mt.SetShininess(model.GL_FRONT, val1)
		mt.SetShininess(model.GL_BACK, val2)
	}
}

func writeShininess(os *io.OsgOstream, pro interface{}) {
	mt := pro.(*model.Material)
	os.Write(mt.ShininessFrontAndBack)
	os.PROPERTY.Name = "Front"
	os.Write(os.PROPERTY)
	os.Write(mt.GetShininess(model.GL_FRONT))
	os.PROPERTY.Name = "Back"
	os.Write(os.PROPERTY)
	os.Write(mt.GetShininess(model.GL_BACK))
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
		return &mt
	}
	wrap := io.NewObjectWrapper("Material", fn, "osg::Object osg::StateAttribute osg::Material")
	ser1 := io.NewEnumSerializer("ColorMode", getColorMode, setColorMode)
	ser1.Add("AMBIENT", model.AMBIENT)
	ser1.Add("DIFFUSE", model.DIFFUSE)
	ser1.Add("SPECULAR", model.SPECULAR)
	ser1.Add("EMISSION", model.EMISSION)
	ser1.Add("AMBIENT_AND_DIFFUSE", model.AMBIENT_AND_DIFFUSE)
	ser1.Add("MTLOFF", model.MTLOFF)
	wrap.AddSerializer(&ser1, io.RW_ENUM)

	ser2 := io.NewUserSerializer("Ambient", checkAmbient, readAmbient, writeAmbient)
	ser3 := io.NewUserSerializer("Diffuse", checkDiffuse, readDiffuse, writeDiffuse)
	ser4 := io.NewUserSerializer("Specular", checkSpecular, readSpecular, writeSpecular)
	ser5 := io.NewUserSerializer("Emission", checkEmission, readEmission, writeEmission)
	ser6 := io.NewUserSerializer("Shininess", checkShininess, readShininess, writeShininess)

	wrap.AddSerializer(&ser2, io.RW_USER)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)
	wrap.AddSerializer(&ser6, io.RW_USER)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
