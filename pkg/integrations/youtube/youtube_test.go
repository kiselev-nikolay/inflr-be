package youtube_test

import (
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/youtube"
	"github.com/stretchr/testify/require"
)

type tableTest struct {
	value    string
	expected string
}

var testLinks = []tableTest{
	{"https://www.youtube.com/channel/UCwkTfp14Sj7o6q9_8ADJpnA", "UCwkTfp14Sj7o6q9_8ADJpnA"},
	{"https://www.youtube.com/channel/UC-lHJZR3Gqxm24_Vd_AJ5Yw", "UC-lHJZR3Gqxm24_Vd_AJ5Yw"},
	{"https://www.youtube.com/user/PewDiePie", ""},
	{"https://www.youtube.com/user/Vsauce", ""},
	{"https://www.youtube.com/c/vsauce1", ""},
	{"https://www.youtube.com/c/%D0%A1%D0%BA%D1%80%D1%8B%D1%82%D1%8B%D0%B9%D1%81%D0%BC%D1%8B%D1%81%D0%BB", ""},
}

func TestGetYTIDFromLink(t *testing.T) {
	for i, ve := range testLinks {
		t.Run("Link#"+strconv.Itoa(i), func(t *testing.T) {
			ytLink, _ := url.Parse(ve.value)
			value, err := youtube.GetYTIDFromLink(ytLink)
			if ve.expected == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, ve.expected, value)
		})
	}
}

func TestYTResponseReader(t *testing.T) {
	require := require.New(t)
	b, _ := os.Open("./test/res.json")
	info, err := youtube.ReadYTResponse(b)
	require.NoError(err)
	require.NotEmpty(info)
	require.Equal("PewDiePie", info.Title)
	require.Equal("https://yt3.ggpht.com/ytc/AAUvwnga3eXKkQgGU-3j1_jccZ0K9m6MbjepV0ksd7eBEw=s800-c-k-c0x00ffffff-no-rj", info.ImageURL)
	require.EqualValues(110000000, info.Subs)
	require.EqualValues(27225286026, info.Views)
	require.EqualValues(4318, info.Videos)
}
