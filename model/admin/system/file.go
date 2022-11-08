package system

import (
	"gin-skeleton/helper"
	"gin-skeleton/model"
)

const (
	FileProviderAli     = "ali"     // 提供商 - 阿里云
	FileProviderQiniu   = "qiniu"   // 提供商 - 七牛云
	FileProviderTencent = "tencent" // 提供商 - 腾讯云
)

// File 文件表
type File struct {
	model.BaseModel
	FileName  string `json:"fileName"`  // 文件名
	FileSize  int64  `json:"fileSize"`  // 文件大小，单位：kb
	FileType  string `json:"fileType"`  // 文件类型
	FileUrl   string `json:"fileUrl"`   // 文件地址
	Thumbnail string `json:"thumbnail"` // 缩略图地址
	Provider  string `json:"provider"`  // 提供商：qiniu-七牛云，ali-阿里云，tencent-腾讯云
	Username  string `json:"username"`  // 用户名
	Nickname  string `json:"nickname"`  // 昵称
	Remark    string `json:"remark"`    // 备注
}

// FileConfig 文件配置
type FileConfig struct {
	Provider   string   `json:"provider"`   // 提供商：qiniu-七牛云，ali-阿里云，tencent-腾讯云
	Bucket     string   `json:"bucket"`     // 空间
	AccessKey  string   `json:"accessKey"`  // 密匙
	SecretKey  string   `json:"secretKey"`  // 密钥
	Domain     string   `json:"domain"`     // 域名
	ThumbConf  string   `json:"thumbConf"`  // 缩略图配置
	AllowTypes []string `json:"allowTypes"` // 允许上传的类型
	Remark     string   `json:"remark"`     // 备注
}

// NewFile 初始化对象
func NewFile() *File {
	return &File{}
}

// GetFileList 获取文件列表
func (f *File) GetFileList(keyword, uploader string, createdDate []string, pageInfo model.BasePageParams) (pr *model.BasePageResult[File], err error) {
	fileModel := helper.GormDefaultDb.Model(NewFile())
	if keyword != "" {
		fileModel.Where("file_name like ? or remark like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if uploader != "" {
		fileModel.Where("username like ? or nickname like ?", "%"+uploader+"%", "%"+uploader+"%")
	}
	if len(createdDate) == 2 && createdDate[0] != "" && createdDate[1] != "" {
		fileModel.Where("created_at >= ? and created_at <= ?", createdDate[0]+" 00:00:00", createdDate[1]+" 23:59:59")
	}

	pr = &model.BasePageResult[File]{Items: make([]*File, 0), Total: 0}
	err = fileModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.Count).Count(&pr.Total).Error
	return
}

// FindFileInfo 查询文件信息
func (f *File) FindFileInfo(id int64) (file *File, err error) {
	err = helper.GormDefaultDb.First(&file, id).Error
	return
}
