package proxyhttp

import (
	"net/http"
	"time"
)

// ProxyClient 代理客户端
type ProxyClient struct {
	Client     *http.Client
	Expiration time.Time
}

// ProxyFetcher 代理获取器接口，支持自定义实现
type ProxyFetcher interface {
	// Fetch 返回指定数量的 ProxyClient
	Fetch(count int) ([]*ProxyClient, error)
}
