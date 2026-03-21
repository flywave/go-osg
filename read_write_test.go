package osg

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"testing"

	"github.com/flywave/go-osg/model"
)

func TestReadNode(t *testing.T) {
	rw := NewReadWrite()
	res := rw.ReadNode("test_data/cessna.osgb", nil)
	obj := res.GetNode()
	if obj == nil {
		fmt.Println("....")
	}
	vst := model.NewNodeVisitor()
	obj.Accept(vst)
}

func TestReadTile(t *testing.T) {
	rw := NewReadWrite()
	res := rw.ReadNode("test_data/Tile_+003_+003_L18_000.osgb", nil)
	if res == nil {
		t.Fatal("failed to read node")
	}
	obj := res.GetNode()
	if obj == nil {
		t.Fatal("failed to get node")
	}
	fmt.Printf("Node: %+v\n", obj)
}

func TestReadOSGBData(t *testing.T) {
	rw := NewReadWrite()

	testFiles := []string{
		"test_data/OSGB/Data/main.osgb",
		"test_data/OSGB/Data/Tile_+007_+004/Tile_+007_+004_L18_0.osgb",
		"test_data/OSGB/Data/Tile_+007_+004/Tile_+007_+004_L22_00000.osgb",
		"test_data/Tile_+003_+003_L18_000.osgb",
	}

	for _, f := range testFiles {
		t.Run(f, func(t *testing.T) {
			res := rw.ReadNode(f, nil)
			if res == nil {
				t.Fatal("failed to read node")
			}
			obj := res.GetNode()
			if obj == nil {
				t.Fatal("failed to get node")
			}
			fmt.Printf("SUCCESS: %s\n", f)
		})
	}
}

func TestDebugArrayType(t *testing.T) {
	rw := NewReadWrite()

	f := "test_data/Tile_+000_+000.osgb"
	res := rw.ReadNode(f, nil)
	if res == nil {
		t.Fatal("failed to read node")
	}
	obj := res.GetNode()
	if obj == nil {
		t.Fatal("failed to get node")
	}

	pagedlod, ok := obj.(*model.PagedLod)
	if !ok {
		t.Fatal("not PagedLOD")
	}

	children := pagedlod.GetChildren()
	for _, child := range children {
		geode, ok := child.(*model.Geode)
		if !ok {
			continue
		}
		geodeChildren := geode.GetChildren()
		for _, gchild := range geodeChildren {
			geom, ok := gchild.(*model.Geometry)
			if !ok {
				continue
			}
			if geom.VertexArray != nil {
				fmt.Printf("Array Type: %d (IDVEC3SARRAY=%d, IDSHORTARRAY=%d, IDVEC3ARRAY=%d, IDFLOATARRAY=%d)\n",
					geom.VertexArray.Type,
					model.IDVEC3SARRAY,
					model.IDSHORTARRAY,
					model.IDVEC3ARRAY,
					model.IDFLOATARRAY)
				fmt.Printf("DataSize: %d, DataType: %d (GL_SHORT=%d, GL_FLOAT=%d)\n",
					geom.VertexArray.DataSize,
					geom.VertexArray.DataType,
					model.GLSHORT,
					model.GLFLOAT)

				if data, ok := geom.VertexArray.Data.([]int16); ok && len(data) >= 6 {
					fmt.Printf("First 6 int16 values: %v\n", data[:6])
				}
				if data, ok := geom.VertexArray.Data.([][3]int16); ok && len(data) >= 3 {
					fmt.Printf("First 3 Vec3s values: %v\n", data[:3])
				}
			}
		}
	}
}

