package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gin-skeleton/helper"

	"gorm.io/gorm"
)

// LocalTime 自定义时间格式
type LocalTime time.Time

// BaseModel 基础Model
type BaseModel struct {
	ID        int64     `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	CreatedAt LocalTime `json:"createdAt"`
	UpdatedAt LocalTime `json:"updatedAt"`
}

// BaseIdParams ID请求参数
type BaseIdParams struct {
	ID int64 `json:"id" Remark:"ID" binding:"required,gt=0"`
}

// BaseIdResult ID结果响应
type BaseIdResult struct {
	ID int64 `json:"id"`
}

// BaseOrderByParams 排序请求参数
type BaseOrderByParams struct {
	Field string `json:"field,omitempty" remark:"排序字段"`
	Order string `json:"order,omitempty" remark:"排序方向"`
}

// BasePageParams 分页请求参数
type BasePageParams struct {
	BaseOrderByParams
	Page     int `json:"page" remark:"页码" binding:"required,gt=0"`
	PageSize int `json:"pageSize" remark:"条数" binding:"required,gt=0"`
}

// BasePageResult 分页结果响应
type BasePageResult[T any] struct {
	Items []*T  `json:"items"`
	Total int64 `json:"total"`
}

// Paginate 分页数据
func Paginate(pageInfo BasePageParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageInfo.Page <= 0 {
			pageInfo.Page = 1
		}
		if pageInfo.PageSize <= 0 {
			pageInfo.PageSize = 10
		}
		if pageInfo.PageSize > 100 {
			pageInfo.PageSize = 100
		}

		offset := (pageInfo.Page - 1) * pageInfo.PageSize
		return db.Offset(offset).Limit(pageInfo.PageSize).Scopes(OrderBy(pageInfo.BaseOrderByParams))
	}
}

// OrderBy 排序处理
func OrderBy(oderByInfo BaseOrderByParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if oderByInfo.Field == "" || oderByInfo.Order == "" {
			return db
		}

		// 将排序字段由小驼峰转化为下划线格式
		// 将排序参数尾部的end字符去掉，因为前端传过来的排序字段为"ascend"、"descend"
		field := helper.Camel2Case(oderByInfo.Field)
		order := strings.TrimRight(oderByInfo.Order, "end")
		return db.Order(fmt.Sprintf("%s %s", field, order))
	}
}

// CancelPaginate 取消分页
func CancelPaginate(db *gorm.DB) *gorm.DB {
	return db.Offset(-1).Limit(-1)
}

// Scan 重写查询方法
func (lt *LocalTime) Scan(v interface{}) error {
	if val, ok := v.(time.Time); ok {
		*lt = LocalTime(val)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Value 重写写入方法
func (lt LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(lt).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(lt), nil
}

// MarshalJSON 时间格式化处理方法
func (lt *LocalTime) MarshalJSON() ([]byte, error) {
	if time.Time(*lt).IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*lt).Format(helper.ToDateTimeString))), nil
}
