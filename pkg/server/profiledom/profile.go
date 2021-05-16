package profiledom

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/controllers"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/models"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile/views"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

func Connect(router *gin.Engine, prefix string, repo repository.Repo) {
	model := models.New(repo)
	ctrl := controllers.New(model)
	view := views.New(model)

	router.POST(prefix+"/new", ctrl.New)
	router.POST(prefix+"/add-yt", ctrl.AddYoutube)
	router.GET(prefix+"/get", view.Get)
}