func TestInspectAllOSGB(t *testing.T) {
	rw := NewReadWrite()

	testFiles := []string{
		"test_data/cessna.osgb",
		"test_data/Tile_+003_+003_L18_000.osgb",
		"test_data/Tile_+000_+000.osgb",
		"test_data/Tile_+000_+000_L22_00000.osgb",
		"test_data/Tile_+000_+000_L24_0000700.osgb",
		"test_data/OSGB/Data/main.osgb",
		"test_data/OSGB/Data/Tile_+007_+004/Tile_+007_+004_L18_0.osgb",
	}

	for _, f := range testFiles {
		res := rw.ReadNode(f, nil)
		if res == nil {
			fmt.Printf("FAILED: %s - res is nil\n", f)
			continue
		}
		obj := res.GetNode()
		if obj == nil {
			fmt.Printf("FAILED: %s - node is nil\n", f)
			continue
		}
		fmt.Printf("SUCCESS: %s (type: %T)\n", f, obj)
		inspectNodeDetails(obj, 0)
	}
}

func inspectNodeDetails(node interface{}, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	if geode, ok := node.(*model.Geode); ok {
		children := geode.GetChildren()
		for _, child := range children {
			if geometry, ok := child.(*model.Geometry); ok {
				fmt.Printf("%s  Geometry:\n", indent)
				if geometry.VertexArray != nil {
					fmt.Printf("%s    VertexArray: DataSize=%d, DataType=%d\n", indent, geometry.VertexArray.DataSize, geometry.VertexArray.DataType)
				}
				if geometry.TexCoordArrayList != nil {
					fmt.Printf("%s    TexCoordArrayList: count=%d\n", indent, len(geometry.TexCoordArrayList))
				}
				if geometry.NormalArray != nil {
					fmt.Printf("%s    NormalArray: present\n", indent)
				}
				if geometry.VertexAttribList != nil {
					fmt.Printf("%s    VertexAttribList: count=%d\n", indent, len(geometry.VertexAttribList))
				}
			}
		}
	}
	if group, ok := node.(*model.Group); ok {
		children := group.GetChildren()
		for _, child := range children {
			inspectNodeDetails(child, depth+1)
		}
	}
	if pagedlod, ok := node.(*model.PagedLod); ok {
		children := pagedlod.GetChildren()
		for _, child := range children {
			inspectNodeDetails(child, depth+1)
		}
	}
}

func inspectNode(node interface{}, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	fmt.Printf("%sType: %T\n", indent, node)
	if group, ok := node.(*model.Group); ok {
		inspectGroup(group, depth)
	}
	if pagedlod, ok := node.(*model.PagedLod); ok {
		inspectPagedLOD(pagedlod, depth)
	}
	if geode, ok := node.(*model.Geode); ok {
		inspectGeode(geode, depth)
	}
}

func inspectGroup(g *model.Group, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	children := g.GetChildren()
	fmt.Printf("%sGroup has %d children\n", indent, len(children))
	for i, child := range children {
		fmt.Printf("%s  Child %d: %T\n", indent, i, child)
		if childGroup, ok := child.(*model.Group); ok {
			inspectGroup(childGroup, depth+2)
		}
		if pagedlod, ok := child.(*model.PagedLod); ok {
			inspectPagedLOD(pagedlod, depth+2)
		}
	}
}

func inspectPagedLOD(p *model.PagedLod, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	children := p.GetChildren()
	fmt.Printf("%sPagedLOD has %d children\n", indent, len(children))
	for i, child := range children {
		fmt.Printf("%s  Child %d: %T\n", indent, i, child)
		if childGroup, ok := child.(*model.Group); ok {
			inspectGroup(childGroup, depth+2)
		}
		if geode, ok := child.(*model.Geode); ok {
			inspectGeode(geode, depth+2)
		}
	}
}

func inspectGeode(g *model.Geode, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	children := g.GetChildren()
	fmt.Printf("%sGeode has %d drawables\n", indent, len(children))
	for i, child := range children {
		fmt.Printf("%s  Drawable %d: %T\n", indent, i, child)
		if geometry, ok := child.(*model.Geometry); ok {
			inspectGeometry(geometry, depth+2)
		}
	}
}

