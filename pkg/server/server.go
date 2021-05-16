package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/landing"
	"github.com/kiselev-nikolay/inflr-be/pkg/api/profile"
	"github.com/kiselev-nikolay/inflr-be/pkg/authware"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
)

const (
	ModeDev = iota
	ModeProduction
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

	repo := &memorystore.MemoryStoreRepo{}
	repo.Connect()

	reactNativePrefix := "/rnai"
	authware.Connect(router, reactNativePrefix+"/auth", repo, key)
	profile.Connect(router, reactNativePrefix+"/profile", repo)
	landing.Connect(router, reactNativePrefix+"/landing", repo)

	return router
}
