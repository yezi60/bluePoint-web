package initialize

import (
	"bluepoint_api/user-web/global"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
)

func InitTrans(locale string) (err error) {
	// 修改gin框架中的validator引擎属性，实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			// 拿到tag字段中json的第一个
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		// 第一个参数是备用的语言环境。后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT)
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale {
		case "en":
			en_translation.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zh_translation.RegisterDefaultTranslations(v, global.Trans)
		default:
			en_translation.RegisterDefaultTranslations(v, global.Trans)
		}
		return

	}
	return
}
