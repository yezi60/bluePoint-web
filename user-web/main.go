package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"bluepoint_api/user-web/global"
	"bluepoint_api/user-web/initialize"
	"bluepoint_api/user-web/utils"
	"bluepoint_api/user-web/utils/register/consul"
	myVaildtor "bluepoint_api/user-web/validator"
)

func main() {
	// 1.初始化logger
	initialize.InitLogger()

	// 2.初始化配置文件
	initialize.InitConfig()

	// 3.初始化router
	router := initialize.Routers()

	// 4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Panicf("initialize.InitTrans failed:%v", err.Error())
		//panic(err)
	}

	// 5. 初始化srv的链接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	// 如果是本地开发环境，端口号固定
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// 6. 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 自定义错误提示
		err := v.RegisterValidation("mobile", myVaildtor.ValidateMobile)

		if err != nil {
			zap.S().Panicf("RegisterValidation failed:%v", err.Error())
		}

		// 自定义翻译
		err = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})

		if err != nil {
			zap.S().Panicf("RegisterTranslation failed:%v", err.Error())
		}

	}

	// 7. 服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	serviceId := uuid.NewV4().String()
	err := register_client.Register(utils.GeIP(), global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败", err.Error())
	}

	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的loggar
		2. 日志是分级别的，debug，info，warn，error，fatal
		3. S函数和L函数很有用，提供了一个全局的安全访问logger的途径
	*/
	zap.S().Infof("启动服务器,端口：%d", global.ServerConfig.Port)

	go func() {
		if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Panic("服务注销失败：", err.Error())
	} else {
		zap.S().Info("服务注销成功：")
	}
}