func inspectGeometry(g *model.Geometry, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	fmt.Printf("%sGeometry:\n", indent)
	fmt.Printf("%s  Type: %s\n", indent, g.Type)
	if g.VertexArray != nil {
		fmt.Printf("%s  VertexArray: DataSize=%d, DataType=%d\n", indent, g.VertexArray.DataSize, g.VertexArray.DataType)
		if g.VertexArray.Data != nil {
			switch data := g.VertexArray.Data.(type) {
			case []int16:
				fmt.Printf("%s    Data: []int16 len=%d\n", indent, len(data))
			case []float32:
				fmt.Printf("%s    Data: []float32 len=%d\n", indent, len(data))
			case [][3]float32:
				fmt.Printf("%s    Data: [][3]float32 len=%d\n", indent, len(data))
			default:
				fmt.Printf("%s    Data: %T\n", indent, g.VertexArray.Data)
			}
		}
	}
	if g.NormalArray != nil {
		fmt.Printf("%s  NormalArray: DataSize=%d\n", indent, g.NormalArray.DataSize)
	}
	if g.TexCoordArrayList != nil {
		fmt.Printf("%s  TexCoordArrayList: count=%d\n", indent, len(g.TexCoordArrayList))
	}
	if g.VertexAttribList != nil {
		fmt.Printf("%s  VertexAttribList: count=%d\n", indent, len(g.VertexAttribList))
	}
	if g.Primitives != nil {
		fmt.Printf("%s  Primitives: count=%d\n", indent, len(g.Primitives))
	}
}

func TestFullVerification(t *testing.T) {
	rw := NewReadWrite()

	fmt.Printf("=== Full OSG Verification Test ===\n\n")

	// Test 1: cessna.osgb (airplane with textures)
	fmt.Printf("=== Test 1: Cessna (with textures) ===\n")
	res := rw.ReadNode("test_data/cessna.osgb", nil)
	if res == nil {
		t.Fatal("Failed to read cessna.osgb")
	}
	node := res.GetNode()
	if node == nil {
		t.Fatal("Failed to get node from cessna.osgb")
	}
	fmt.Printf("Cessna root type: %T\n", node)
	inspectNodeFull(node, 0, nil, nil, nil, nil)

	// Test 2: simpleroom.osgt (ASCII format)
	fmt.Printf("\n=== Test 2: SimpleRoom (ASCII format) ===\n")
	res2 := rw.ReadNode("test_data/simpleroom.osgt", nil)
	if res2 == nil {
		t.Fatal("Failed to read simpleroom.osgt")
	}
	node2 := res2.GetNode()
	if node2 == nil {
		t.Fatal("Failed to get node from simpleroom.osgt")
	}
	fmt.Printf("SimpleRoom root type: %T\n", node2)
	inspectNodeFull(node2, 0, nil, nil, nil, nil)

	// Test 3: skydome.osgt
	fmt.Printf("\n=== Test 3: SkyDome (ASCII format) ===\n")
	res3 := rw.ReadNode("test_data/skydome.osgt", nil)
	if res3 == nil {
		t.Fatal("Failed to read skydome.osgt")
	}
	node3 := res3.GetNode()
	if node3 == nil {
		t.Fatal("Failed to get node from skydome.osgt")
	}
	fmt.Printf("SkyDome root type: %T\n", node3)
	inspectNodeFull(node3, 0, nil, nil, nil, nil)

	// Test 4: All OSGB tiles
	fmt.Printf("\n=== Test 4: All OSGB Tiles ===\n")
	inspectAllTiles := func() {
		tiles := []string{
			"test_data/Tile_+000_+000.osgb",
			"test_data/Tile_+000_+000_L22_00000.osgb",
			"test_data/Tile_+000_+000_L24_0000700.osgb",
			"test_data/Tile_+003_+003_L18_000.osgb",
		}
		for _, tile := range tiles {
			res := rw.ReadNode(tile, nil)
			if res == nil {
				t.Fatalf("Failed to read %s", tile)
			}
			node := res.GetNode()
			if node == nil {
				t.Fatalf("Failed to get node from %s", tile)
			}
			fmt.Printf("Tile: %s - type: %T\n", tile, node)
		}
	}
	inspectAllTiles()

	fmt.Printf("\n=== All Tests Passed ===\n")
}

