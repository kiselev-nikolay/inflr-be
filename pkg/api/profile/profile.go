package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/controllers"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/views"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

func Connect(router *gin.Engine, prefix string, repo repository.Repo) {
	model := models.New(repo)
	view := views.New(model)
	ctrl := controllers.New(model)

	router.GET(prefix+"/get", view.Get)
	router.POST(prefix+"/new", ctrl.Create)
	router.POST(prefix+"/add-yt", ctrl.AddYoutube)
}
