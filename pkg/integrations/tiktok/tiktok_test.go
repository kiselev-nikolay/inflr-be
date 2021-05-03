package tiktok_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/tiktok"
	"github.com/stretchr/testify/require"
)

type tableTest struct {
	value    string
	expected *tiktok.TiktokInfo
}

var testLinks = []tableTest{
	{
		"Саша Пынин (@eho_rw) on TikTok | 319 Likes. 26 Fans. machineand.me Watch the latest video from Саша Пынин (@eho_rw).",
		&tiktok.TiktokInfo{Name: "Саша Пынин", Login: "@eho_rw", Likes: 319, Subs: 26},
	},
	{
		"Painted.zhang (@painted.zhang) on TikTok | 2.5M Likes. 260.4K Fans. Watch the latest video from Painted.zhang (@painted.zhang).",
		&tiktok.TiktokInfo{Name: "Painted.zhang", Login: "@painted.zhang", Likes: 2500000, Subs: 260400},
	},
	{
		"Карина  (@karinakross) on TikTok | 485.7M Likes. 13.4M Fans. Новый трек ТУСЫ 😳🔥👇🏻 Watch the latest video from Карина  (@karinakross).",
		&tiktok.TiktokInfo{Name: "Карина", Login: "@karinakross", Likes: 485700000, Subs: 13400000},
	},
	{
		"Давид Манукян (@dava_m) on TikTok | 247.5M Likes. 10M Fans. 👇🏻 Слушай трек: «Под гитару»💔🧸",
		&tiktok.TiktokInfo{Name: "Давид Манукян", Login: "@dava_m", Likes: 247500000, Subs: 10000000},
	},
	{
		"Wow, Ivleeva?! (@_agentgirl_) on TikTok | 61.5M Likes. 6.2M Fans. 𝖙𝖍𝖎𝖘 𝖎𝖘 𝖌𝖊𝖓𝖎𝖚𝖘 Watch the latest video from Wow, Ivleeva?! (@_agentgirl_).",
		&tiktok.TiktokInfo{Name: "Wow, Ivleeva?!", Login: "@_agentgirl_", Likes: 61500000, Subs: 6200000},
	},
	{
		"Ольга Бузова (@buzova86) on TikTok | 171.6M Likes. 7.5M Fans. Телеведущая, актриса, дизайнер,певица 🎤👸🏻 Так сильно 💔🎵⬇️",
		&tiktok.TiktokInfo{Name: "Ольга Бузова", Login: "@buzova86", Likes: 171600000, Subs: 7500000},
	},
	{
		"Егор Крид  (@egorkreed) on TikTok | 226.8M Likes. 10.6M Fans. ❤️ heART brEak kid 💔 Watch the latest video from Егор Крид  (@egorkreed).",
		&tiktok.TiktokInfo{Name: "Егор Крид", Login: "@egorkreed", Likes: 226800000, Subs: 10600000},
	},
}

func TestWebPageReader(t *testing.T) {
	require := require.New(t)
	b, _ := os.Open("./test/page.html")
	info, err := tiktok.ReadWebPage(b)
	require.NoError(err)
	require.NotEmpty(info)
}

func TestMetaDescriptionReader(t *testing.T) {
	for i, ve := range testLinks {
		t.Run("Link#"+strconv.Itoa(i), func(t *testing.T) {
			require := require.New(t)
			info, err := tiktok.ReadMetaDescription(ve.value)
			require.NoError(err)
			require.EqualValues(info, ve.expected)
		})
	}
}