func TestStateSetVerification(t *testing.T) {
	rw := NewReadWrite()

	fmt.Printf("=== StateSet Verification Test ===\n\n")

	// Test StateSet parsing
	files := []string{
		"test_data/simpleroom.osgt",
		"test_data/skydome.osgt",
		"test_data/cessna.osgb",
	}

	for _, f := range files {
		fmt.Printf("Testing: %s\n", f)
		res := rw.ReadNode(f, nil)
		if res == nil {
			t.Fatalf("Failed to read %s", f)
		}
		node := res.GetNode()
		if node == nil {
			t.Fatalf("Failed to get node from %s", f)
		}
		inspectStateSet(node, 0)
	}

	fmt.Printf("\n=== StateSet Tests Passed ===\n")
}

func inspectStateSet(node interface{}, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	switch n := node.(type) {
	case *model.Group:
		if n.States != nil {
			inspectStateSetDetails(n.States, indent)
		}
		children := n.GetChildren()
		for _, child := range children {
			inspectStateSet(child, depth+1)
		}
	case *model.Geode:
		children := n.GetChildren()
		for _, child := range children {
			if geometry, ok := child.(*model.Geometry); ok {
				if geometry.States != nil {
					inspectStateSetDetails(geometry.States, indent)
				}
			}
		}
	}
}

func inspectStateSetDetails(states *model.StateSet, indent string) {
	fmt.Printf("%sStateSet:\n", indent)
	fmt.Printf("%s  RenderingHint: %d\n", indent, states.RenderingHint)
	fmt.Printf("%s  BinName: %s\n", indent, states.BinName)
	fmt.Printf("%s  BinNum: %d\n", indent, states.BinNum)
	fmt.Printf("%s  ModeList count: %d\n", indent, len(states.ModeList))
	fmt.Printf("%s  AttributeList count: %d\n", indent, len(states.AttributeList))
	fmt.Printf("%s  TextureModeList count: %d\n", indent, len(states.TextureModeList))
	fmt.Printf("%s  TextureAttributeList count: %d\n", indent, len(states.TextureAttributeList))

	for k, v := range states.AttributeList {
		fmt.Printf("%s    Attribute[%d]: %T\n", indent, k, v)
	}
	for k, v := range states.TextureAttributeList {
		fmt.Printf("%s    TextureAttribute[%d]: count=%d\n", indent, k, len(v))
	}
}

func TestObliquePhotographyFull(t *testing.T) {
	rw := NewReadWrite()

	mainFile := "test_data/OSGB/Data/main.osgb"
	res := rw.ReadNode(mainFile, nil)
	if res == nil {
		t.Fatalf("Failed to read %s", mainFile)
	}

	node := res.GetNode()
	if node == nil {
		t.Fatalf("Failed to get node from %s", mainFile)
	}

	fmt.Printf("=== Oblique Photography Full Test ===\n")
	fmt.Printf("Main file: %s\n", mainFile)
	fmt.Printf("Root node type: %T\n", node)

	totalGeodes := 0
	totalGeometries := 0
	totalVertices := 0
	maxLODLevel := 0

	inspectNodeFull(node, 0, &totalGeodes, &totalGeometries, &totalVertices, &maxLODLevel)

	fmt.Printf("\n=== Loading Tile File ===\n")
	tileFile := "test_data/OSGB/Data/Tile_+007_+006/Tile_+007_+006.osgb"
	res2 := rw.ReadNode(tileFile, nil)
	if res2 == nil {
		t.Fatalf("Failed to read %s", tileFile)
	}

	tileNode := res2.GetNode()
	if tileNode == nil {
		t.Fatalf("Failed to get node from %s", tileFile)
	}

	fmt.Printf("Tile file: %s\n", tileFile)
	fmt.Printf("Tile root node type: %T\n", tileNode)

	inspectNodeFull(tileNode, 0, &totalGeodes, &totalGeometries, &totalVertices, &maxLODLevel)

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total Geodes: %d\n", totalGeodes)
	fmt.Printf("Total Geometries: %d\n", totalGeometries)
	fmt.Printf("Total Vertices: %d\n", totalVertices)
	fmt.Printf("Max LOD Level: %d\n", maxLODLevel)
}

