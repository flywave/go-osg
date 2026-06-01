package tiles3d

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/go-osg/model"

	"github.com/flywave/gltf"
	"github.com/flywave/gltf/ext/draco"
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
			Version:    "1.0",
			GenBy:      "go-osg-3dtiles",
			GltfUpAxis: "Z",
		},
		GeometricError: g.opts.GeometricError,
	}

	if tile != nil {
		tileset.Root = g.convertTileToJSON(tile, true)
		if tileset.Root != nil && tile.Content != nil && len(tile.Content.Vertices) < 3 && len(tile.Children) > 0 {
			tileset.Root.GeometricError = 0
			for _, c := range tile.Children {
				if c.GeometricError > tileset.Root.GeometricError {
					tileset.Root.GeometricError = c.GeometricError
				}
			}
		}
	}

	return tileset
}

func (g *TilesetGenerator) convertTileToJSON(tile *Tile, isRoot bool) *TileJSONNode {
	bbox := tile.BoundingBox

	node := &TileJSONNode{
		GeometricError: tile.GeometricError,
		BoundVolume: BoundVolumeJSON{
			Box: bbox[:],
		},
	}

	node.Refine = "REPLACE"
	if tile.Content != nil && tile.Path != "" && len(tile.Content.Vertices) >= 3 {
		node.Content = &ContentJSON{
			URI: "./" + filepath.Base(tile.Path),
		}
	}

	if len(tile.Children) > 0 {
		node.Children = make([]*TileJSONNode, len(tile.Children))
		for i, child := range tile.Children {
			node.Children[i] = g.convertTileToJSON(child, false)
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

	if g.opts != nil && g.opts.EnableDraco {
		quantOpts := map[string]interface{}{
			"quantization": map[string]int{
				"position": g.opts.DracoPositionBits,
				"normal":   g.opts.DracoNormalBits,
				"texcoord": g.opts.DracoTexCoordBits,
			},
		}
		if err := draco.EncodeAll(doc, quantOpts); err != nil {
			return nil, fmt.Errorf("draco encode failed: %w", err)
		}
	}

	b3dm := tile3d.NewB3dm()
	b3dm.Header.Version = 1
	b3dm.Model = doc

	view := tile3d.B3dmFeatureTableView{
		BatchLength: tile.Content.BatchLength,
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

	if g.opts != nil && g.opts.EnableDraco {
		quantOpts := map[string]interface{}{
			"quantization": map[string]int{
				"position": g.opts.DracoPositionBits,
				"normal":   g.opts.DracoNormalBits,
				"texcoord": g.opts.DracoTexCoordBits,
			},
		}
		if err := draco.EncodeAll(doc, quantOpts); err != nil {
			return nil, fmt.Errorf("draco encode failed: %w", err)
		}
	}

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

func osgModeToGLTF(mode int32) gltf.PrimitiveMode {
	switch mode {
	case model.GLPOINTS:
		return gltf.PrimitivePoints
	case model.GLLINES:
		return gltf.PrimitiveLines
	case model.GLLINELOOP:
		return gltf.PrimitiveLineLoop
	case model.GLLINESTRIP:
		return gltf.PrimitiveLineStrip
	case model.GLTRIANGLES:
		return gltf.PrimitiveTriangles
	case model.GLTRIANGLESTRIP:
		return gltf.PrimitiveTriangleStrip
	case model.GLTRIANGLEFAN:
		return gltf.PrimitiveTriangleFan
	case model.GLQUADS, model.GLQUADSTRIP:
		return gltf.PrimitiveTriangles
	default:
		return gltf.PrimitiveTriangles
	}
}

func pickIndexComponentType(maxIndex uint32) gltf.ComponentType {
	if maxIndex <= 255 {
		return gltf.ComponentUbyte
	} else if maxIndex <= 65535 {
		return gltf.ComponentUshort
	}
	return gltf.ComponentUint
}

func indicesToBytes(indices []uint32, ct gltf.ComponentType) []byte {
	switch ct {
	case gltf.ComponentUbyte:
		buf := make([]byte, len(indices))
		for i, idx := range indices {
			buf[i] = byte(idx)
		}
		return buf
	case gltf.ComponentUshort:
		buf := make([]byte, len(indices)*2)
		for i, idx := range indices {
			buf[i*2] = byte(idx)
			buf[i*2+1] = byte(idx >> 8)
		}
		return buf
	default:
		return uint32ToByte(indices)
	}
}

func align4(n int) int {
	return (n + 3) & ^3
}

func (g *B3DMGenerator) buildGLTF(content *TileContent) *gltf.Document {
	doc := gltf.NewDocument()

	var bin bytes.Buffer
	type bufferRegion struct {
		offset   int
		length   int
		stride   int
		target   gltf.Target
	}

	writeRegion := func(data []byte, stride int, target gltf.Target) bufferRegion {
		offset := bin.Len()
		bin.Write(data)
		padding := align4(bin.Len()) - bin.Len()
		bin.Write(make([]byte, padding))
		return bufferRegion{offset: offset, length: len(data), stride: stride, target: target}
	}

	posData := float32ToByte(content.Vertices)
	posRegion := writeRegion(posData, 12, gltf.TargetArrayBuffer)

	posBV := uint32(len(doc.BufferViews))
	doc.BufferViews = append(doc.BufferViews, &gltf.BufferView{
		Buffer:     0,
		ByteOffset: uint32(posRegion.offset),
		ByteLength: uint32(posRegion.length),
		ByteStride: uint32(posRegion.stride),
		Target:     posRegion.target,
	})

	posAcc := &gltf.Accessor{
		BufferView:    &posBV,
		ComponentType: gltf.ComponentFloat,
		Count:         uint32(len(content.Vertices) / 3),
		Type:          gltf.AccessorVec3,
		Min: []float32{
			float32(content.BoundingBox[0] - content.BoundingBox[3]),
			float32(content.BoundingBox[1] - content.BoundingBox[7]),
			float32(content.BoundingBox[2] - content.BoundingBox[11]),
		},
		Max: []float32{
			float32(content.BoundingBox[0] + content.BoundingBox[3]),
			float32(content.BoundingBox[1] + content.BoundingBox[7]),
			float32(content.BoundingBox[2] + content.BoundingBox[11]),
		},
	}
	doc.Accessors = append(doc.Accessors, posAcc)

	normalAccessorIdx := -1
	if len(content.Normals) > 0 {
		normData := float32ToByte(content.Normals)
		normRegion := writeRegion(normData, 12, gltf.TargetArrayBuffer)
		normBV := uint32(len(doc.BufferViews))
		doc.BufferViews = append(doc.BufferViews, &gltf.BufferView{
			Buffer:     0,
			ByteOffset: uint32(normRegion.offset),
			ByteLength: uint32(normRegion.length),
			ByteStride: uint32(normRegion.stride),
			Target:     normRegion.target,
		})
		doc.Accessors = append(doc.Accessors, &gltf.Accessor{
			BufferView:    &normBV,
			ComponentType: gltf.ComponentFloat,
			Count:         uint32(len(content.Normals) / 3),
			Type:          gltf.AccessorVec3,
		})
		normalAccessorIdx = len(doc.Accessors) - 1
	}

	texCoordAccessorIdx := -1
	if len(content.TexCoords) > 0 {
		tcData := float32ToByte(content.TexCoords)
		tcRegion := writeRegion(tcData, 8, gltf.TargetArrayBuffer)
		tcBV := uint32(len(doc.BufferViews))
		doc.BufferViews = append(doc.BufferViews, &gltf.BufferView{
			Buffer:     0,
			ByteOffset: uint32(tcRegion.offset),
			ByteLength: uint32(tcRegion.length),
			ByteStride: uint32(tcRegion.stride),
			Target:     tcRegion.target,
		})
		doc.Accessors = append(doc.Accessors, &gltf.Accessor{
			BufferView:    &tcBV,
			ComponentType: gltf.ComponentFloat,
			Count:         uint32(len(content.TexCoords) / 2),
			Type:          gltf.AccessorVec2,
		})
		texCoordAccessorIdx = len(doc.Accessors) - 1
	}

	makePrimitive := func(indices []uint32, mode int32) *gltf.Primitive {
		prim := &gltf.Primitive{
			Attributes: map[string]uint32{
				"POSITION": 0,
			},
			Mode: osgModeToGLTF(mode),
		}
		if normalAccessorIdx >= 0 {
			prim.Attributes["NORMAL"] = uint32(normalAccessorIdx)
		}
		if texCoordAccessorIdx >= 0 {
			prim.Attributes["TEXCOORD_0"] = uint32(texCoordAccessorIdx)
		}
		if len(indices) > 0 {
			maxIdx := uint32(0)
			for _, idx := range indices {
				if idx > maxIdx {
					maxIdx = idx
				}
			}
			ct := pickIndexComponentType(maxIdx)
			idxBytes := indicesToBytes(indices, ct)
			idxRegion := writeRegion(idxBytes, 0, gltf.TargetElementArrayBuffer)
			idxBV := uint32(len(doc.BufferViews))
			doc.BufferViews = append(doc.BufferViews, &gltf.BufferView{
				Buffer:     0,
				ByteOffset: uint32(idxRegion.offset),
				ByteLength: uint32(idxRegion.length),
				Target:     idxRegion.target,
			})
			doc.Accessors = append(doc.Accessors, &gltf.Accessor{
				BufferView:    &idxBV,
				ComponentType: ct,
				Count:         uint32(len(indices)),
				Type:          gltf.AccessorScalar,
			})
			prim.Indices = gltf.Index(uint32(len(doc.Accessors) - 1))
		}
		return prim
	}

	var primitives []*gltf.Primitive
	if len(content.Primitives) > 0 {
		for _, p := range content.Primitives {
			primitives = append(primitives, makePrimitive(p.Indices, p.Mode))
		}
	} else if len(content.Indices) > 0 {
		primitives = append(primitives, makePrimitive(content.Indices, content.Mode))
	} else {
		primitives = append(primitives, makePrimitive(nil, content.Mode))
	}

	imageBVOffsets := make([]uint32, 0)
	if len(content.Textures) > 0 {
		for i, texData := range content.Textures {
			off := bin.Len()
			bin.Write(texData)
			padding := align4(bin.Len()) - bin.Len()
			bin.Write(make([]byte, padding))
			imgBV := uint32(len(doc.BufferViews))
			doc.BufferViews = append(doc.BufferViews, &gltf.BufferView{
				Buffer:     0,
				ByteOffset: uint32(off),
				ByteLength: uint32(len(texData)),
			})
			imageBVOffsets = append(imageBVOffsets, imgBV)

			doc.Images = append(doc.Images, &gltf.Image{
				URI:        "",
				BufferView: &imgBV,
				MimeType:   "image/jpeg",
			})
			doc.Samplers = append(doc.Samplers, &gltf.Sampler{
				MagFilter: gltf.MagLinear,
				MinFilter: gltf.MinLinearMipMapLinear,
				WrapS:     gltf.WrapRepeat,
				WrapT:     gltf.WrapRepeat,
			})
			doc.Textures = append(doc.Textures, &gltf.Texture{
				Name:    fmt.Sprintf("texture_%d", i),
				Source:  gltf.Index(uint32(len(doc.Images) - 1)),
				Sampler: gltf.Index(uint32(len(doc.Samplers) - 1)),
			})
		}

		doc.Materials = append(doc.Materials, &gltf.Material{
			Name: "TexturedMaterial",
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
				BaseColorFactor: &[4]float32{1.0, 1.0, 1.0, 1.0},
				MetallicFactor:  &[]float32{0.0}[0],
				RoughnessFactor: &[]float32{1.0}[0],
				BaseColorTexture: &gltf.TextureInfo{
					Index: uint32(len(doc.Textures) - 1),
				},
			},
			DoubleSided: false,
		})
		matID := uint32(len(doc.Materials) - 1)
		for _, prim := range primitives {
			prim.Material = &matID
		}
	}

	doc.Buffers = append(doc.Buffers, &gltf.Buffer{
		Data:       bin.Bytes(),
		ByteLength: uint32(bin.Len()),
	})

	doc.Meshes = append(doc.Meshes, &gltf.Mesh{
		Name:       "Mesh",
		Primitives: primitives,
	})

	doc.Nodes = append(doc.Nodes, &gltf.Node{
		Mesh: gltf.Index(0),
	})

	doc.Scene = gltf.Index(0)
	if len(doc.Scenes) > 0 {
		doc.Scenes[0].Nodes = []uint32{0}
	}

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
	visitedTiles  map[string]bool
	visitedDirs   map[string]bool
}

func NewConverter(opts *ConverterOptions) *Converter {
	geoidConv := NewGeoidConverter(opts.GeoidModel, opts.GeoidDataPath)

	coordTrans := NewCoordinateTransformer(opts.SourceSRS, opts.TargetSRS)
	coordTrans.SetGeoidConverter(geoidConv)

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
		visitedTiles:  make(map[string]bool),
		visitedDirs:   make(map[string]bool),
	}
}

func (c *Converter) Convert(inputPath, outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	c.basePath = filepath.Dir(inputPath)

	c.detectMetadata(inputPath)

	node, err := LoadOSGB(inputPath)
	if err != nil {
		return err
	}

	tile := c.convertNodeToTile(node, filepath.Base(inputPath))
	tile.SourcePath = inputPath
	fmt.Printf("DEBUG Convert: root tile ID=%s, SourcePath=%s, bbox center = (%f, %f, %f)\n",
		tile.ID, tile.SourcePath, tile.BoundingBox[0], tile.BoundingBox[1], tile.BoundingBox[2])

	if c.opts.MaxLOD < 0 {
		c.loadPagedLODs(tile, outputPath)
	} else {
		c.loadPagedLODsLimited(tile, outputPath, 0, c.opts.MaxLOD)
	}

	c.extendTileBox(tile)
	fmt.Printf("DEBUG Convert: after extendTileBox, root tile bbox center = (%f, %f, %f)\n",
		tile.BoundingBox[0], tile.BoundingBox[1], tile.BoundingBox[2])

	c.calcGeometricError(tile)

	tileset := c.tilesetGen.Generate(tile)

	if c.coordTrans.HasGeoReference() {
		center := c.coordTrans.GetCenter()
		latRad := center[1] * math.Pi / 180.0
		lonRad := center[0] * math.Pi / 180.0
		height := center[2]

		ecef := c.coordTrans.ToECEFFromLatLon(latRad, lonRad, height)

		sinLat := math.Sin(latRad)
		cosLat := math.Cos(latRad)
		sinLon := math.Sin(lonRad)
		cosLon := math.Cos(lonRad)

		transform := []float64{
			-sinLon, cosLon, 0, 0,
			-sinLat * cosLon, -sinLat * sinLon, cosLat, 0,
			cosLat * cosLon, cosLat * sinLon, sinLat, 0,
			ecef[0], ecef[1], ecef[2], 1,
		}

		if tileset.Root != nil {
			tileset.Root.Transform = transform
		}
		fmt.Printf("DEBUG Convert: added transform matrix, ECEF=(%f, %f, %f)\n", ecef[0], ecef[1], ecef[2])
	}

	if err := c.tilesetGen.Write(tileset, filepath.Join(outputPath, "tileset.json")); err != nil {
		return err
	}

	return c.writeTiles(tile, outputPath)
}

func (c *Converter) detectMetadata(inputPath string) {
	searchPaths := []string{
		filepath.Dir(inputPath),
		filepath.Join(filepath.Dir(inputPath), "OSGB"),
		filepath.Join(filepath.Dir(inputPath), ".."),
		filepath.Join(filepath.Dir(inputPath), "..", "OSGB"),
		filepath.Join(filepath.Dir(inputPath), "..", ".."),
		filepath.Join(filepath.Dir(inputPath), "..", "..", "OSGB"),
		filepath.Join(filepath.Dir(inputPath), "..", "..", ".."),
		filepath.Join(filepath.Dir(inputPath), "..", "..", "..", "OSGB"),
	}

	var metadataPath string
	for _, p := range searchPaths {
		if p != "" {
			path, err := FindMetadataFile(p)
			if err == nil {
				metadataPath = path
				break
			}
		}
	}

	if metadataPath == "" {
		fmt.Printf("DEBUG detectMetadata: no metadata.xml found, using provided options\n")
		if c.opts.CenterLongitude != 0 || c.opts.CenterLatitude != 0 {
			c.coordTrans.SetCenter(c.opts.CenterLongitude, c.opts.CenterLatitude, c.opts.CenterHeight)
		}
		return
	}

	fmt.Printf("DEBUG detectMetadata: found metadata.xml at %s\n", metadataPath)

	metadata, err := ParseMetadataXML(metadataPath)
	if err != nil {
		fmt.Printf("DEBUG detectMetadata: failed to parse metadata.xml: %v\n", err)
		if c.opts.CenterLongitude != 0 || c.opts.CenterLatitude != 0 {
			c.coordTrans.SetCenter(c.opts.CenterLongitude, c.opts.CenterLatitude, c.opts.CenterHeight)
		}
		return
	}

	fmt.Printf("DEBUG detectMetadata: SRS=%q, SRSOrigin=%q\n", metadata.SRS, metadata.SRSOrigin)

	if c.opts.SourceSRS != "" {
		fmt.Printf("DEBUG detectMetadata: using user-provided SourceSRS=%q\n", c.opts.SourceSRS)
		metadata.SRS = c.opts.SourceSRS
	}

	err = c.coordTrans.SetGeoReferenceFromMetadata(metadata, c.opts.GeoidDataPath)
	if err != nil {
		fmt.Printf("DEBUG detectMetadata: failed to set geo reference: %v\n", err)
		if c.opts.CenterLongitude != 0 || c.opts.CenterLatitude != 0 {
			c.coordTrans.SetCenter(c.opts.CenterLongitude, c.opts.CenterLatitude, c.opts.CenterHeight)
		}
		return
	}

	fmt.Printf("DEBUG detectMetadata: geo reference set successfully\n")
	fmt.Printf("DEBUG detectMetadata: center=(%f, %f, %f)\n",
		c.coordTrans.GetCenter()[0], c.coordTrans.GetCenter()[1], c.coordTrans.GetCenter()[2])
	fmt.Printf("DEBUG detectMetadata: originOffset=(%f, %f, %f)\n",
		c.coordTrans.GetOriginOffset()[0], c.coordTrans.GetOriginOffset()[1], c.coordTrans.GetOriginOffset()[2])
}

var loadCount = 0

func (c *Converter) loadPagedLODs(tile *Tile, outputPath string) {
	if tile == nil {
		return
	}

	children := c.getPagedLODChildren(tile)
	fmt.Printf("DEBUG loadPagedLODs: tile=%s, children=%d\n", tile.ID, len(children))

	for _, childPath := range children {
		fullPath := childPath
		if !filepath.IsAbs(fullPath) {
			if tile.SourcePath != "" {
				tileDir := filepath.Dir(tile.SourcePath)
				fullPath = filepath.Join(tileDir, childPath)
			} else if c.basePath != "" {
				fullPath = filepath.Join(c.basePath, childPath)
			}
		}

		fmt.Printf("DEBUG loadPagedLODs: Loading %s -> %s\n", childPath, fullPath)

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

		loadCount++
		if loadCount%10 == 0 {
			fmt.Printf("Progress: Loaded %d files...\n", loadCount)
		}

		childTile := c.convertNodeToTile(node, filepath.Base(childPath))
		childTile.Path = filepath.Base(childPath) + ".b3dm"
		childTile.SourcePath = fullPath

		tile.Children = append(tile.Children, childTile)

		c.loadPagedLODs(childTile, outputPath)
	}
}

func (c *Converter) loadPagedLODsLimited(tile *Tile, outputPath string, currentLevel, maxLevel int) {
	if tile == nil || currentLevel >= maxLevel {
		return
	}

	fmt.Printf(">>> loadPagedLODsLimited: START tile=%s, level=%d/%d\n", tile.ID, currentLevel, maxLevel)

	children := c.getPagedLODChildren(tile)
	fmt.Printf(">>> loadPagedLODsLimited: AFTER getPagedLODChildren tile=%s, found %d children\n",
		tile.ID, len(children))

	for i, childPath := range children {
		if c.visitedTiles[childPath] {
			fmt.Printf("DEBUG: Skipping visited tile: %s\n", childPath)
			continue
		}

		fullPath := childPath
		if !filepath.IsAbs(fullPath) {
			if tile.SourcePath != "" {
				tileDir := filepath.Dir(tile.SourcePath)
				fullPath = filepath.Join(tileDir, childPath)
			} else if c.basePath != "" {
				fullPath = filepath.Join(c.basePath, childPath)
			}
		}
		fmt.Printf(">>> [%d/%d] Loading child: %s -> %s\n", i+1, len(children), childPath, fullPath)

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

		fmt.Printf(">>> Loaded successfully: %s\n", childPath)

		c.visitedTiles[childPath] = true

		childTile := c.convertNodeToTile(node, filepath.Base(childPath))
		childTile.Path = filepath.Base(childPath) + ".b3dm"
		childTile.SourcePath = fullPath

		tile.Children = append(tile.Children, childTile)

		c.loadPagedLODsLimited(childTile, outputPath, currentLevel+1, maxLevel)
	}
}

func (c *Converter) getPagedLODChildren(tile *Tile) []string {
	if tile == nil {
		return nil
	}
	if tile.Node == nil {
		fmt.Printf("DEBUG getPagedLODChildren: tile=%s has nil Node!\n", tile.ID)
		return nil
	}

	var result []string

	// First, get children from PagedLOD
	pagedLODChildren := c.getPagedLODChildrenFromNode(tile)

	// Then, scan directory for sibling tiles at the same LOD level
	// Only scan for L21-L25 levels (middle LODs)
	lodLevel := extractLODLevel(tile.ID)
	isMidLevel := false
	if lodLevel != "" {
		// Parse LOD level number
		if len(lodLevel) >= 2 && lodLevel[0] == 'L' {
			levelNum := 0
			fmt.Sscanf(lodLevel[1:], "%d", &levelNum)
			// Only scan siblings for mid-level LODs (L21-L25)
			if levelNum >= 21 && levelNum <= 25 {
				isMidLevel = true
			}
		}
	}

	var dirChildren []string
	if isMidLevel {
		dirChildren = c.scanDirectoryForSiblingsBylod(tile)
	}

	// Merge both lists (avoid duplicates)
	childMap := make(map[string]bool)
	for _, child := range pagedLODChildren {
		if child != "" {
			childMap[child] = true
		}
	}
	for _, child := range dirChildren {
		if child != "" && !childMap[child] {
			childMap[child] = true
		}
	}

	for child := range childMap {
		result = append(result, child)
	}

	fmt.Printf("DEBUG getPagedLODChildren: tile=%s, pagedLOD=%d, dir=%d, total=%d\n",
		tile.ID, len(pagedLODChildren), len(dirChildren), len(result))

	return result
}

func (c *Converter) scanDirectoryForSiblings(tile *Tile) []string {
	if tile == nil || c.basePath == "" {
		return nil
	}

	// Get the directory of the current tile (use basePath directly since that's where the tile was loaded from)
	tileDir := c.basePath
	fmt.Printf("DEBUG scanDirectoryForSiblings: tile=%s, tileDir=%q, basePath=%q\n",
		tile.ID, tileDir, c.basePath)
	if tileDir == "" {
		return nil
	}

	// Resolve full path
	fullTileDir := tileDir
	if !filepath.IsAbs(fullTileDir) {
		fullTileDir = filepath.Join(c.basePath, tileDir)
	}

	// Skip if already scanned this directory
	if c.visitedDirs[fullTileDir] {
		fmt.Printf("DEBUG scanDirectoryForSiblings: dir already scanned: %s\n", fullTileDir)
		return nil
	}
	c.visitedDirs[fullTileDir] = true

	// Check if directory exists
	info, err := os.Stat(fullTileDir)
	if err != nil || !info.IsDir() {
		return nil
	}

	// Read all osgb files in the directory
	entries, err := os.ReadDir(fullTileDir)
	if err != nil {
		return nil
	}

	// Extract current tile name for exclusion
	currentTileName := filepath.Base(tile.Path)

	var siblings []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".osgb") {
			continue
		}
		// Skip self
		if name == currentTileName {
			continue
		}

		// Calculate relative path from basePath
		relPath, err := filepath.Rel(c.basePath, filepath.Join(fullTileDir, name))
		if err != nil {
			continue
		}
		relPath = filepath.ToSlash(relPath)
		siblings = append(siblings, relPath)
	}

	if len(siblings) > 0 {
		fmt.Printf("DEBUG scanDirectoryForSiblings: tile=%s, dir=%s, found %d siblings\n",
			tile.ID, fullTileDir, len(siblings))
	}

	return siblings
}

