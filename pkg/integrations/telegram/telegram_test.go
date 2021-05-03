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
	require.Equal("Тёмная сторона", info.Title)
	require.Equal("Аркадий Морейнис @amoreynis https://www.instagram.com/temnografika/", info.Description)
	require.NotEmpty(info.ImageURL)
	require.EqualValues(84275, info.Subs)
}
