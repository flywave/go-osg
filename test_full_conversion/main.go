package main

import (
	"fmt"
	"os"
	"time"

	"github.com/flywave/go-osg/tiles3d"
	_ "github.com/flywave/go-proj"
)

func init() {
	os.Setenv("PROJ_DATA", "/Users/xuning/Work/go-proj/proj_data")
	os.Setenv("PROJ_LIB", "/Users/xuning/Work/go-proj/proj_data")
}

func main() {
	inputPath := "/Users/xuning/Work/go-osg/test_data/0131/Data/main.osgb"
	outputPath := "/tmp/osgb_0131_output"

	startTime := time.Now()

	fmt.Printf("========================================\n")
	fmt.Printf("完整转换测试\n")
	fmt.Printf("========================================\n")
	fmt.Printf("输入路径: %s\n", inputPath)
	fmt.Printf("输出路径: %s\n", outputPath)
	fmt.Printf("开始时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("========================================\n\n")

	opts := tiles3d.DefaultConverterOptions()
	opts.EnableTexture = true
	opts.EnableUnlit = true
	opts.MaxLOD = -1 // 转换所有LOD层级，
	opts.GeoidModel = "none"
	opts.GeoidDataPath = ""

	// 坐标 518078, 4080366 看起来像是CGCS2000坐标系
	// 根据Y坐标的百万位数字，可能是3度带或6度带
	// 尝试EPSG:4545 (CGCS2000 / 3-degree Gauss-Kruger zone 36)
	opts.SourceSRS = "EPSG:4545"

	err := os.RemoveAll(outputPath)
	if err != nil {
		fmt.Printf("清理输出目录失败: %v\n", err)
	}

	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	fmt.Println("开始转换...")
	fmt.Println()

	result, err := tiles3d.OSGBTo3DTiles(inputPath, outputPath, opts)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("\n========================================\n")
	fmt.Printf("转换成功!\n")
	fmt.Printf("========================================\n")
	fmt.Printf("JSON长度: %d 字节\n", len(result.JSON))
	fmt.Printf("包围盒: %v\n", result.BoundingBox)
	fmt.Printf("结束时间: %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("总耗时: %v\n", duration)
	fmt.Printf("========================================\n\n")

	tilesetPath := outputPath + "/tileset.json"
	if _, err := os.Stat(tilesetPath); err == nil {
		fmt.Println("✓ tileset.json 创建成功")
	}

	files, _ := os.ReadDir(outputPath)
	fmt.Printf("\n输出文件统计:\n")
	fmt.Printf("  总文件数: %d\n", len(files))

	b3dmCount := 0
	jsonCount := 0
	for _, f := range files {
		if f.Name()[len(f.Name())-5:] == ".b3dm" {
			b3dmCount++
		}
		if f.Name()[len(f.Name())-5:] == ".json" {
			jsonCount++
		}
	}
	fmt.Printf("  B3DM文件: %d\n", b3dmCount)
	fmt.Printf("  JSON文件: %d\n", jsonCount)

	if b3dmCount > 0 {
		fmt.Printf("\n前10个文件:\n")
		count := 0
		for _, f := range files {
			if count >= 10 {
				break
			}
			info, _ := f.Info()
			fmt.Printf("  - %s (%d bytes)\n", f.Name(), info.Size())
			count++
		}
	}
}
