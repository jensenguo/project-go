// Package http http请求封装
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	Post(ctx context.Context, url string, req, rsp interface{}) error
}

type client struct {
	httpClient *http.Client
}

type option func(c *http.Client) error

// NewClient 新建http client
func NewClient(address string, opts ...option) Client {
	c := &http.Client{}
	for _, opt := range opts {
		opt(c)
	}
	return &client{httpClient: c}
}

func OptionWithTimeout(timeout time.Duration) option {
	return func(c *http.Client) error {
		c.Timeout = timeout
		return nil
	}
}

// Post 发送POST请求，req必须是json结构体
func (c *client) Post(ctx context.Context, url string, req, rsp interface{}) error {
	// 判断req，rsp必须是指针
	breq, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("req is not json struct fail, err: %v", err)
	}
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
