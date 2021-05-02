package memorystore

import (
	"errors"
	"sync"

	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters"
)

type MemoryStoreRepo struct {
	data  map[string]map[string]repository_adapters.Item
	mutex sync.Mutex
}

func (ms *MemoryStoreRepo) Connect() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.data = make(map[string]map[string]repository_adapters.Item)
}

func (ms *MemoryStoreRepo) Send(collection string, item *repository_adapters.Item) error {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]repository_adapters.Item)
	}
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]repository_adapters.Item)
	}
	defer ms.mutex.Unlock()
	ms.data[collection][item.Key] = *item
	return nil
}

func (ms *MemoryStoreRepo) Find(collection string, item *repository_adapters.Item) error {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]repository_adapters.Item)
	}
	defer ms.mutex.Unlock()
	v, ok := ms.data[collection][item.Key]
	if ok {
		*item = v
		return nil
	}
	return errors.New("Not found")
}

func (ms *MemoryStoreRepo) Delete(collection string, item *repository_adapters.Item) error {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]repository_adapters.Item)
	}
	defer ms.mutex.Unlock()
	delete(ms.data[collection], item.Key)
	return nil
}

func (ms *MemoryStoreRepo) List(collection string) ([]*repository_adapters.Item, error) {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]repository_adapters.Item)
	}
	defer ms.mutex.Unlock()
	items := make([]*repository_adapters.Item, 0)
	for _, item := range ms.data[collection] {
		items = append(items, &item)
	}
	return items, nil
}
