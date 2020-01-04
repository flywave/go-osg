package osg

import (
	"strings"

	"github.com/flywave/go-osg/model"
)

const INT32_MAX = int32(^uint32(0) >> 1)

var manager *objectWrapperManager

type ObjectWrapperAssociate struct {
	FirstVersion int32
	LastVersion  int32
	Name         string
}

type CreateInstanceFuncType func() interface{}

type ObjectWrapper struct {
	CreateInstanceFunc                   CreateInstanceFuncType
	Domain                               string
	Name                                 string
	Associates                           []*ObjectWrapperAssociate
	Serializers                          []interface{}
	BackupSerializers                    []interface{}
	TypeList                             []SerType
	Version                              int32
	IsAssociatesRevisionsInheritanceDone bool
}

func NewObjectWrapper(name string, fn CreateInstanceFuncType, associates string) *ObjectWrapper {
	ow := &ObjectWrapper{Name: name, CreateInstanceFunc: fn, Version: 0, IsAssociatesRevisionsInheritanceDone: false}
	ow.SplitAssociates(associates, " ")
	return ow
}

func NewObjectWrapper2(name string, domain string, fn CreateInstanceFuncType, associates string) *ObjectWrapper {
	ow := &ObjectWrapper{Name: name, Domain: domain, CreateInstanceFunc: fn, Version: 0, IsAssociatesRevisionsInheritanceDone: false}
	ow.SplitAssociates(associates, " ")
	return ow
}

func (wp *ObjectWrapper) SplitAssociates(str string, separator string) {
	list := strings.Split(str, separator)
	if separator == "" {
		separator = " "
	}

	for _, l := range list {
		if l != separator {
			owa := ObjectWrapperAssociate{Name: l, FirstVersion: 0, LastVersion: INT32_MAX}
			wp.Associates = append(wp.Associates, &owa)
		}
	}
}

func (wp *ObjectWrapper) CreateInstance() interface{} {
	if wp.CreateInstanceFunc == nil {
		return nil
	}
	return wp.CreateInstanceFunc()
}

func (wp *ObjectWrapper) AddSerializer(s interface{}, t SerType) {
	s.(Serializer).SetFirstVersion(wp.Version)
	wp.Serializers = append(wp.Serializers, s)
	wp.TypeList = append(wp.TypeList, t)
}

func (wp *ObjectWrapper) MarkSerializerAsAdded(name string) {
	for _, as := range wp.Associates {
		if as.Name == name {
			as.FirstVersion = wp.Version
			return
		}
	}
}

func (wp *ObjectWrapper) MarkSerializerAsRemoved(name string) {
	for _, s := range wp.Serializers {
		ser := s.(Serializer)
		if ser.GetSerializerName() == name {
			ser.SetLastVersion(wp.Version - 1)
		}
	}
}

func (wp *ObjectWrapper) GetSerializer(name string) interface{} {
	for _, s := range wp.Serializers {
		ser := s.(Serializer)
		if ser.GetSerializerName() == name {
			return s
		}
	}

	for _, as := range wp.Associates {
		w := GetObjectWrapperManager().FindWrap(as.Name)
		if w == nil {
			continue
		}
		for _, s := range w.Serializers {
			ser := s.(Serializer)
			if ser.GetSerializerName() == name {
				return s
			}
		}
	}
	return nil
}

func (wp *ObjectWrapper) GetSerializerAndType(name string, ty *SerType) interface{} {
	for i, s := range wp.Serializers {
		ser := s.(Serializer)
		if ser.GetSerializerName() == name {
			*ty = wp.TypeList[i]
			return s
		}
	}

	for _, as := range wp.Associates {
		w := GetObjectWrapperManager().FindWrap(as.Name)
		if w == nil {
			continue
		}
		for _, s := range w.Serializers {
			ser := s.(Serializer)
			if ser.GetSerializerName() == name {
				*ty = w.TypeList[0]
				return s
			}
		}
	}
	*ty = RWUNDEFINED
	return nil
}

func (wp *ObjectWrapper) Read(is *OsgIstream, obj interface{}) {
	inputVersion := is.GetFileVersion(wp.Domain)
	for _, s := range wp.Serializers {
		ser := s.(Serializer)
		fv := ser.GetFirstVersion()
		lv := ser.GetLastVersion()
		srw := ser.SupportsReadWrite()
		if fv <= inputVersion && inputVersion <= lv && srw {
			s := Serializer(ser)
			s.Read(is, obj)
		}
	}
}

