package model

// TestModel 测试
type TestModel struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Age       int        `json:"age"`
	CreatedAt FormatTime `json:"created_at"`
	UpdatedAt FormatTime `json:"updated_at"`
}

// TableName 表名
func (TestModel) TableName() string {
	return "test"
}

// NewTestModel 创建
func NewTestModel() *TestModel {
	return &TestModel{}
}
