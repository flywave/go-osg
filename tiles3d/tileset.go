package tiles3d

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/go-osg/model"

	"github.com/flywave/gltf"
	"github.com/flywave/go-3dtile"
)

type TilesetGenerator struct {
	opts *ConverterOptions
}

func NewTilesetGenerator(opts *ConverterOptions) *TilesetGenerator {
	return &TilesetGenerator{opts: opts}
}

func (g *TilesetGenerator) Generate(tile *Tile) *TileJSON {
	tileset := &TileJSON{
		Asset: AssetJSON{
			Version: "1.0",
			GenBy:   "go-osg-3dtiles",
		},
		GeometricError: g.opts.GeometricError,
	}

	if tile != nil {
		tileset.Root = g.convertTileToJSON(tile)
	}

	return tileset
}

func (g *TilesetGenerator) convertTileToJSON(tile *Tile) *TileJSONNode {
	node := &TileJSONNode{
		GeometricError: tile.GeometricError,
		Refine:         "REPLACE",
		BoundVolume: BoundVolumeJSON{
			Box: tile.BoundingBox[:],
		},
	}

	if tile.Content != nil && tile.Path != "" {
		node.Content = &ContentJSON{
			URI: filepath.Base(tile.Path),
		}
	}

	if len(tile.Children) > 0 {
		node.Children = make([]*TileJSONNode, len(tile.Children))
		for i, child := range tile.Children {
			node.Children[i] = g.convertTileToJSON(child)
		}
	}

	return node
}

