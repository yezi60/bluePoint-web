package config

type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key" yaml:"key"`
}

type MessageConfig struct {
	ApiKey    string `mapstructure:"key" json:"key" yaml:"key"`
	ApiSecret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Expire    int    `mapstructure:"expire" json:"expire" yaml:"expire"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name" yaml:"name"`
	Port        int           `mapstructure:"port" json:"port" yaml:"port"`
	Host        string        `mapstructure:"host" json:"host" yaml:"host"`
	Tags        []string      `mapstructure:"tags" json:"tags" yaml:"tags"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv" yaml:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	SmsInfo     MessageConfig `mapstructure:"sms" json:"sms" yaml:"sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis" yaml:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul" yaml:"consul"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	NameSpace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}
