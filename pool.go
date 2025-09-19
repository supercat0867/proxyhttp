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
	for i, c := range p.clients {
		if c == client {
			p.clients = append(p.clients[:i], p.clients[i+1:]...)
			break
		}
	}
}

// GetProxyClient 从代理池中获取一个代理
func (p *HttpPool) GetProxyClient() (*ProxyClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		if len(p.clients) == 0 {
			newClients, err := p.fetcher.Fetch(p.max)
			if err != nil {
				return nil, err
			}
			if len(newClients) == 0 {
				return nil, errors.New("no proxy clients fetched")
			}
			p.clients = append(p.clients, newClients...)
		}

		index := rand.Intn(len(p.clients))
		proxyClient := p.clients[index]

		if proxyClient.Expiration.Before(time.Now()) {
			p.removeProxyClient(proxyClient)
		} else {
			return proxyClient, nil
		}
	}
}
