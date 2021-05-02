package telegram_test

import (
	"os"
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/telegram"
	"github.com/stretchr/testify/require"
)

func TestWebPageReader(t *testing.T) {
	require := require.New(t)
	b, _ := os.Open("./test/page.html")
	info, err := telegram.ReadWebPage(b)
	require.NoError(err)
	require.NotEmpty(info)
	require.NotEmpty(info.Image.String())
	require.Empty(info.Image.String())
}