func (g *TilesetGenerator) Write(tileset *TileJSON, path string) error {
	data, err := json.MarshalIndent(tileset, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type B3DMGenerator struct {
	opts *ConverterOptions
}

func NewB3DMGenerator(opts *ConverterOptions) *B3DMGenerator {
	return &B3DMGenerator{opts: opts}
}

func (g *B3DMGenerator) Generate(tile *Tile) ([]byte, error) {
	if tile.Content == nil {
		return nil, fmt.Errorf("tile content is nil")
	}

	doc := g.buildGLTF(tile.Content)

	b3dm := tile3d.NewB3dm()
	b3dm.Model = doc

	view := tile3d.B3dmFeatureTableView{
		BatchLength: tile.Content.BatchLength,
	}

	if len(tile.Content.BoundingBox) == 6 {
		center := []float64{
			(tile.Content.BoundingBox[0] + tile.Content.BoundingBox[3]) / 2,
			(tile.Content.BoundingBox[1] + tile.Content.BoundingBox[4]) / 2,
			(tile.Content.BoundingBox[2] + tile.Content.BoundingBox[5]) / 2,
		}
		view.RtcCenter = center
	}

	b3dm.SetFeatureTable(view)

	var buf bytes.Buffer
	err := b3dm.Write(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *B3DMGenerator) GenerateGLB(tile *Tile) ([]byte, error) {
	if tile.Content == nil {
		return nil, fmt.Errorf("tile content is nil")
	}

	doc := g.buildGLTF(tile.Content)

	var buf bytes.Buffer
	encoder := gltf.NewEncoder(&buf)
	encoder.AsBinary = true
	err := encoder.Encode(doc)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func float32ToByte(data []float32) []byte {
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

func uint32ToByte(data []uint32) []byte {
	result := make([]byte, len(data)*4)
	for i, v := range data {
		result[i*4] = byte(v)
		result[i*4+1] = byte(v >> 8)
		result[i*4+2] = byte(v >> 16)
		result[i*4+3] = byte(v >> 24)
	}
	return result
}

func (g *B3DMGenerator) buildGLTF(content *TileContent) *gltf.Document {
	doc := gltf.NewDocument()

	vertexData := content.Vertices
	vertexBuffer := &gltf.Buffer{
		Data:       float32ToByte(vertexData),
		ByteLength: uint32(len(vertexData) * 4),
	}
	vertexBufferID := uint32(len(doc.Buffers))
	doc.Buffers = append(doc.Buffers, vertexBuffer)

	vertexBufferView := &gltf.BufferView{
		Buffer:     vertexBufferID,
		ByteOffset: 0,
		ByteLength: uint32(len(vertexData) * 4),
		ByteStride: 12,
	}
	vertexBufferViewID := uint32(len(doc.BufferViews))
	doc.BufferViews = append(doc.BufferViews, vertexBufferView)

	vertexAccessor := &gltf.Accessor{
		BufferView:    &vertexBufferViewID,
		ComponentType: gltf.ComponentFloat,
		Count:         uint32(len(content.Vertices) / 3),
		Type:          gltf.AccessorVec3,
		Max:           []float32{float32(content.BoundingBox[3]), float32(content.BoundingBox[4]), float32(content.BoundingBox[5])},
		Min:           []float32{float32(content.BoundingBox[0]), float32(content.BoundingBox[1]), float32(content.BoundingBox[2])},
	}
	doc.Accessors = append(doc.Accessors, vertexAccessor)

	var normalAccessor *gltf.Accessor
	if len(content.Normals) > 0 {
		normalBuffer := &gltf.Buffer{
			Data:       float32ToByte(content.Normals),
			ByteLength: uint32(len(content.Normals) * 4),
		}
		normalBufferID := uint32(len(doc.Buffers))
		doc.Buffers = append(doc.Buffers, normalBuffer)

		normalBufferView := &gltf.BufferView{
			Buffer:     normalBufferID,
			ByteOffset: 0,
			ByteLength: uint32(len(content.Normals) * 4),
			ByteStride: 12,
		}
		normalBufferViewID := uint32(len(doc.BufferViews))
		doc.BufferViews = append(doc.BufferViews, normalBufferView)

		normalAccessor = &gltf.Accessor{
			BufferView:    &normalBufferViewID,
			ComponentType: gltf.ComponentFloat,
			Count:         uint32(len(content.Normals) / 3),
			Type:          gltf.AccessorVec3,
		}
		doc.Accessors = append(doc.Accessors, normalAccessor)
	}

	var texCoordAccessor *gltf.Accessor
	if len(content.TexCoords) > 0 {
		texCoordBuffer := &gltf.Buffer{
			Data:       float32ToByte(content.TexCoords),
			ByteLength: uint32(len(content.TexCoords) * 4),
		}
		texCoordBufferID := uint32(len(doc.Buffers))
		doc.Buffers = append(doc.Buffers, texCoordBuffer)

		texCoordBufferView := &gltf.BufferView{
			Buffer:     texCoordBufferID,
			ByteOffset: 0,
			ByteLength: uint32(len(content.TexCoords) * 4),
			ByteStride: 8,
		}
		texCoordBufferViewID := uint32(len(doc.BufferViews))
		doc.BufferViews = append(doc.BufferViews, texCoordBufferView)

		texCoordAccessor = &gltf.Accessor{
			BufferView:    &texCoordBufferViewID,
			ComponentType: gltf.ComponentFloat,
			Count:         uint32(len(content.TexCoords) / 2),
			Type:          gltf.AccessorVec2,
		}
		doc.Accessors = append(doc.Accessors, texCoordAccessor)
	}

	primitive := &gltf.Primitive{
		Attributes: map[string]uint32{
			"POSITION": 0,
		},
		Mode: gltf.PrimitiveTriangles,
	}

	if normalAccessor != nil {
		normalAccessorID := uint32(len(doc.Accessors) - 1)
		primitive.Attributes["NORMAL"] = normalAccessorID
	}

	if texCoordAccessor != nil {
		texCoordAccessorID := uint32(len(doc.Accessors) - 1)
		primitive.Attributes["TEXCOORD_0"] = texCoordAccessorID
	}

	if len(content.Indices) > 0 {
		indexBuffer := &gltf.Buffer{
			Data:       uint32ToByte(content.Indices),
			ByteLength: uint32(len(content.Indices) * 4),
		}
		indexBufferID := uint32(len(doc.Buffers))
		doc.Buffers = append(doc.Buffers, indexBuffer)

		indexBufferView := &gltf.BufferView{
			Buffer:     indexBufferID,
			ByteOffset: 0,
			ByteLength: uint32(len(content.Indices) * 4),
		}
		indexBufferViewID := uint32(len(doc.BufferViews))
		doc.BufferViews = append(doc.BufferViews, indexBufferView)

		indexAccessor := &gltf.Accessor{
			BufferView:    &indexBufferViewID,
			ComponentType: gltf.ComponentUint,
			Count:         uint32(len(content.Indices)),
			Type:          gltf.AccessorScalar,
		}
		doc.Accessors = append(doc.Accessors, indexAccessor)

		primitive.Indices = gltf.Index(uint32(len(doc.Accessors) - 1))
	}

	mesh := &gltf.Mesh{
		Name:       "Mesh",
		Primitives: []*gltf.Primitive{primitive},
	}
	meshID := uint32(len(doc.Meshes))
	doc.Meshes = append(doc.Meshes, mesh)

	node := &gltf.Node{
		Mesh: &meshID,
	}
	doc.Nodes = append(doc.Nodes, node)

	sceneIndex := uint32(len(doc.Nodes) - 1)
	doc.Scene = &sceneIndex

	baseColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	metallic := float32(0.0)
	roughness := float32(1.0)

	material := &gltf.Material{
		Name: "DefaultMaterial",
		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorFactor: &baseColor,
			MetallicFactor:  &metallic,
			RoughnessFactor: &roughness,
		},
		DoubleSided: false,
	}
	doc.Materials = append(doc.Materials, material)

	doc.ExtensionsUsed = []string{"KHR_materials_unlit"}

	return doc
}

func (g *B3DMGenerator) Write(tile *Tile, path string) error {
	data, err := g.Generate(tile)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type Converter struct {
	opts          *ConverterOptions
	coordTrans    *CoordinateTransformer
	geoidConv     *GeoidConverter
	geomConverter *GeometryConverter
	tilesetGen    *TilesetGenerator
	b3dmGen       *B3DMGenerator
	basePath      string
}

func NewConverter(opts *ConverterOptions) *Converter {
	coordTrans := NewCoordinateTransformer(opts.SourceSRS, opts.TargetSRS)
	coordTrans.SetCenter(opts.CenterLongitude, opts.CenterLatitude, opts.CenterHeight)

	geoidConv := NewGeoidConverter(opts.GeoidModel, opts.GeoidDataPath)

	geomConverter := NewGeometryConverter(opts, coordTrans, geoidConv)

	tilesetGen := NewTilesetGenerator(opts)

	b3dmGen := NewB3DMGenerator(opts)

	return &Converter{
		opts:          opts,
		coordTrans:    coordTrans,
		geoidConv:     geoidConv,
		geomConverter: geomConverter,
		tilesetGen:    tilesetGen,
		b3dmGen:       b3dmGen,
		basePath:      "",
	}
}

func (c *Converter) Convert(inputPath, outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	c.basePath = filepath.Dir(inputPath)

	node, err := LoadOSGB(inputPath)
	if err != nil {
		return err
	}

	tile := c.convertNodeToTile(node, filepath.Base(inputPath))

	c.extendTileBox(tile)

	if c.opts.MaxLOD < 0 {
		c.loadPagedLODs(tile, outputPath)
	} else {
		c.loadPagedLODsLimited(tile, outputPath, 0, c.opts.MaxLOD)
	}

	c.calcGeometricError(tile)

	tileset := c.tilesetGen.Generate(tile)

	if err := c.tilesetGen.Write(tileset, filepath.Join(outputPath, "tileset.json")); err != nil {
		return err
	}

	return c.writeTiles(tile, outputPath)
}

func (c *Converter) loadPagedLODs(tile *Tile, outputPath string) {
	if tile == nil {
		return
	}

	children := c.getPagedLODChildren(tile)
	for _, childPath := range children {
		fullPath := childPath
		if !filepath.IsAbs(fullPath) && c.basePath != "" {
			fullPath = filepath.Join(c.basePath, childPath)
		}

		node, err := LoadOSGB(fullPath)
		if err != nil {
			fmt.Printf("Warning: Failed to load %s: %v\n", fullPath, err)
			continue
		}

		childTile := c.convertNodeToTile(node, filepath.Base(childPath))
		childTile.Path = filepath.Base(childPath) + ".b3dm"

		tile.Children = append(tile.Children, childTile)

		c.loadPagedLODs(childTile, outputPath)
	}
}

func (c *Converter) loadPagedLODsLimited(tile *Tile, outputPath string, currentLevel, maxLevel int) {
	if tile == nil || currentLevel >= maxLevel {
		return
	}

	children := c.getPagedLODChildren(tile)
	for _, childPath := range children {
		fullPath := childPath
		if !filepath.IsAbs(fullPath) && c.basePath != "" {
			fullPath = filepath.Join(c.basePath, childPath)
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Warning: Recovered from panic loading %s: %v\n", fullPath, r)
			}
		}()

		node, err := LoadOSGB(fullPath)
		if err != nil {
			fmt.Printf("Warning: Failed to load %s: %v\n", fullPath, err)
			continue
		}

		childTile := c.convertNodeToTile(node, filepath.Base(childPath))
		childTile.Path = filepath.Base(childPath) + ".b3dm"

		tile.Children = append(tile.Children, childTile)

		c.loadPagedLODsLimited(childTile, outputPath, currentLevel+1, maxLevel)
	}
}

func (c *Converter) getPagedLODChildren(tile *Tile) []string {
	if tile == nil || tile.Node == nil {
		return nil
	}

	var result []string

	processPagedLod := func(plod *model.PagedLod) {
		for _, perRange := range plod.PerRangeDataList {
			if perRange.FileName != "" {
				childPath := strings.ReplaceAll(perRange.FileName, "\\", string(filepath.Separator))
				if !filepath.IsAbs(childPath) && plod.DataBasePath != "" {
					childPath = filepath.Join(plod.DataBasePath, childPath)
				}
				result = append(result, childPath)
			}
		}
	}

	switch n := tile.Node.(type) {
	case *model.PagedLod:
		processPagedLod(n)
	case *model.Lod:
	case *model.Group:
		for _, child := range n.GetChildren() {
			if child == nil {
				continue
			}
			if plod, ok := child.(*model.PagedLod); ok {
				processPagedLod(plod)
			}
		}
	case *model.MatrixTransform:
		for _, child := range n.GetChildren() {
			if child == nil {
				continue
			}
			if plod, ok := child.(*model.PagedLod); ok {
				processPagedLod(plod)
			}
		}
	case *model.PositionAttitudeTransform:
		for _, child := range n.GetChildren() {
			if child == nil {
				continue
			}
			if plod, ok := child.(*model.PagedLod); ok {
				processPagedLod(plod)
			}
		}
	}

	return result
}

func (c *Converter) extendTileBox(tile *Tile) {
	if tile == nil {
		return
	}

	for _, child := range tile.Children {
		c.extendTileBox(child)
		if len(child.BoundingBox) == 6 && len(tile.BoundingBox) == 6 {
			for i := 0; i < 3; i++ {
				if child.BoundingBox[i] < tile.BoundingBox[i] {
					tile.BoundingBox[i] = child.BoundingBox[i]
				}
				if child.BoundingBox[i+3] > tile.BoundingBox[i+3] {
					tile.BoundingBox[i+3] = child.BoundingBox[i+3]
				}
			}
		}
	}
}

func (c *Converter) calcGeometricError(tile *Tile) {
	if tile == nil {
		return
	}

	if len(tile.Children) == 0 {
		if len(tile.BoundingBox) == 6 {
			tile.GeometricError = calculateGeometricError(tile.BoundingBox)
		}
	} else {
		maxChildError := 0.0
		for _, child := range tile.Children {
			c.calcGeometricError(child)
			if child.GeometricError > maxChildError {
				maxChildError = child.GeometricError
			}
		}
		tile.GeometricError = maxChildError * 2.0
	}
}

func (c *Converter) convertNodeToTile(node interface{}, id string) *Tile {
	content := c.geomConverter.Convert(node)

	tile := &Tile{
		ID:       id,
		Path:     id + ".b3dm",
		Children: []*Tile{},
		Node:     node,
	}

	if content != nil {
		boundingBox := content.BoundingBox
		if len(boundingBox) != 6 {
			boundingBox = calculateBoundingBox(content.Vertices)
		}
		tile.BoundingBox = boundingBox
		tile.GeometricError = calculateGeometricError(boundingBox)
		tile.Content = content
	}

	return tile
}

func (c *Converter) writeTiles(tile *Tile, outputPath string) error {
	if tile.Content != nil && tile.Path != "" {
		if err := c.b3dmGen.Write(tile, filepath.Join(outputPath, tile.Path)); err != nil {
			return err
		}
	}

	for _, child := range tile.Children {
		if err := c.writeTiles(child, outputPath); err != nil {
			return err
		}
	}

	return nil
}

func calculateBoundingBox(vertices []float32) [6]float64 {
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

func calculateGeometricError(bbox [6]float64) float64 {
	dx := bbox[3] - bbox[0]
	dy := bbox[4] - bbox[1]
	dz := bbox[5] - bbox[2]
	return (dx + dy + dz) / 6.0 * 2.0
}

type ConvertResult struct {
	JSON        string
	BoundingBox [6]float64
}

func OSGBToGLB(inputPath string, opts *ConverterOptions) ([]byte, error) {
	if opts == nil {
		opts = DefaultConverterOptions()
	}

	coordTrans := NewCoordinateTransformer(opts.SourceSRS, opts.TargetSRS)
	coordTrans.SetCenter(opts.CenterLongitude, opts.CenterLatitude, opts.CenterHeight)

	geoidConv := NewGeoidConverter(opts.GeoidModel, opts.GeoidDataPath)

	geomConverter := NewGeometryConverter(opts, coordTrans, geoidConv)

	node, err := LoadOSGB(inputPath)
	if err != nil {
		return nil, err
	}

	content := geomConverter.Convert(node)
	if content == nil {
		return nil, fmt.Errorf("failed to extract geometry from %s", inputPath)
	}

	b3dmGen := NewB3DMGenerator(opts)

	tile := &Tile{
		Content:     content,
		BoundingBox: content.BoundingBox,
	}

	if len(tile.BoundingBox) != 6 {
		tile.BoundingBox = calculateBoundingBox(content.Vertices)
	}

	glbData, err := b3dmGen.GenerateGLB(tile)
	if err != nil {
		return nil, err
	}

	return glbData, nil
}

func OSGBToGLBFile(inputPath, outputPath string, opts *ConverterOptions) error {
	data, err := OSGBToGLB(inputPath, opts)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, data, 0644)
}

func (c *Converter) OSGBToB3DM(inputPath string) ([]byte, *Tile, error) {
	node, err := LoadOSGB(inputPath)
	if err != nil {
		return nil, nil, err
	}

	tile := c.convertNodeToTile(node, filepath.Base(inputPath))

	data, err := c.b3dmGen.Generate(tile)
	if err != nil {
		return nil, tile, err
	}

	return data, tile, nil
}

func OSGBTo3DTiles(inputPath, outputPath string, opts *ConverterOptions) (*ConvertResult, error) {
	if opts == nil {
		opts = DefaultConverterOptions()
	}

	converter := NewConverter(opts)

	if err := converter.Convert(inputPath, outputPath); err != nil {
		return nil, err
	}

	jsonPath := filepath.Join(outputPath, "tileset.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	tilePath := filepath.Join(outputPath, "tileset.json")
	tileData, err := os.ReadFile(tilePath)
	if err != nil {
		return nil, err
	}

	var tileset TileJSON
	if err := json.Unmarshal(tileData, &tileset); err != nil {
		return nil, err
	}

	result := &ConvertResult{
		JSON: string(data),
	}

	if tileset.Root != nil && len(tileset.Root.BoundVolume.Box) >= 6 {
		copy(result.BoundingBox[:], tileset.Root.BoundVolume.Box[:6])
	}

	return result, nil
}
