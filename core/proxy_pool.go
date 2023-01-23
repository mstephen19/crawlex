package core

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type proxyHandler func(*http.Request) (*url.URL, error)

type Proxy struct {
	raw        string
	handler    proxyHandler
	rating     int
	lock       *sync.Mutex
	httpClient *http.Client
}

func (proxy *Proxy) markBad() int {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()

	proxy.rating++
	return proxy.rating
}

func (proxy *Proxy) markGood() int {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()

	if proxy.rating > 0 {
		proxy.rating--
	}
	return proxy.rating
}

func NewProxy(raw string, timeout time.Duration) (proxy *Proxy, err error) {
	proxyUrl, err := url.Parse(raw)
	if err != nil {
		return
	}
	handler := http.ProxyURL(proxyUrl)
	proxy = &Proxy{
		raw:     raw,
		handler: handler,
		lock:    &sync.Mutex{},
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy: handler,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	return
}

type ProxyPool struct {
	proxies   map[*Proxy]struct{}
	maxRating int
	lock      *sync.Mutex
}

func NewProxyPool(maxRating int, proxies ...*Proxy) *ProxyPool {
	proxyMap := map[*Proxy]struct{}{}
	for _, proxy := range proxies {
		proxyMap[proxy] = struct{}{}
	}

	return &ProxyPool{
		proxies:   proxyMap,
		lock:      &sync.Mutex{},
		maxRating: maxRating,
	}
}

func (pool *ProxyPool) IsEmpty() bool {
	return len(pool.proxies) == 0
}

func (pool *ProxyPool) RandomProxy() (proxy *Proxy) {
	if len(pool.proxies) == 0 {
		return nil
	}

	index := rand.Intn(len(pool.proxies))
	count := 0
	for prx := range pool.proxies {
		if count == index {
			proxy = prx
			return
		}

		count++
	}
	return
}

func (pool *ProxyPool) MarkBad(proxy *Proxy) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	_, ok := pool.proxies[proxy]
	if !ok {
		return
	}

	rating := proxy.markBad()
	if rating >= pool.maxRating {
		delete(pool.proxies, proxy)
	}
}

func (pool *ProxyPool) MarkGood(proxy *Proxy) {
	_, ok := pool.proxies[proxy]
	if !ok {
		return
	}

	proxy.markGood()
}