func (wp *ObjectWrapper) Write(os *OsgOstream, obj interface{}) {
	inputVersion := os.GetFileVersion(wp.Domain)
	for _, s := range wp.Serializers {
		ser := s.(Serializer)
		if ser.GetFirstVersion() <= inputVersion &&
			inputVersion <= ser.GetLastVersion() && ser.SupportsGetSet() {
			s := Serializer(ser)
			s.Write(os, obj)
		}
	}
}

func (wp *ObjectWrapper) ReadSchema(properties []string, types []SerType) bool {
	if len(wp.BackupSerializers) != 0 {
		wp.BackupSerializers = wp.Serializers
	}
	wp.Serializers = wp.Serializers[0:0]
	size := len(properties)
	serializersSize := len(wp.BackupSerializers)

	for i := 0; i < size; i++ {
		if serializersSize < i {
			break
		}
		prop := properties[i]
		ser := wp.BackupSerializers[i].(Serializer)
		if prop == ser.GetSerializerName() {
			wp.Serializers = append(wp.Serializers, wp.BackupSerializers[i])
		} else {
			for _, s := range wp.Serializers {
				ser := s.(Serializer)
				if prop != ser.GetSerializerName() {
					continue
				}
				wp.Serializers = append(wp.Serializers, ser)
			}
		}
	}
	return size == len(wp.Serializers)
}

func (wp *ObjectWrapper) WriteSchema(properties []string, types []SerType) {
	ssize := len(wp.Serializers)
	tsize := len(wp.TypeList)
	i := 0
	for {
		if ssize-1 <= i || tsize-1 <= i {
			break
		}
		s := wp.Serializers[i]
		ser := s.(Serializer)
		t := wp.TypeList[i]
		if ser.SupportsGetSet() {
			properties = append(properties, ser.GetSerializerName())
			types = append(types, t)
		}
		i++
	}
}

func (wp *ObjectWrapper) ResetSchema() {
	if len(wp.BackupSerializers) > 0 {
		wp.Serializers = wp.BackupSerializers
	}
}

type AddPropFuncType func(obj *ObjectWrapper)
type AddPropCustomFuncType func(str string, obj *ObjectWrapper)

func NewRegisterCustomWrapperProxy(instfunc CreateInstanceFuncType, domain string, name string, associates string) {
	wrap := NewObjectWrapper2(name, domain, instfunc, associates)
	wrap.CreateInstanceFunc = instfunc
	wrap.Name = name
	wrap.Domain = domain
	wrap.SplitAssociates(associates, " ")
	GetObjectWrapperManager().AddWrap(wrap)
}

type objectWrapperManager struct {
	Wraps       map[string]*ObjectWrapper
	Compressors map[string]*CompressorStream
	GlobalMap   map[string]*IntLookup
}

