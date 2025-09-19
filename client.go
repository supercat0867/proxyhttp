package proxyhttp

import (
	"net/http"
	"time"
)

// Client 代理客户端
type Client struct {
	pool *HttpPool
}

// NewClient 创建代理客户端
func NewClient(pool *HttpPool) *Client {
	return &Client{pool: pool}
}

// Do 使用代理池发请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	proxyClient, err := c.pool.GetProxyClient()
	if err != nil {
		return nil, err
	}
	return proxyClient.Client.Do(req)
}

// DoWithRetry 带重试机制的请求
// retries: 最大重试次数
// interval: 每次重试的间隔时间
func (c *Client) DoWithRetry(req *http.Request, retries int, interval time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= retries; attempt++ {
		proxyClient, getErr := c.pool.GetProxyClient()
		if getErr != nil {
			err = getErr
		} else {
			resp, err = proxyClient.Client.Do(req)
			if err == nil {
				return resp, nil
			}
		}

		if attempt < retries {
			time.Sleep(interval)
		}
	}

	return resp, err
}
