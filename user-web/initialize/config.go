package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"bluepoint_api/user-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()

	// 设置文件路径
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息：%v", global.NacosConfig)

	// 从nacos中读取配置信息
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NameSpace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfig := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.Host,
			ContextPath: "/nacos",
			Port:        uint64(global.NacosConfig.Port),
			Scheme:      "http",
		},
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})

	if err != nil {
		panic(err)
	}

	configClient, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(configClient), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败：%s", err.Error())
	}

	zap.S().Infof("配置信息：%v", global.ServerConfig)

	// 监听配置文件是否更改
	err = client.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Infof("配置文件更改！")
			zap.S().Infof("dataId: %v group: %v 配置信息：%v", dataId, group, global.ServerConfig)
		},
	})

}

func InitConfigOld() {

	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()

	// 设置文件路径
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息：%v", global.ServerConfig)

	// viper的功能，动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("config file changed :%v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})

}
