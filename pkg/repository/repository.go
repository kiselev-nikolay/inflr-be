package repository

import (
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"

	"github.com/kiselev-nikolay/inflr-be/pkg/repository/firestore"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
)

type Repo interface {
	Send(collection string, item *irepository.Item) error
	Find(collection string, item *irepository.Item) error
	Delete(collection string, item *irepository.Item) error
	List(collection string) ([]*irepository.Item, error)
}

func init() {
	var _ Repo = (*firestore.FireStoreRepo)(nil)
	var _ Repo = (*memorystore.MemoryStoreRepo)(nil)
}
