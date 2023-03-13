package nacos

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"

	_ "github.com/go-kratos/kratos/v2/encoding/form"
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
	_ "github.com/go-kratos/kratos/v2/encoding/yaml"

	config "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	sourceconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jensenguo/project-go/utils/coroutine"
)

// LoadAndWatch 加载nacos配置并监听
func (c *client) LoadAndWatch(group, dataID, key string, typ reflect.Type) error {
	ctx := context.Background()
	// 首次加载配置
	configSource := config.NewConfigSource(c.configClient, config.WithGroup(group), config.WithDataID(dataID))
	configs, err := configSource.Load()
	if err != nil {
		return fmt.Errorf("config source load fail, err: %v", err)
	}
	// 保存到本地内存
	if err = c.storeToLocal(group, dataID, key, typ, configs); err != nil {
		return fmt.Errorf("store to local fail, err: %v", err)
	}
	// 监听配置
	if err = c.watch(ctx, configSource, group, dataID, key, typ); err != nil {
		return fmt.Errorf("watch fail, err: %v", err)
	}
	return nil
}

// Get 获取配置信息，返回指向配置的指针*T
func (c *client) Get(group, dataID, key string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cname := c.getName(group, dataID, key)
	storage, ok := c.configs[cname]
	if !ok {
		return nil, fmt.Errorf("find config %s fail", cname)
	}
	val, ok := storage.Load().(reflect.Value)
	if !ok {
		return nil, fmt.Errorf("storage type invalid")
	}
	return val.Interface(), nil
}

func (c *client) watch(ctx context.Context, sc sourceconfig.Source, group, dataID, key string, typ reflect.Type) error {
	// watch定时器配置
	watcher, err := sc.Watch()
	if err != nil {
		return fmt.Errorf("watch fail, err: %v", err)
	}
	coroutine.Go(func() {
		for {
			configs, err := watcher.Next()
			if err != nil {
				log.Errorf("watch store to local fail, err: %v", err)
				continue // 不返回错误
			}
			if err := c.storeToLocal(group, dataID, key, typ, configs); err != nil {
				log.Errorf("watch store to local fail, err: %v", err)
				continue
			}
		}
	})
	return nil
}

func (c *client) storeToLocal(group, dataID, key string, typ reflect.Type, configs []*sourceconfig.KeyValue) error {
	// 兼容typ是指针的情况，由于reflect.New(T)返回指向T的指针，如果T本身是指针，返回值是**T不符合习惯用法
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	config, err := c.findConfig(configs, key)
	if err != nil {
		return err
	}
	codec := encoding.GetCodec(config.Format)
	if codec == nil {
		return fmt.Errorf("unsuport config type %s fail", config.Format)
	}
	val := reflect.New(typ)
	if err := codec.Unmarshal(config.Value, val.Interface()); err != nil {
		return fmt.Errorf("unmarshal %v fail, err: %v", typ.Kind(), err)
	}
	storage := &atomic.Value{}
	storage.Store(val)

	c.lock.Lock()
	defer c.lock.Unlock()
	c.configs[c.getName(group, dataID, key)] = storage
	return nil
}

func (c *client) getName(group, dataID, key string) string {
	return fmt.Sprintf("%s_%s_%s", group, dataID, key)
}

func (c *client) findConfig(configs []*sourceconfig.KeyValue, key string) (*sourceconfig.KeyValue, error) {
	for _, config := range configs {
		log.Errorf("config %+v", config)
		if config.Key == key {
			return config, nil
		}
	}
	return nil, fmt.Errorf("find key %s config fail", key)
}
