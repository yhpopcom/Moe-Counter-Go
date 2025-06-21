package controller

import (
	"io/fs"
)

// 计数器请求参数结构
type CounterRequest struct {
	Name     string  `form:"name"`
	Theme    string  `form:"theme"`
	Length   int     `form:"length" default:"7"`
	Scale    float64 `form:"scale" default:"0"`
	Offset   int     `form:"offset" default:"0"`
	Align    string  `form:"align" default:"left"`
	Pixelate string  `form:"pixelate" default:"off"`
	Darkmode string  `form:"darkmode" default:"auto"`
	Base     int     `form:"base" default:"10"`
	Num      string  `form:"num"`
}

// 获取可用主题列表
func ListThemes(publicFS fs.FS) ([]string, error) {
	// 打开主题目录
	themeDir := "assets/theme"
	dir, err := fs.Sub(publicFS, themeDir)
	if err != nil {
		return nil, err
	}

	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		return nil, err
	}

	themes := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			themes = append(themes, entry.Name())
		}
	}
	return themes, nil
}
