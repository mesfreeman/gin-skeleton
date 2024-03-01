package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhts "github.com/go-playground/validator/v10/translations/zh"
)

// 翻译器
var trans ut.Translator

// InitTranslator 初始化校验翻译器
func InitTranslator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			// 读取请求字段的中文名，更好地兼容中文提示
			remark := strings.SplitN(field.Tag.Get("remark"), ",", 2)[0]
			if remark != "" && remark != "-" {
				return remark
			}

			// 找不到中文备注的话，使用请求参数的字段名
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}

			return name
		})

		zhObj := zh.New()
		trans, _ = ut.New(zhObj, zhObj).GetTranslator("zh")
		_ = zhts.RegisterDefaultTranslations(v, trans)
	}
}

// Translate 翻译校验器错误信息为中文
func Translate(err error) string {
	// 如果是校验器的错误，则只返回第一个错误信息
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, err := range errs.Translate(trans) {
			return err
		}
	}

	// 其它类型错误，原样返回
	return err.Error()
}
