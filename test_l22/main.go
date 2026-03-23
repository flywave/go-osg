package main

import (
	"fmt"

	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func main() {
	inputPath := "/Users/xuning/Work/go-osg/test_data/0131/Data/Tile_+001_+000/Tile_+001_+000_L22_00020.osgb"

	fmt.Printf("========================================\n")
	fmt.Printf("Reading OSGB file: %s\n", inputPath)
	fmt.Printf("========================================\n\n")

	rw := osg.NewReadWrite()
	res := rw.ReadNode(inputPath, nil)
	if res == nil {
		fmt.Printf("Failed to read OSGB file\n")
		return
	}

	node := res.GetNode()
	if node == nil {
		fmt.Printf("Node is nil after reading!\n")
		return
	}
	fmt.Printf("Node type: %T\n\n", node)

	// 递归查找所有Geometry对象
	findGeometries(node, 0)
}

func findGeometries(node interface{}, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	switch n := node.(type) {
	case *model.Group:
		children := n.GetChildren()
		for _, child := range children {
			fmt.Printf("%sGroup child: %T\n", prefix, child)
			findGeometries(child, indent+1)
		}
	case *model.Geode:
		children := n.GetChildren()
		for _, child := range children {
			fmt.Printf("%sGeode child: %T\n", prefix, child)
			if geom, ok := child.(*model.Geometry); ok {
				printGeometryInfo(geom, prefix+"    ")
			}
		}
	case *model.PagedLod:
		children := n.GetChildren()
		for _, child := range children {
			findGeometries(child, indent+1)
		}
	case *model.MatrixTransform:
		children := n.GetChildren()
		for _, child := range children {
			findGeometries(child, indent+1)
		}
	}
}

func printGeometryInfo(geom *model.Geometry, prefix string) {
	fmt.Printf("\n%s========================================\n", prefix)
	fmt.Printf("%sGeometry Info:\n", prefix)
	fmt.Printf("%s========================================\n", prefix)

	if geom.VertexArray == nil {
		fmt.Printf("%sVertexArray is nil\n", prefix)
		return
	}

	fmt.Printf("%sVertexArray properties:\n", prefix)
	fmt.Printf("%s  Type: %d (expected 16 for Vec3Array)\n", prefix, geom.VertexArray.Type)
	fmt.Printf("%s  DataType: %d (expected 5126 for GLFLOAT)\n", prefix, geom.VertexArray.DataType)
	fmt.Printf("%s  DataSize: %d (expected 3 for Vec3)\n", prefix, geom.VertexArray.DataSize)
	fmt.Printf("%s  Binding: %d\n", prefix, geom.VertexArray.Binding)
	fmt.Printf("%s  Normalize: %v\n", prefix, geom.VertexArray.Normalize)
	fmt.Printf("%s  Data type: %T\n", prefix, geom.VertexArray.Data)

	if geom.VertexArray.Data == nil {
		fmt.Printf("%sData is nil\n", prefix)
		return
	}

	switch data := geom.VertexArray.Data.(type) {
	case []float32:
		fmt.Printf("%sFloat32 array, length: %d\n", prefix, len(data))
		fmt.Printf("%sFirst 9 values (3 vertices):\n", prefix)
		for j := 0; j < 9 && j < len(data); j += 3 {
			fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n", prefix, j/3, data[j], data[j+1], data[j+2])
		}
		if len(data) > 12 {
			fmt.Printf("%sLast 9 values:\n", prefix)
			start := len(data) - 9
			for j := 0; j < 9 && start+j < len(data); j += 3 {
				fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n", prefix, (len(data)-9)/3+j/3, data[start+j], data[start+j+1], data[start+j+2])
			}
		}
	case [][3]float32:
		fmt.Printf("%s[3]float32 array, length: %d\n", prefix, len(data))
		fmt.Printf("%sFirst 3 vertices:\n", prefix)
		for j := 0; j < 3 && j < len(data); j++ {
			fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n", prefix, j, data[j][0], data[j][1], data[j][2])
		}
		if len(data) > 3 {
			fmt.Printf("%sLast 3 vertices:\n", prefix)
			start := len(data) - 3
			for j := 0; j < 3 && start+j < len(data); j++ {
				fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n", prefix, start+j, data[start+j][0], data[start+j][1], data[start+j][2])
			}
		}
	default:
		fmt.Printf("%sUnknown data type: %T\n", prefix, data)
	}

	fmt.Printf("%sPrimitives: %d\n", prefix, len(geom.Primitives))
}
