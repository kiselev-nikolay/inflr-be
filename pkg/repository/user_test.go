package repository_test

import (
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters/memorystore"
	"github.com/stretchr/testify/require"
)

func TestUserModel(t *testing.T) {
	require := require.New(t)
	repo := memorystore.MemoryStoreRepo{}
	repo.Connect()
	um := repository.NewUserModel(&repo)

	err := um.Send("test", &repository.User{
		Login: "test",
	})
	require.NoError(err)
	user, err := um.Find("test")
	require.NoError(err)
	require.Equal("test", user.Login)

	users, err := um.List()
	require.NoError(err)
	require.Len(users, 1)
	require.Equal("test", users[0].Login)

	um.Delete("test")
	users, err = um.List()
	require.NoError(err)
	require.Empty(users)
}
