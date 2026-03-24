package serialization_test

import (
	"testing"

	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

const testFilePath = "../test_data/0131/Data/Tile_+001_+000/Tile_+001_+000_L22_00020.osgb"

// TestReadGroup tests reading osg::Group with children
func TestReadGroup(t *testing.T) {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(testFilePath, nil)
	if res == nil {
		t.Fatalf("Failed to read OSGB file")
	}

	node := res.GetNode()
	if node == nil {
		t.Fatalf("Node is nil after reading")
	}

	group, ok := node.(*model.Group)
	if !ok {
		t.Fatalf("Expected *model.Group, got %T", node)
	}

	children := group.GetChildren()
	if len(children) == 0 {
		t.Errorf("Group has no children")
	}

	t.Logf("Group has %d children", len(children))
}

// TestReadPagedLod tests reading osg::PagedLOD
func TestReadPagedLod(t *testing.T) {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(testFilePath, nil)
	if res == nil {
		t.Fatalf("Failed to read OSGB file")
	}

	node := res.GetNode()
	group, ok := node.(*model.Group)
	if !ok {
		t.Fatalf("Expected *model.Group, got %T", node)
	}

	children := group.GetChildren()
	if len(children) == 0 {
		t.Fatalf("Group has no children")
	}

	pagedLod, ok := children[0].(*model.PagedLod)
	if !ok {
		t.Fatalf("Expected first child to be *model.PagedLod, got %T", children[0])
	}
	t.Logf("PagedLod found, has %d children", len(pagedLod.Lod.Group.Children))
}

// TestReadGeometry tests reading Geometry with VertexArray and TexCoordArrayList
func TestReadGeometry(t *testing.T) {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(testFilePath, nil)
	if res == nil {
		t.Fatalf("Failed to read OSGB file")
	}

	geom := findFirstGeometry(res.GetNode())
	if geom == nil {
		t.Fatal("No geometry found in file")
	}

	t.Logf("Found Geometry")

	// Check VertexArray
	if geom.VertexArray == nil {
		t.Fatal("VertexArray is nil")
	}

	t.Logf("VertexArray properties:")
	t.Logf("  Type: %d (expected 16 for Vec3Array)", geom.VertexArray.Type)
	t.Logf("  DataType: %d (expected 5126 for GLFLOAT)", geom.VertexArray.DataType)
	t.Logf("  DataSize: %d (expected 3 for Vec3)", geom.VertexArray.DataSize)
	t.Logf("  Binding: %d", geom.VertexArray.Binding)
	t.Logf("  Normalize: %v", geom.VertexArray.Normalize)
	t.Logf("  Data type: %T", geom.VertexArray.Data)

	if geom.VertexArray.Data == nil {
		t.Fatal("VertexArray.Data is nil - this is the main bug!")
	}

	// Check data content
	switch data := geom.VertexArray.Data.(type) {
	case []float32:
		t.Logf("Vertex count: %d", len(data)/3)
		if len(data) >= 9 {
			t.Logf("First vertex: (%.6f, %.6f, %.6f)", data[0], data[1], data[2])
		}
	case [][3]float32:
		t.Logf("Vertex count: %d", len(data))
		if len(data) > 0 {
			t.Logf("First vertex: (%.6f, %.6f, %.6f)", data[0][0], data[0][1], data[0][2])
		}
	default:
		t.Errorf("Unexpected VertexArray.Data type: %T", data)
	}

	// Check TexCoordArrayList
	if len(geom.TexCoordArrayList) > 0 && geom.TexCoordArrayList[0] != nil {
		texArray := geom.TexCoordArrayList[0]
		t.Logf("TexCoordArray Type: %d", texArray.Type)
		t.Logf("TexCoordArray DataType: %d", texArray.DataType)
		t.Logf("TexCoordArray DataSize: %d", texArray.DataSize)
		if texArray.Data == nil {
			t.Error("TexCoordArray.Data is nil")
		}
	}

	// Check Primitives
	if len(geom.Primitives) > 0 {
		t.Logf("Primitive count: %d", len(geom.Primitives))
	}
}

// TestArraySerialization tests that arrays are properly serialized for FileVersion >= 112
func TestArraySerialization(t *testing.T) {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(testFilePath, nil)
	if res == nil {
		t.Fatalf("Failed to read OSGB file")
	}

	geom := findFirstGeometry(res.GetNode())
	if geom == nil {
		t.Fatal("No geometry found")
	}

	// Test Vec3Array (VertexArray)
	testArrayProperties(t, geom.VertexArray, "VertexArray", 16, 5126, 3)

	// Test Vec2Array (TexCoordArray)
	if len(geom.TexCoordArrayList) > 0 {
		testArrayProperties(t, geom.TexCoordArrayList[0], "TexCoordArray", 16, 5126, 2)
	}
}

func testArrayProperties(t *testing.T, arr *model.Array, name string, expectedType int, expectedDataType int32, expectedDataSize int32) {
	if arr == nil {
		t.Fatalf("%s is nil", name)
	}

	if int(arr.Type) != expectedType {
		t.Errorf("%s.Type = %d, expected %d", name, arr.Type, expectedType)
	}

	if arr.DataType != expectedDataType {
		t.Errorf("%s.DataType = %d, expected %d", name, arr.DataType, expectedDataType)
	}

	if arr.DataSize != expectedDataSize {
		t.Errorf("%s.DataSize = %d, expected %d", name, arr.DataSize, expectedDataSize)
	}

	if arr.Data == nil {
		t.Errorf("%s.Data is nil", name)
		return
	}

	// Check data length
	switch data := arr.Data.(type) {
	case []float32:
		if len(data)%int(expectedDataSize) != 0 {
			t.Errorf("%s data length %d is not a multiple of DataSize %d", name, len(data), expectedDataSize)
		} else {
			t.Logf("%s has %d elements ([]float32)", name, len(data)/int(expectedDataSize))
		}
	case [][2]float32:
		t.Logf("%s has %d elements ([][2]float32)", name, len(data))
	case [][3]float32:
		t.Logf("%s has %d elements ([][3]float32)", name, len(data))
	case [][4]float32:
		t.Logf("%s has %d elements ([][4]float32)", name, len(data))
	default:
		t.Errorf("%s.Data has unexpected type %T", name, data)
	}
}

func findFirstGeometry(node interface{}) *model.Geometry {
	switch n := node.(type) {
	case *model.Group:
		children := n.GetChildren()
		for _, child := range children {
			if geom := findFirstGeometry(child); geom != nil {
				return geom
			}
		}
	case *model.PagedLod:
		children := n.Lod.Group.GetChildren()
		for _, child := range children {
			if geom := findFirstGeometry(child); geom != nil {
				return geom
			}
		}
	case *model.Geode:
		children := n.GetChildren()
		for _, child := range children {
			if geom, ok := child.(*model.Geometry); ok {
				return geom
			}
		}
	}
	return nil
}
