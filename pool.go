package proxyhttp

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// HttpPool 代理池
type HttpPool struct {
	mu      sync.RWMutex
	clients []*ProxyClient // 代理
	max     int            // 最大代理数
	fetcher ProxyFetcher   // 代理获取方法
}

func NewHttpPool(max int, fetcher ProxyFetcher) *HttpPool {
	return &HttpPool{
		max:     max,
		fetcher: fetcher,
	}
}

// removeProxyClient 从代理池中移除指定代理
func (p *HttpPool) removeProxyClient(client *ProxyClient) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, c := range p.clients {
		if c == client {
			p.clients = append(p.clients[:i], p.clients[i+1:]...)
			break
		}
	}
}

// GetProxyClient 从代理池中获取一个代理
func (p *HttpPool) GetProxyClient() (*ProxyClient, error) {
	for {
		p.mu.Lock()

		// 代理池空时去拉取代理
		if len(p.clients) == 0 {
			p.mu.Unlock()
			newClients, err := p.fetcher.Fetch(p.max)
			if err != nil {
				return nil, err
			}
			if len(newClients) == 0 {
				return nil, errors.New("no proxy clients fetched")
			}

			p.mu.Lock()
			p.clients = append(p.clients, newClients...)
			p.mu.Unlock()
			continue
		}

		// 随机获取一个代理
		index := rand.Intn(len(p.clients))
		proxyClient := p.clients[index]

		// 判断代理是否过期
		if proxyClient.Expiration.Before(time.Now()) {
			p.mu.Unlock()
			p.removeProxyClient(proxyClient)
			continue
		}

		p.mu.Unlock()
		return proxyClient, nil
	}
}
