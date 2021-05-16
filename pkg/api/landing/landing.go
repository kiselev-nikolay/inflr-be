package landing

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/controllers"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing/views"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

func Connect(router *gin.Engine, prefix string, repo repository.Repo) {
	model := models.New(repo)
	view := views.New(model)
	ctrl := controllers.New(model)

	router.GET("/l/:user/:landing", view.Get)
	router.POST(prefix+"/new", ctrl.Create)
	router.POST(prefix+"/add-card", ctrl.AddCard)
}
