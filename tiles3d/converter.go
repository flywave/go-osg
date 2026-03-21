package tiles3d

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/flywave/go-draco"
	"github.com/flywave/go-meshopt"
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

type GeometryConverter struct {
	opts       *ConverterOptions
	coordTrans *CoordinateTransformer
	geoidConv  *GeoidConverter
}

func NewGeometryConverter(opts *ConverterOptions, coordTrans *CoordinateTransformer, geoidConv *GeoidConverter) *GeometryConverter {
	return &GeometryConverter{
		opts:       opts,
		coordTrans: coordTrans,
		geoidConv:  geoidConv,
	}
}

func (c *GeometryConverter) Convert(node interface{}) *TileContent {
	content := &TileContent{}

	switch n := node.(type) {
	case *model.Group:
		children := n.GetChildren()
		for _, child := range children {
			childContent := c.Convert(child)
			if childContent != nil {
				content.Merge(childContent)
			}
		}
	case *model.Lod:
		children := n.GetChildren()
		if len(children) > 0 {
			content = c.Convert(children[0])
		}
	case *model.PagedLod:
		children := n.GetChildren()
		if len(children) > 0 {
			content = c.Convert(children[0])
		}
	case *model.MatrixTransform:
		children := n.GetChildren()
		for _, child := range children {
			childContent := c.Convert(child)
			if childContent != nil {
				c.applyMatrixTransform(childContent, n.Matrix)
				content.Merge(childContent)
			}
		}
	case *model.Geode:
		children := n.GetChildren()
		for _, child := range children {
			if geometry, ok := child.(*model.Geometry); ok {
				geomContent := c.extractGeometry(geometry)
				if geomContent != nil {
					content.Merge(geomContent)
				}
			}
		}
	}

	return content
}

func (c *GeometryConverter) extractGeometry(geom *model.Geometry) *TileContent {
	content := &TileContent{}

	if geom.VertexArray == nil || geom.VertexArray.Data == nil {
		return nil
	}

	vertices := c.extractVertices(geom.VertexArray)
	if len(vertices) == 0 {
		return nil
	}

	bbox := c.calculateBoundingBox(vertices)
	content.BoundingBox = bbox
	content.Vertices = vertices

	if geom.NormalArray != nil && geom.NormalArray.Data != nil {
		content.Normals = c.extractNormals(geom.NormalArray)
	}

	if len(geom.TexCoordArrayList) > 0 && geom.TexCoordArrayList[0] != nil {
		content.TexCoords = c.extractTexCoords(geom.TexCoordArrayList[0])
	}

	if geom.Primitives != nil && len(geom.Primitives) > 0 {
		content.Indices = c.extractIndices(geom.Primitives[0])
	}

	content.BatchLength = len(vertices) / 3

	return content
}

func (c *GeometryConverter) applyMatrixTransform(content *TileContent, matrix [4][4]float32) {
	if len(content.Vertices) == 0 {
		return
	}

	for i := 0; i < len(content.Vertices); i += 3 {
		x := float64(content.Vertices[i])
		y := float64(content.Vertices[i+1])
		z := float64(content.Vertices[i+2])

		m00 := float64(matrix[0][0])
		m01 := float64(matrix[0][1])
		m02 := float64(matrix[0][2])
		m03 := float64(matrix[0][3])
		m10 := float64(matrix[1][0])
		m11 := float64(matrix[1][1])
		m12 := float64(matrix[1][2])
		m13 := float64(matrix[1][3])
		m20 := float64(matrix[2][0])
		m21 := float64(matrix[2][1])
		m22 := float64(matrix[2][2])
		m23 := float64(matrix[2][3])

		newX := m00*x + m01*y + m02*z + m03
		newY := m10*x + m11*y + m12*z + m13
		newZ := m20*x + m21*y + m22*z + m23

		content.Vertices[i] = float32(newX)
		content.Vertices[i+1] = float32(newY)
		content.Vertices[i+2] = float32(newZ)
	}

	if len(content.Normals) > 0 {
		for i := 0; i < len(content.Normals); i += 3 {
			x := float64(content.Normals[i])
			y := float64(content.Normals[i+1])
			z := float64(content.Normals[i+2])

			m00 := float64(matrix[0][0])
			m01 := float64(matrix[0][1])
			m02 := float64(matrix[0][2])
			m10 := float64(matrix[1][0])
			m11 := float64(matrix[1][1])
			m12 := float64(matrix[1][2])
			m20 := float64(matrix[2][0])
			m21 := float64(matrix[2][1])
			m22 := float64(matrix[2][2])

			newX := m00*x + m01*y + m02*z
			newY := m10*x + m11*y + m12*z
			newZ := m20*x + m21*y + m22*z

			content.Normals[i] = float32(newX)
			content.Normals[i+1] = float32(newY)
			content.Normals[i+2] = float32(newZ)
		}
	}

	content.BoundingBox = calculateBoundingBox(content.Vertices)
}

