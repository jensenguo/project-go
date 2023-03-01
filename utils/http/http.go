// Package http http请求封装
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

// Client http client接口定义
type Client interface {
	// Post 发送POST请求
	Post(ctx context.Context, path string, req, rsp interface{}) error
}

type client struct {
	httpClient *http.Client
	scheme     string
	host       string
}

type option func(c *http.Client) error

// NewClient 新建http client
func NewClient(scheme, host string, opts ...option) Client {
	c := &http.Client{}
	for _, opt := range opts {
		opt(c)
	}
	return &client{
		httpClient: c,
		scheme:     scheme,
		host:       host,
	}
}

// OptionWithTimeout 设置超时时间
func OptionWithTimeout(timeout time.Duration) option {
	return func(c *http.Client) error {
		c.Timeout = timeout
		return nil
	}
}

// Post 发送POST请求
func (c *client) Post(ctx context.Context, path string, req, rsp interface{}) error {
	// 判断rsp必须是指针
	if reflect.TypeOf(rsp).Kind() != reflect.Ptr {
		return fmt.Errorf("rsp must be a pointer")
	}
	breq, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("req is not json struct fail, err: %v", err)
	}
	url := fmt.Sprintf("%s://%s%s", c.scheme, c.host, path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(breq))
	if err != nil {
		return fmt.Errorf("new request with context fail, err: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do http req fail, err: %v", err)
	}
	defer httpRsp.Body.Close()

	bodyRsp, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("read rsp fail, err: %v", err)
	}
	if err := json.Unmarshal(bodyRsp, rsp); err != nil {
		return fmt.Errorf("unmarshal fail, err: %v", err)
	}
	return nil
}
