package passwords_test

import (
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"

	"github.com/stretchr/testify/require"
)

func TestPassworder(t *testing.T) {
	require := require.New(t)
	pw := passwords.Passworder{
		KeySecret: []byte("test"),
		MinLen:    3,
	}
	_, err := pw.Hash("123")
	require.Error(err)
	require.Equal(passwords.ErrTooShort, err)

	h, err := pw.Hash("meow")
	require.NoError(err)
	meowIsCorrect, err := pw.IsCorrect(h, "meow")
	require.NoError(err)
	require.True(meowIsCorrect)
	catIsCorrect, err := pw.IsCorrect(h, "cat")
	require.NoError(err)
	require.False(catIsCorrect)
	dogIsCorrect, err := pw.IsCorrect(h, "dog")
	require.NoError(err)
	require.False(dogIsCorrect)

	pw.KeySecret = []byte("test2")
	meowIsCorrect, err = pw.IsCorrect(h, "meow")
	require.NoError(err)
	require.False(meowIsCorrect)
}
