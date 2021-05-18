package model

import (
	"time"
)

type TestModel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 创建
func NewTestModel() *TestModel {
	return &TestModel{}
}
