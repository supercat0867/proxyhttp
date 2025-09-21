package proxyhttp

import (
	"fmt"
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

// GetPool 获取代理池
func (c *Client) GetPool() *HttpPool {
	return c.pool
}

// Do 使用代理池发请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	proxyClient, err := c.pool.GetProxyClient()
	if err != nil {
		return nil, err
	}

	resp, err := proxyClient.Client.Do(req)
	if err != nil {
		// 出现错误就移除代理
		c.pool.removeProxyClient(proxyClient)
	}
	return resp, err
}

// DoWithRetry 带重试机制的请求
// retries: 最大重试次数
// interval: 每次重试的间隔时间
func (c *Client) DoWithRetry(req *http.Request, retries int, interval time.Duration) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= retries; attempt++ {
		proxyClient, err := c.pool.GetProxyClient()
		if err != nil {
			lastErr = fmt.Errorf("failed to get proxy client: %v", err)
			time.Sleep(interval)
			continue
		}

		// 使用代理发送请求
		resp, err := proxyClient.Client.Do(req)
		if err == nil {
			return resp, nil
		}

		// 请求失败，记录错误并移除该代理
		lastErr = fmt.Errorf("failed to do request: %v", err)
		c.pool.removeProxyClient(proxyClient)

		if attempt < retries {
			time.Sleep(interval)
		}
	}

	return nil, lastErr
}
