package main

import (
	"fmt"
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func main() {
	inputPath := "/Users/xuning/Work/go-osg/test_data/0131/Data/Tile_+001_+000/Tile_+001_+000_L22_00020.osgb"

	fmt.Printf("Reading OSGB file: %s\n", inputPath)

	rw := osg.NewReadWrite()
	res := rw.ReadNode(inputPath, nil)
	if res == nil {
		fmt.Printf("Failed to read OSGB file\n")
		return
	}

	node := res.GetNode()
	fmt.Printf("Node type: %T\n", node)

	// Find first Geometry with VertexArray
	findFirstGeometryWithArray(node, 0)
}

func findFirstGeometryWithArray(node interface{}, depth int) {
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "  "
	}

	switch n := node.(type) {
	case *model.Group:
		children := n.GetChildren()
		for _, child := range children {
			fmt.Printf("%sGroup child: %T\n", prefix, child)
			findFirstGeometryWithArray(child, depth+1)
		}
	case *model.Geode:
		children := n.GetChildren()
		for _, child := range children {
			fmt.Printf("%sGeode child: %T\n", prefix, child)
			if geom, ok := child.(*model.Geometry); ok {
				printGeometryInfo(geom, prefix+"  ")
				if geom.VertexArray != nil && geom.VertexArray.Data != nil {
					fmt.Printf("%s*** Found Geometry with VertexArray ***\n", prefix)
					return
				}
			}
		}
	case *model.PagedLod:
		children := n.GetChildren()
		for _, child := range children {
			findFirstGeometryWithArray(child, depth+1)
		}
	case *model.MatrixTransform:
		children := n.GetChildren()
		for _, child := range children {
			findFirstGeometryWithArray(child, depth+1)
		}
	}
}

func printGeometryInfo(geom *model.Geometry, prefix string) {
	fmt.Printf("%s========================================\n", prefix)
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

	fmt.Printf("%sData length: %d\n", prefix, getArrayLength(geom.VertexArray.Data))

	switch data := geom.VertexArray.Data.(type) {
	case []float32:
		fmt.Printf("%sFloat32 array, length: %d\n", prefix, len(data))
		fmt.Printf("%sFirst 6 values (2 vertices):\n", prefix)
		for j := 0; j < 6 && j < len(data); j += 3 {
			fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n",
				prefix, j/3, data[j], data[j+1], data[j+2])
		}
		if len(data) > 12 {
			fmt.Printf("%sLast 6 values:\n", prefix)
			start := len(data) - 6
			for j := 0; j < 6 && start+j < len(data); j += 3 {
				fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n",
					prefix, start/3+j/3, data[start+j], data[start+j+1], data[start+j+2])
			}
		}
	case [][3]float32:
		fmt.Printf("%s[3]float32 array, length: %d\n", prefix, len(data))
		fmt.Printf("%sFirst 2 vertices:\n", prefix)
		for j := 0; j < 2 && j < len(data); j++ {
			fmt.Printf("%s  [%d] x=%.6f, y=%.6f, z=%.6f\n", prefix, j, data[j][0], data[j][1], data[j][2])
		}
	default:
		fmt.Printf("%sUnknown data type: %T\n", prefix, geom.VertexArray.Data)
	}
}

func getArrayLength(data interface{}) int {
	switch v := data.(type) {
	case []float32:
		return len(v)
	case []float64:
		return len(v)
	case []int8:
		return len(v)
	case []uint8:
		return len(v)
	case []int16:
		return len(v)
	case []uint16:
		return len(v)
	case []int32:
		return len(v)
	case []uint32:
		return len(v)
	case []int64:
		return len(v)
	case []uint64:
		return len(v)
	case [][3]float32:
		return len(v)
	case [][2]float32:
		return len(v)
	case [][4]float32:
		return len(v)
	default:
		return 0
	}
}