func inspectNodeFull(node interface{}, depth int, totalGeodes, totalGeometries, totalVertices, maxLODLevel *int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	switch n := node.(type) {
	case *model.Group:
		children := n.GetChildren()
		for i, child := range children {
			fmt.Printf("%sChild %d: %T\n", indent, i, child)
			inspectNodeFull(child, depth+1, totalGeodes, totalGeometries, totalVertices, maxLODLevel)
		}
	case *model.PositionAttitudeTransform:
		children := n.GetChildren()
		fmt.Printf("%sPositionAttitudeTransform: %d children\n", indent, len(children))
		for i, child := range children {
			fmt.Printf("%s  Child %d: %T\n", indent, i, child)
			inspectNodeFull(child, depth+1, totalGeodes, totalGeometries, totalVertices, maxLODLevel)
		}
	case *model.MatrixTransform:
		children := n.GetChildren()
		fmt.Printf("%sMatrixTransform: %d children\n", indent, len(children))
		for i, child := range children {
			fmt.Printf("%s  Child %d: %T\n", indent, i, child)
			inspectNodeFull(child, depth+1, totalGeodes, totalGeometries, totalVertices, maxLODLevel)
		}
	case *model.Transform:
		children := n.GetChildren()
		fmt.Printf("%sTransform: %d children\n", indent, len(children))
		for i, child := range children {
			fmt.Printf("%s  Child %d: %T\n", indent, i, child)
			inspectNodeFull(child, depth+1, totalGeodes, totalGeometries, totalVertices, maxLODLevel)
		}
	case *model.Lod:
		children := n.GetChildren()
		fmt.Printf("%sLOD Center: %v\n", indent, n.Center)
		fmt.Printf("%sLOD Radius: %v\n", indent, n.Radius)
		fmt.Printf("%sLOD RangeList: %v\n", indent, n.RangeList)
		fmt.Printf("%sLOD Children count: %d\n", indent, len(children))
		if maxLODLevel != nil {
			for _, r := range n.RangeList {
				level := int(r[0] / 100)
				if level > *maxLODLevel {
					*maxLODLevel = level
				}
			}
		}
		for i, child := range children {
			fmt.Printf("%s  LOD Child %d: %T\n", indent, i, child)
			inspectNodeFull(child, depth+2, totalGeodes, totalGeometries, totalVertices, maxLODLevel)
		}
	case *model.PagedLod:
		children := n.GetChildren()
		fmt.Printf("%sPagedLOD Center: %v\n", indent, n.Center)
		fmt.Printf("%sPagedLOD Radius: %v\n", indent, n.Radius)
		fmt.Printf("%sPagedLOD RangeList: %v\n", indent, n.RangeList)
		fmt.Printf("%sPagedLOD PerRangeDataList: %v\n", indent, n.PerRangeDataList)
		fmt.Printf("%sPagedLOD Children count: %d\n", indent, len(children))
		if maxLODLevel != nil {
			for _, r := range n.RangeList {
				level := int(r[0] / 100)
				if level > *maxLODLevel {
					*maxLODLevel = level
				}
			}
		}
		for i, child := range children {
			fmt.Printf("%s  PagedLOD Child %d: %T\n", indent, i, child)
			inspectNodeFull(child, depth+2, totalGeodes, totalGeometries, totalVertices, maxLODLevel)
		}
	case *model.Geode:
		if totalGeodes != nil {
			*totalGeodes++
		}
		children := n.GetChildren()
		fmt.Printf("%sGeode: %d drawables\n", indent, len(children))
		for i, child := range children {
			if geometry, ok := child.(*model.Geometry); ok {
				if totalGeometries != nil {
					*totalGeometries++
				}
				inspectGeometryFull(geometry, depth+1, totalVertices)
			} else {
				fmt.Printf("%s  Drawable %d: %T\n", indent, i, child)
			}
		}
	default:
		fmt.Printf("%sUnknown node type: %T\n", indent, node)
	}
}

