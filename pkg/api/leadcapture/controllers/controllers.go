package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/apierrors"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/leadcapture/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
)

func New(model *models.Model) *Ctrl {
	return &Ctrl{Model: model}
}

type Ctrl struct {
	Model *models.Model
}

type CreateReq struct {
	LandingKey string            `json:"landingKey" bind:"required"`
	Title      string            `json:"title" bind:"required"`
	Email      string            `json:"email" bind:"required"`
	Phone      string            `json:"phone" bind:"required"`
	Forms      map[string]string `json:"forms" bind:"required"`
}

func (c *Ctrl) Create(g *gin.Context) {
	if !authware.RequiredLoginPassed(g) {
		return
	}
	u := authware.GetUserFromContext(g)
	body := &CreateReq{}
	if err := g.BindJSON(body); err != nil {
		apierrors.MissingRequiredData.Send(g)
		return
	}

	_, findErr := c.Model.Find(u, body.Key)
	if findErr == nil {
		apierrors.AlreadyHave.Send(g)
		return
	}

	newLanding := &models.Landing{
		Key:   body.Key,
		Title: body.Title,
		Cards: []models.Card{},
	}

	sendErr := c.Model.Send(u, body.Key, newLanding)
	if sendErr != nil {
		apierrors.CannotCreate.Send(g)
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status": "created",
	})
}
