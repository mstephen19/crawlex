package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Store[T any] interface {
	Push(...T) error
	Stop()
}

type basicStore[T any] struct {
	items   []T
	ticker  *time.Ticker
	lock    *sync.Mutex
	group   *sync.WaitGroup
	stopped bool
}

func ensureStorageFolder() {
	if _, err := os.Stat("./storage"); !os.IsNotExist(err) {
		return
	}

	os.Mkdir("./storage", os.ModePerm)
}

func addItemToStorage[T any](name string, item T) (err error) {
	bytes, err := json.Marshal(item)
	if err != nil {
		return
	}

	os.WriteFile(fmt.Sprintf(`./storage/%s.json`, name), bytes, 0644)
	return
}

func NewBasicStore[T any](interval int64) (store Store[T]) {
	store = &basicStore[T]{
		items:  make([]T, 0),
		ticker: time.NewTicker(time.Duration(interval)),
		lock:   &sync.Mutex{},
		group:  &sync.WaitGroup{},
	}

	store.(*basicStore[T]).start()

	return
}

func (store *basicStore[T]) Wait() {
	store.group.Wait()
}

func (store *basicStore[T]) Stop() {
	store.stopped = true
	store.ticker.Stop()
	store.group.Wait()
}

func (store *basicStore[T]) start() {
	ensureStorageFolder()

	store.group.Add(1)
	go func() {
		defer func() {
			store.group.Done()
		}()

	loop:
		for {
			if !store.stopped {
				continue loop
			}
			<-store.ticker.C

			store.lock.Lock()
			items := store.items
			store.items = make([]T, 0)
			store.group.Add(len(items))
			store.lock.Unlock()

			for _, item := range items {
				go func(name string, item T) {
					defer store.group.Done()
					addItemToStorage(name, item)
				}(fmt.Sprint(rand.Intn(int(time.Now().UnixMilli()))), item)
			}

			if store.stopped {
				return
			}
		}
	}()
}

func (store *basicStore[T]) Push(items ...T) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.items = append(store.items, items...)
	return nil
}
