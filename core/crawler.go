package core

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const DefaultMaxConcurrency int = 50
const DefaultMaxProxyRating int = 5
const DefaultRequestTimeoutSeconds int = 60

type CrawlerConfig struct {
	MaxConcurrency        int
	Handler               HandlerFunc
	Proxies               []string
	MaxProxyRating        int
	RequestTimeoutSeconds int
}

type Crawler struct {
	manager        *RequestManager
	handler        HandlerFunc
	MaxConcurrency int
	active         int
	maxed          chan struct{}
	lock           *sync.Mutex
	running        bool
	group          *sync.WaitGroup
	proxyPool      *ProxyPool
	timeout        int
	defaultClient  *http.Client
}

func NewCrawler(config *CrawlerConfig) *Crawler {
	if config.MaxConcurrency <= 0 {
		config.MaxConcurrency = DefaultMaxConcurrency
	}

	if config.MaxProxyRating <= 0 {
		config.MaxProxyRating = DefaultMaxProxyRating
	}

	if config.RequestTimeoutSeconds <= 0 {
		config.RequestTimeoutSeconds = DefaultRequestTimeoutSeconds
	}

	proxies := make([]*Proxy, len(config.Proxies))
	for i, raw := range config.Proxies {
		proxy, err := NewProxy(raw, time.Second*time.Duration(config.RequestTimeoutSeconds))
		if err != nil {
			log.Fatal(err)
		}

		proxies[i] = proxy
	}

	return &Crawler{
		manager:        NewRequestManager(),
		handler:        config.Handler,
		MaxConcurrency: config.MaxConcurrency,
		timeout:        config.RequestTimeoutSeconds,
		maxed:          make(chan struct{}),
		lock:           &sync.Mutex{},
		group:          &sync.WaitGroup{},
		proxyPool:      NewProxyPool(config.MaxProxyRating, proxies...),
		defaultClient: &http.Client{
			Timeout: time.Second * time.Duration(config.RequestTimeoutSeconds),
		},
	}
}

func (crawler *Crawler) Enqueue(requests ...*RequestOptions) (err error) {
	for _, opts := range requests {
		err = crawler.manager.Push(opts)
	}

	return
}

func (crawler *Crawler) incr() {
	crawler.lock.Lock()
	defer crawler.lock.Unlock()

	crawler.active++
}

func (crawler *Crawler) decr() {
	crawler.lock.Lock()
	defer crawler.lock.Unlock()

	crawler.active--

	// If the current active count is less than the max
	// concurrency, send a message on the channel notifying
	// that another request can go ahead.
	if crawler.active < crawler.MaxConcurrency {
		select {
		case crawler.maxed <- struct{}{}:
		default:
			return
		}
	}
}

func (crawler *Crawler) activeCount() int {
	crawler.lock.Lock()
	defer crawler.lock.Unlock()

	return crawler.active
}

func (crawler *Crawler) Run(requests ...*RequestOptions) {
	if crawler.running {
		log.Fatal(`Cannot have two calls of "Run" occurring at the same time.`)
	}

	crawler.running = true

	crawler.group.Add(1)

	go func() {
		defer crawler.group.Done()

		for {
			// If the max concurrency has already been reached, block
			// until the concurrency has gone down.
			if crawler.activeCount() >= crawler.MaxConcurrency {
				<-crawler.maxed
			}

			crawler.incr()
			opts, ok := crawler.manager.Shift()
			if !ok {
				crawler.decr()
				if crawler.activeCount() == 0 {
					return
				}
				continue
			}

			crawler.group.Add(1)
			go func() {
				defer func() {
					crawler.group.Done()
					crawler.decr()
				}()

				var response *http.Response
				var err error
				var proxy *Proxy

				if !opts.SkipRequest {
					client := crawler.defaultClient
					proxy = crawler.proxyPool.RandomProxy()
					if proxy != nil {
						client = proxy.httpClient
					}
					response, err = MakeRequest(opts, client)
					defer func() {
						if response == nil || response.Body == nil {
							return
						}

						response.Body.Close()
					}()
				}

				crawler.handler(&HandlerContext{
					Options:  opts,
					Response: response,
					crawler:  crawler,
					proxy:    proxy,
				}, err)
			}()
		}
	}()

	err := crawler.Enqueue(requests...)
	if err != nil {
		log.Fatal(err)
	}

	crawler.group.Wait()
	crawler.running = false
}
