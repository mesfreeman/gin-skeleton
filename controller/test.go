package controller

import (
	"errors"
	"gin-skeleton/helper/response"
	"gin-skeleton/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 测试动作
func Test(c *gin.Context) {
	// 开发者可以在这里加上自己的任意的测试代码，但是测试代码不应提交到仓库中！

	result := map[string]string{
		"welcome":  "hello world!",
		"time":     time.Now().Format("2006-01-02 15:04:05"),
		"clientIp": c.ClientIP(),
		"mode":     viper.GetString("app.mode"),
	}

	response.SuccessJSON(result, "success", c)
}

// 添加
func Add(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	age, _ := strconv.Atoi(c.DefaultQuery("age", "0"))

	// 各种检查逻辑（比如：参数检查、关联性检查、内部检查等）

	test, err := model.AddTest(name, age)
	if err != nil {
		response.LogicExceptionJSON("添加失败："+err.Error(), c)
		return
	}
	response.SuccessJSON(gin.H{"id": test.ID}, "添加成功", c)
}

// 删除
func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))

	// 各种检查逻辑（比如：ID大于零、数据是否存在等）

	err := model.DeleteTest(id)
	if err != nil {
		response.LogicExceptionJSON("删除失败："+err.Error(), c)
		return
	}
	response.SuccessJSON(gin.H{"id": id}, "删除成功", c)
}

// 修改
func Modify(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	name := c.DefaultQuery("name", "")

	// 判断是否存在
	_, err := model.FindTestById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("要修改的数据不存在", c)
		return
	}

	// 修改
	if err := model.ModifyTest(id, model.Test{Name: name}); err != nil {
		response.LogicExceptionJSON("修改失败："+err.Error(), c)
		return
	}

	response.SuccessJSON(gin.H{"id": id}, "修改成功", c)
}

// 查询
func View(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))

	// 判断是否存在
	test, err := model.FindTestById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("要查询的数据不存在", c)
		return
	}

	response.SuccessJSON(test, "查询成功", c)
}