const (
	meterToLat = 0.000000157891
	meterToLon = 0.000000156785
)

func meterToLatDeg(m float64) float64 {
	return m * meterToLat
}

func meterToLonDeg(m float64, latRad float64) float64 {
	return m * meterToLon / math.Cos(latRad)
}

func (c *GeometryConverter) extractVertices(arr *model.Array) []float32 {
	if arr == nil || arr.Data == nil {
		return nil
	}

	var vertices []float32

	switch data := arr.Data.(type) {
	case []float32:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		copy(vertices, data)
	case [][3]float32:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data)*3)
		for i, v := range data {
			vertices[i*3] = v[0]
			vertices[i*3+1] = v[1]
			vertices[i*3+2] = v[2]
		}
	case []float64:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		for i, v := range data {
			vertices[i] = float32(v)
		}
	case []int16:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		for i, v := range data {
			vertices[i] = float32(v)
		}
	case []uint16:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		for i, v := range data {
			vertices[i] = float32(v)
		}
	default:
		return nil
	}

	if len(vertices) < 3 {
		return nil
	}

	if len(vertices)%3 != 0 {
		vertices = vertices[:len(vertices)/3*3]
	}

	center := c.coordTrans.GetCenter()
	hasOriginOffset := center[0] != 0 || center[1] != 0 || center[2] != 0
	hasProjTransformation := c.coordTrans.HasProjection()
	isGeographicOutput := c.coordTrans.IsGeographicOutput() && !c.coordTrans.IsECEFOutput()

	// First add origin offset
	if hasOriginOffset {
		for i := 0; i < len(vertices); i += 3 {
			vertices[i] += float32(center[0])
			vertices[i+1] += float32(center[1])
			vertices[i+2] += float32(center[2])
		}
	}

	// Then apply PROJ transformation if target is set
	if hasProjTransformation && c.coordTrans.GetTargetSRS() != "" {
		for i := 0; i < len(vertices); i += 3 {
			// EPSG:4548 is (Easting, Northing) = (X, Y)
			// PROJ expects (lon, lat, height) for geographic, or (X, Y, Z) for projected
			point := [3]float64{float64(vertices[i]), float64(vertices[i+1]), float64(vertices[i+2])}

			point = c.coordTrans.Transform(point)

			if isGeographicOutput {
				point[0] = point[0] * 180.0 / math.Pi
				point[1] = point[1] * 180.0 / math.Pi
			}

			vertices[i] = float32(point[0])
			vertices[i+1] = float32(point[1])
			vertices[i+2] = float32(point[2])
		}
	}

	return vertices
}

func (c *GeometryConverter) extractNormals(arr *model.Array) []float32 {
	if arr.Data == nil {
		return nil
	}

	switch data := arr.Data.(type) {
	case []float32:
		return data
	case [][3]float32:
		normals := make([]float32, len(data)*3)
		for i, v := range data {
			normals[i*3] = v[0]
			normals[i*3+1] = v[1]
			normals[i*3+2] = v[2]
		}
		return normals
	}

	return nil
}

func (c *GeometryConverter) extractTexCoords(arr *model.Array) []float32 {
	if arr.Data == nil {
		return nil
	}

	switch data := arr.Data.(type) {
	case []float32:
		return data
	case [][2]float32:
		texCoords := make([]float32, len(data)*2)
		for i, v := range data {
			texCoords[i*2] = v[0]
			texCoords[i*2+1] = v[1]
		}
		return texCoords
	}

	return nil
}

func (c *GeometryConverter) calculateBoundingBox(vertices []float32) [6]float64 {
	if len(vertices) == 0 {
		return [6]float64{0, 0, 0, 0, 0, 0}
	}

	minX, maxX := float64(vertices[0]), float64(vertices[0])
	minY, maxY := float64(vertices[1]), float64(vertices[1])
	minZ, maxZ := float64(vertices[2]), float64(vertices[2])

	for i := 0; i < len(vertices); i += 3 {
		if float64(vertices[i]) < minX {
			minX = float64(vertices[i])
		}
		if float64(vertices[i]) > maxX {
			maxX = float64(vertices[i])
		}
		if float64(vertices[i+1]) < minY {
			minY = float64(vertices[i+1])
		}
		if float64(vertices[i+1]) > maxY {
			maxY = float64(vertices[i+1])
		}
		if float64(vertices[i+2]) < minZ {
			minZ = float64(vertices[i+2])
		}
		if float64(vertices[i+2]) > maxZ {
			maxZ = float64(vertices[i+2])
		}
	}

	return [6]float64{minX, minY, minZ, maxX, maxY, maxZ}
}

