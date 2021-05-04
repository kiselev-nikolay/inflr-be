package authdom

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
)

func Connect(router *gin.Engine, prefix string, repo repository.Repo, key string) {
	pw := passwords.Passworder{KeySecret: []byte(key)}
	um := user.NewUserModel(repo)
	c := &authware.Config{
		Key:        []byte(key),
		UserModel:  *um,
		Passworder: pw,
	}
	router.Use(authware.NewAuthware(c))

	router.POST(prefix+"/token", authware.NewTokenHandler(c))
	router.POST(prefix+"/new", authware.NewRegisterHandler(c))
	router.GET(prefix+"/test", testHandler)
}

func testHandler(g *gin.Context) {
	u := authware.GetUserFromContext(g)
	if u == nil {
		g.Status(http.StatusUnauthorized)
		return
	}
	g.Status(http.StatusOK)
}
