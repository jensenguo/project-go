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

	nconfig "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jensenguo/project-go/utils/coroutine"
)

// LoadAndWatch 加载nacos配置并监听，注意dataID必须是带格式的文件名，例如xxx.yaml
func (c *client) LoadAndWatch(group, dataID string, typ reflect.Type) error {
	ctx := context.Background()
	// 首次加载配置
	source := nconfig.NewConfigSource(c.configClient, nconfig.WithGroup(group), nconfig.WithDataID(dataID))
	configs, err := source.Load()
	if err != nil {
		return fmt.Errorf("config source load fail, err: %v", err)
	}
	// 保存到本地内存
	if err = c.storeToLocal(group, dataID, typ, configs); err != nil {
		return fmt.Errorf("store to local fail, err: %v", err)
	}
	// 监听配置
	if err = c.watch(ctx, source, group, dataID, typ); err != nil {
		return fmt.Errorf("watch fail, err: %v", err)
	}
	return nil
}

// GetConfig 获取配置信息，返回指向配置的指针*T
func (c *client) GetConfig(group, dataID string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cname := c.getName(group, dataID)
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

func (c *client) watch(ctx context.Context, source kconfig.Source, group, dataID string, typ reflect.Type) error {
	// watch定时器配置
	watcher, err := source.Watch()
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
			if err := c.storeToLocal(group, dataID, typ, configs); err != nil {
				log.Errorf("watch store to local fail, err: %v", err)
				continue
			}
		}
	})
	return nil
}

func (c *client) storeToLocal(group, dataID string, typ reflect.Type, configs []*kconfig.KeyValue) error {
	// 兼容typ是指针的情况，由于reflect.New(T)返回指向T的指针，如果T本身是指针，返回值是**T不符合习惯用法
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// 正常这里只会返回一项
	if len(configs) != 1 {
		return fmt.Errorf("configs size invalid")
	}
	config := configs[0]
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
	c.configs[c.getName(group, dataID)] = storage
	return nil
}

func (c *client) getName(group, dataID string) string {
	return fmt.Sprintf("%s_%s", group, dataID)
}