func (c *TileContent) Merge(other *TileContent) {
	if other == nil {
		return
	}

	c.Vertices = append(c.Vertices, other.Vertices...)
	c.Normals = append(c.Normals, other.Normals...)
	c.TexCoords = append(c.TexCoords, other.TexCoords...)
	c.Indices = append(c.Indices, other.Indices...)

	if c.BatchLength == 0 {
		c.BoundingBox = other.BoundingBox
	} else {
		c.BatchLength += other.BatchLength
	}
}

func LoadOSGB(path string) (interface{}, error) {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(path, nil)
	if res == nil {
		return nil, fmt.Errorf("failed to read %s", path)
	}
	return res.GetNode(), nil
}

func (c *GeometryConverter) extractMaterial(geode *model.Geode) []*Material {
	materials := make([]*Material, 0)
	children := geode.GetChildren()
	for _, child := range children {
		if drawable, ok := child.(model.DrawableInterface); ok {
			geom, ok := drawable.(*model.Geometry)
			if !ok {
				continue
			}
			if stateset := geom.GetStates(); stateset != nil {
				mat := c.extractMaterialFromStateSet(stateset)
				if mat != nil {
					materials = append(materials, mat)
				}
			}
		}
	}
	return materials
}

func (c *GeometryConverter) extractMaterialFromStateSet(stateset *model.StateSet) *Material {
	mat := &Material{
		BaseColorFactor: [4]float32{1.0, 1.0, 1.0, 1.0},
		MetallicFactor:  0.0,
		RoughnessFactor: 1.0,
		DoubleSided:     false,
		Unlit:           true,
	}

	if stateset.AttributeList == nil {
		return mat
	}

	for _, attrPair := range stateset.AttributeList {
		if attrPair == nil {
			continue
		}
		if material, ok := attrPair.First.(*model.Material); ok {
			diffuse := material.GetDiffuse(model.GLFRONTANDBACK)
			if diffuse != nil {
				mat.BaseColorFactor = *diffuse
			}
			specular := material.GetSpecular(model.GLFRONTANDBACK)
			if specular != nil && (specular[0] > 0 || specular[1] > 0 || specular[2] > 0) {
				mat.MetallicFactor = 0.9
				mat.RoughnessFactor = 0.1
			}
			emission := material.GetEmission(model.GLFRONTANDBACK)
			if emission != nil && (emission[0] > 0 || emission[1] > 0 || emission[2] > 0) {
				mat.BaseColorFactor[0] = emission[0]
				mat.BaseColorFactor[1] = emission[1]
				mat.BaseColorFactor[2] = emission[2]
			}
		}
	}
	return mat
}

func (c *GeometryConverter) extractTexture(geode *model.Geode, basePath string) []*Texture {
	textures := make([]*Texture, 0)
	children := geode.GetChildren()
	texIndex := 0
	for _, child := range children {
		if drawable, ok := child.(model.DrawableInterface); ok {
			geom, ok := drawable.(*model.Geometry)
			if !ok {
				continue
			}
			if stateset := geom.GetStates(); stateset != nil {
				if tex := c.extractTextureFromStateSet(stateset, basePath, texIndex); tex != nil {
					textures = append(textures, tex)
					texIndex++
				}
			}
		}
	}
	return textures
}

func (c *GeometryConverter) extractTextureFromStateSet(stateset *model.StateSet, basePath string, index int) *Texture {
	if stateset.TextureAttributeList == nil || len(stateset.TextureAttributeList) == 0 {
		return nil
	}

	for unit, attrList := range stateset.TextureAttributeList {
		if unit > 0 {
			break
		}
		for _, attrPair := range attrList {
			if attrPair == nil {
				continue
			}
			if tex, ok := attrPair.First.(*model.Texture); ok {
				if tex.Image != nil {
					texData := c.extractImageData(tex.Image, basePath)
					if texData != nil {
						return &Texture{
							ID:   fmt.Sprintf("texture_%d", index),
							Data: texData,
							Mime: getImageMimeType(tex.Image),
						}
					}
				}
			}
		}
	}
	return nil
}

