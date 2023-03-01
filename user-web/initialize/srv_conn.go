package initialize

import (
	"bluepoint_api/user-web/global"
	"bluepoint_api/user-web/proto"
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrvInfo.Name), // 这是consul的请求地址，whoami是服务的名字，tag可以不写
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	) // 配置负载均衡策略，目前只支持轮训策略)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

// 初始化 grpc 服务(旧版)
func InitSrvConn2() {
	// 从注册中心获取到用户服务的信息 - 包括ip和端口号
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	var userSrvHost string
	var userSrvPort int

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(
		fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}

	for _, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	// 拨号连接用户grpc服务器 跨域的问题 -- 后端解决
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		userSrvHost, userSrvPort), grpc.WithInsecure())

	// 1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡可以做

	// 2. 已经事先创立好了连接，不用后续再次进行tcp三次握手

	// 3. 一个连接多个goroutine共用，性能 - 连接池
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
