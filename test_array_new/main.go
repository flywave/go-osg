package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Testing Array data creation...\n")

	// Test 1: Create a float32 Vec3 array
	data := [][3]float32{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
		{7.0, 8.0, 9.0},
	}

	fmt.Printf("Created [3]float32 array with %d vertices\n", len(data))
	for i, v := range data {
		fmt.Printf("  [%d] x=%.6f, y=%.6f, z=%.6f\n", i, v[0], v[1], v[2])
	}

	fmt.Printf("\nTest 2: Create a model.Array from this data")
	arr := &model.Array{
		Type:     model.Vec3ArrayType,
		DataType: model.GLFLOAT,
		DataSize: 3,
		Data:     &data,
	}

	fmt.Printf("Array properties:\n")
	fmt.Printf("  Type: %d (expected 16 for Vec3Array)\n", arr.Type)
	fmt.Printf("  DataType: %d (expected 5126 for GLFLOAT)\n", arr.DataType)
	fmt.Printf("  DataSize: %d (expected 3 for Vec3)\n", arr.DataSize)
	fmt.Printf("  Data: %v\n", arr.Data)

	fmt.Printf("\nTest 3: Print first vertex\n")
	if arr.Data != nil {
		if arr.Data, ok := arr.Data.(*[3]float32); ok {
			fmt.Printf("First vertex: (%.6f, %.6f, %.6f)\n", arr.Data[0][0], arr.Data[0][1], arr.Data[0][2])
		}
	} else {
		fmt.Printf("Data is nil or type mismatch\n")
	}
}

	fmt.Printf("Created [3]float32 array with %d vertices\n", len(data))
	for i, v := range data {
		fmt.Printf("  [%d] x=%.6f, y=%.6f, z=%.6f\n", i, v[0], v[1], v[2])
	}

	fmt.Printf("\nTest 2: Create a model.Array from this data")
}
}
