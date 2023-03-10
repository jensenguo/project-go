// Package registry 服务注册中心
package registry

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	// nacosHost 注册中心地址
	nacosHost = flag.String("nacos_host", "127.0.0.1", "nacos server host")
	// nacosHost 注册中心端口号
	nacosPort = flag.Uint64("nacos_port", 8848, "nacos server port")
)

// DefaultRegistryClient 默认注册中心client对象
var DefaultRegistryClient = newRegistryClient()

// newRegistryClient 新建registry client实例
func newRegistryClient() *nacos.Registry {
	flag.Parse()
	defaultClientConfig := constant.NewClientConfig(
		constant.WithTimeoutMs(5000),
		constant.WithLogLevel("debug"),
	)
	defaultServerConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			*nacosHost,
			*nacosPort,
			constant.WithContextPath("/nacos"),
		),
	}
	nc, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  defaultClientConfig,
			ServerConfigs: defaultServerConfigs,
		},
	)
	if err != nil {
		panic(fmt.Errorf("new name client failed, err: %v.", err)) // 新建nacos失败，服务直接启动失败
	}
	return nacos.New(nc)
}

// RegisterServer 服务注册接口
func RegisterServer(serviceInstance *registry.ServiceInstance) error {
	if err := DefaultRegistryClient.Register(context.Background(), serviceInstance); err != nil {
		return fmt.Errorf("register failed, err: %v.", err)
	}
	return nil
}
