package global

import (
	"bluepoint_api/user-web/config"
	"bluepoint_api/user-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	// 全局配置
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	// nacos配置
	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	// 翻译器
	Trans ut.Translator

	// grpc client
	UserSrvClient proto.UserClient
)
