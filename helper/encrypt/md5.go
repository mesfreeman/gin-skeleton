package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5 MD5加密
func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
