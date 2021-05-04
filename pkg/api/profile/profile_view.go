package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/apierrors"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/telegram"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/tiktok"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/youtube"
)

type View struct {
	Model *Model
}

func NewView(model *Model) *View {
	return &View{Model: model}
}

type ProfileRes struct {
	Name         string                  `json:"name"`
	About        string                  `json:"about"`
	Availability int                     `json:"availability"`
	Country      string                  `json:"country"`
	Links        []string                `json:"links"`
	Telegram     []telegram.TelegramInfo `json:"telegram"`
	Tiktok       []tiktok.TiktokInfo     `json:"tiktok"`
	Youtube      []youtube.YoutubeInfo   `json:"youtube"`
}

func (ctrl *View) Get(g *gin.Context) {
	if !authware.RequiredLoginPassed(g) {
		return
	}
	u := authware.GetUserFromContext(g)

	p, err := ctrl.Model.Find(u.Login)
	if err != nil {
		apierrors.NotFound.Send(g)
		return
	}
	links := make([]string, len(p.Links))
	for i, url := range p.Links {
		links[i] = url.String()
	}
	g.JSON(http.StatusOK, gin.H{
		"status": "found",
		"profile": ProfileRes{
			Name:         p.Name,
			About:        p.About,
			Availability: p.Availability,
			Country:      p.Country.Code,
			Telegram:     p.Telegram,
			Tiktok:       p.Tiktok,
			Youtube:      p.Youtube,
			Links:        links,
		},
	})
}
