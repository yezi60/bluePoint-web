package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// 单机式的，存放在内存中
var store = base64Captcha.DefaultMemStore

// GetCaptcha生成验证码的handlefunc
func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	// 驱动和存储的地方
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()

	if err != nil {
		zap.S().Errorw("生成验证码错误,:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})

}