func (c *Converter) scanDirectoryForSiblingsBylod(tile *Tile) []string {
	if tile == nil || c.basePath == "" {
		return nil
	}

	// Extract LOD level from tile ID (e.g., "L21" from "Tile_+002_+000_L21_00000.osgb")
	lodLevel := extractLODLevel(tile.ID)
	if lodLevel == "" {
		return nil
	}

	// Use basePath as the directory to scan
	tileDir := c.basePath
	if tileDir == "" {
		return nil
	}

	// Skip if already scanned this LOD level (only scan once per LOD level)
	scanKey := tileDir + ":" + lodLevel
	if c.visitedDirs[scanKey] {
		fmt.Printf("DEBUG scanDirectoryForSiblingsBylod: already scanned %s\n", scanKey)
		return nil
	}
	c.visitedDirs[scanKey] = true

	// Limit: don't scan siblings for very detailed levels (L22+) to avoid explosion
	if lodLevel == "L22" || lodLevel == "L23" || lodLevel == "L24" {
		fmt.Printf("DEBUG scanDirectoryForSiblingsBylod: skipping high LOD level %s to avoid explosion\n", lodLevel)
		return nil
	}

	// Read all osgb files in the directory
	entries, err := os.ReadDir(tileDir)
	if err != nil {
		return nil
	}

	// Extract current tile name for exclusion (strip .b3dm extension if present)
	currentTileName := filepath.Base(tile.Path)
	currentTileName = strings.TrimSuffix(currentTileName, ".b3dm")

	var siblings []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".osgb") {
			continue
		}
		// Skip self
		if name == currentTileName {
			continue
		}
		// Only include tiles with the same LOD level
		if !strings.Contains(name, "_"+lodLevel+"_") {
			continue
		}

		// Calculate relative path from basePath
		relPath, err := filepath.Rel(c.basePath, filepath.Join(tileDir, name))
		if err != nil {
			continue
		}
		relPath = filepath.ToSlash(relPath)
		siblings = append(siblings, relPath)
	}

	if len(siblings) > 0 {
		fmt.Printf("DEBUG scanDirectoryForSiblingsBylod: tile=%s, lod=%s, found %d siblings\n",
			tile.ID, lodLevel, len(siblings))
	}

	return siblings
}

