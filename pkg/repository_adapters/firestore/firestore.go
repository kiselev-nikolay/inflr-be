package firestore

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FireStoreRepo struct {
	firestoreClient *firestore.Client
	caches          map[string]map[string]repository_adapters.Item
	mutex           sync.Mutex
	cancelChannel   chan struct{}
}

type FireStoreRepoConf struct {
	ProjectID, CredentialsPath string
}

func (fs *FireStoreRepo) Connect(conf FireStoreRepoConf) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, conf.ProjectID, option.WithCredentialsFile(conf.CredentialsPath))
	if err != nil {
		return err
	}
	fs.caches = make(map[string]map[string]repository_adapters.Item)
	fs.firestoreClient = client
	return nil
}

func (fs *FireStoreRepo) ResetCache() {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.caches = make(map[string]map[string]repository_adapters.Item)
}

func (fs *FireStoreRepo) Quit() {
	fs.firestoreClient.Close()
	close(fs.cancelChannel)
}

func (fs *FireStoreRepo) Send(collection string, item *repository_adapters.Item) error {
	docRef := fs.firestoreClient.Collection(collection).Doc(item.Key)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := docRef.Set(ctx, item)
	if err != nil {
		return err
	}
	fs.mutex.Lock()
	col, ok := fs.caches[collection]
	if !ok {
		col = make(map[string]repository_adapters.Item)
		fs.caches[collection] = col
	}
	col[item.Key] = *item
	fs.mutex.Unlock()
	return nil
}

func (fs *FireStoreRepo) Find(collection string, item *repository_adapters.Item) error {
	fs.mutex.Lock()
	col, ok := fs.caches[collection]
	if !ok {
		col = make(map[string]repository_adapters.Item)
		fs.caches[collection] = col
	}
	cachedValue, ok := col[item.Key]
	fs.mutex.Unlock()
	*item = cachedValue
	if ok {
		return nil
	}
	docRef := fs.firestoreClient.Collection(collection).Doc(item.Key)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	doc, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	err = doc.DataTo(item)
	if err != nil {
		return err
	}
	fs.mutex.Lock()
	fs.caches[collection][item.Key] = *item
	fs.mutex.Unlock()
	return nil
}

func (fs *FireStoreRepo) Delete(collection string, item *repository_adapters.Item) error {
	docRef := fs.firestoreClient.Collection(collection).Doc(item.Key)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}
	fs.mutex.Lock()
	col, ok := fs.caches[collection]
	if !ok {
		col = make(map[string]repository_adapters.Item)
		fs.caches[collection] = col
	}
	delete(col, item.Key)
	fs.mutex.Unlock()
	return nil
}

func (fs *FireStoreRepo) List(collection string) ([]*repository_adapters.Item, error) {
	items := make([]*repository_adapters.Item, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	documentIterator := fs.firestoreClient.Collection(collection).Documents(ctx)
	documents, err := documentIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		item := &repository_adapters.Item{}
		err := document.DataTo(item)
		if err != nil {
			log.Println("Broken node:", document.Ref.Path)
		}
		items = append(items, item)
	}
	return items, nil
}

func IsNotFoundError(err error) bool {
	return status.Code(err) == codes.NotFound
}
