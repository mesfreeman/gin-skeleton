package helper

import (
	"crypto/md5"
	"encoding/hex"
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

// MD5 加密
func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