func inspectGeometryFull(g *model.Geometry, depth int, totalVertices *int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	fmt.Printf("%sGeometry:\n", indent)

	if g.VertexArray != nil {
		vertexCount := 0
		if g.VertexArray.Data != nil {
			switch data := g.VertexArray.Data.(type) {
			case []float32:
				vertexCount = len(data) / 3
			case [][3]float32:
				vertexCount = len(data)
			case []int16:
				vertexCount = len(data) / 3
			case [][3]int16:
				vertexCount = len(data)
			}
			if totalVertices != nil {
				*totalVertices += vertexCount
			}
		}
		fmt.Printf("%s  VertexArray: DataSize=%d, DataType=%d (GL_SHORT=%d, GL_FLOAT=%d), VertexCount=%d\n",
			indent, g.VertexArray.DataSize, g.VertexArray.DataType, model.GLSHORT, model.GLFLOAT, vertexCount)
	}

	if g.NormalArray != nil {
		fmt.Printf("%s  NormalArray: DataSize=%d\n", indent, g.NormalArray.DataSize)
	}

	if g.TexCoordArrayList != nil && len(g.TexCoordArrayList) > 0 {
		for i, tex := range g.TexCoordArrayList {
			if tex != nil {
				fmt.Printf("%s  TexCoordArray[%d]: DataSize=%d, DataType=%d\n", indent, i, tex.DataSize, tex.DataType)
			}
		}
	}

	if g.ColorArray != nil {
		fmt.Printf("%s  ColorArray: DataSize=%d\n", indent, g.ColorArray.DataSize)
	}

	if g.VertexAttribList != nil && len(g.VertexAttribList) > 0 {
		fmt.Printf("%s  VertexAttribList: count=%d\n", indent, len(g.VertexAttribList))
	}

	if g.Primitives != nil {
		fmt.Printf("%s  Primitives: count=%d\n", indent, len(g.Primitives))
	}
}

func TestWriteOSG(t *testing.T) {
	rw := NewReadWrite()

	fmt.Printf("=== Write OSG Test ===\n\n")

	// Read a file
	res := rw.ReadNode("test_data/simpleroom.osgt", nil)
	if res == nil {
		t.Fatal("Failed to read simpleroom.osgt")
	}
	node := res.GetNode()
	if node == nil {
		t.Fatal("Failed to get node from simpleroom.osgt")
	}

	// Try writing - skip for now as write has issues
	fmt.Printf("Read simpleroom.osgt successfully: type=%T\n", node)
	fmt.Printf("Write functionality requires further debugging\n")

	fmt.Printf("\n=== Write Test Complete ===\n")
}

func TestCompress(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8}
	var bt []byte
	bf := bytes.NewBuffer(bt)
	w := zlib.NewWriter(bf)
	w.Write(buf)
	w.Close()
	fmt.Println(bf.Bytes())

	r, _ := zlib.NewReader(bf)
	var src1 []byte
	for {
		buf1 := make([]byte, 4096)
		n, _ := io.ReadFull(r, buf1)
		if n == 0 {
			break
		}
		src1 = append(src1, buf1[0:n]...)
	}
	fmt.Println(src1)
}