func extractLODLevel(tileID string) string {
	// Extract LOD level like "L21" from "Tile_+002_+000_L21_00000.osgb"
	// Pattern: Tile_+XXX_+YYY_L<level>_<index>
	parts := strings.Split(tileID, "_")
	for i, part := range parts {
		if len(part) >= 2 && part[0] == 'L' {
			// Check if following parts are numbers
			if i+1 < len(parts) {
				return part
			}
		}
	}
	return ""
}

func (c *Converter) getPagedLODChildrenFromNode(tile *Tile) []string {
	if tile == nil || tile.Node == nil {
		return nil
	}

	var result []string

	fmt.Printf("DEBUG getPagedLODChildren: tile=%s, Node type=%T\n", tile.ID, tile.Node)

	processPagedLod := func(plod *model.PagedLod) {
		fmt.Printf("DEBUG getPagedLODChildren: found PagedLod with %d ranges\n", len(plod.PerRangeDataList))
		for i, perRange := range plod.PerRangeDataList {
			fmt.Printf("DEBUG: range[%d] filename=%q, dbpath=%q\n", i, perRange.FileName, plod.DataBasePath)
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
	}

	if len(tile.Children) == 0 {
		return
	}

	fmt.Printf("DEBUG extendTileBox: tile=%s, children=%d, tile.BoundingBox=%v\n",
		tile.ID, len(tile.Children), tile.BoundingBox[:3])
	if len(tile.Children) > 0 {
		fmt.Printf("DEBUG extendTileBox: first child bbox center = (%f, %f, %f)\n",
			tile.Children[0].BoundingBox[0], tile.Children[0].BoundingBox[1], tile.Children[0].BoundingBox[2])
	}

	minX := math.MaxFloat64
	maxX := -math.MaxFloat64
	minY := math.MaxFloat64
	maxY := -math.MaxFloat64
	minZ := math.MaxFloat64
	maxZ := -math.MaxFloat64

	hasValidChild := false

	for _, child := range tile.Children {
		if len(child.BoundingBox) != 12 {
			continue
		}

		hasValidChild = true

		cminX := child.BoundingBox[0] - child.BoundingBox[3]
		cmaxX := child.BoundingBox[0] + child.BoundingBox[3]
		cminY := child.BoundingBox[1] - child.BoundingBox[7]
		cmaxY := child.BoundingBox[1] + child.BoundingBox[7]
		cminZ := child.BoundingBox[2] - child.BoundingBox[11]
		cmaxZ := child.BoundingBox[2] + child.BoundingBox[11]

		if cminX < minX {
			minX = cminX
		}
		if cmaxX > maxX {
			maxX = cmaxX
		}
		if cminY < minY {
			minY = cminY
		}
		if cmaxY > maxY {
			maxY = cmaxY
		}
		if cminZ < minZ {
			minZ = cminZ
		}
		if cmaxZ > maxZ {
			maxZ = cmaxZ
		}
	}

	if !hasValidChild {
		return
	}

	tile.BoundingBox[0] = (maxX + minX) / 2
	tile.BoundingBox[1] = (maxY + minY) / 2
	tile.BoundingBox[2] = (maxZ + minZ) / 2
	xHalf := (maxX - minX) / 2
	yHalf := (maxY - minY) / 2
	zHalf := (maxZ - minZ) / 2
	// extend bbox by 10% margin per C++ reference (TileBox::extend(0.2))
	if xHalf > 0 { xHalf *= 1.1 }
	if yHalf > 0 { yHalf *= 1.1 }
	if zHalf > 0 { zHalf *= 1.1 }
	tile.BoundingBox[3] = xHalf
	tile.BoundingBox[7] = yHalf
	tile.BoundingBox[11] = zHalf
	fmt.Printf("DEBUG extendTileBox: tile=%s computed bbox center = (%f, %f, %f)\n",
		tile.ID, tile.BoundingBox[0], tile.BoundingBox[1], tile.BoundingBox[2])
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

	if content != nil {
		content.Vertices, content.Normals, content.TexCoords, content.Indices =
			c.geomConverter.optimizeMesh(content.Vertices, content.Normals, content.TexCoords, content.Indices)
	}

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
	if tile.Content != nil && tile.Path != "" && len(tile.Content.Vertices) >= 3 {
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

func calculateBoundingBox(vertices []float32) [12]float64 {
	if len(vertices) == 0 {
		return [12]float64{0, 0, 0, 0.01, 0, 0, 0, 0.01, 0, 0, 0, 0.01}
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

	centerX := (maxX + minX) / 2
	centerY := (maxY + minY) / 2
	centerZ := (maxZ + minZ) / 2

	xHalf := (maxX - minX) / 2
	yHalf := (maxY - minY) / 2
	zHalf := (maxZ - minZ) / 2

	if xHalf < 0.01 {
		xHalf = 0.01
	}
	if yHalf < 0.01 {
		yHalf = 0.01
	}
	if zHalf < 0.01 {
		zHalf = 0.01
	}

	return [12]float64{
		centerX, centerY, centerZ,
		xHalf, 0, 0,
		0, yHalf, 0,
		0, 0, zHalf,
	}
}

func calculateGeometricError(bbox [12]float64) float64 {
	xHalf := bbox[3]
	yHalf := bbox[7]
	zHalf := bbox[11]
	dx := xHalf * 2
	dy := yHalf * 2
	dz := zHalf * 2
	maxDim := dx
	if dy > maxDim {
		maxDim = dy
	}
	if dz > maxDim {
		maxDim = dz
	}
	return maxDim / 2.0
}

type ConvertResult struct {
	JSON        string
	BoundingBox [12]float64
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