func newObjectWrapperManager() *objectWrapperManager {
	if manager != nil {
		return manager
	}
	obj := objectWrapperManager{Wraps: make(map[string]*ObjectWrapper), Compressors: make(map[string]*CompressorStream), GlobalMap: make(map[string]*IntLookup)}
	lk := NewIntLookup()
	obj.GlobalMap["GL"] = lk
	lk.Add("GLALPHATEST", model.GLALPHATEST)
	lk.Add("GLBLEND", model.GLBLEND)
	lk.Add("GLCOLORLOGICOP", model.GLCOLORLOGICOP)
	lk.Add("GLCOLORMATERIAL", model.GLCOLORMATERIAL)
	lk.Add("GLCULLFACE", model.GLCULLFACE)
	lk.Add("GLDEPTHTEST", model.GLDEPTHTEST)
	lk.Add("GLFOG", model.GLFOG)
	lk.Add("GLFRAGMENTPROGRAMARB", model.GLFRAGMENTPROGRAMARB)
	lk.Add("GLLINESTIPPLE", model.GLLINESTIPPLE)
	lk.Add("GLPOINTSMOOTH", model.GLPOINTSMOOTH)
	lk.Add("GLPOINTSPRITEARB", model.GLPOINTSPRITEARB)
	lk.Add("GLPOLYGONOFFSETFILL", model.GLPOLYGONOFFSETFILL)
	lk.Add("GLPOLYGONOFFSETLINE", model.GLPOLYGONOFFSETLINE)
	lk.Add("GLPOLYGONOFFSETPOINT", model.GLPOLYGONOFFSETPOINT)
	lk.Add("GLPOLYGONSTIPPLE", model.GLPOLYGONSTIPPLE)
	lk.Add("GLSCISSORTEST", model.GLSCISSORTEST)
	lk.Add("GLSTENCILTEST", model.GLSTENCILTEST)
	lk.Add("GLSTENCILTESTTWOSIDE", model.GLSTENCILTESTTWOSIDE)
	lk.Add("GLVERTEXPROGRAMARB", model.GLVERTEXPROGRAMARB)

	lk.Add("GLCOLORSUM", model.GLCOLORSUM)
	lk.Add("GLLIGHTING", model.GLLIGHTING)
	lk.Add("GLNORMALIZE", model.GLNORMALIZE)
	lk.Add("GLRESCALENORMAL", model.GLRESCALENORMAL)

	lk.Add("GLTEXTURE1D", model.GLTEXTURE1D)
	lk.Add("GLTEXTURE2D", model.GLTEXTURE2D)
	lk.Add("GLTEXTURE3D", model.GLTEXTURE3D)
	lk.Add("GLTEXTURECUBEMAP", model.GLTEXTURECUBEMAP)
	lk.Add("GLTEXTURERECTANGLE", model.GLTEXTURERECTANGLE)
	lk.Add("GLTEXTUREGENQ", model.GLTEXTUREGENQ)
	lk.Add("GLTEXTUREGENR", model.GLTEXTUREGENR)
	lk.Add("GLTEXTUREGENS", model.GLTEXTUREGENS)
	lk.Add("GLTEXTUREGENT", model.GLTEXTUREGENT)

	lk.Add("GLCLIPPLANE0", model.GLCLIPPLANE0)
	lk.Add("GLCLIPPLANE1", model.GLCLIPPLANE1)
	lk.Add("GLCLIPPLANE2", model.GLCLIPPLANE2)
	lk.Add("GLCLIPPLANE3", model.GLCLIPPLANE3)
	lk.Add("GLCLIPPLANE4", model.GLCLIPPLANE4)
	lk.Add("GLCLIPPLANE5", model.GLCLIPPLANE5)

	lk.Add("GLLIGHT0", model.GLLIGHT0)
	lk.Add("GLLIGHT1", model.GLLIGHT1)
	lk.Add("GLLIGHT2", model.GLLIGHT2)
	lk.Add("GLLIGHT3", model.GLLIGHT3)
	lk.Add("GLLIGHT4", model.GLLIGHT4)
	lk.Add("GLLIGHT5", model.GLLIGHT5)
	lk.Add("GLLIGHT6", model.GLLIGHT6)
	lk.Add("GLLIGHT7", model.GLLIGHT7)

	lk.Add("GLVERTEXPROGRAMPOINTSIZE", model.GLVERTEXPROGRAMPOINTSIZE)
	lk.Add("GLVERTEXPROGRAMTWOSIDE", model.GLVERTEXPROGRAMTWOSIDE)

	// Functions
	lk.Add("NEVER", model.GLNEVER)
	lk.Add("LESS", model.GLLESS)
	lk.Add("EQUAL", model.GLEQUAL)
	lk.Add("LEQUAL", model.GLLEQUAL)
	lk.Add("GREATER", model.GLGREATER)
	lk.Add("NOTEQUAL", model.GLNOTEQUAL)
	lk.Add("GEQUAL", model.GLGEQUAL)
	lk.Add("ALWAYS", model.GLALWAYS)

	// Texture environment states
	lk.Add("REPLACE", model.GLREPLACE)
	lk.Add("MODULATE", model.GLMODULATE)
	lk.Add("Add", model.GLADD)
	lk.Add("AddSIGNED", model.GLADDSIGNEDARB)
	lk.Add("INTERPOLATE", model.GLINTERPOLATEARB)
	lk.Add("SUBTRACT", model.GLSUBTRACTARB)
	lk.Add("DOT3RGB", model.GLDOT3RGBARB)
	lk.Add("DOT3RGBA", model.GLDOT3RGBAARB)

	lk.Add("CONSTANT", model.GLCONSTANTARB)
	lk.Add("PRIMARYCOLOR", model.GLPRIMARYCOLORARB)
	lk.Add("PREVIOUS", model.GLPREVIOUSARB)
	lk.Add("TEXTURE", model.GLTEXTURE)
	lk.Add("TEXTURE0", model.GLTEXTURE0)
	lk.Add("TEXTURE1", model.GLTEXTURE0+1)
	lk.Add("TEXTURE2", model.GLTEXTURE0+2)
	lk.Add("TEXTURE3", model.GLTEXTURE0+3)
	lk.Add("TEXTURE4", model.GLTEXTURE0+4)
	lk.Add("TEXTURE5", model.GLTEXTURE0+5)
	lk.Add("TEXTURE6", model.GLTEXTURE0+6)
	lk.Add("TEXTURE7", model.GLTEXTURE0+7)

	// Texture clamp modes
	lk.Add("CLAMP", model.GLCLAMP)
	lk.Add("CLAMPTOEDGE", model.GLCLAMPTOEDGE)
	lk.Add("CLAMPTOBORDER", model.GLCLAMPTOBORDERARB)
	lk.Add("REPEAT", model.GLREPEAT)
	lk.Add("MIRROR", model.GLMIRROREDREPEATIBM)

	// Texture filter modes
	lk.Add("LINEAR", model.GLLINEAR)
	lk.Add("LINEARMIPMAPLINEAR", model.GLLINEARMIPMAPLINEAR)
	lk.Add("LINEARMIPMAPNEAREST", model.GLLINEARMIPMAPNEAREST)
	lk.Add("NEAREST", model.GLNEAREST)
	lk.Add("NEARESTMIPMAPLINEAR", model.GLNEARESTMIPMAPLINEAR)
	lk.Add("NEARESTMIPMAPNEAREST", model.GLNEARESTMIPMAPNEAREST)

	// Texture formats
	lk.Add("GLINTENSITY", model.GLINTENSITY)
	lk.Add("GLLUMINANCE", model.GLLUMINANCE)
	lk.Add("GLALPHA", model.GLALPHA)
	lk.Add("GLLUMINANCEALPHA", model.GLLUMINANCEALPHA)
	lk.Add("GLRGB", model.GLRGB)
	lk.Add("GLRGBA", model.GLRGBA)
	lk.Add("GLCOMPRESSEDALPHAARB", model.GLCOMPRESSEDALPHAARB)
	lk.Add("GLCOMPRESSEDLUMINANCEARB", model.GLCOMPRESSEDLUMINANCEARB)
	lk.Add("GLCOMPRESSEDINTENSITYARB", model.GLCOMPRESSEDINTENSITYARB)
	lk.Add("GLCOMPRESSEDLUMINANCEALPHAARB",
		model.GLCOMPRESSEDLUMINANCEALPHAARB)
	lk.Add("GLCOMPRESSEDRGBARB", model.GLCOMPRESSEDRGBARB)
	lk.Add("GLCOMPRESSEDRGBAARB", model.GLCOMPRESSEDRGBAARB)
	lk.Add("GLCOMPRESSEDRGBS3TCDXT1EXT",
		model.GLCOMPRESSEDRGBS3TCDXT1EXT)
	lk.Add("GLCOMPRESSEDRGBAS3TCDXT1EXT",
		model.GLCOMPRESSEDRGBAS3TCDXT1EXT)
	lk.Add("GLCOMPRESSEDRGBAS3TCDXT3EXT",
		model.GLCOMPRESSEDRGBAS3TCDXT3EXT)
	lk.Add("GLCOMPRESSEDRGBAS3TCDXT5EXT",
		model.GLCOMPRESSEDRGBAS3TCDXT5EXT)
	lk.Add("GLCOMPRESSEDRGBPVRTC4BPPV1IMG",
		model.GLCOMPRESSEDRGBPVRTC4BPPV1IMG)
	lk.Add("GLCOMPRESSEDRGBPVRTC2BPPV1IMG",
		model.GLCOMPRESSEDRGBPVRTC2BPPV1IMG)
	lk.Add("GLCOMPRESSEDRGBAPVRTC4BPPV1IMG",
		model.GLCOMPRESSEDRGBAPVRTC4BPPV1IMG)
	lk.Add("GLCOMPRESSEDRGBAPVRTC2BPPV1IMG",
		model.GLCOMPRESSEDRGBAPVRTC2BPPV1IMG)
	lk.Add("GLETC1RGB8OES", model.GLETC1RGB8OES)
	lk.Add("GLCOMPRESSEDRGB8ETC2", model.GLCOMPRESSEDRGB8ETC2)
	lk.Add("GLCOMPRESSEDSRGB8ETC2", model.GLCOMPRESSEDSRGB8ETC2)
	lk.Add("GLCOMPRESSEDRGB8PUNCHTHROUGHALPHA1ETC2",
		model.GLCOMPRESSEDRGB8PUNCHTHROUGHALPHA1ETC2)
	lk.Add("GLCOMPRESSEDSRGB8PUNCHTHROUGHALPHA1ETC2",
		model.GLCOMPRESSEDSRGB8PUNCHTHROUGHALPHA1ETC2)
	lk.Add("GLCOMPRESSEDRGBA8ETC2EAC", model.GLCOMPRESSEDRGBA8ETC2EAC)
	lk.Add("GLCOMPRESSEDSRGB8ALPHA8ETC2EAC",
		model.GLCOMPRESSEDSRGB8ALPHA8ETC2EAC)
	lk.Add("GLCOMPRESSEDR11EAC", model.GLCOMPRESSEDR11EAC)
	lk.Add("GLCOMPRESSEDSIGNEDR11EAC", model.GLCOMPRESSEDSIGNEDR11EAC)
	lk.Add("GLCOMPRESSEDRG11EAC", model.GLCOMPRESSEDRG11EAC)
	lk.Add("GLCOMPRESSEDSIGNEDRG11EAC", model.GLCOMPRESSEDSIGNEDRG11EAC)

	// Texture source types
	lk.Add("GLBYTE", model.GLBYTE)
	lk.Add("GLSHORT", model.GLSHORT)
	lk.Add("GLINT", model.GLINT)
	lk.Add("GLFLOAT", model.GLFLOAT)
	lk.Add("GLDOUBLE", model.GLDOUBLE)
	lk.Add("GLUNSIGNEDBYTE", model.GLUNSIGNEDBYTE)
	lk.Add("GLUNSIGNEDSHORT", model.GLUNSIGNEDSHORT)
	lk.Add("GLUNSIGNEDINT", model.GLUNSIGNEDINT)

	// Blend values
	lk.Add("DSTALPHA", model.GLDSTALPHA)
	lk.Add("DSTCOLOR", model.GLDSTCOLOR)
	lk.Add("ONE", model.GLONE)
	lk.Add("ONEMINUSDSTALPHA", model.GLONEMINUSDSTALPHA)
	lk.Add("ONEMINUSDSTCOLOR", model.GLONEMINUSDSTCOLOR)
	lk.Add("ONEMINUSSRCALPHA", model.GLONEMINUSSRCALPHA)
	lk.Add("ONEMINUSSRCCOLOR", model.GLONEMINUSSRCCOLOR)
	lk.Add("SRCALPHA", model.GLSRCALPHA)
	lk.Add("SRCALPHASATURATE", model.GLSRCALPHASATURATE)
	lk.Add("SRCCOLOR", model.GLSRCCOLOR)
	lk.Add("CONSTANTCOLOR", model.GLCONSTANTCOLOR)
	lk.Add("ONEMINUSCONSTANTCOLOR", model.GLONEMINUSCONSTANTCOLOR)
	lk.Add("CONSTANTALPHA", model.GLCONSTANTALPHA)
	lk.Add("ONEMINUSCONSTANTALPHA", model.GLONEMINUSCONSTANTALPHA)
	lk.Add("ZERO", model.GLZERO)

	// Fog coordinate sources
	lk.Add("COORDINATE", model.GLFOGCOORDINATE)
	lk.Add("DEPTH", model.GLFRAGMENTDEPTH)

	// Hint targets
	lk.Add("FOGHINT", model.GLFOGHINT)
	lk.Add("GENERATEMIPMAPHINT", model.GLGENERATEMIPMAPHINTSGIS)
	lk.Add("LINESMOOTHHINT", model.GLLINESMOOTHHINT)
	lk.Add("PERSPECTIVECORRECTIONHINT", model.GLPERSPECTIVECORRECTIONHINT)
	lk.Add("POINTSMOOTHHINT", model.GLPOINTSMOOTHHINT)
	lk.Add("POLYGONSMOOTHHINT", model.GLPOLYGONSMOOTHHINT)
	lk.Add("TEXTURECOMPRESSIONHINT", model.GLTEXTURECOMPRESSIONHINTARB)
	lk.Add("FRAGMENTSHADERDERIVATIVEHINT",
		model.GLFRAGMENTSHADERDERIVATIVEHINT)

	// Polygon modes
	lk.Add("POINT", model.GLPOINT)
	lk.Add("LINE", model.GLLINE)
	lk.Add("FILL", model.GLFILL)

	// Misc
	lk.Add("BACK", model.GLBACK)
	lk.Add("FRONT", model.GLFRONT)
	lk.Add("FRONTANDBACK", model.GLFRONTANDBACK)
	lk.Add("FIXEDONLY", model.GLFIXEDONLY)
	lk.Add("FASTEST", model.GLFASTEST)
	lk.Add("NICEST", model.GLNICEST)
	lk.Add("DONTCARE", model.GLDONTCARE)

	arrayTable := NewIntLookup()
	obj.GlobalMap["ArrayType"] = arrayTable

	arrayTable.Add("ByteArray", model.IDBYTEARRAY)
	arrayTable.Add("UByteArray", model.IDUBYTEARRAY)
	arrayTable.Add("ShortArray", model.IDSHORTARRAY)
	arrayTable.Add("UShortArray", model.IDUSHORTARRAY)
	arrayTable.Add("IntArray", model.IDINTARRAY)
	arrayTable.Add("UIntArray", model.IDUINTARRAY)
	arrayTable.Add("FloatArray", model.IDFLOATARRAY)
	arrayTable.Add("DoubleArray", model.IDDOUBLEARRAY)

	arrayTable.Add("Vec2bArray", model.IDVEC2BARRAY)
	arrayTable.Add("Vec3bArray", model.IDVEC3BARRAY)
	arrayTable.Add("Vec4bArray", model.IDVEC4BARRAY)
	arrayTable.Add("Vec2ubArray", model.IDVEC2UBARRAY)
	arrayTable.Add("Vec3ubArray", model.IDVEC3UBARRAY)
	arrayTable.Add("Vec4ubArray", model.IDVEC4UBARRAY)
	arrayTable.Add("Vec2sArray", model.IDVEC2SARRAY)
	arrayTable.Add("Vec3sArray", model.IDVEC3SARRAY)
	arrayTable.Add("Vec4sArray", model.IDVEC4SARRAY)
	arrayTable.Add("Vec2usArray", model.IDVEC2USARRAY)
	arrayTable.Add("Vec3usArray", model.IDVEC3USARRAY)
	arrayTable.Add("Vec4usArray", model.IDVEC4USARRAY)
	arrayTable.Add("Vec2fArray", model.IDVEC2ARRAY)
	arrayTable.Add("Vec3fArray", model.IDVEC3ARRAY)
	arrayTable.Add("Vec4fArray", model.IDVEC4ARRAY)
	arrayTable.Add("Vec2dArray", model.IDVEC2DARRAY)
	arrayTable.Add("Vec3dArray", model.IDVEC3DARRAY)
	arrayTable.Add("Vec4dArray", model.IDVEC4DARRAY)

	arrayTable.Add("Vec2iArray", model.IDVEC2IARRAY)
	arrayTable.Add("Vec3iArray", model.IDVEC3IARRAY)
	arrayTable.Add("Vec4iArray", model.IDVEC4IARRAY)
	arrayTable.Add("Vec2uiArray", model.IDVEC2UIARRAY)
	arrayTable.Add("Vec3uiArray", model.IDVEC3UIARRAY)
	arrayTable.Add("Vec4uiArray", model.IDVEC4UIARRAY)

	primitiveTable := NewIntLookup()
	obj.GlobalMap["PrimitiveType"] = primitiveTable

	primitiveTable.Add("DrawArrays", model.IDDRAWARRAYS)
	primitiveTable.Add("DrawArraysLength", model.IDDRAWARRAYLENGTH)
	primitiveTable.Add("DrawElementsUByte", model.IDDRAWELEMENTSUBYTE)
	primitiveTable.Add("DrawElementsUShort", model.IDDRAWELEMENTSUSHORT)
	primitiveTable.Add("DrawElementsUInt", model.IDDRAWELEMENTSUINT)

	primitiveTable.Add("GLPOINTS", model.GLPOINTS)
	primitiveTable.Add("GLLINES", model.GLLINES)
	primitiveTable.Add("GLLINESTRIP", model.GLLINESTRIP)
	primitiveTable.Add("GLLINELOOP", model.GLLINELOOP)
	primitiveTable.Add("GLTRIANGLES", model.GLTRIANGLES)
	primitiveTable.Add("GLTRIANGLESTRIP", model.GLTRIANGLESTRIP)
	primitiveTable.Add("GLTRIANGLEFAN", model.GLTRIANGLEFAN)
	primitiveTable.Add("GLQUADS", model.GLQUADS)
	primitiveTable.Add("GLQUADSTRIP", model.GLQUADSTRIP)
	primitiveTable.Add("GLPOLYGON", model.GLPOLYGON)

	primitiveTable.Add2("GLLINESADJACENCYEXT", "GLLINESADJACENCY",
		model.GLLINESADJACENCY)
	primitiveTable.Add2("GLLINESTRIPADJACENCYEXT", "GLLINESTRIPADJACENCY",
		model.GLLINESTRIPADJACENCY)
	primitiveTable.Add2("GLTRIANGLESADJACENCYEXT", "GLTRIANGLESADJACENCY",
		model.GLTRIANGLESADJACENCY)
	primitiveTable.Add2("GLTRIANGLESTRIPADJACENCYEXT",
		"GLTRIANGLESTRIPADJACENCY",
		model.GLTRIANGLESTRIPADJACENCY)

	primitiveTable.Add("GLPATCHES", model.GLPATCHES)
	return &obj
}

