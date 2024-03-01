package tool

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"gin-skeleton/model/admin/system"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var _ Storager = (*qiniuStorage)(nil)
var _ Storager = (*aliStorage)(nil)
var _ Storager = (*tencentStorage)(nil)

// Storager 存储接口
type Storager interface {
	// PutFileByIo 以文件IO的方式上传
	PutFileByIo(source string, fileName string, fileSize int64, fileIo io.Reader) (url string, err error)

	// PutFileByUrl 抓取URL文件并上传
	PutFileByUrl(source string, fileName string, sourceUrl string) (url string, err error)

	// DeleteFile 删除文件（注：如果配置有变更，会导致删除失败）
	DeleteFile(fileUrl string) (err error)
}

// 七牛存储
type qiniuStorage struct {
	Bucket    string // 空间
	AccessKey string // 密匙
	SecretKey string // 密钥
	Domain    string // 域名
}

// 阿里存储
type aliStorage struct {
	Bucket    string // 空间
	AccessKey string // 密匙
	SecretKey string // 密钥
	Domain    string // 域名
}

// 腾讯存储
type tencentStorage struct {
	Bucket    string // 空间
	AccessKey string // 密匙
	SecretKey string // 密钥
	Domain    string // 域名
}

// NewStorage 初始化存储对象
func NewStorage(cf system.FileConfig) Storager {
	switch cf.Provider {
	case system.FileProviderQiniu:
		return &qiniuStorage{Bucket: cf.Bucket, AccessKey: cf.AccessKey, SecretKey: cf.SecretKey, Domain: cf.Domain}
	case system.FileProviderAli:
		return &aliStorage{Bucket: cf.Bucket, AccessKey: cf.AccessKey, SecretKey: cf.SecretKey, Domain: cf.Domain}
	default:
		return &tencentStorage{Bucket: cf.Bucket, AccessKey: cf.AccessKey, SecretKey: cf.SecretKey, Domain: cf.Domain}
	}
}

// PutFileByIo 以文件IO的方式上传
func (s *qiniuStorage) PutFileByIo(source string, fileName string, fileSize int64, fileIo io.Reader) (url string, err error) {
	upResult := storage.PutRet{}
	putPolicy := storage.PutPolicy{Scope: s.Bucket}
	upToken := putPolicy.UploadToken(qbox.NewMac(s.AccessKey, s.SecretKey))
	uploader := storage.NewFormUploader(&storage.Config{})
	err = uploader.Put(context.Background(), &upResult, upToken, generateFilePath(source, fileName), fileIo, fileSize, &storage.PutExtra{})
	if err != nil {
		return
	}
	url = s.Domain + "/" + upResult.Key
	return
}

// PutFileByUrl 抓取URL文件并上传
func (s *qiniuStorage) PutFileByUrl(source string, fileName string, sourceUrl string) (url string, err error) {
	manager := storage.NewBucketManager(qbox.NewMac(s.AccessKey, s.SecretKey), &storage.Config{})
	fetch, err := manager.Fetch(sourceUrl, s.Bucket, generateFilePath(source, fileName))
	if err != nil {
		return
	}
	url = s.Domain + "/" + fetch.Key
	return
}

// DeleteFile 删除文件
func (s *qiniuStorage) DeleteFile(fileUrl string) (err error) {
	manager := storage.NewBucketManager(qbox.NewMac(s.AccessKey, s.SecretKey), &storage.Config{})
	err = manager.Delete(s.Bucket, strings.TrimPrefix(fileUrl, s.Domain+"/"))
	return
}

// PutFileByIo 以文件IO的方式上传
func (s *aliStorage) PutFileByIo(source string, fileName string, fileSize int64, fileIo io.Reader) (url string, err error) {
	// @todo 待实现
	return
}

// PutFileByUrl 抓取URL文件并上传
func (s *aliStorage) PutFileByUrl(source string, fileName string, sourceUrl string) (url string, err error) {
	// @todo 待实现
	return
}

// DeleteFile 删除文件
func (s *aliStorage) DeleteFile(fileUrl string) (err error) {
	// @todo 待实现
	return
}

// PutFileByIo 以文件IO的方式上传
func (s *tencentStorage) PutFileByIo(source string, fileName string, fileSize int64, fileIo io.Reader) (url string, err error) {
	// @todo 待实现
	return
}

// PutFileByUrl 抓取URL文件并上传
func (s *tencentStorage) PutFileByUrl(source string, fileName string, sourceUrl string) (url string, err error) {
	// @todo 待实现
	return
}

// DeleteFile 删除文件
func (s *tencentStorage) DeleteFile(fileUrl string) (err error) {
	// @todo 待实现
	return
}

// 生成文件地址
func generateFilePath(source, filename string) string {
	return fmt.Sprintf("%s/%s/%s%s", source, time.Now().Format("20060102"), time.Now().Format("150405"), filename)
}
