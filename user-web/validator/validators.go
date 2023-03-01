package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidateMobile 自定义验证手机格式的验证器
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()

	// 使用正则表达式进行判断
	pattern := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	ok, err := regexp.MatchString(pattern, mobile)
	if err != nil {
		return false
	}
	return ok
}