func GetObjectWrapperManager() *objectWrapperManager {
	if manager == nil {
		manager = newObjectWrapperManager()
	}
	return manager
}

func (man *objectWrapperManager) AddWrap(wrap *ObjectWrapper) {
	if wrap == nil {
		return
	}
	manager.Wraps["osg::"+strings.ToLower(wrap.Name)] = wrap
}

func (man *objectWrapperManager) RemoveWrap(wrap *ObjectWrapper) {
	nm := strings.ToLower(wrap.Name)
	delete(manager.Wraps, nm)
}

func (man *objectWrapperManager) FindWrap(str string) *ObjectWrapper {
	nm := strings.ToLower(str)
	w, ok := manager.Wraps[nm]
	if ok {
		return w
	}
	return nil
}

func (man *objectWrapperManager) FindLookup(group string) *IntLookup {
	lk, ok := man.GlobalMap[group]
	if !ok {
		lk, ok = man.GlobalMap["GL"]
	}
	return lk
}

func (man *objectWrapperManager) AddCompressor(st *CompressorStream) {
	man.Compressors[st.Name] = st
}

func (man *objectWrapperManager) RemoveCompressor(st *CompressorStream) {
	if st == nil {
		return
	}
	_, ok := man.Compressors[st.Name]
	if ok {
		delete(man.Compressors, st.Name)
	}
}

func (man *objectWrapperManager) FindCompressor(st string) *CompressorStream {
	lk, ok := man.Compressors[st]
	if ok {
		return lk
	}
	return nil
}

type UpdateWrapperVersionProxy struct {
	Wrap        *ObjectWrapper
	LastVersion int32
}

func (uv *UpdateWrapperVersionProxy) SetLastVersion() {
	uv.Wrap.Version = uv.LastVersion

}

func AddUpdateWrapperVersionProxy(w *ObjectWrapper, v int32) *UpdateWrapperVersionProxy {
	prox := UpdateWrapperVersionProxy{Wrap: w, LastVersion: w.Version}
	w.Version = v
	return &prox
}
