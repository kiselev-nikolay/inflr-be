package profiledom

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

func Connect(router *gin.Engine, prefix string, repo repository.Repo) {
	model := profile.NewModel(repo)
	ctrl := profile.NewController(model)
	view := profile.NewView(model)

	router.POST(prefix+"/new", ctrl.New)
	router.POST(prefix+"/add-yt", ctrl.AddYoutube)
	router.GET(prefix+"/get", view.Get)
}
