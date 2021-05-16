package controllers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/apierrors"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/youtube"
)

func New(model *models.Model) *Ctrl {
	return &Ctrl{Model: model}
}

type Ctrl struct {
	Model *models.Model
}

type NewReq struct {
	Name string `json:"name" bind:"required"`
}

func (ctrl *Ctrl) New(g *gin.Context) {
	if !authware.RequiredLoginPassed(g) {
		return
	}
	u := authware.GetUserFromContext(g)
	body := &NewReq{}
	if err := g.BindJSON(body); err != nil {
		apierrors.MissingRequiredData.Send(g)
		return
	}

	_, err := ctrl.Model.Find(u.Login)
	if err == nil {
		apierrors.AlreadyHave.Send(g)
		return
	}

	p := &models.Profile{}
	p.Bio.Name = body.Name
	err = ctrl.Model.Send(u.Login, p)
	if err != nil {
		apierrors.CannotCreate.Send(g)
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status": "created",
	})
}

type AddYoutubeReq struct {
	Link string `json:"link" bind:"required"`
}

func (ctrl *Ctrl) AddYoutube(g *gin.Context) {
	if !authware.RequiredLoginPassed(g) {
		return
	}
	u := authware.GetUserFromContext(g)
	body := &AddYoutubeReq{}
	if err := g.BindJSON(body); err != nil {
		apierrors.MissingRequiredData.Send(g)
		return
	}

	p, err := ctrl.Model.Find(u.Login)
	if err != nil {
		apierrors.NotFound.Send(g)
		return
	}

	ytl, err := url.Parse(body.Link)
	if err != nil {
		apierrors.WrongData.Send(g)
		return
	}
	ytid, err := youtube.GetYTIDFromLink(ytl)
	if err != nil {
		apierrors.WrongData.Send(g)
		return
	}
	ytInfo, err := youtube.GetInfo(ytid)
	if err != nil {
		apierrors.IntegrationFail.Send(g)
		return
	}
	if p.Youtube == nil {
		p.Youtube = make(map[string]youtube.YoutubeInfo, 1)
		p.Youtube[ytid] = *ytInfo
	} else {
		p.Youtube[ytid] = *ytInfo
	}

	err = ctrl.Model.Send(u.Login, p)
	if err != nil {
		apierrors.CannotCreate.Send(g)
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status": "added",
	})
}
