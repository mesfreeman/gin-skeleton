package model

import (
	"gin-skeleton/helper"
	"time"
)

type Test struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddTest 添加测试数据
func AddTest(name string, age int) (test Test, err error) {
	test = Test{
		Name: name,
		Age:  age,
	}
	err = helper.GormDefaultDb.Create(&test).Error
	return
}

// DeleteTest 删除测试数据
func DeleteTest(id int) (err error) {
	var test Test
	err = helper.GormDefaultDb.Where("id = ?", id).Delete(&test).Error
	return
}

// ModifyTest
func ModifyTest(id int, test Test) (err error) {
	err = helper.GormDefaultDb.Where("id = ?", id).Updates(test).Error
	return
}

// FindTestById 查询测试数据
func FindTestById(id int) (test Test, err error) {
	err = helper.GormDefaultDb.Where("id = ?", id).First(&test).Error
	return
}
