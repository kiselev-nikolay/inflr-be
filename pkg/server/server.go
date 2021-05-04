package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository/memorystore"
	"github.com/kiselev-nikolay/inflr-be/pkg/server/authdom"
	"github.com/kiselev-nikolay/inflr-be/pkg/server/profiledom"
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

	repo := &memorystore.MemoryStoreRepo{}
	repo.Connect()

	reactNativePrefix := "/rnai"
	authdom.Connect(router, reactNativePrefix+"/auth", repo, key)
	profiledom.Connect(router, reactNativePrefix+"/profile", repo)

	return router
}
