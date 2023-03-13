// Package nacos 服务注册中心、配置中心
package nacos

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// envNacosHostPort 存放nacos地址环境变量
var envNacosHostPort = "TEM_REGISTRY_ADDRESS"

type client struct {
	nameClient   naming_client.INamingClient
	configClient config_client.IConfigClient
	configs      map[string]*atomic.Value // 配置名和配置缓存对应关系
	lock         sync.RWMutex
}

// DefaultClient 默认nacos client，务必设置环境变量TEM_REGISTRY_ADDRESS，否则适用会panic！
var DefaultClient = func() *client {
	nacosHostPort := os.Getenv(envNacosHostPort)
	client, err := NewClient(nacosHostPort, nil, nil)
	if err != nil {
		log.Errorf("new default nacos client fail, err: %v", err)
		return nil // 默认client初始化失败不返回错误，
	}
	return client
}()

// NewClient 新建nacos client实例
func NewClient(hostport string, serverOpts []constant.ServerOption,
	clientOpts []constant.ClientOption) (*client, error) {
	ip, strPort, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, fmt.Errorf("split host port fail, err: %v.", err)
	}
	port, err := strconv.ParseUint(strPort, 10, 64)
	if err != nil {
		return nil, err
	}
	defaultServerOpts := []constant.ServerOption{
		constant.WithContextPath("/nacos"),
	}
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, append(defaultServerOpts, serverOpts...)...),
	}
	defaultCLientOpts := []constant.ClientOption{
		constant.WithTimeoutMs(5000),
		constant.WithLogLevel("debug"),
	}
	clientConfig := constant.NewClientConfig(append(defaultCLientOpts, clientOpts...)...)
	nc, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new name client fail, err: %v.", err)
	}
	cc, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new config client fail, err: %v.", err)
	}
	return &client{
		nameClient:   nc,
		configClient: cc,
		configs:      make(map[string]*atomic.Value),
	}, nil
}
