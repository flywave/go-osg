package main

import (
	"fmt"
	"osg"
	"github.com/flywave/go-osg/model"
)

func main() {
	fmt.Println("Testing...")

	data := make([][3]float32, 4)
	data[0] = [1.0, 2.0, 3.0]
	data[1] = [4.0, 5.0, 6.0]
	data[2] = [7.0, 8.0, 9.0]
	data[3] = [10.0, 11.0, 12.0]

	arr := &model.Array{
		Type:     model.Vec3ArrayType,
		DataType: model.GLFLOAT,
		DataSize: 3,
		Data:     data,
	}

	fmt.Printf("Array created: Type=%d, DataType=%d, DataSize=%d\n", arr.Type, arr.DataType, arr.DataSize)
	fmt.Printf("Data length: %d\n", len(arr.Data))
	fmt.Printf("Data pointer: %v\n", arr.Data)

	if arr.Data == nil {
		fmt.Println("Data is nil!")
		return
	}

	if float32Data, ok := arr.Data.([][3]float32); ok {
		fmt.Printf("Type assertion failed: Data is %[3]float32\n", float32Data)
		fmt.Printf("First 3 vertices:\n")
		for i := 0; i < 3 && i < len(float32Data); i++ {
			fmt.Printf("  [%d] x=%.6f, y=%.6f, z=%.6f\n", i, float32Data[i][0], float32Data[i][1], float32Data[i][2])
		}
	}
}
