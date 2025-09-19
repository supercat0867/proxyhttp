# proxyhttp

一个支持代理池的 HTTP 客户端，支持自定义代理获取器，免维护代理池，发送请求时随机从代理池中取代理客户端转发请求。

## 安装

``` bash
go get github.com/supercat0867/proxyhttp
```

## 使用示例

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/supercat0867/proxyhttp"
)

// ExampleFetcher 自定义代理获取器的实现
type ExampleFetcher struct {
	ProxyURLs []string
}

// Fetch 返回指定数量的 ProxyClient
func (f *ExampleFetcher) Fetch(count int) ([]*proxyhttp.ProxyClient, error) {
	var clients []*proxyhttp.ProxyClient

	for i, raw := range f.ProxyURLs {
		if i >= count {
			break
		}
		proxyURL, err := url.Parse(raw)
		if err != nil {
			continue
		}
		client := &proxyhttp.ProxyClient{
			Client: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
				Timeout: 10 * time.Second,
			},
			// 假设 5 分钟有效
			Expiration: time.Now().Add(5 * time.Minute),
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func main() {
	// 自定义代理获取器
	fetcher := &ExampleFetcher{
		ProxyURLs: []string{
			"http://user:pass@127.0.0.1:8080",
			"http://user:pass@127.0.0.1:8081",
		},
	}

	pool := proxyhttp.NewHttpPool(5, fetcher)
	client := proxyhttp.NewClient(pool)

	// 请求示例
	req, err := http.NewRequest("GET", "https://httpbin.org/ip", nil)
	body, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer body.Body.Close()

	resp, err := io.ReadAll(body.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
}
```