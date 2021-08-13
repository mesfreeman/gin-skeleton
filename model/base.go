package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// FormatTime 自定义Time格式
type FormatTime struct {
	Time time.Time
}

// MarshalJSON 重写 MarshaJSON 方法，实现格式化时间
func (f FormatTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", f.Time.Format("2006-01-02 15:04:05"))
	return []byte(output), nil
}

// Value 重写 Value 方法，自定义时间格式写入数据库
func (f FormatTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if f.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return f.Time, nil
}

// Scan 重写 Scan 方法，读取时自动转换
func (f *FormatTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*f = FormatTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
