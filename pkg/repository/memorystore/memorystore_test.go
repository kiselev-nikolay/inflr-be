package memorystore_test

import (
	"testing"

	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"

	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/stretchr/testify/require"
)

func TestMemstorage(t *testing.T) {
	require := require.New(t)
	repo := memorystore.MemoryStoreRepo{}
	repo.Connect()
	sentItem := &irepository.Item{
		Key:   "1",
		Value: "hello",
	}
	repo.Send("test", sentItem)
	foundItem := &irepository.Item{
		Key: "1",
	}
	repo.Find("test", foundItem)
	require.Equal(sentItem.Value, foundItem.Value)

	repo.Send("test", sentItem)
	testItems, _ := repo.List("test")
	require.Len(testItems, 1)

	anotherSentItem := &irepository.Item{
		Key:   "2",
		Value: "hello",
	}
	repo.Send("test", anotherSentItem)
	testItems, _ = repo.List("test")
	require.Len(testItems, 2)

	repo.Delete("test", sentItem)
	testItems, _ = repo.List("test")
	require.Len(testItems, 1)
	require.EqualValues(anotherSentItem, testItems[0])
}
