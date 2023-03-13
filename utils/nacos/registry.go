package nacos

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
)

// RegisterServer 服务注册接口
func (c *client) RegisterServer(serviceInstance *registry.ServiceInstance, opts ...nacos.Option) error {
	if err := c.GetRegistrarClient(opts...).Register(context.Background(), serviceInstance); err != nil {
		return fmt.Errorf("register server fail, err: %v.", err)
	}
	return nil
}

// GetRegistrarClient 获取服务注册client
func (c *client) GetRegistrarClient(opts ...nacos.Option) registry.Registrar {
	return nacos.New(c.nameClient, opts...)
}

// GetDiscoveryClient 获取服务发现client
func (c *client) GetDiscoveryClient(opts ...nacos.Option) registry.Discovery {
	return nacos.New(c.nameClient, opts...)
}
