package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/passwords"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository_adapters/memorystore"
)

const (
	ModeDev        = iota
	ModeProduction = iota
)

const key = "Kh4Hy=bKRZ^fkq!RE7P8cBx=KLAb#nU^4Es$7srGHdH8@g79q2"

func GetRouter(mode int) *gin.Engine {
	switch mode {
	case ModeDev:
		gin.SetMode(gin.DebugMode)
	case ModeProduction:
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	repo := memorystore.MemoryStoreRepo{}
	repo.Connect()

	um := repository.NewUserModel(&repo)

	pw := passwords.Passworder{KeySecret: []byte(key)}

	c := &authware.Config{
		Key:        []byte(key),
		UserModel:  *um,
		Passworder: pw,
	}
	router.Use(authware.NewAuthware(c))
	router.POST("/token", authware.NewTokenHandler(c))
	router.POST("/register", authware.NewRegisterHandler(c))
	router.GET("/test", func(g *gin.Context) {
		u := authware.GetUserFromContext(g)
		if u == nil {
			g.Status(http.StatusUnauthorized)
			return
		}
		g.Status(http.StatusOK)
	})

	return router
}
