package memorystore

import (
	"errors"
	"sync"

	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
)

type MemoryStoreRepo struct {
	data  map[string]map[string]irepository.Item
	mutex sync.Mutex
}

func (ms *MemoryStoreRepo) Connect() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.data = make(map[string]map[string]irepository.Item)
}

func (ms *MemoryStoreRepo) Send(collection string, item *irepository.Item) error {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]irepository.Item)
	}
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]irepository.Item)
	}
	defer ms.mutex.Unlock()
	ms.data[collection][item.Key] = *item
	return nil
}

func (ms *MemoryStoreRepo) Find(collection string, item *irepository.Item) error {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]irepository.Item)
	}
	defer ms.mutex.Unlock()
	v, ok := ms.data[collection][item.Key]
	if ok {
		*item = v
		return nil
	}
	return errors.New("Not found")
}

func (ms *MemoryStoreRepo) Delete(collection string, item *irepository.Item) error {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]irepository.Item)
	}
	defer ms.mutex.Unlock()
	delete(ms.data[collection], item.Key)
	return nil
}

func (ms *MemoryStoreRepo) List(collection string) ([]*irepository.Item, error) {
	ms.mutex.Lock()
	if _, ok := ms.data[collection]; !ok {
		ms.data[collection] = make(map[string]irepository.Item)
	}
	defer ms.mutex.Unlock()
	items := make([]*irepository.Item, 0)
	for _, item := range ms.data[collection] {
		items = append(items, &item)
	}
	return items, nil
}
