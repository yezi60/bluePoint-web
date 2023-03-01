package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

type Registry struct {
	Host string
	Port int
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

// Register 服务注册
func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	// consul的ip与port

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)

	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s", // 官方默认应该是30s和10m
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	return nil
}

func (r *Registry) DeRegister(serviceId string) error {

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	// consul的ip与port

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	err = client.Agent().ServiceDeregister(serviceId)

	return err
}
