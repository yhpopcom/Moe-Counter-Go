package utils

import (
	"encoding/base64"
	"image"
	_ "image/gif"
	_ "image/png"
	"io/fs"
	"strconv"
	"strings"
)

// 图片信息结构
type ImageInfo struct {
	Width  int
	Height int
	Data   string // base64编码的图片数据
}

// 拼接图片 - 生成计数器SVG（使用文件系统）
func CombineImages(count uint, publicFS fs.FS, themeName string, length int, scale float64, offset int, align string, pixelate string, darkmode string) (string, error) {
	// 为可选参数设置默认值
	if length <= 0 {
		length = 7 // 默认长度
	}
	if scale <= 0 {
		scale = 1.0 // 默认缩放
	}
	// 对齐方式默认居中
	if align == "" {
		align = "center"
	}
	// 像素化模式默认关闭
	if pixelate == "" {
		pixelate = "off"
	}
	// 暗黑模式处理已迁移到前端，此处不再使用该参数

	// 将计数转换为数字字符串并处理位数
	countStr := strconv.FormatUint(uint64(count), 10)
	// 长度不足时左侧补零，超过时保留实际长度
	if len(countStr) < length {
		countStr = strings.Repeat("0", length-len(countStr)) + countStr
	}

	// 构建主题路径（使用正斜杠分隔符）
	themePath := "assets/theme/" + themeName

	// 确定图片格式
	ext := ".png"
	if _, err := fs.Stat(publicFS, themePath+"/0.gif"); err == nil {
		ext = ".gif"
	}

	// 加载所有图片信息
	imageMap := make(map[string]*ImageInfo)

	// 加载数字图片
	for _, char := range "0123456789" {
		imgPath := themePath + "/" + string(char) + ext
		info, err := loadImageInfo(publicFS, imgPath)
		if err != nil {
			return "", err
		}
		imageMap[string(char)] = info
	}

	// 加载特殊图片
	specialFiles := []string{"start", "end"}
	for _, file := range specialFiles {
		imgPath := themePath + "/" + file + ext
		if _, err := fs.Stat(publicFS, imgPath); err == nil {
			info, err := loadImageInfo(publicFS, imgPath)
			if err != nil {
				return "", err
			}
			imageMap[file] = info
		}
	}

	// 构建字符序列
	chars := []string{}
	if _, exists := imageMap["start"]; exists {
		chars = append(chars, "start")
	}
	chars = append(chars, strings.Split(countStr, "")...)
	if _, exists := imageMap["end"]; exists {
		chars = append(chars, "end")
	}

	// 计算总宽度和最大高度
	totalWidth := 0
	maxHeight := 0
	for i, char := range chars {
		if info, exists := imageMap[char]; exists {
			scaledWidth := int(float64(info.Width) * scale)
			totalWidth += scaledWidth
			if i < len(chars)-1 {
				totalWidth += offset
			}
			if scaledHeight := int(float64(info.Height) * scale); scaledHeight > maxHeight {
				maxHeight = scaledHeight
			}
		}
	}

	// 根据对齐方式计算起始位置
	startX := 0
	if align == "center" {
		startX = (totalWidth - (totalWidth - offset*(len(chars)-1))) / 2
	} else if align == "right" {
		startX = totalWidth - (totalWidth - offset*(len(chars)-1))
	}

	// 生成SVG
	svg := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
	width="` + strconv.Itoa(totalWidth) + `" height="` + strconv.Itoa(maxHeight) + `" 
	viewBox="0 0 ` + strconv.Itoa(totalWidth) + ` ` + strconv.Itoa(maxHeight) + `">`

	// 仅在黑暗模式开启时添加CSS样式
	if darkmode == "on" {
		svg += `
	<style>
	  svg {
		image-rendering: pixelated;
		filter: brightness(.6);
	  }
	</style>`
	}

	xPos := startX
	for _, char := range chars {
		if info, exists := imageMap[char]; exists {
			width := int(float64(info.Width) * scale)
			height := int(float64(info.Height) * scale)
			yPos := (maxHeight - height) / 2 // 垂直居中

			// 添加像素化效果（占位）
			if pixelate == "on" {
				// 此处添加像素化处理逻辑
			}

			svg += `
	<image x="` + strconv.Itoa(xPos) + `" y="` + strconv.Itoa(yPos) + `" 
		width="` + strconv.Itoa(width) + `" height="` + strconv.Itoa(height) + `" 
		xlink:href="` + info.Data + `" />`

			xPos += width + offset
		}
	}

	svg += `
</svg>`
	return svg, nil
}

// 从文件系统加载图片信息
func loadImageInfo(publicFS fs.FS, path string) (*ImageInfo, error) {
	// 读取图片文件
	data, err := fs.ReadFile(publicFS, path)
	if err != nil {
		return nil, err
	}

	// 获取图片尺寸
	file, err := publicFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	// 获取MIME类型
	mimeType := "image/png"
	if strings.HasSuffix(path, ".gif") {
		mimeType = "image/gif"
	}

	// 生成base64数据URI
	base64Data := base64.StdEncoding.EncodeToString(data)
	dataURI := "data:" + mimeType + ";base64," + base64Data

	return &ImageInfo{
		Width:  config.Width,
		Height: config.Height,
		Data:   dataURI,
	}, nil
}
