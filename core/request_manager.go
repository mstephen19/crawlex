package core

import (
	"sync"
)

type RequestManager struct {
	requests []*RequestOptions
	lock     *sync.Mutex
}

func NewRequestManager() *RequestManager {
	return &RequestManager{
		lock: &sync.Mutex{},
	}
}

func (manager *RequestManager) Push(opts *RequestOptions) (err error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	err = CleanRequestOptions(opts)
	if err != nil {
		return
	}

	manager.requests = append(manager.requests, opts)
	return
}

func (manager *RequestManager) Shift() (item *RequestOptions, ok bool) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	ok = true

	if len(manager.requests) == 0 {
		ok = false
		return
	}

	item = shift(&manager.requests)
	return
}