func (c *GeometryConverter) extractImageData(img *model.Image, basePath string) []byte {
	if img == nil {
		return nil
	}

	if len(img.FileName) > 0 {
		texPath := img.FileName
		if !filepath.IsAbs(texPath) && len(basePath) > 0 {
			texPath = filepath.Join(basePath, texPath)
		}
		data, err := os.ReadFile(texPath)
		if err == nil {
			return data
		}
	}

	if img.Data != nil && len(img.Data) > 0 {
		return img.Data
	}

	return nil
}

func getImageMimeType(img *model.Image) string {
	if len(img.FileName) > 0 {
		ext := filepath.Ext(img.FileName)
		switch ext {
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".png":
			return "image/png"
		case ".ktx", ".ktx2":
			return "image/ktx2"
		}
	}
	return "image/jpeg"
}

func (c *GeometryConverter) extractIndices(prim interface{}) []uint32 {
	if prim == nil {
		return nil
	}

	switch p := prim.(type) {
	case *model.DrawElementsUByte:
		indices := make([]uint32, 0)
		if p.Data != nil {
			for _, idx := range p.Data {
				indices = append(indices, uint32(idx))
			}
		}
		mode := p.Mode
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	case *model.DrawElementsUShort:
		indices := make([]uint32, 0)
		if p.Data != nil {
			for _, idx := range p.Data {
				indices = append(indices, uint32(idx))
			}
		}
		mode := p.Mode
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	case *model.DrawElementsUInt:
		indices := make([]uint32, 0)
		if p.Data != nil {
			for _, idx := range p.Data {
				indices = append(indices, idx)
			}
		}
		mode := p.Mode
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	case *model.DrawArrays:
		mode := p.Mode
		count := int(p.Count)
		first := int(p.First)
		indices := make([]uint32, count)
		for i := 0; i < count; i++ {
			indices[i] = uint32(first + i)
		}
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	}
	return nil
}

func triangulateQuadLike(indices []uint32, mode int32) []uint32 {
	if len(indices) < 4 {
		return indices
	}

	var result []uint32

	if mode == model.QUADS {
		if len(indices)%4 != 0 {
			indices = indices[:len(indices)/4*4]
		}
		quadCount := len(indices) / 4
		result = make([]uint32, 0, quadCount*6)
		for q := 0; q < quadCount; q++ {
			base := q * 4
			result = append(result, indices[base], indices[base+1], indices[base+2])
			result = append(result, indices[base], indices[base+2], indices[base+3])
		}
		return result
	}

	if mode == model.QUADSTRIP {
		if len(indices)%2 != 0 {
			indices = indices[:len(indices)/2*2]
		}
		pairCount := len(indices) / 2
		if pairCount < 2 {
			return indices
		}
		result = make([]uint32, 0, (pairCount-1)*6)
		for i := 0; i < pairCount-1; i++ {
			base := i * 2
			if base+3 >= len(indices) {
				break
			}
			a := indices[base]
			b := indices[base+1]
			c := indices[base+2]
			d := indices[base+3]
			result = append(result, a, b, c, b, d, c)
		}
		return result
	}

	return indices
}

func (c *GeometryConverter) simplifyMesh(vertices, normals, texcoords []float32, indices []uint32) ([]float32, []float32, []float32, []uint32) {
	if c.opts == nil || !c.opts.EnableSimplify || c.opts.SimplifyRatio >= 1.0 {
		return vertices, normals, texcoords, indices
	}

	mesh := &meshopt.Mesh{}
	vertexCount := len(vertices) / 3
	indexCount := len(indices)

	if indexCount == 0 || vertexCount == 0 {
		return vertices, normals, texcoords, indices
	}

	mesh.Indices = make([][3]uint32, indexCount/3)
	for i := 0; i < indexCount/3; i++ {
		mesh.Indices[i] = [3]uint32{indices[i*3], indices[i*3+1], indices[i*3+2]}
	}

	mesh.Nodes = append(mesh.Nodes, meshopt.Node{
		Type: meshopt.AttributeTypePosition,
		Data: c.float32ToBytes(vertices),
	})

	if len(normals) > 0 {
		mesh.Nodes = append(mesh.Nodes, meshopt.Node{
			Type: meshopt.AttributeTypeNormal,
			Data: c.float32ToBytes(normals),
		})
	}

	if len(texcoords) > 0 {
		mesh.Nodes = append(mesh.Nodes, meshopt.Node{
			Type: meshopt.AttributeTypeTexCoord,
			Data: c.float32ToBytes(texcoords),
		})
	}

	settings := &meshopt.Settings{
		SimplifyThreshold:  c.opts.SimplifyRatio,
		SimplifyAggressive: false,
		CompressMore:       false,
	}

	err := meshopt.ProcessMesh(mesh, settings)
	if err != nil {
		return vertices, normals, texcoords, indices
	}

	newIndices := make([]uint32, len(mesh.Indices)*3)
	for i, tri := range mesh.Indices {
		newIndices[i*3] = tri[0]
		newIndices[i*3+1] = tri[1]
		newIndices[i*3+2] = tri[2]
	}

	return vertices, normals, texcoords, newIndices
}

