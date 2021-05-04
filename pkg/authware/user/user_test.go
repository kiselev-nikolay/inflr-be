package user_test

import (
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/stretchr/testify/require"
)

func TestUserModel(t *testing.T) {
	require := require.New(t)
	repo := memorystore.MemoryStoreRepo{}
	repo.Connect()
	um := user.NewUserModel(&repo)

	err := um.Send("test", &user.User{
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
