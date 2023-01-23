package core

import (
	"log"
	"sync"
)

const DefaultMaxConcurrency int = 50

type CrawlerConfig struct {
	MaxConcurrency int
	Handler        HandlerFunc
}

type Crawler struct {
	manager        *RequestManager
	handler        HandlerFunc
	MaxConcurrency int
	active         int
	maxed          chan struct{}
	lock           *sync.Mutex
	running        bool
}

func NewCrawler(config *CrawlerConfig) *Crawler {
	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = DefaultMaxConcurrency
	}

	return &Crawler{
		manager:        NewRequestManager(),
		handler:        config.Handler,
		MaxConcurrency: config.MaxConcurrency,
		maxed:          make(chan struct{}),
		lock:           &sync.Mutex{},
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
	group := sync.WaitGroup{}

	group.Add(1)

	go func() {
		defer group.Done()

		for {
			// If the max concurrency has already been reached, block
			// until the concurrency has gone down.
			if crawler.activeCount() >= crawler.MaxConcurrency {
				<-crawler.maxed
			}

			crawler.incr()
			group.Add(1)
			opts, ok := crawler.manager.Shift()
			if !ok {
				group.Done()
				crawler.decr()
				if crawler.activeCount() == 0 {
					return
				}
				continue
			}

			go func() {
				defer func() {
					group.Done()
					crawler.decr()
				}()

				response, err := MakeRequest(opts)
				defer response.Body.Close()

				crawler.handler(&HandlerContext{
					Options:  opts,
					Response: response,
					Crawler:  crawler,
				}, err)
			}()
		}
	}()

	err := crawler.Enqueue(requests...)
	if err != nil {
		log.Fatal(err)
	}

	group.Wait()
	crawler.running = false
}
