package views

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/models"
)

type View struct {
	Model *models.Model
}

func New(model *models.Model) *View {
	return &View{Model: model}
}

func (c *View) Get(g *gin.Context) {
	// Todo render template
}
