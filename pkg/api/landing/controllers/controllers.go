package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/apierrors"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
)

func New(model *models.Model) *Ctrl {
	return &Ctrl{Model: model}
}

type Ctrl struct {
	Model *models.Model
}

type CreateReq struct {
	Key   string `json:"key" bind:"required"`
	Title string `json:"title" bind:"required"`
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

type AddCardReq struct {
	LandingKey string `json:"landingKey" bind:"required"`
	Title      string `json:"title" bind:"required"`
	Text       string `json:"text" bind:"required"`
	Link       string `json:"link" bind:"required"`
}

func (c *Ctrl) AddCard(g *gin.Context) {
	if !authware.RequiredLoginPassed(g) {
		return
	}
	u := authware.GetUserFromContext(g)
	body := &AddCardReq{}
	if err := g.BindJSON(body); err != nil {
		apierrors.MissingRequiredData.Send(g)
		return
	}

	landing, findErr := c.Model.Find(u, body.LandingKey)
	if findErr != nil {
		apierrors.NotFound.Send(g)
		return
	}

	if landing.Cards == nil {
		landing.Cards = make([]models.Card, 0, 1)
	}
	landing.Cards = append(landing.Cards, models.Card{
		Title: body.Title,
		Text:  body.Text,
		Link:  body.Link,
	})

	sendErr := c.Model.Send(u, body.LandingKey, landing)
	if sendErr != nil {
		apierrors.CannotCreate.Send(g)
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}
