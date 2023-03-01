package api

import (
	"bluepoint_api/user-web/forms"
	"bluepoint_api/user-web/global"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/haoweitech/sms-go-sdk/sms"
	"go.uber.org/zap"
)

// GenerateSmsCode 生成指定长度的验证码
func GenerateSmsCode(width int) string {
	numeric := [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().Unix())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// SendSms 发送验证码 阿里云
func SendSms(ctx *gin.Context) {

	sendMmsForm := forms.SendMmsForm{}

	if err := ctx.ShouldBindJSON(&sendMmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	params := make(map[string]interface{}, 0)

	params["code"] = GenerateSmsCode(6) // 发送时生成的验证码

	client := sms.NewClient()
	client.SetAppId(global.ServerConfig.SmsInfo.ApiKey)        // 平台中的用户id
	client.SetSecretKey(global.ServerConfig.SmsInfo.ApiSecret) // 对应的用户密钥

	request := sms.NewRequest()
	request.SetMethod("sms.message.send")
	request.SetBizContent(sms.TemplateMessage{
		Mobile:     []string{sendMmsForm.Mobile}, // 要发送的手机号
		Type:       0,
		Sign:       "闪速码", // 签名
		TemplateId: "ST_2020101100000005",
		SendTime:   "",
		Params:     params,
	})

	buf, err := client.Execute(request)
	if err != nil {
		zap.S().Errorf("短信服务出现问题: %v", err.Error())
		return
	}
	zap.S().Infof("短信服务：%v", buf)

	// 将验证码保存起来 key-value 手机号:验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	rdb.Set(context.Background(), sendMmsForm.Mobile, params["code"], time.Second*time.Duration(global.ServerConfig.SmsInfo.Expire))

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
