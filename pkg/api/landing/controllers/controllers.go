package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/models"
)

func New(model *models.Model) *Ctrl {
	return &Ctrl{Model: model}
}

type Ctrl struct {
	Model *models.Model
}

func (c *Ctrl) Create(g *gin.Context) {}
