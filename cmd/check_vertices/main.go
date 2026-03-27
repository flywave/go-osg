package main

import (
	"fmt"
	"math"
	"os"
)

func float32fromBytes(b []byte) float32 {
	bits := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	return math.Float32frombits(bits)
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("File size: %d bytes\n", len(data))

	vec3Offsets := []int{}
	for i := 0; i < len(data)-4; i++ {
		if string(data[i:i+4]) == "Vec3" {
			vec3Offsets = append(vec3Offsets, i)
		}
	}

	fmt.Printf("Found %d 'Vec3' patterns\n", len(vec3Offsets))

	for _, offset := range vec3Offsets {
		fmt.Printf("\nAnalyzing at offset %d:", offset)
		fmt.Printf("\n  Context: % x", data[offset-10:offset+30])

		sizeOffset := offset - 4
		if sizeOffset < 0 {
			continue
		}
		size := int(data[sizeOffset]) | int(data[sizeOffset+1])<<8 | int(data[sizeOffset+2])<<16 | int(data[sizeOffset+3])<<24
		fmt.Printf("\n  Size at offset %d: %d", sizeOffset, size)

		if size > 0 && size < 100000 {
			dataStart := sizeOffset + 4
			fmt.Printf("\n  Data at offset %d:", dataStart)
			for v := 0; v < size && v < 4; v++ {
				base := dataStart + v*12
				if base+12 <= len(data) {
					x := float32fromBytes(data[base : base+4])
					y := float32fromBytes(data[base+4 : base+8])
					z := float32fromBytes(data[base+8 : base+12])
					fmt.Printf("\n    V[%d]: (%.3f, %.3f, %.3f)", v, x, y, z)
				}
			}
		}

		if len(vec3Offsets) > 3 {
			break
		}
	}
}
