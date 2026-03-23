package main

import (
	"fmt"
	"os"
)

func main() {
	inputPath := "/Users/xuning/Work/go-osg/test_data/0131/Data/Tile_+000_+000/Tile_+000_+000.osgb"

	fmt.Printf("Reading OSGB file: %s\n", inputPath)

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// 读取前200个字节
	buf := make([]byte, 200)
	n, err := file.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	fmt.Printf("Read %d bytes\n", n)
	fmt.Printf("\nFirst 200 bytes:\n")
	for i := 0; i < n && i < 200; i++ {
		fmt.Printf("%02x ", buf[i])
		if (i+1)%16 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")

	// 尝试读取为ASCII
	fmt.Printf("\nFirst 200 bytes as ASCII:\n")
	for i := 0; i < n && i < 200; i++ {
		if buf[i] >= 32 && buf[i] <= 126 {
			fmt.Printf("%c", buf[i])
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
}
