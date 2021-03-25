package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// 获取项目根目录
func GetRootPath() string {
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("get project root path faild: %s", err))
	}
	return rootPath
}