func (c *GeometryConverter) float32ToBytes(data []float32) []byte {
	result := make([]byte, len(data)*4)
	for i, v := range data {
		bits := float32ToBits(v)
		result[i*4] = byte(bits)
		result[i*4+1] = byte(bits >> 8)
		result[i*4+2] = byte(bits >> 16)
		result[i*4+3] = byte(bits >> 24)
	}
	return result
}

func float32ToBits(f float32) uint32 {
	bits := uint32(0)
	copy((*[4]byte)(unsafe.Pointer(&bits))[:], (*[4]byte)(unsafe.Pointer(&f))[:])
	return bits
}

func (c *GeometryConverter) compressMeshDraco(vertices, normals, texcoords []float32, indices []uint32) ([]byte, map[string]int, error) {
	if c.opts == nil || !c.opts.EnableDraco {
		return nil, nil, nil
	}

	meshBuilder := draco.NewMeshBuilder()

	posQuantization := int32(c.opts.DracoPositionBits)
	if posQuantization == 0 {
		posQuantization = 11
	}
	normQuantization := int32(c.opts.DracoNormalBits)
	if normQuantization == 0 {
		normQuantization = 10
	}
	texQuantization := int32(c.opts.DracoTexCoordBits)
	if texQuantization == 0 {
		texQuantization = 12
	}

	vertexCount := len(vertices) / 3
	meshBuilder.Start(vertexCount)

	posData := make([]float32, len(vertices))
	copy(posData, vertices)
	meshBuilder.SetAttribute(vertexCount, posData, draco.GAT_POSITION)

	if len(normals) > 0 {
		normData := make([]float32, len(normals))
		copy(normData, normals)
		meshBuilder.SetAttribute(vertexCount, normData, draco.GAT_NORMAL)
	}

	if len(texcoords) > 0 {
		texData := make([]float32, len(texcoords))
		copy(texData, texcoords)
		meshBuilder.SetAttribute(vertexCount, texData, draco.GAT_TEX_COORD)
	}

	mesh := meshBuilder.GetMesh()

	encoder := draco.NewEncoder()
	encoder.SetAttributeQuantization(draco.GAT_POSITION, posQuantization)
	encoder.SetAttributeQuantization(draco.GAT_NORMAL, normQuantization)
	encoder.SetAttributeQuantization(draco.GAT_TEX_COORD, texQuantization)

	err, compressed := encoder.EncodeMesh(mesh)
	if err != nil {
		return nil, nil, err
	}

	attrIds := map[string]int{
		"POSITION": 0,
	}

	if len(normals) > 0 {
		attrIds["NORMAL"] = 1
	}
	if len(texcoords) > 0 {
		attrIds["TEXCOORD_0"] = 2
	}

	return compressed, attrIds, nil
}

func (c *GeometryConverter) optimizeMesh(vertices, normals, texcoords []float32, indices []uint32) ([]float32, []float32, []float32, []uint32) {
	if c.opts == nil {
		return vertices, normals, texcoords, indices
	}

	if c.opts.EnableSimplify && c.opts.SimplifyRatio < 1.0 {
		return c.simplifyMesh(vertices, normals, texcoords, indices)
	}

	return vertices, normals, texcoords, indices
}

func (c *GeometryConverter) GetDatabasePath(node interface{}) string {
	switch n := node.(type) {
	case *model.PagedLod:
		return n.DataBasePath
	}
	return ""
}

func (c *GeometryConverter) GetFileNames(node interface{}) []string {
	switch n := node.(type) {
	case *model.PagedLod:
		names := make([]string, 0)
		children := n.GetChildren()
		for i := 0; i < len(children); i++ {
			if i < len(n.PerRangeDataList) {
				if fname := n.PerRangeDataList[i].FileName; fname != "" {
					names = append(names, fname)
				}
			}
		}
		return names
	}
	return nil
}
