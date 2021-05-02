package repository

import (
	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters/firestore"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters/memorystore"
)

type Repo interface {
	Send(collection string, item *repository_adapters.Item) error
	Find(collection string, item *repository_adapters.Item) error
	Delete(collection string, item *repository_adapters.Item) error
	List(collection string) ([]*repository_adapters.Item, error)
}

func init() {
	var _ Repo = (*firestore.FireStoreRepo)(nil)
	var _ Repo = (*memorystore.MemoryStoreRepo)(nil)
}
